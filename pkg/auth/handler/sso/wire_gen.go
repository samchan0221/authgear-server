// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package sso

import (
	"github.com/gorilla/mux"
	"github.com/skygeario/skygear-server/pkg/auth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/audit"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authn"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/hook"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/loginid"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/mfa"
	pq3 "github.com/skygeario/skygear-server/pkg/auth/dependency/mfa/pq"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/passwordhistory/pq"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/principal"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/principal/oauth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/principal/password"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/session"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/session/redis"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/sso"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/urlprefix"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/userprofile"
	"github.com/skygeario/skygear-server/pkg/core/async"
	auth2 "github.com/skygeario/skygear-server/pkg/core/auth"
	pq2 "github.com/skygeario/skygear-server/pkg/core/auth/authinfo/pq"
	"github.com/skygeario/skygear-server/pkg/core/config"
	"github.com/skygeario/skygear-server/pkg/core/db"
	"github.com/skygeario/skygear-server/pkg/core/handler"
	"github.com/skygeario/skygear-server/pkg/core/logging"
	"github.com/skygeario/skygear-server/pkg/core/mail"
	"github.com/skygeario/skygear-server/pkg/core/sms"
	"github.com/skygeario/skygear-server/pkg/core/time"
	"github.com/skygeario/skygear-server/pkg/core/validation"
	"net/http"
)

// Injectors from wire.go:

func newAuthHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	provider := urlprefix.NewProvider(r)
	authHandlerHTMLProvider := sso.ProvideAuthHandlerHTMLProvider(provider)
	ssoProvider := sso.ProvideSSOProvider(context, tenantConfiguration)
	requestID := auth.ProvideLoggingRequestID(r)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	timeProvider := time.NewProvider()
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	store := pq.ProvidePasswordHistoryStore(timeProvider, sqlBuilder, sqlExecutor)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	passwordProvider := password.ProvidePasswordProvider(sqlBuilder, sqlExecutor, timeProvider, store, factory, tenantConfiguration, reservedNameChecker)
	oauthProvider := oauth.ProvideOAuthProvider(sqlBuilder, sqlExecutor)
	v := auth.ProvidePrincipalProviders(oauthProvider, passwordProvider)
	identityProvider := principal.ProvideIdentityProvider(sqlBuilder, sqlExecutor, v)
	authenticateProcess := authn.ProvideAuthenticateProcess(factory, timeProvider, passwordProvider, oauthProvider, identityProvider)
	passwordChecker := audit.ProvidePasswordChecker(tenantConfiguration, store)
	loginIDChecker := loginid.ProvideLoginIDChecker(tenantConfiguration, reservedNameChecker)
	authinfoStore := pq2.ProvideStore(sqlBuilderFactory, sqlExecutor)
	userprofileStore := userprofile.ProvideStore(timeProvider, sqlBuilder, sqlExecutor)
	contextGetter := auth2.ProvideAuthContextGetter(context)
	hookProvider := hook.ProvideHookProvider(sqlBuilder, sqlExecutor, requestID, tenantConfiguration, provider, contextGetter, txContext, timeProvider, authinfoStore, userprofileStore, passwordProvider, factory)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, txContext, requestID, tenantConfiguration, executor)
	signupProcess := authn.ProvideSignupProcess(passwordChecker, loginIDChecker, identityProvider, passwordProvider, oauthProvider, timeProvider, authinfoStore, userprofileStore, hookProvider, tenantConfiguration, provider, queue)
	oAuthCoordinator := &authn.OAuthCoordinator{
		Authn:  authenticateProcess,
		Signup: signupProcess,
	}
	mfaStore := pq3.ProvideStore(tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider)
	client := sms.ProvideSMSClient(tenantConfiguration)
	sender := mail.ProvideMailSender(tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	mfaSender := mfa.ProvideMFASender(tenantConfiguration, client, sender, engine)
	mfaProvider := mfa.ProvideMFAProvider(mfaStore, tenantConfiguration, timeProvider, mfaSender)
	sessionStore := redis.ProvideStore(context, tenantConfiguration, timeProvider, factory)
	eventStore := redis.ProvideEventStore(context, tenantConfiguration)
	sessionProvider := session.ProvideSessionProvider(r, sessionStore, eventStore, tenantConfiguration)
	authnSessionProvider := authn.ProvideSessionProvider(mfaProvider, sessionProvider, tenantConfiguration, timeProvider, authinfoStore, userprofileStore, identityProvider, hookProvider)
	authnProvider := &authn.Provider{
		OAuth:   oAuthCoordinator,
		Authn:   authenticateProcess,
		Signup:  signupProcess,
		Session: authnSessionProvider,
	}
	loginIDNormalizerFactory := loginid.ProvideLoginIDNormalizerFactory(tenantConfiguration)
	oAuthProviderFactory := sso.ProvideOAuthProviderFactory(tenantConfiguration, provider, timeProvider, loginIDNormalizerFactory)
	oAuthProvider := provideOAuthProviderFromRequestVars(r, oAuthProviderFactory)
	handler := provideAuthHandler(txContext, tenantConfiguration, authHandlerHTMLProvider, ssoProvider, authnProvider, oAuthProvider)
	return handler
}

func newAuthResultHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	contextGetter := auth2.ProvideAuthContextGetter(context)
	requestID := auth.ProvideLoggingRequestID(r)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	requireAuthz := handler.NewRequireAuthzFactory(contextGetter, factory)
	provider := time.NewProvider()
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	store := pq.ProvidePasswordHistoryStore(provider, sqlBuilder, sqlExecutor)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	passwordProvider := password.ProvidePasswordProvider(sqlBuilder, sqlExecutor, provider, store, factory, tenantConfiguration, reservedNameChecker)
	oauthProvider := oauth.ProvideOAuthProvider(sqlBuilder, sqlExecutor)
	v := auth.ProvidePrincipalProviders(oauthProvider, passwordProvider)
	identityProvider := principal.ProvideIdentityProvider(sqlBuilder, sqlExecutor, v)
	authenticateProcess := authn.ProvideAuthenticateProcess(factory, provider, passwordProvider, oauthProvider, identityProvider)
	passwordChecker := audit.ProvidePasswordChecker(tenantConfiguration, store)
	loginIDChecker := loginid.ProvideLoginIDChecker(tenantConfiguration, reservedNameChecker)
	authinfoStore := pq2.ProvideStore(sqlBuilderFactory, sqlExecutor)
	userprofileStore := userprofile.ProvideStore(provider, sqlBuilder, sqlExecutor)
	urlprefixProvider := urlprefix.NewProvider(r)
	hookProvider := hook.ProvideHookProvider(sqlBuilder, sqlExecutor, requestID, tenantConfiguration, urlprefixProvider, contextGetter, txContext, provider, authinfoStore, userprofileStore, passwordProvider, factory)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, txContext, requestID, tenantConfiguration, executor)
	signupProcess := authn.ProvideSignupProcess(passwordChecker, loginIDChecker, identityProvider, passwordProvider, oauthProvider, provider, authinfoStore, userprofileStore, hookProvider, tenantConfiguration, urlprefixProvider, queue)
	oAuthCoordinator := &authn.OAuthCoordinator{
		Authn:  authenticateProcess,
		Signup: signupProcess,
	}
	mfaStore := pq3.ProvideStore(tenantConfiguration, sqlBuilder, sqlExecutor, provider)
	client := sms.ProvideSMSClient(tenantConfiguration)
	sender := mail.ProvideMailSender(tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	mfaSender := mfa.ProvideMFASender(tenantConfiguration, client, sender, engine)
	mfaProvider := mfa.ProvideMFAProvider(mfaStore, tenantConfiguration, provider, mfaSender)
	sessionStore := redis.ProvideStore(context, tenantConfiguration, provider, factory)
	eventStore := redis.ProvideEventStore(context, tenantConfiguration)
	sessionProvider := session.ProvideSessionProvider(r, sessionStore, eventStore, tenantConfiguration)
	authnSessionProvider := authn.ProvideSessionProvider(mfaProvider, sessionProvider, tenantConfiguration, provider, authinfoStore, userprofileStore, identityProvider, hookProvider)
	authnProvider := &authn.Provider{
		OAuth:   oAuthCoordinator,
		Authn:   authenticateProcess,
		Signup:  signupProcess,
		Session: authnSessionProvider,
	}
	validator := auth.ProvideValidator(m)
	ssoProvider := sso.ProvideSSOProvider(context, tenantConfiguration)
	httpHandler := provideAuthResultHandler(txContext, requireAuthz, authnProvider, validator, ssoProvider)
	return httpHandler
}

func newLinkHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	contextGetter := auth2.ProvideAuthContextGetter(context)
	requestID := auth.ProvideLoggingRequestID(r)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	requireAuthz := handler.NewRequireAuthzFactory(contextGetter, factory)
	validator := auth.ProvideValidator(m)
	provider := sso.ProvideSSOProvider(context, tenantConfiguration)
	timeProvider := time.NewProvider()
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	store := pq.ProvidePasswordHistoryStore(timeProvider, sqlBuilder, sqlExecutor)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	passwordProvider := password.ProvidePasswordProvider(sqlBuilder, sqlExecutor, timeProvider, store, factory, tenantConfiguration, reservedNameChecker)
	oauthProvider := oauth.ProvideOAuthProvider(sqlBuilder, sqlExecutor)
	v := auth.ProvidePrincipalProviders(oauthProvider, passwordProvider)
	identityProvider := principal.ProvideIdentityProvider(sqlBuilder, sqlExecutor, v)
	authenticateProcess := authn.ProvideAuthenticateProcess(factory, timeProvider, passwordProvider, oauthProvider, identityProvider)
	passwordChecker := audit.ProvidePasswordChecker(tenantConfiguration, store)
	loginIDChecker := loginid.ProvideLoginIDChecker(tenantConfiguration, reservedNameChecker)
	authinfoStore := pq2.ProvideStore(sqlBuilderFactory, sqlExecutor)
	userprofileStore := userprofile.ProvideStore(timeProvider, sqlBuilder, sqlExecutor)
	urlprefixProvider := urlprefix.NewProvider(r)
	hookProvider := hook.ProvideHookProvider(sqlBuilder, sqlExecutor, requestID, tenantConfiguration, urlprefixProvider, contextGetter, txContext, timeProvider, authinfoStore, userprofileStore, passwordProvider, factory)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, txContext, requestID, tenantConfiguration, executor)
	signupProcess := authn.ProvideSignupProcess(passwordChecker, loginIDChecker, identityProvider, passwordProvider, oauthProvider, timeProvider, authinfoStore, userprofileStore, hookProvider, tenantConfiguration, urlprefixProvider, queue)
	oAuthCoordinator := &authn.OAuthCoordinator{
		Authn:  authenticateProcess,
		Signup: signupProcess,
	}
	mfaStore := pq3.ProvideStore(tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider)
	client := sms.ProvideSMSClient(tenantConfiguration)
	sender := mail.ProvideMailSender(tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	mfaSender := mfa.ProvideMFASender(tenantConfiguration, client, sender, engine)
	mfaProvider := mfa.ProvideMFAProvider(mfaStore, tenantConfiguration, timeProvider, mfaSender)
	sessionStore := redis.ProvideStore(context, tenantConfiguration, timeProvider, factory)
	eventStore := redis.ProvideEventStore(context, tenantConfiguration)
	sessionProvider := session.ProvideSessionProvider(r, sessionStore, eventStore, tenantConfiguration)
	authnSessionProvider := authn.ProvideSessionProvider(mfaProvider, sessionProvider, tenantConfiguration, timeProvider, authinfoStore, userprofileStore, identityProvider, hookProvider)
	authnProvider := &authn.Provider{
		OAuth:   oAuthCoordinator,
		Authn:   authenticateProcess,
		Signup:  signupProcess,
		Session: authnSessionProvider,
	}
	loginIDNormalizerFactory := loginid.ProvideLoginIDNormalizerFactory(tenantConfiguration)
	oAuthProviderFactory := sso.ProvideOAuthProviderFactory(tenantConfiguration, urlprefixProvider, timeProvider, loginIDNormalizerFactory)
	oAuthProvider := provideOAuthProviderFromRequestVars(r, oAuthProviderFactory)
	httpHandler := provideLinkHandler(txContext, requireAuthz, validator, contextGetter, provider, authnProvider, oAuthProvider)
	return httpHandler
}

func newLoginHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	contextGetter := auth2.ProvideAuthContextGetter(context)
	requestID := auth.ProvideLoggingRequestID(r)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	requireAuthz := handler.NewRequireAuthzFactory(contextGetter, factory)
	validator := auth.ProvideValidator(m)
	provider := sso.ProvideSSOProvider(context, tenantConfiguration)
	timeProvider := time.NewProvider()
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	store := pq.ProvidePasswordHistoryStore(timeProvider, sqlBuilder, sqlExecutor)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	passwordProvider := password.ProvidePasswordProvider(sqlBuilder, sqlExecutor, timeProvider, store, factory, tenantConfiguration, reservedNameChecker)
	oauthProvider := oauth.ProvideOAuthProvider(sqlBuilder, sqlExecutor)
	v := auth.ProvidePrincipalProviders(oauthProvider, passwordProvider)
	identityProvider := principal.ProvideIdentityProvider(sqlBuilder, sqlExecutor, v)
	authenticateProcess := authn.ProvideAuthenticateProcess(factory, timeProvider, passwordProvider, oauthProvider, identityProvider)
	passwordChecker := audit.ProvidePasswordChecker(tenantConfiguration, store)
	loginIDChecker := loginid.ProvideLoginIDChecker(tenantConfiguration, reservedNameChecker)
	authinfoStore := pq2.ProvideStore(sqlBuilderFactory, sqlExecutor)
	userprofileStore := userprofile.ProvideStore(timeProvider, sqlBuilder, sqlExecutor)
	urlprefixProvider := urlprefix.NewProvider(r)
	hookProvider := hook.ProvideHookProvider(sqlBuilder, sqlExecutor, requestID, tenantConfiguration, urlprefixProvider, contextGetter, txContext, timeProvider, authinfoStore, userprofileStore, passwordProvider, factory)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, txContext, requestID, tenantConfiguration, executor)
	signupProcess := authn.ProvideSignupProcess(passwordChecker, loginIDChecker, identityProvider, passwordProvider, oauthProvider, timeProvider, authinfoStore, userprofileStore, hookProvider, tenantConfiguration, urlprefixProvider, queue)
	oAuthCoordinator := &authn.OAuthCoordinator{
		Authn:  authenticateProcess,
		Signup: signupProcess,
	}
	mfaStore := pq3.ProvideStore(tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider)
	client := sms.ProvideSMSClient(tenantConfiguration)
	sender := mail.ProvideMailSender(tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	mfaSender := mfa.ProvideMFASender(tenantConfiguration, client, sender, engine)
	mfaProvider := mfa.ProvideMFAProvider(mfaStore, tenantConfiguration, timeProvider, mfaSender)
	sessionStore := redis.ProvideStore(context, tenantConfiguration, timeProvider, factory)
	eventStore := redis.ProvideEventStore(context, tenantConfiguration)
	sessionProvider := session.ProvideSessionProvider(r, sessionStore, eventStore, tenantConfiguration)
	authnSessionProvider := authn.ProvideSessionProvider(mfaProvider, sessionProvider, tenantConfiguration, timeProvider, authinfoStore, userprofileStore, identityProvider, hookProvider)
	authnProvider := &authn.Provider{
		OAuth:   oAuthCoordinator,
		Authn:   authenticateProcess,
		Signup:  signupProcess,
		Session: authnSessionProvider,
	}
	loginIDNormalizerFactory := loginid.ProvideLoginIDNormalizerFactory(tenantConfiguration)
	oAuthProviderFactory := sso.ProvideOAuthProviderFactory(tenantConfiguration, urlprefixProvider, timeProvider, loginIDNormalizerFactory)
	oAuthProvider := provideOAuthProviderFromRequestVars(r, oAuthProviderFactory)
	httpHandler := provideLoginHandler(txContext, requireAuthz, validator, provider, authnProvider, oAuthProvider)
	return httpHandler
}

// wire.go:

func provideOAuthProviderFromRequestVars(r *http.Request, spf *sso.OAuthProviderFactory) sso.OAuthProvider {
	vars := mux.Vars(r)
	return spf.NewOAuthProvider(vars["provider"])
}

func provideAuthHandler(
	tx db.TxContext,
	cfg *config.TenantConfiguration,
	hp sso.AuthHandlerHTMLProvider,
	sp sso.Provider,
	ap AuthHandlerAuthnProvider,
	op sso.OAuthProvider,
) http.Handler {
	h := &AuthHandler{
		TxContext:               tx,
		TenantConfiguration:     cfg,
		AuthHandlerHTMLProvider: hp,
		SSOProvider:             sp,
		AuthnProvider:           ap,
		OAuthProvider:           op,
	}
	return h
}

func provideAuthResultHandler(
	tx db.TxContext,
	requireAuthz handler.RequireAuthz,
	ap AuthResultAuthnProvider,
	v *validation.Validator,
	sp sso.Provider,
) http.Handler {
	h := &AuthResultHandler{
		TxContext:     tx,
		AuthnProvider: ap,
		Validator:     v,
		SSOProvider:   sp,
	}
	return requireAuthz(h, h)
}

func provideLinkHandler(
	tx db.TxContext,
	requireAuthz handler.RequireAuthz,
	v *validation.Validator,
	ac auth2.ContextGetter,
	sp sso.Provider,
	ap LinkAuthnProvider,
	op sso.OAuthProvider,
) http.Handler {
	h := &LinkHandler{
		TxContext:     tx,
		Validator:     v,
		AuthContext:   ac,
		SSOProvider:   sp,
		AuthnProvider: ap,
		OAuthProvider: op,
	}
	return requireAuthz(h, h)
}

func provideLoginHandler(
	tx db.TxContext,
	requireAuthz handler.RequireAuthz,
	v *validation.Validator,
	sp sso.Provider,
	ap LoginAuthnProvider,
	op sso.OAuthProvider,
) http.Handler {
	h := &LoginHandler{
		TxContext:     tx,
		Validator:     v,
		SSOProvider:   sp,
		AuthnProvider: ap,
		OAuthProvider: op,
	}
	return requireAuthz(h, h)
}
