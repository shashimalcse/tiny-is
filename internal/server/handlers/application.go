package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/shashimalcse/tiny-is/internal/application"
	app_models "github.com/shashimalcse/tiny-is/internal/application/models"
	"github.com/shashimalcse/tiny-is/internal/server/models"
)

type ApplicationHandler struct {
	applicationService *application.ApplicationService
}

func NewApplicationHandler(applicationService *application.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		applicationService: applicationService,
	}
}

func (handler ApplicationHandler) GetApplications(w http.ResponseWriter, r *http.Request) {
	applications, err := handler.applicationService.GetApplications()
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
	application, err := handler.applicationService.GetApplicationByID(applicationId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.GetApplicationResponse(application))
}

func (handler ApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {

	var applicationRequest models.ApplicationCreateRequest
	err := json.NewDecoder(r.Body).Decode(&applicationRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	appId := uuid.New().String()
	clientId, err := handler.applicationService.GenerateClientId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	clientSecret, err := handler.applicationService.GenerateClientSecreat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	application := app_models.Application{
		Id:           appId,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectUri:  applicationRequest.RedirectUri,
		GrantTypes:   strings.Join(applicationRequest.GrantTypes, " "),
	}
	err = handler.applicationService.CreateApplication(application)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
