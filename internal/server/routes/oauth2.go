package routes

import (
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/oauth2"
	"github.com/shashimalcse/tiny-is/internal/server/handlers"
	tinyhttp "github.com/shashimalcse/tiny-is/internal/server/http"
	"github.com/shashimalcse/tiny-is/internal/server/middlewares"
)

func RegisterOAuth2Routes(mux *tinyhttp.TinyServeMux, oauth2Service oauth2.OAuth2Service) {
	handler := handlers.NewOAuth2Handler(oauth2Service)
	authorizeHandler := middlewares.ChainMiddleware(handler.Authorize, middlewares.ErrorMiddleware())
	tokenHandler := middlewares.ChainMiddleware(handler.Token, middlewares.ErrorMiddleware())
	revokeHandler := middlewares.ChainMiddleware(handler.Revoke, middlewares.ErrorMiddleware())
	metadataHandler := middlewares.ChainMiddleware(handler.Metadata, middlewares.ErrorMiddleware())
	mux.HandleFunc("GET /authorize", func(w http.ResponseWriter, r *http.Request) { authorizeHandler(w, r) })
	mux.HandleFunc("POST /token", func(w http.ResponseWriter, r *http.Request) { tokenHandler(w, r) })
	mux.HandleFunc("POST /revoke", func(w http.ResponseWriter, r *http.Request) { revokeHandler(w, r) })
	mux.HandleFunc("GET /.well-known/oauth-authorization-server", func(w http.ResponseWriter, r *http.Request) { metadataHandler(w, r) })
}
