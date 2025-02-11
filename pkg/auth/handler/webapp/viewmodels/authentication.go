package viewmodels

import (
	"github.com/authgear/authgear-server/pkg/lib/authn"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity"
	"github.com/authgear/authgear-server/pkg/lib/interaction"
)

// Ideally we should use type alias to present
// LoginPageTextLoginIDVariant and LoginPageTextLoginIDInputType
// But they may be passed to localize which does not support type alias of builtin types.

const (
	LoginPageTextLoginIDVariantNone            = "none"
	LoginPageTextLoginIDVariantEamilOrUsername = "email_or_username"
	LoginPageTextLoginIDVariantEmail           = "email"
	LoginPageTextLoginIDVariantUsername        = "username"
)

const (
	LoginPageTextLoginIDInputTypeText  = "text"
	LoginPageTextLoginIDInputTypeEmail = "email"
)

type IdentityCandidatesGetter interface {
	GetIdentityCandidates() []identity.Candidate
}

type AuthenticationViewModel struct {
	IdentityCandidates            []identity.Candidate
	IdentityCount                 int
	LoginPageLoginIDHasPhone      bool
	LoginPageTextLoginIDVariant   string
	LoginPageTextLoginIDInputType string
}

func NewAuthenticationViewModelWithGraph(graph *interaction.Graph) AuthenticationViewModel {
	var node IdentityCandidatesGetter
	if !graph.FindLastNode(&node) {
		panic("webapp: no node with identity candidates found")
	}

	return NewAuthenticationViewModelWithCandidates(node.GetIdentityCandidates())
}

func NewAuthenticationViewModelWithCandidates(candidates []identity.Candidate) AuthenticationViewModel {
	hasEmail := false
	hasUsername := false
	hasPhone := false
	identityCount := 0

	for _, c := range candidates {
		typ, _ := c[identity.CandidateKeyType].(string)
		if typ == string(authn.IdentityTypeLoginID) {
			loginIDType, _ := c[identity.CandidateKeyLoginIDType].(string)
			switch loginIDType {
			case "phone":
				c["login_id_input_type"] = "phone"
				hasPhone = true
			case "email":
				c["login_id_input_type"] = "email"
				hasEmail = true
			default:
				c["login_id_input_type"] = "text"
				hasUsername = true
			}
		}

		identityID := c[identity.CandidateKeyIdentityID].(string)
		if identityID != "" {
			identityCount++
		}
	}

	var loginPageTextLoginIDVariant string
	var loginPageTextLoginIDInputType string
	if hasEmail {
		if hasUsername {
			loginPageTextLoginIDVariant = LoginPageTextLoginIDVariantEamilOrUsername
			loginPageTextLoginIDInputType = LoginPageTextLoginIDInputTypeText
		} else {
			loginPageTextLoginIDVariant = LoginPageTextLoginIDVariantEmail
			loginPageTextLoginIDInputType = LoginPageTextLoginIDInputTypeEmail
		}
	} else {
		if hasUsername {
			loginPageTextLoginIDVariant = LoginPageTextLoginIDVariantUsername
			loginPageTextLoginIDInputType = LoginPageTextLoginIDInputTypeText
		} else {
			loginPageTextLoginIDVariant = LoginPageTextLoginIDVariantNone
			loginPageTextLoginIDInputType = LoginPageTextLoginIDInputTypeText
		}
	}

	return AuthenticationViewModel{
		IdentityCandidates:            candidates,
		IdentityCount:                 identityCount,
		LoginPageLoginIDHasPhone:      hasPhone,
		LoginPageTextLoginIDVariant:   loginPageTextLoginIDVariant,
		LoginPageTextLoginIDInputType: loginPageTextLoginIDInputType,
	}
}
