package httputil

import (
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

type CookieDef struct {
	// NameSuffix means the cookie could have prefix.
	NameSuffix string
	Path       string
	// Domain is omitted because it is controlled somewhere else.
	// Domain            string
	AllowScriptAccess bool
	SameSite          http.SameSite
	MaxAge            *int
}

func UpdateCookie(w http.ResponseWriter, cookie *http.Cookie) {
	header := w.Header()
	resp := http.Response{Header: header}
	cookies := resp.Cookies()
	updated := false
	for i, c := range cookies {
		if c.Name == cookie.Name && c.Domain == cookie.Domain && c.Path == cookie.Path {
			cookies[i] = cookie
			updated = true
		}
	}
	if !updated {
		cookies = append(cookies, cookie)
	}
	setCookies := make([]string, len(cookies))
	for i, c := range cookies {
		setCookies[i] = c.String()
	}
	header["Set-Cookie"] = setCookies
}

// CookieDomainFromETLDPlusOneWithoutPort derives host from r.
// If host has port, the port is removed.
// If ETLD+1 cannot be derived, an empty string is returned.
// The return value never have port.
func CookieDomainFromETLDPlusOneWithoutPort(host string) string {
	// Trim the port if it is present.
	// We have to trim the port first.
	// Passing host:port to EffectiveTLDPlusOne confuses it.
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		ipv6Str := host[1 : len(host)-1]
		if ipv6 := net.ParseIP(ipv6Str); ipv6 != nil {
			return ""
		}
	}

	if ipv4or6 := net.ParseIP(host); ipv4or6 != nil {
		return ""
	}

	host, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		return ""
	}

	return host
}

type CookieManager struct {
	Request      *http.Request
	TrustProxy   bool
	CookiePrefix string
	CookieDomain string
}

func (f *CookieManager) fixupCookie(cookie *http.Cookie) {
	host := GetHost(f.Request, f.TrustProxy)
	proto := GetProto(f.Request, f.TrustProxy)

	cookie.Secure = proto == "https"
	if cookie.Domain == "" {
		cookie.Domain = CookieDomainFromETLDPlusOneWithoutPort(host)
	}

	if cookie.SameSite == http.SameSiteNoneMode &&
		!ShouldSendSameSiteNone(f.Request.UserAgent(), cookie.Secure) {
		cookie.SameSite = 0
	}
}

// CookieName returns the full name, that is, CookiePrefix followed by NameSuffix.
func (f *CookieManager) CookieName(def *CookieDef) string {
	return f.CookiePrefix + def.NameSuffix
}

// GetCookie is wrapper around http.Request.Cookie, taking care of cookie name.
func (f *CookieManager) GetCookie(r *http.Request, def *CookieDef) (*http.Cookie, error) {
	cookieName := f.CookieName(def)
	return r.Cookie(cookieName)
}

// ValueCookie generates a cookie that when set, the cookie is set to the specified value.
func (f *CookieManager) ValueCookie(def *CookieDef, value string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     f.CookieName(def),
		Path:     def.Path,
		Domain:   f.CookieDomain,
		HttpOnly: !def.AllowScriptAccess,
		SameSite: def.SameSite,
	}

	cookie.Value = value
	if def.MaxAge != nil {
		cookie.MaxAge = *def.MaxAge
	}

	f.fixupCookie(cookie)

	return cookie
}

// ClearCookie generates a cookie that when set, the cookie is clear.
func (f *CookieManager) ClearCookie(def *CookieDef) *http.Cookie {
	cookie := &http.Cookie{
		Name:     f.CookieName(def),
		Path:     def.Path,
		Domain:   f.CookieDomain,
		HttpOnly: !def.AllowScriptAccess,
		SameSite: def.SameSite,
		Expires:  time.Unix(0, 0),
	}

	f.fixupCookie(cookie)

	return cookie
}
