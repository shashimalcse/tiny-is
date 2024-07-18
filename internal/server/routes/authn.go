package routes

import (
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/authn"
	"github.com/shashimalcse/tiny-is/internal/server/handlers"
	tinyhttp "github.com/shashimalcse/tiny-is/internal/server/http"
	"github.com/shashimalcse/tiny-is/internal/server/middlewares"
)

func RegisterAuthnRoutes(mux *tinyhttp.TinyServeMux, authnService authn.AuthnService) {
	handler := handlers.NewAuthnHandler(authnService)
	loginHandler := middlewares.ChainMiddleware(handler.Login, middlewares.ErrorMiddleware())
	getLoginFormHandler := middlewares.ChainMiddleware(handler.GetLoginForm, middlewares.ErrorMiddleware())
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) { loginHandler(w, r) })
	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) { getLoginFormHandler(w, r) })
}
