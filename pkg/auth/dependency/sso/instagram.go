package sso

import (
	"github.com/skygeario/skygear-server/pkg/core/config"
)

const (
	instagramAuthorizationURL string = "https://api.instagram.com/oauth/authorize"
	// nolint: gosec
	instagramTokenURL    string = "https://api.instagram.com/oauth/access_token"
	instagramUserInfoURL string = "https://api.instagram.com/v1/users/self"
)

type InstagramImpl struct {
	OAuthConfig    config.OAuthConfiguration
	ProviderConfig config.OAuthProviderConfiguration
}

type instagramAuthInfoProcessor struct {
	defaultAuthInfoProcessor
}

func newInstagramAuthInfoProcessor() instagramAuthInfoProcessor {
	return instagramAuthInfoProcessor{}
}

func (f *InstagramImpl) GetAuthURL(params GetURLParams) (string, error) {
	p := authURLParams{
		oauthConfig:    f.OAuthConfig,
		providerConfig: f.ProviderConfig,
		state:          NewState(params),
		baseURL:        instagramAuthorizationURL,
	}
	return authURL(p)
}

func (f *InstagramImpl) DecodeState(encodedState string) (*State, error) {
	return DecodeState(f.OAuthConfig.StateJWTSecret, encodedState)
}

func (f *InstagramImpl) GetAuthInfo(r OAuthAuthorizationResponse) (authInfo AuthInfo, err error) {
	return f.NonOpenIDConnectGetAuthInfo(r)
}

func (f *InstagramImpl) NonOpenIDConnectGetAuthInfo(r OAuthAuthorizationResponse) (authInfo AuthInfo, err error) {
	p := newInstagramAuthInfoProcessor()
	h := getAuthInfoRequest{
		oauthConfig:    f.OAuthConfig,
		providerConfig: f.ProviderConfig,
		code:           r.Code,
		accessTokenURL: instagramTokenURL,
		userProfileURL: instagramUserInfoURL,
		processor:      p,
	}
	return h.getAuthInfo()
}

func (i instagramAuthInfoProcessor) DecodeUserInfo(userProfile map[string]interface{}) (info ProviderUserInfo) {
	// Check GET /users/self response
	// https://www.instagram.com/developer/endpoints/users/
	data, ok := userProfile["data"].(map[string]interface{})
	if !ok {
		return
	}

	info.ID, _ = data["id"].(string)
	info.Email, _ = data["email"].(string)
	return
}

func (f *InstagramImpl) ExternalAccessTokenGetAuthInfo(accessTokenResp AccessTokenResp) (authInfo AuthInfo, err error) {
	p := newInstagramAuthInfoProcessor()
	h := getAuthInfoRequest{
		oauthConfig:    f.OAuthConfig,
		providerConfig: f.ProviderConfig,
		accessTokenURL: instagramTokenURL,
		userProfileURL: instagramUserInfoURL,
		processor:      p,
	}
	return h.getAuthInfoByAccessTokenResp(accessTokenResp)
}

var (
	_ OAuthProvider                   = &InstagramImpl{}
	_ NonOpenIDConnectProvider        = &InstagramImpl{}
	_ ExternalAccessTokenFlowProvider = &InstagramImpl{}
)
