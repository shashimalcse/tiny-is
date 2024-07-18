package routes

import (
	"github.com/shashimalcse/tiny-is/internal/application"
	"github.com/shashimalcse/tiny-is/internal/authn"
	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/config"
	"github.com/shashimalcse/tiny-is/internal/oauth2"
	"github.com/shashimalcse/tiny-is/internal/oauth2/token"
	"github.com/shashimalcse/tiny-is/internal/organization"
	tinyhttp "github.com/shashimalcse/tiny-is/internal/server/http"
	"github.com/shashimalcse/tiny-is/internal/session"
	"github.com/shashimalcse/tiny-is/internal/user"
)

func NewRouter(cfg *config.Config, cacheService cache.CacheService, sessionStore session.SessionStore, organizationService organization.OrganizationService, applicationService application.ApplicationService, userService user.UserService, tokenService token.TokenService) *tinyhttp.TinyServeMux {
	mux := tinyhttp.NewTinyServeMux(organizationService)

	RegisterOAuth2Routes(mux, oauth2.NewOAuth2Service(cacheService, tokenService, applicationService))
	RegisterAuthnRoutes(mux, authn.NewAuthnService(cacheService, sessionStore, userService))
	RegisterApplicationRoutes(mux, cfg, applicationService)
	RegisterUserRoutes(mux, cfg, userService)
	return mux
}
