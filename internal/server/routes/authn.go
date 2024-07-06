package routes

import (
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/authn"
	"github.com/shashimalcse/tiny-is/internal/server/handlers"
	"github.com/shashimalcse/tiny-is/internal/server/utils"
)

func RegisterAuthnRoutes(mux *utils.OrgServeMux, authn *authn.Authn) {
	handler := handlers.NewAuthnHandler(authn)

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		handler.Login(w, r)
	})
	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		handler.LoginForm(w, r)
	})
}
