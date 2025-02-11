package oauth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/lib/session"
	"github.com/authgear/authgear-server/pkg/lib/session/access"
	"github.com/authgear/authgear-server/pkg/lib/session/idpsession"
	"github.com/authgear/authgear-server/pkg/util/clock"
	"github.com/authgear/authgear-server/pkg/util/httputil"
)

type ResolverSessionProvider interface {
	AccessWithID(id string, accessEvent access.Event) (*idpsession.IDPSession, error)
}

type AccessTokenDecoder interface {
	DecodeAccessToken(encodedToken string) (tok string, isHash bool, err error)
}

type ResolverCookieManager interface {
	GetCookie(r *http.Request, def *httputil.CookieDef) (*http.Cookie, error)
}

type Resolver struct {
	OAuthConfig        *config.OAuthConfig
	TrustProxy         config.TrustProxy
	Authorizations     AuthorizationStore
	AccessGrants       AccessGrantStore
	OfflineGrants      OfflineGrantStore
	AppSessions        AppSessionStore
	AccessTokenDecoder AccessTokenDecoder
	Sessions           ResolverSessionProvider
	Cookies            ResolverCookieManager
	SessionCookieDef   session.CookieDef
	Clock              clock.Clock
}

func (re *Resolver) Resolve(rw http.ResponseWriter, r *http.Request) (session.Session, error) {
	// The resolve function has the following outcomes:
	// - (nil, nil) which means no session was found and no error.
	// - (nil, err) in which err is ErrInvalidSession, which means the session was found but was invalid.
	//              in which err is something else, which means some unexpected error occurred.
	// - (s, nil)  which means a session was resolved successfully.
	//
	// Here we want to try the next resolve function iff the outcome is (nil, nil).
	funcs := []func(*http.Request) (session.Session, error){
		re.resolveHeader,
		re.resolveCookie,
	}

	for _, f := range funcs {
		s, err := f(r)
		if err != nil {
			return nil, err
		}
		if s != nil {
			return s, nil
		}
	}

	return nil, nil
}

func (re *Resolver) resolveHeader(r *http.Request) (session.Session, error) {
	token := parseAuthorizationHeader(r)
	if token == "" {
		// No bearer token in Authorization header. Simply proceed.
		return nil, nil
	}

	tok, isHash, err := re.AccessTokenDecoder.DecodeAccessToken(token)
	if err != nil {
		return nil, session.ErrInvalidSession
	}

	var tokenHash string
	if isHash {
		tokenHash = tok
	} else {
		tokenHash = HashToken(token)
	}

	grant, err := re.AccessGrants.GetAccessGrant(tokenHash)
	if errors.Is(err, ErrGrantNotFound) {
		return nil, session.ErrInvalidSession
	} else if err != nil {
		return nil, err
	}

	_, err = re.Authorizations.GetByID(grant.AuthorizationID)
	if errors.Is(err, ErrAuthorizationNotFound) {
		// Authorization does not exists (e.g. revoked)
		return nil, session.ErrInvalidSession
	} else if err != nil {
		return nil, err
	}

	var authSession session.Session
	event := access.NewEvent(re.Clock.NowUTC(), r, bool(re.TrustProxy))

	switch grant.SessionKind {
	case GrantSessionKindSession:
		s, err := re.Sessions.AccessWithID(grant.SessionID, event)
		if errors.Is(err, idpsession.ErrSessionNotFound) {
			return nil, session.ErrInvalidSession
		} else if err != nil {
			return nil, err
		}
		authSession = s

	case GrantSessionKindOffline:
		g, err := re.OfflineGrants.GetOfflineGrant(grant.SessionID)
		if errors.Is(err, ErrGrantNotFound) {
			return nil, session.ErrInvalidSession
		} else if err != nil {
			return nil, err
		}

		expiry, err := ComputeOfflineGrantExpiryWithClients(g, re.OAuthConfig)
		if errors.Is(err, ErrGrantNotFound) {
			return nil, session.ErrInvalidSession
		} else if err != nil {
			return nil, err
		}

		g, err = re.OfflineGrants.AccessWithID(g.ID, event, expiry)
		if err != nil {
			return nil, err
		}

		authSession = g
	default:
		panic("oauth: resolving unknown grant session kind")
	}

	return authSession, nil
}

func (re *Resolver) resolveCookie(r *http.Request) (session.Session, error) {
	cookie, err := re.Cookies.GetCookie(r, re.SessionCookieDef.Def)
	if err != nil {
		// No session cookie. Simply proceed.
		return nil, nil
	}

	aSession, err := re.AppSessions.GetAppSession(HashToken(cookie.Value))
	if errors.Is(err, ErrGrantNotFound) {
		return nil, session.ErrInvalidSession
	} else if err != nil {
		return nil, err
	}

	offlineGrant, err := re.OfflineGrants.GetOfflineGrant(aSession.OfflineGrantID)
	if errors.Is(err, ErrGrantNotFound) {
		return nil, session.ErrInvalidSession
	} else if err != nil {
		return nil, err
	}

	authz, err := re.Authorizations.GetByID(offlineGrant.AuthorizationID)
	if errors.Is(err, ErrAuthorizationNotFound) {
		// Authorization does not exists (e.g. revoked)
		return nil, session.ErrInvalidSession
	} else if err != nil {
		return nil, err
	} else if !authz.IsAuthorized([]string{FullAccessScope}) {
		// App sessions must have full user access to be valid
		return nil, session.ErrInvalidSession
	}

	expiry, err := ComputeOfflineGrantExpiryWithClients(offlineGrant, re.OAuthConfig)
	if errors.Is(err, ErrGrantNotFound) {
		return nil, session.ErrInvalidSession
	} else if err != nil {
		return nil, err
	}

	event := access.NewEvent(re.Clock.NowUTC(), r, bool(re.TrustProxy))
	offlineGrant, err = re.OfflineGrants.AccessWithID(offlineGrant.ID, event, expiry)
	if err != nil {
		return nil, err
	}

	return offlineGrant, nil
}

func parseAuthorizationHeader(r *http.Request) (token string) {
	authorization := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(authorization) != 2 {
		return
	}

	scheme := authorization[0]
	if strings.ToLower(scheme) != "bearer" {
		return
	}

	return authorization[1]
}
