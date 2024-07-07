package grant_handlers

import (
	"context"

	"github.com/shashimalcse/tiny-is/internal/oauth2/models"
	server_models "github.com/shashimalcse/tiny-is/internal/server/models"
)

type GrantHandler interface {
	HandleGrant(ctx context.Context, oauth2TokenContext models.OAuth2TokenContext) (server_models.TokenResponse, error)
}
