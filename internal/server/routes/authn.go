package routes

import (
	"github.com/shashimalcse/tiny-is/internal/authn"
	"github.com/shashimalcse/tiny-is/internal/server/handlers"
	"github.com/shashimalcse/tiny-is/internal/server/middlewares"
	"github.com/shashimalcse/tiny-is/internal/server/utils"
)

func RegisterAuthnRoutes(mux *utils.OrgServeMux, authnService authn.AuthnService) {
	handler := handlers.NewAuthnHandler(authnService)

	mux.HandleFunc("POST /login", middlewares.ErrorMiddleware(handler.Login))
	mux.HandleFunc("GET /login", middlewares.ErrorMiddleware(handler.GetLoginForm))
}
