package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/application"
	app_models "github.com/shashimalcse/tiny-is/internal/application/models"
	"github.com/shashimalcse/tiny-is/internal/server/models"
)

type ApplicationHandler struct {
	applicationService application.ApplicationService
}

func NewApplicationHandler(applicationService application.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		applicationService: applicationService,
	}
}

func (handler ApplicationHandler) GetApplications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orgId := r.Header.Get("org_id")
	if orgId == "" {
		http.Error(w, "Organization not found!", http.StatusNotFound)
		return
	}
	applications, err := handler.applicationService.GetApplications(ctx, orgId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.GetApplicationResponses(applications))
}

func (handler ApplicationHandler) GetApplicationByID(w http.ResponseWriter, r *http.Request) {

	applicationId := r.URL.Query().Get("id")
	orgId := r.Header.Get("org_id")
	if orgId == "" {
		http.Error(w, "Organization not found!", http.StatusNotFound)
		return
	}
	ctx := r.Context()
	application, err := handler.applicationService.GetApplicationByID(ctx, applicationId, orgId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.GetApplicationResponse(application))
}

func (handler ApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {

	orgId := r.Header.Get("org_id")
	if orgId == "" {
		http.Error(w, "Organization not found!", http.StatusNotFound)
		return
	}
	var applicationRequest models.ApplicationCreateRequest
	err := json.NewDecoder(r.Body).Decode(&applicationRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	application := app_models.Application{
		Name:           applicationRequest.Name,
		RedirectUris:   applicationRequest.RedirectUris,
		GrantTypes:     applicationRequest.GrantTypes,
		OrganizationId: orgId,
	}
	ctx := r.Context()
	err = handler.applicationService.CreateApplication(ctx, application)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
