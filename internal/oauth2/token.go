package oauth2

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/shashimalcse/tiny-is/internal/server/models"
)

func (o OAuth2) GetAccessToken(oauth2AuthroizeContext models.OAuth2AuthorizeContext) (string, error) {

	claims := jwt.MapClaims{
		"sub": oauth2AuthroizeContext.AuthenticatedUser.Id,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := accessToken.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
