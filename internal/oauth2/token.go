package oauth2

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/shashimalcse/tiny-is/internal/server/models"
)

func (o OAuth2) GetAccessToken(oauth2AuthroizeContext models.OAuth2AuthorizeContext) (string, error) {

	claims, err := o.GetClaimsForAccessToken(oauth2AuthroizeContext.AuthenticatedUser.Id, "tiny-is")
	if err != nil {
		return "", err
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := accessToken.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (o OAuth2) GetRefreshToken(oauth2AuthroizeContext models.OAuth2AuthorizeContext) (string, error) {

	claims, err := o.GetClaimsForRefreshTokenToken(oauth2AuthroizeContext.AuthenticatedUser.Id, "tiny-is", oauth2AuthroizeContext.OAuth2AuthorizeRequest.ClientId)
	if err != nil {
		return "", err
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := accessToken.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (o OAuth2) GetClaimsForAccessToken(sub, issuer string) (jwt.MapClaims, error) {

	expiresAt := time.Now().Add(time.Minute * 60).Unix()
	iat := time.Now().Unix()
	nbf := time.Now().Unix()
	jti, err := uuid.NewUUID()
	if err != nil {
		return jwt.MapClaims{}, err
	}
	claims := jwt.MapClaims{
		"sub": sub,
		"iss": issuer,
		"exp": expiresAt,
		"iat": iat,
		"nbf": nbf,
		"jti": jti.String(),
	}
	return claims, nil
}

func (o OAuth2) GetClaimsForRefreshTokenToken(sub, issuer, client_id string) (jwt.MapClaims, error) {

	expiresAt := time.Now().Add(time.Minute * 60 * 24 * 30).Unix()
	iat := time.Now().Unix()
	nbf := time.Now().Unix()
	jti, err := uuid.NewUUID()
	if err != nil {
		return jwt.MapClaims{}, err // Return an empty string and the error if generation fails
	}
	claims := jwt.MapClaims{
		"sub":       sub,
		"iss":       issuer,
		"exp":       expiresAt,
		"iat":       iat,
		"nbf":       nbf,
		"jti":       jti.String(),
		"client_id": client_id,
	}
	return claims, nil
}
