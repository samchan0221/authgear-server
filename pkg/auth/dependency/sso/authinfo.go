package sso

import (
	"crypto/subtle"
	"fmt"

	"github.com/skygeario/skygear-server/pkg/core/config"
	"github.com/skygeario/skygear-server/pkg/core/hash"
)

type getAuthInfoRequest struct {
	oauthConfig    config.OAuthConfiguration
	providerConfig config.OAuthProviderConfiguration
	accessTokenURL string
	userProfileURL string
	processor      AuthInfoProcessor
}

func (h getAuthInfoRequest) getAuthInfo(r OAuthAuthorizationResponse) (authInfo AuthInfo, err error) {
	state, err := DecodeState(h.oauthConfig.StateJWTSecret, r.State)
	if err != nil {
		return
	}

	if subtle.ConstantTimeCompare([]byte(state.Nonce), []byte(hash.SHA256String(r.Nonce))) != 1 {
		err = fmt.Errorf("invalid nonce")
		return
	}

	// compare nonce
	authInfo = AuthInfo{
		ProviderConfig: h.providerConfig,
	}

	accessTokenResp, err := fetchAccessTokenResp(
		r.Code,
		h.accessTokenURL,
		h.oauthConfig,
		h.providerConfig,
	)
	if err != nil {
		return
	}

	err = accessTokenResp.Validate()
	if err != nil {
		return
	}
	authInfo.ProviderAccessTokenResp = accessTokenResp

	return h.getAuthInfoByAccessTokenResp(accessTokenResp)
}

func (h getAuthInfoRequest) getAuthInfoByAccessTokenResp(accessTokenResp AccessTokenResp) (authInfo AuthInfo, err error) {
	authInfo = AuthInfo{
		ProviderConfig: h.providerConfig,
		// validated accessTokenResp
		ProviderAccessTokenResp: accessTokenResp,
	}

	userProfile, err := fetchUserProfile(accessTokenResp, h.userProfileURL)
	if err != nil {
		return
	}
	authInfo.ProviderRawProfile = userProfile
	authInfo.ProviderUserInfo = h.processor.DecodeUserInfo(userProfile)

	return
}
