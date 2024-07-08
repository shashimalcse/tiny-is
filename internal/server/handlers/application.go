package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/application"
	app_models "github.com/shashimalcse/tiny-is/internal/application/models"
	"github.com/shashimalcse/tiny-is/internal/server/middlewares"
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

func (handler ApplicationHandler) GetApplications(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	orgId := r.Header.Get("org_id")
	if orgId == "" {
		return middlewares.NewAPIError(http.StatusNotFound, "Organization not found!")
	}
	applications, err := handler.applicationService.GetApplications(ctx, orgId)
	if err != nil {
		return middlewares.NewAPIError(http.StatusInternalServerError, err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.GetApplicationResponses(applications))
	return nil
}

func (handler ApplicationHandler) GetApplicationByID(w http.ResponseWriter, r *http.Request) error {

	applicationId := r.URL.Query().Get("id")
	orgId := r.Header.Get("org_id")
	if orgId == "" {
		return middlewares.NewAPIError(http.StatusNotFound, "Organization not found!")
	}
	ctx := r.Context()
	application, err := handler.applicationService.GetApplicationByID(ctx, applicationId, orgId)
	if err != nil {
		return middlewares.NewAPIError(http.StatusInternalServerError, err.Error())
	}
	if application.Id == "" {
		return middlewares.NewAPIError(http.StatusNotFound, "Application not found!")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.GetApplicationResponse(application))
	return nil
}

func (handler ApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) error {

	orgId := r.Header.Get("org_id")
	if orgId == "" {
		return middlewares.NewAPIError(http.StatusNotFound, "Organization not found!")
	}
	var applicationRequest models.ApplicationCreateRequest
	err := json.NewDecoder(r.Body).Decode(&applicationRequest)
	if err != nil {
		return middlewares.NewAPIError(http.StatusBadRequest, "Invalid request payload")
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
		return middlewares.NewAPIError(http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}
