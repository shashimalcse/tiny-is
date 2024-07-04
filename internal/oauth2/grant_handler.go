package oauth2

import (
	"errors"
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/server/models"
)

func (oauth2 OAuth2) AuthorizationCodeGrant(r *http.Request) (models.TokenResponse, error) {

	code := r.Form.Get("code")
	authroizeContext, found := oauth2.CacheService.GetCodeFromCache(code)
	if !found {
		return models.TokenResponse{}, errors.New("invalid code")
	}

	tokenString, err := oauth2.GetAccessToken(authroizeContext)
	if err != nil {
		return models.TokenResponse{}, err
	}
	tokenResponse := models.TokenResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
	}
	return tokenResponse, nil
}
