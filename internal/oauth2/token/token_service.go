package token

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	authn_models "github.com/shashimalcse/tiny-is/internal/authn/models"
	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/oauth2/models"
	server_models "github.com/shashimalcse/tiny-is/internal/server/models"
)

type TokenService interface {
	GenerateAccessToken(ctx context.Context, oauth2AuthroizeContext models.OAuth2AuthorizeContext, UserData map[string]string) (string, error)
	GenerateRefreshToken(ctx context.Context, oauth2AuthroizeContext models.OAuth2AuthorizeContext, UserData map[string]string) (string, error)
	ValidateRefreshToken(ctx context.Context, tokenString string) (models.OAuth2AuthorizeContext, error)
	RevokeToken(ctx context.Context, tokenString string)
}

type tokenService struct {
	cacheService    *cache.CacheService
	tokenRepository TokenRepository
	signingKey      []byte
}

func NewTokenService(cacheService *cache.CacheService, tokenRepository TokenRepository, signingKey []byte) TokenService {
	return &tokenService{
		cacheService:    cacheService,
		tokenRepository: tokenRepository,
		signingKey:      signingKey,
	}
}

func (s *tokenService) GenerateAccessToken(ctx context.Context, oauth2AuthroizeContext models.OAuth2AuthorizeContext, UserData map[string]string) (string, error) {

	claims, err := GetClaimsForAccessToken(oauth2AuthroizeContext.AuthenticatedUser.Id, "tiny-is")
	if err != nil {
		return "", err
	}

	// s.tokenRepository.PersistToken(ctx, claims["jti"].(string), claims["sub"].(string), oauth2AuthroizeContext.OAuth2AuthorizeRequest.ClientId, oauth2AuthroizeContext.OAuth2AuthorizeRequest.OrganizationId, claims["iat"].(int64), claims["exp"].(int64))

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := accessToken.SignedString(s.signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *tokenService) GenerateRefreshToken(ctx context.Context, oauth2AuthroizeContext models.OAuth2AuthorizeContext, UserData map[string]string) (string, error) {

	claims, err := GetClaimsForRefreshTokenToken(oauth2AuthroizeContext.AuthenticatedUser.Id, "tiny-is", oauth2AuthroizeContext.OAuth2AuthorizeRequest.ClientId)
	if err != nil {
		return "", err
	}

	s.tokenRepository.PersistToken(ctx, claims["jti"].(string), claims["sub"].(string), oauth2AuthroizeContext.OAuth2AuthorizeRequest.ClientId, oauth2AuthroizeContext.OAuth2AuthorizeRequest.OrganizationId, claims["iat"].(int64), claims["exp"].(int64))

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := accessToken.SignedString(s.signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *tokenService) ValidateRefreshToken(ctx context.Context, tokenString string) (models.OAuth2AuthorizeContext, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return models.OAuth2AuthorizeContext{}, errors.New("unexpected signing method")
		}

		return s.signingKey, nil
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
		jti, ok := claims["jti"].(string)
		if !ok {
			return models.OAuth2AuthorizeContext{}, errors.New("jti not found in refresh token")
		}
		isRefreshTokenExists, err := s.tokenRepository.IsTokenExists(ctx, jti)
		if err != nil {
			return models.OAuth2AuthorizeContext{}, err
		}
		if !isRefreshTokenExists {
			return models.OAuth2AuthorizeContext{}, errors.New("refresh token not found! It may be revoked or expired")
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

func (s *tokenService) RevokeToken(ctx context.Context, tokenString string) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return models.OAuth2AuthorizeContext{}, errors.New("unexpected signing method")
		}
		return s.signingKey, nil
	})

	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if jti, ok := claims["jti"].(string); ok {
			s.tokenRepository.DeleteToken(ctx, jti)
		}
	}
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
