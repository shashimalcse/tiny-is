package models

import (
	"github.com/shashimalcse/tiny-is/internal/application/models"
)

type ApplicationResponse struct {
	Id           string   `json:"id"`
	ClientId     string   `json:"client_id,omitempty"`
	ClientSecret string   `json:"client_secret,omitempty"`
	RedirectUri  string   `json:"redirect_uri,omitempty"`
	GrantTypes   []string `json:"grant_types,omitempty"`
}

type ApplicationCreateRequest struct {
	RedirectUri string   `json:"redirect_uri,omitempty"`
	GrantTypes  []string `json:"grant_types,omitempty"`
}

func GetApplicationResponse(application models.Application) ApplicationResponse {
	return ApplicationResponse{
		Id:           application.Id,
		ClientId:     application.ClientId,
		ClientSecret: application.ClientSecret,
		RedirectUri:  application.RedirectUri,
		GrantTypes:   application.GetGrantTypes(),
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
