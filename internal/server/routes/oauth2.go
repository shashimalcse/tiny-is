package routes

import (
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/oauth2"
	"github.com/shashimalcse/tiny-is/internal/server/handlers"
	"github.com/shashimalcse/tiny-is/internal/server/utils"
)

func RegisterOAuth2Routes(mux *utils.OrgServeMux, oauth2Service oauth2.OAuth2Service) {
	handler := handlers.NewOAuth2Handler(oauth2Service)

	mux.HandleFunc("GET /authorize", func(w http.ResponseWriter, r *http.Request) {
		handler.Authorize(w, r)
	})
	mux.HandleFunc("POST /token", func(w http.ResponseWriter, r *http.Request) {
		handler.Token(w, r)
	})
}
