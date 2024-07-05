package routes

import (
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/application"
	"github.com/shashimalcse/tiny-is/internal/authn"
	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/oauth2"
	"github.com/shashimalcse/tiny-is/internal/session"
	"github.com/shashimalcse/tiny-is/internal/user"
)

func NewRouter(cacheService *cache.CacheService, sessionStore *session.SessionStore, applicationService *application.ApplicationService, userService *user.UserService) *http.ServeMux {
	mux := http.NewServeMux()

	RegisterOAuth2Routes(mux, oauth2.NewOAuth2(cacheService, applicationService))
	RegisterAuthnRoutes(mux, authn.NewAuthn(cacheService, sessionStore, userService))
	RegisterApplicationRoutes(mux, applicationService)
	return mux
}
