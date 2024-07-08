package models

import (
	"github.com/shashimalcse/tiny-is/internal/application/models"
)

type ApplicationResponse struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	ClientId     string   `json:"client_id,omitempty"`
	ClientSecret string   `json:"client_secret,omitempty"`
	RedirectUris []string `json:"redirect_uris,omitempty"`
	GrantTypes   []string `json:"grant_types,omitempty"`
}

type ApplicationCreateRequest struct {
	Name         string   `json:"name"`
	RedirectUris []string `json:"redirect_uris,omitempty"`
	GrantTypes   []string `json:"grant_types,omitempty"`
}

type ApplicationUpdateRequest struct {
	Name         string   `json:"name,omitempty"`
	RedirectUris []string `json:"redirect_uris,omitempty"`
	GrantTypes   []string `json:"grant_types,omitempty"`
}

func GetApplicationResponse(application models.Application) ApplicationResponse {
	return ApplicationResponse{
		Id:           application.Id,
		Name:         application.Name,
		ClientId:     application.ClientId,
		ClientSecret: application.ClientSecret,
		RedirectUris: application.RedirectUris,
		GrantTypes:   application.GrantTypes,
	}
}

func GetApplicationResponses(applications []models.Application) []ApplicationResponse {

	if applications == nil {
		return []ApplicationResponse{}
	}
	var applicationResponses []ApplicationResponse
	for _, application := range applications {
		applicationResponses = append(applicationResponses, GetApplicationResponse(application))
	}
	return applicationResponses
}
