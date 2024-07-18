package utils

import (
	"context"

	"github.com/shashimalcse/tiny-is/internal/application"
	"github.com/shashimalcse/tiny-is/internal/config"
	"github.com/shashimalcse/tiny-is/internal/organization"
	"github.com/shashimalcse/tiny-is/internal/organization/models"
	"github.com/shashimalcse/tiny-is/internal/user"
	user_models "github.com/shashimalcse/tiny-is/internal/user/models"
)

func InitServer(cfg *config.Config, organizationService organization.OrganizationService, applicationService application.ApplicationService, userService user.UserService) error {

	//check if super organization exists
	exists, err := organizationService.IsOrganizationExistByName(context.Background(), cfg.SuperOrganization.Name)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	super_organization := models.Organization{
		Name: cfg.SuperOrganization.Name,
	}
	err = organizationService.CreateOrganization(context.Background(), super_organization)
	if err != nil {
		return err
	}
	super_org, err := organizationService.GetOrganizationByName(context.Background(), super_organization.Name)
	if err != nil {
		return err
	}
	admin := user_models.User{
		Username:       cfg.SuperOrganization.Admin.Username,
		Password:       cfg.SuperOrganization.Admin.Password,
		Email:          "admin",
		OrganizationId: super_org.Id,
	}
	err = userService.CreateUser(context.Background(), admin)
	if err != nil {
		return err
	}
	return nil
}
