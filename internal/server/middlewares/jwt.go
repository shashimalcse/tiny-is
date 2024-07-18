package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shashimalcse/tiny-is/internal/config"
)

func JWTMiddleware(cfg *config.Config) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			publicPaths := map[string]bool{
				"/login":     true,
				"/authorize": true,
			}
			if publicPaths[r.URL.Path] {
				next(w, r)
			}
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				return NewAPIError(http.StatusUnauthorized, "missing Authorization header")
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
				return NewAPIError(http.StatusUnauthorized, "invalid Authorization header format")
			}

			tokenString := bearerToken[1]
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte("secret"), nil
			})

			if err != nil {
				return NewAPIError(http.StatusUnauthorized, "invalid token: "+err.Error())
			}

			if !token.Valid {
				return NewAPIError(http.StatusUnauthorized, "invalid token")
			}
			ctx := context.WithValue(r.Context(), "claims", claims)
			return next(w, r.WithContext(ctx))
		}

	}
}
