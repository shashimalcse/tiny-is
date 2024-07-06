package organization

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/shashimalcse/tiny-is/internal/organization/models"
)

type OrganizationRepository interface {
	GetOrganizationByName(ctx context.Context, name string) (models.Organization, error)
	CreateOrganization(ctx context.Context, organization models.Organization) error
}

type organizationRepository struct {
	db *sqlx.DB
}

func NewOrganizationRepository(db *sqlx.DB) OrganizationRepository {
	return &organizationRepository{
		db: db,
	}
}

func (r *organizationRepository) GetOrganizationByName(ctx context.Context, name string) (models.Organization, error) {
	var organization models.Organization
	err := r.db.Get(&organization, "SELECT id, name FROM organization WHERE name=$1", name)
	if err != nil {
		return models.Organization{}, err
	}
	return organization, nil
}

func (r *organizationRepository) CreateOrganization(ctx context.Context, organization models.Organization) error {
	_, err := r.db.Exec("INSERT INTO organization (name) VALUES ($1)", organization.Name)
	if err != nil {
		return err
	}
	return nil
}
