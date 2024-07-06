package grant_handlers

import (
	"errors"
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/oauth2/token"
	"github.com/shashimalcse/tiny-is/internal/server/models"
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

func (gh *RefreshTokenGrantHandler) HandleGrant(r *http.Request) (models.TokenResponse, error) {

	refresh_token := r.Form.Get("refresh_token")
	if refresh_token == "" {
		return models.TokenResponse{}, errors.New("refresh token not found")
	}
	ctx := r.Context()
	authroizeContext, err := gh.tokenService.ValidateToken(ctx, refresh_token)
	if err != nil {
		return models.TokenResponse{}, err
	}
	tokenString, err := gh.tokenService.GenerateAccessToken(ctx, authroizeContext, map[string]string{})
	if err != nil {
		return models.TokenResponse{}, err
	}
	tokenResponse := models.TokenResponse{
		AccessToken:  tokenString,
		RefreshToken: refresh_token,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}
	return tokenResponse, nil
}
