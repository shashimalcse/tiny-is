package token

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	authn_models "github.com/shashimalcse/tiny-is/internal/authn/models"
	"github.com/shashimalcse/tiny-is/internal/oauth2/models"
	server_models "github.com/shashimalcse/tiny-is/internal/server/models"
)

type TokenService interface {
	GenerateAccessToken(ctx context.Context, oauth2AuthroizeContext models.OAuth2AuthorizeContext, UserData map[string]string) (string, error)
	GenerateRefreshToken(ctx context.Context, oauth2AuthroizeContext models.OAuth2AuthorizeContext, UserData map[string]string) (string, error)
	ValidateToken(ctx context.Context, tokenString string) (models.OAuth2AuthorizeContext, error)
}

type tokenService struct {
	signingKey []byte
}

func NewTokenService(signingKey []byte) TokenService {
	return &tokenService{
		signingKey: signingKey,
	}
}

func (ts *tokenService) GenerateAccessToken(ctx context.Context, oauth2AuthroizeContext models.OAuth2AuthorizeContext, UserData map[string]string) (string, error) {

	claims, err := GetClaimsForAccessToken(oauth2AuthroizeContext.AuthenticatedUser.Id, "tiny-is")
	if err != nil {
		return "", err
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := accessToken.SignedString(ts.signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (ts *tokenService) GenerateRefreshToken(ctx context.Context, oauth2AuthroizeContext models.OAuth2AuthorizeContext, UserData map[string]string) (string, error) {

	claims, err := GetClaimsForRefreshTokenToken(oauth2AuthroizeContext.AuthenticatedUser.Id, "tiny-is", oauth2AuthroizeContext.OAuth2AuthorizeRequest.ClientId)
	if err != nil {
		return "", err
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := accessToken.SignedString(ts.signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (ts *tokenService) ValidateToken(ctx context.Context, tokenString string) (models.OAuth2AuthorizeContext, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return models.OAuth2AuthorizeContext{}, errors.New("unexpected signing method")
		}

		return ts.signingKey, nil
	})

	if err != nil {
		return models.OAuth2AuthorizeContext{}, errors.New("invalid refresh token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return models.OAuth2AuthorizeContext{}, errors.New("refresh token has expired")
			}
		} else {
			return models.OAuth2AuthorizeContext{}, errors.New("invalid expiration claim")
		}

		clientID, ok := claims["client_id"].(string)
		if !ok {
			return models.OAuth2AuthorizeContext{}, errors.New("client ID not found in refresh token")
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			return models.OAuth2AuthorizeContext{}, errors.New("sub not found in refresh token")
		}

		authroizeContext := models.OAuth2AuthorizeContext{
			OAuth2AuthorizeRequest: server_models.OAuth2AuthorizeRequest{
				ClientId: clientID,
			},
			AuthenticatedUser: authn_models.AuthenticatedUser{
				Id: sub,
			},
		}
		return authroizeContext, nil
	}
	return models.OAuth2AuthorizeContext{}, errors.New("invalid token")
}

func GetClaimsForAccessToken(sub, issuer string) (jwt.MapClaims, error) {

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

func GetClaimsForRefreshTokenToken(sub, issuer, client_id string) (jwt.MapClaims, error) {

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
