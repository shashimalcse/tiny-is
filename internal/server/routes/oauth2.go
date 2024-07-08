package routes

import (
	"github.com/shashimalcse/tiny-is/internal/oauth2"
	"github.com/shashimalcse/tiny-is/internal/server/handlers"
	"github.com/shashimalcse/tiny-is/internal/server/middlewares"
	"github.com/shashimalcse/tiny-is/internal/server/utils"
)

func RegisterOAuth2Routes(mux *utils.OrgServeMux, oauth2Service oauth2.OAuth2Service) {
	handler := handlers.NewOAuth2Handler(oauth2Service)

	mux.HandleFunc("GET /authorize", middlewares.ErrorMiddleware(handler.Authorize))
	mux.HandleFunc("POST /token", middlewares.ErrorMiddleware(handler.Token))
}
