package oauth2

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	authn_models "github.com/shashimalcse/tiny-is/internal/authn/models"
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
	refreshTokenString, err := oauth2.GetRefreshToken(authroizeContext)
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

func (oauth2 OAuth2) RefreshTokenGrant(r *http.Request) (models.TokenResponse, error) {

	refresh_token := r.Form.Get("refresh_token")
	if refresh_token == "" {
		return models.TokenResponse{}, errors.New("refresh token not found")
	}

	token, err := jwt.Parse(refresh_token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte("secret"), nil
	})

	if err != nil {
		return models.TokenResponse{}, errors.New("invalid refresh token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return models.TokenResponse{}, errors.New("refresh token has expired")
			}
		} else {
			return models.TokenResponse{}, errors.New("invalid expiration claim")
		}

		clientID, ok := claims["client_id"].(string)
		if !ok {
			return models.TokenResponse{}, errors.New("client ID not found in refresh token")
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			return models.TokenResponse{}, errors.New("sub not found in refresh token")
		}

		authroizeContext := models.OAuth2AuthorizeContext{
			OAuth2AuthorizeRequest: models.OAuth2AuthorizeRequest{
				ClientId: clientID,
			},
			AuthenticatedUser: authn_models.AuthenticatedUser{
				Id: sub,
			},
		}

		tokenString, err := oauth2.GetAccessToken(authroizeContext)
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

	return models.TokenResponse{}, errors.New("invalid token")

}
