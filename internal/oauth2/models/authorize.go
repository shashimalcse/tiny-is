package models

import (
	"github.com/shashimalcse/tiny-is/internal/authn/models"
	server_models "github.com/shashimalcse/tiny-is/internal/server/models"
)

type OAuth2AuthorizeContext struct {
	OAuth2AuthorizeRequest server_models.OAuth2AuthorizeRequest `json:"oauth2_authorize_request"`
	AuthenticatedUser      models.AuthenticatedUser             `json:"authenticated_user"`
}

type OAuth2TokenContext struct {
	OAuth2TokenRequest server_models.OAuth2TokenRequest `json:"oauth2_token_request"`
}
