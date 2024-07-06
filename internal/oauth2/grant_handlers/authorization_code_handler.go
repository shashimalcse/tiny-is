package grant_handlers

import (
	"errors"
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/oauth2/token"
	"github.com/shashimalcse/tiny-is/internal/server/models"
)

type AuthorizationCodeGrantHandler struct {
	cacheService *cache.CacheService
	tokenService token.TokenService
}

func NewAuthorizationCodeGrantHandler(cacheService *cache.CacheService, tokenService token.TokenService) *AuthorizationCodeGrantHandler {
	return &AuthorizationCodeGrantHandler{
		cacheService: cacheService,
		tokenService: tokenService,
	}
}

func (gh *AuthorizationCodeGrantHandler) HandleGrant(r *http.Request) (models.TokenResponse, error) {

	code := r.Form.Get("code")
	authroizeContext, found := gh.cacheService.GetOAuth2AuthorizeContextFromCacheByAuthCode(code)
	if !found {
		return models.TokenResponse{}, errors.New("invalid code")
	}
	ctx := r.Context()
	tokenString, err := gh.tokenService.GenerateAccessToken(ctx, authroizeContext, map[string]string{})
	if err != nil {
		return models.TokenResponse{}, err
	}
	refreshTokenString, err := gh.tokenService.GenerateRefreshToken(ctx, authroizeContext, map[string]string{})
	if err != nil {
		return models.TokenResponse{}, err
	}
	tokenResponse := models.TokenResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}
	return tokenResponse, nil
}
