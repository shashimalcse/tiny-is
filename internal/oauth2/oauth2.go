package oauth2

import (
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/application"
	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/server/models"
)

type OAuth2 struct {
	CacheService       *cache.CacheService
	ApplicationService *application.ApplicationService
	GrantTypes         map[string]func(*http.Request) (models.TokenResponse, error)
}

func NewOAuth2(cacheService *cache.CacheService, applicationService *application.ApplicationService) *OAuth2 {
	oauth2 := &OAuth2{
		CacheService: cacheService,
		GrantTypes:   make(map[string]func(*http.Request) (models.TokenResponse, error)),
	}
	oauth2.ApplicationService = applicationService
	oauth2.GrantTypes["authorization_code"] = oauth2.AuthorizationCodeGrant
	return oauth2
}
