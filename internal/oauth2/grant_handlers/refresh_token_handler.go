package grant_handlers

import (
	"context"
	"errors"

	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/oauth2/models"
	"github.com/shashimalcse/tiny-is/internal/oauth2/token"
	server_models "github.com/shashimalcse/tiny-is/internal/server/models"
)

type RefreshTokenGrantHandler struct {
	cacheService *cache.CacheService
	tokenService token.TokenService
}

func NewRefreshTokenGrantHandler(cacheService *cache.CacheService, tokenService token.TokenService) *RefreshTokenGrantHandler {
	return &RefreshTokenGrantHandler{
		cacheService: cacheService,
		tokenService: tokenService,
	}
}

func (gh *RefreshTokenGrantHandler) HandleGrant(ctx context.Context, oauth2TokenContext models.OAuth2TokenContext) (server_models.TokenResponse, error) {

	refresh_token := oauth2TokenContext.OAuth2TokenRequest.RefreshToken
	if refresh_token == "" {
		return server_models.TokenResponse{}, errors.New("invalid_refresh_token")
	}
	authroizeContext, err := gh.tokenService.ValidateRefreshToken(ctx, refresh_token)
	if err != nil {
		return server_models.TokenResponse{}, errors.New("invalid_refresh_token")
	}
	tokenString, err := gh.tokenService.GenerateAccessToken(ctx, authroizeContext, map[string]string{})
	if err != nil {
		return server_models.TokenResponse{}, err
	}
	tokenResponse := server_models.TokenResponse{
		AccessToken:  tokenString,
		RefreshToken: refresh_token,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}
	return tokenResponse, nil
}
