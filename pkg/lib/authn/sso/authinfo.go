package sso

import (
	"github.com/authgear/authgear-server/pkg/lib/authn/identity"
	"github.com/authgear/authgear-server/pkg/lib/config"
)

// AuthInfo contains auth info from HandleAuthzResp
type AuthInfo struct {
	ProviderConfig          config.OAuthSSOProviderConfig
	ProviderRawProfile      map[string]interface{}
	ProviderAccessTokenResp interface{}
	ProviderUserInfo        ProviderUserInfo
}

type ProviderUserInfo struct {
	ID string
	// Email is normalized.
	Email string
	// PhoneNumber is in E.164 format.
	PhoneNumber string
	// PreferredUsername is populated when Email or PhoneNumber is not applicable.
	PreferredUsername string
}

func (i ProviderUserInfo) ClaimsValue() map[string]interface{} {
	claimsValue := map[string]interface{}{}
	if i.Email != "" {
		claimsValue[identity.StandardClaimEmail] = i.Email
	}
	if i.PhoneNumber != "" {
		claimsValue[identity.StandardClaimPhoneNumber] = i.PhoneNumber
	}
	if i.PreferredUsername != "" {
		claimsValue[identity.StandardClaimPreferredUsername] = i.PreferredUsername
	}
	return claimsValue
}

type OAuthAuthorizationResponse struct {
	Code  string
	State string
	Scope string
}
