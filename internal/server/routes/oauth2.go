package routes

import (
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/oauth2"
	"github.com/shashimalcse/tiny-is/internal/server/handlers"
)

func RegisterOAuth2Routes(mux *http.ServeMux, oauth2 *oauth2.OAuth2) {
	handler := handlers.NewOAuth2Handler(oauth2)

	mux.HandleFunc("GET /authorize", func(w http.ResponseWriter, r *http.Request) {
		handler.Authorize(w, r)
	})
	mux.HandleFunc("POST /token", func(w http.ResponseWriter, r *http.Request) {
		handler.Token(w, r)
	})
}
