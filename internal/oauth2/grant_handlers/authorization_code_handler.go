package grant_handlers

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/oauth2/models"
	"github.com/shashimalcse/tiny-is/internal/oauth2/token"
	server_models "github.com/shashimalcse/tiny-is/internal/server/models"
)

type AuthorizationCodeGrantHandler struct {
	cacheService cache.CacheService
	tokenService token.TokenService
}

func NewAuthorizationCodeGrantHandler(cacheService cache.CacheService, tokenService token.TokenService) *AuthorizationCodeGrantHandler {
	return &AuthorizationCodeGrantHandler{
		cacheService: cacheService,
		tokenService: tokenService,
	}
}

func (gh *AuthorizationCodeGrantHandler) HandleGrant(ctx context.Context, oauth2TokenContext models.OAuth2TokenContext) (server_models.TokenResponse, error) {
	authorizeContext, found := gh.cacheService.GetOAuth2AuthorizeContextFromCacheByAuthCode(oauth2TokenContext.OAuth2TokenRequest.Code)
	if !found {
		return server_models.TokenResponse{}, errors.New("invalid_code")
	}
	// handle pkce
	if authorizeContext.OAuth2AuthorizeRequest.CodeChallenge != "" {
		if authorizeContext.OAuth2AuthorizeRequest.CodeChallengeMethod == "" || authorizeContext.OAuth2AuthorizeRequest.CodeChallengeMethod == "plain" {
			if authorizeContext.OAuth2AuthorizeRequest.CodeChallenge != oauth2TokenContext.OAuth2TokenRequest.CodeVerifier {
				return server_models.TokenResponse{}, errors.New("invalid_code_verifier")
			}
		} else if authorizeContext.OAuth2AuthorizeRequest.CodeChallengeMethod == "S256" {
			h := sha256.New()
			h.Write([]byte(oauth2TokenContext.OAuth2TokenRequest.CodeVerifier))
			codeChallenge := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
			if authorizeContext.OAuth2AuthorizeRequest.CodeChallenge != codeChallenge {
				return server_models.TokenResponse{}, errors.New("invalid_code_verifier")
			}
		} else {
			return server_models.TokenResponse{}, errors.New("invalid_code_challenge_method")
		}
	}
	tokenString, err := gh.tokenService.GenerateAccessToken(ctx, authorizeContext, map[string]string{})
	if err != nil {
		return server_models.TokenResponse{}, err
	}
	refreshTokenString, err := gh.tokenService.GenerateRefreshToken(ctx, authorizeContext, map[string]string{})
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
