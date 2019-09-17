package middleware

import (
	"net/http"
	"strconv"

	"github.com/skygeario/skygear-server/pkg/core/auth"
	coreHttp "github.com/skygeario/skygear-server/pkg/core/http"
	coreMiddleware "github.com/skygeario/skygear-server/pkg/core/middleware"
	"github.com/skygeario/skygear-server/pkg/core/model"
)

// AuthInfoMiddleware injects auth info headers into the request
// if x-skygear-access-token is present in the request.
type AuthInfoMiddleware struct {
	AuthContext auth.ContextGetter `dependency:"AuthContextGetter"`
}

// AuthInfoMiddlewareFactory creates AuthInfoMiddleware per request.
type AuthInfoMiddlewareFactory struct{}

// NewInjectableMiddleware implements InjectableMiddlewareFactory.
func (f AuthInfoMiddlewareFactory) NewInjectableMiddleware() coreMiddleware.InjectableMiddleware {
	return &AuthInfoMiddleware{}
}

// Handle implements InjectableMiddleware.
func (m *AuthInfoMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		model.SetAccessKey(r, m.AuthContext.AccessKey())

		// Remove untrusted headers first.
		r.Header.Del(coreHttp.HeaderUserID)
		r.Header.Del(coreHttp.HeaderUserVerified)
		r.Header.Del(coreHttp.HeaderUserDisabled)
		r.Header.Del(coreHttp.HeaderSessionIdentityType)
		r.Header.Del(coreHttp.HeaderSessionAuthenticatorType)

		// TODO(mfa): If refresh token is enabled and the session is invalid,
		// do not forward the request and write `x-skygear-try-refresh-token: true`
		authInfo, _ := m.AuthContext.AuthInfo()
		if authInfo != nil {
			id := authInfo.ID
			disabled := authInfo.Disabled
			verified := authInfo.Verified

			r.Header.Set(coreHttp.HeaderUserID, id)
			r.Header.Set(coreHttp.HeaderUserVerified, strconv.FormatBool(verified))
			r.Header.Set(coreHttp.HeaderUserDisabled, strconv.FormatBool(disabled))
		}
		sess, _ := m.AuthContext.Session()
		if sess != nil {
			ptype := sess.PrincipalType
			if ptype != "" {
				r.Header.Set(coreHttp.HeaderSessionIdentityType, string(ptype))
			}
			atype := sess.AuthenticatorType
			if atype != "" {
				r.Header.Set(coreHttp.HeaderSessionAuthenticatorType, string(atype))
			}
		}

		next.ServeHTTP(w, r)
	})
}
