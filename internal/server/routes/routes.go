package routes

import (
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/application"
	"github.com/shashimalcse/tiny-is/internal/authn"
	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/oauth2"
	"github.com/shashimalcse/tiny-is/internal/user"
)

func NewRouter(cacheService *cache.CacheService, applicationService *application.ApplicationService, userService *user.UserService) *http.ServeMux {
	mux := http.NewServeMux()

	RegisterOAuth2Routes(mux, oauth2.NewOAuth2(cacheService, applicationService))
	RegisterAuthnRoutes(mux, authn.NewAuthn(cacheService, userService))
	RegisterApplicationRoutes(mux, applicationService)
	return mux
}
