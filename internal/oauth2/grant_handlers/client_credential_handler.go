package grant_handlers

import (
	"context"

	authn_models "github.com/shashimalcse/tiny-is/internal/authn/models"
	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/oauth2/models"
	"github.com/shashimalcse/tiny-is/internal/oauth2/token"
	server_models "github.com/shashimalcse/tiny-is/internal/server/models"
)

type ClientCredetialGrantHandler struct {
	cacheService cache.CacheService
	tokenService token.TokenService
}

func NewClientCredetialGrantHandler(cacheService cache.CacheService, tokenService token.TokenService) *ClientCredetialGrantHandler {
	return &ClientCredetialGrantHandler{
		cacheService: cacheService,
		tokenService: tokenService,
	}
}

func (gh *ClientCredetialGrantHandler) HandleGrant(ctx context.Context, oauth2TokenContext models.OAuth2TokenContext) (server_models.TokenResponse, error) {
	authroizeContext := models.OAuth2AuthorizeContext{
		AuthenticatedUser: authn_models.AuthenticatedUser{
			Id: oauth2TokenContext.OAuth2TokenRequest.ClientId,
		},
	}
	tokenString, err := gh.tokenService.GenerateAccessToken(ctx, authroizeContext, map[string]string{})
	if err != nil {
		return server_models.TokenResponse{}, err
	}
	tokenResponse := server_models.TokenResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}
	return tokenResponse, nil
}
