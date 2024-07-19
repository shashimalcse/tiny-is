package routes

import (
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/application"
	"github.com/shashimalcse/tiny-is/internal/config"
	"github.com/shashimalcse/tiny-is/internal/security"
	"github.com/shashimalcse/tiny-is/internal/server/handlers"
	tinyhttp "github.com/shashimalcse/tiny-is/internal/server/http"
	"github.com/shashimalcse/tiny-is/internal/server/middlewares"
)

func RegisterApplicationRoutes(mux *tinyhttp.TinyServeMux, cfg *config.Config, keyManager *security.KeyManager, applicationService application.ApplicationService) {
	handler := handlers.NewApplicationHandler(applicationService)
	getApplicationsHandler := middlewares.ChainMiddleware(handler.GetApplications, middlewares.ErrorMiddleware(), middlewares.JWTMiddleware(cfg, keyManager))
	createApplicationHandler := middlewares.ChainMiddleware(handler.CreateApplication, middlewares.ErrorMiddleware(), middlewares.JWTMiddleware(cfg, keyManager))
	updateApplicationHandler := middlewares.ChainMiddleware(handler.UpdateApplication, middlewares.ErrorMiddleware(), middlewares.JWTMiddleware(cfg, keyManager))
	deleteApplicationHandler := middlewares.ChainMiddleware(handler.DeleteApplication, middlewares.ErrorMiddleware(), middlewares.JWTMiddleware(cfg, keyManager))
	mux.HandleFunc("GET /applications", func(w http.ResponseWriter, r *http.Request) { getApplicationsHandler(w, r) })
	mux.HandleFunc("POST /applications", func(w http.ResponseWriter, r *http.Request) { createApplicationHandler(w, r) })
	mux.HandleFunc("PUT /applications/{id}", func(w http.ResponseWriter, r *http.Request) { updateApplicationHandler(w, r) })
	mux.HandleFunc("DELETE /applications/{id}", func(w http.ResponseWriter, r *http.Request) { deleteApplicationHandler(w, r) })
}
