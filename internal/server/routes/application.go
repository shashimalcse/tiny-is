package routes

import (
	"github.com/shashimalcse/tiny-is/internal/application"
	"github.com/shashimalcse/tiny-is/internal/server/handlers"
	"github.com/shashimalcse/tiny-is/internal/server/middlewares"
	"github.com/shashimalcse/tiny-is/internal/server/utils"
)

func RegisterApplicationRoutes(mux *utils.OrgServeMux, applicationService application.ApplicationService) {
	handler := handlers.NewApplicationHandler(applicationService)

	mux.HandleFunc("GET /applications", middlewares.ErrorMiddleware(handler.GetApplications))
	mux.HandleFunc("GET /applications/{id}", middlewares.ErrorMiddleware(handler.GetApplicationByID))
	mux.HandleFunc("POST /applications", middlewares.ErrorMiddleware(handler.CreateApplication))
}
