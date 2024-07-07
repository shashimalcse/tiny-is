package routes

import (
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/application"
	"github.com/shashimalcse/tiny-is/internal/server/handlers"
	"github.com/shashimalcse/tiny-is/internal/server/utils"
)

func RegisterApplicationRoutes(mux *utils.OrgServeMux, applicationService application.ApplicationService) {
	handler := handlers.NewApplicationHandler(applicationService)

	mux.HandleFunc("GET /applications", func(w http.ResponseWriter, r *http.Request) {
		handler.GetApplications(w, r)
	})
	mux.HandleFunc("GET /applications/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetApplicationByID(w, r)
	})
	mux.HandleFunc("POST /applications", func(w http.ResponseWriter, r *http.Request) {
		handler.CreateApplication(w, r)
	})
}
