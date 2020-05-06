// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package loginid

import (
	"github.com/skygeario/skygear-server/pkg/auth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/audit"
	auth2 "github.com/skygeario/skygear-server/pkg/auth/dependency/auth"
	redis3 "github.com/skygeario/skygear-server/pkg/auth/dependency/auth/redis"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/bearertoken"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/oob"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/password"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/recoverycode"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/totp"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/hook"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/identity/loginid"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/identity/oauth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/interaction"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/interaction/adaptors"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/interaction/flows"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/interaction/redis"
	oauth2 "github.com/skygeario/skygear-server/pkg/auth/dependency/oauth"
	handler2 "github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/handler"
	pq3 "github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/pq"
	redis2 "github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/redis"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oidc"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/passwordhistory/pq"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/session"
	redis4 "github.com/skygeario/skygear-server/pkg/auth/dependency/session/redis"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/urlprefix"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/userprofile"
	"github.com/skygeario/skygear-server/pkg/core/async"
	pq2 "github.com/skygeario/skygear-server/pkg/core/auth/authinfo/pq"
	"github.com/skygeario/skygear-server/pkg/core/db"
	"github.com/skygeario/skygear-server/pkg/core/handler"
	"github.com/skygeario/skygear-server/pkg/core/logging"
	"github.com/skygeario/skygear-server/pkg/core/time"
	"github.com/skygeario/skygear-server/pkg/core/validation"
	"net/http"
)

// Injectors from wire.go:

func newAddLoginIDHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	validator := auth.ProvideValidator(m)
	context := auth.ProvideContext(r)
	requestID := auth.ProvideLoggingRequestID(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	requireAuthz := handler.NewRequireAuthzFactory(factory)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	provider := time.NewProvider()
	store := redis.ProvideStore(context, tenantConfiguration, provider)
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	loginidProvider := loginid.ProvideProvider(sqlBuilder, sqlExecutor, provider, tenantConfiguration, reservedNameChecker)
	oauthProvider := oauth.ProvideProvider(sqlBuilder, sqlExecutor, provider)
	identityAdaptor := &adaptors.IdentityAdaptor{
		LoginID: loginidProvider,
		OAuth:   oauthProvider,
	}
	passwordhistoryStore := pq.ProvidePasswordHistoryStore(provider, sqlBuilder, sqlExecutor)
	passwordChecker := audit.ProvidePasswordChecker(tenantConfiguration, passwordhistoryStore)
	passwordProvider := password.ProvideProvider(sqlBuilder, sqlExecutor, provider, factory, passwordhistoryStore, passwordChecker, tenantConfiguration)
	totpProvider := totp.ProvideProvider(sqlBuilder, sqlExecutor, provider, tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	urlprefixProvider := urlprefix.NewProvider(r)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, txContext, requestID, tenantConfiguration, executor)
	oobProvider := oob.ProvideProvider(tenantConfiguration, sqlBuilder, sqlExecutor, provider, engine, urlprefixProvider, queue)
	bearertokenProvider := bearertoken.ProvideProvider(sqlBuilder, sqlExecutor, provider, tenantConfiguration)
	recoverycodeProvider := recoverycode.ProvideProvider(sqlBuilder, sqlExecutor, provider, tenantConfiguration)
	authenticatorAdaptor := &adaptors.AuthenticatorAdaptor{
		Password:     passwordProvider,
		TOTP:         totpProvider,
		OOBOTP:       oobProvider,
		BearerToken:  bearertokenProvider,
		RecoveryCode: recoverycodeProvider,
	}
	authinfoStore := pq2.ProvideStore(sqlBuilderFactory, sqlExecutor)
	userprofileStore := userprofile.ProvideStore(provider, sqlBuilder, sqlExecutor)
	hookProvider := hook.ProvideHookProvider(context, sqlBuilder, sqlExecutor, requestID, tenantConfiguration, txContext, provider, authinfoStore, userprofileStore, loginidProvider, factory)
	userProvider := interaction.ProvideUserProvider(authinfoStore, userprofileStore, provider, hookProvider, urlprefixProvider, queue, tenantConfiguration)
	interactionProvider := interaction.ProvideProvider(store, provider, factory, identityAdaptor, authenticatorAdaptor, userProvider, oobProvider, tenantConfiguration)
	authorizationStore := &pq3.AuthorizationStore{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	grantStore := redis2.ProvideGrantStore(context, factory, tenantConfiguration, sqlBuilder, sqlExecutor, provider)
	eventStore := redis3.ProvideEventStore(context, tenantConfiguration)
	accessEventProvider := auth2.AccessEventProvider{
		Store: eventStore,
	}
	sessionStore := redis4.ProvideStore(context, tenantConfiguration, provider, factory)
	authAccessEventProvider := &auth2.AccessEventProvider{
		Store: eventStore,
	}
	sessionProvider := session.ProvideSessionProvider(r, sessionStore, authAccessEventProvider, tenantConfiguration)
	idTokenIssuer := oidc.ProvideIDTokenIssuer(tenantConfiguration, urlprefixProvider, authinfoStore, userprofileStore, provider)
	tokenGenerator := _wireTokenGeneratorValue
	tokenHandler := handler2.ProvideTokenHandler(r, tenantConfiguration, factory, authorizationStore, grantStore, grantStore, grantStore, accessEventProvider, sessionProvider, idTokenIssuer, tokenGenerator, provider)
	insecureCookieConfig := auth.ProvideSessionInsecureCookieConfig(m)
	cookieConfiguration := session.ProvideSessionCookieConfiguration(r, insecureCookieConfig, tenantConfiguration)
	userController := flows.ProvideUserController(authinfoStore, userprofileStore, tokenHandler, cookieConfiguration, sessionProvider, hookProvider, provider, tenantConfiguration)
	authAPIFlow := &flows.AuthAPIFlow{
		Interactions:   interactionProvider,
		UserController: userController,
	}
	httpHandler := provideAddLoginIDHandler(validator, requireAuthz, txContext, authAPIFlow)
	return httpHandler
}

var (
	_wireTokenGeneratorValue = handler2.TokenGenerator(oauth2.GenerateToken)
)

// wire.go:

func provideAddLoginIDHandler(
	v *validation.Validator,
	requireAuthz handler.RequireAuthz,
	tx db.TxContext,
	f AddLoginIDInteractionFlow,
) http.Handler {
	h := &AddLoginIDHandler{
		Validator:    v,
		TxContext:    tx,
		Interactions: f,
	}
	return requireAuthz(h, h)
}
