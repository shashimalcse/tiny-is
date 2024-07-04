package models

import (
	"github.com/shashimalcse/tiny-is/internal/application/models"
)

type Application struct {
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

func GetApplicationResponse(application models.Application) Application {
	return Application{
		Id:           application.Id,
		ClientId:     application.ClientId,
		ClientSecret: application.ClientSecret,
		RedirectUri:  application.RedirectUri,
		GrantTypes:   application.GetGrantTypes(),
	}
}

func GetApplicationResponses(applications []models.Application) []Application {

	if applications == nil {
		return []Application{}
	}
	var applicationResponses []Application
	for _, application := range applications {
		applicationResponses = append(applicationResponses, GetApplicationResponse(application))
	}
	return applicationResponses
}
