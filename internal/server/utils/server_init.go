package utils

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/shashimalcse/tiny-is/internal/application"
	"github.com/shashimalcse/tiny-is/internal/config"
	"github.com/shashimalcse/tiny-is/internal/organization"
	"github.com/shashimalcse/tiny-is/internal/organization/models"
	"github.com/shashimalcse/tiny-is/internal/user"
	user_models "github.com/shashimalcse/tiny-is/internal/user/models"
)

func addGrantTypes(db *sqlx.DB) error {
	_, err := db.Exec("INSERT INTO grant_types (name) VALUES ('authorization_code')")
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO grant_types (name) VALUES ('refresh_token')")
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO grant_types (name) VALUES ('client_credentials')")
	if err != nil {
		return err
	}
	return nil
}

func InitServer(cfg *config.Config, db *sqlx.DB, organizationService organization.OrganizationService, applicationService application.ApplicationService, userService user.UserService) error {

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
	super_org, err := organizationService.CreateOrganization(context.Background(), super_organization)
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
	err = addGrantTypes(db)
	if err != nil {
		return err
	}
	return nil
}
