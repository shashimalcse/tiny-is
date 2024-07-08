package grant_handlers

import (
	"context"
	"errors"

	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/oauth2/models"
	"github.com/shashimalcse/tiny-is/internal/oauth2/token"
	server_models "github.com/shashimalcse/tiny-is/internal/server/models"
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

func (gh *AuthorizationCodeGrantHandler) HandleGrant(ctx context.Context, oauth2TokenContext models.OAuth2TokenContext) (server_models.TokenResponse, error) {

	authroizeContext, found := gh.cacheService.GetOAuth2AuthorizeContextFromCacheByAuthCode(oauth2TokenContext.OAuth2TokenRequest.Code)
	if !found {
		return server_models.TokenResponse{}, errors.New("invalid_code")
	}
	tokenString, err := gh.tokenService.GenerateAccessToken(ctx, authroizeContext, map[string]string{})
	if err != nil {
		return server_models.TokenResponse{}, err
	}
	refreshTokenString, err := gh.tokenService.GenerateRefreshToken(ctx, authroizeContext, map[string]string{})
	if err != nil {
		return server_models.TokenResponse{}, err
	}
	tokenResponse := server_models.TokenResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}
	return tokenResponse, nil
}
