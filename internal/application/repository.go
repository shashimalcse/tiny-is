package application

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/shashimalcse/tiny-is/internal/application/models"
)

type ApplicationRepository interface {
	GetApplications(ctx context.Context, orgId string) ([]models.Application, error)
	GetApplicationByID(ctx context.Context, id, orgId string) (models.Application, error)
	CreateApplication(ctx context.Context, application models.Application) error
	UpdateApplication(ctx context.Context, id string, updateApplication models.Application) error
	DeleteApplication(ctx context.Context, id, orgId string) error
	ValidateClientId(ctx context.Context, clientId, orgId string) (bool, error)
	ValidateClientSecret(ctx context.Context, clientId, clientSecret, orgId string) (bool, error)
	ValidateRedirectUri(ctx context.Context, clientId, redirectUri, orgId string) (bool, error)
}

type applicationRepository struct {
	db *sqlx.DB
}

func NewApplicationRepository(db *sqlx.DB) ApplicationRepository {
	return &applicationRepository{
		db: db,
	}
}

func (r *applicationRepository) GetApplications(ctx context.Context, orgId string) ([]models.Application, error) {
	var applications []models.Application
	err := r.db.Select(&applications, "SELECT * FROM application WHERE organization_id=$1", orgId)
	if err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *applicationRepository) GetApplicationByID(ctx context.Context, id, orgId string) (models.Application, error) {
	var application models.Application
	err := r.db.Get(&application, "SELECT * FROM application WHERE id=$1 AND organization_id=$2", id, orgId)
	if err != nil {
		return models.Application{}, err
	}
	grantTypes, err := r.GetApplicationGrant(ctx, id)
	if err != nil {
		return models.Application{}, err
	}
	application.GrantTypes = grantTypes
	return application, nil
}

func (r *applicationRepository) GetApplicationGrant(ctx context.Context, applicationID string) ([]string, error) {
	var grantTypes []string
	query := `
		SELECT gt.*
		FROM grant_type gt
		INNER JOIN client_grant_type cgt ON gt.id = cgt.grant_type_id
		WHERE cgt.application_id = $1
	`
	err := r.db.Select(&grantTypes, query, applicationID)
	if err != nil {
		return nil, err
	}
	return grantTypes, nil
}

func (r *applicationRepository) CreateApplication(ctx context.Context, application models.Application) error {
	_, err := r.db.NamedExec("INSERT INTO application (id, name, organization_id, client_id, client_secret, redirect_uris) VALUES (:id, :name, :organization_id, :client_id, :client_secret, :redirect_uris)", map[string]interface{}{
		"id":              application.Id,
		"name":            application.Name,
		"organization_id": application.OrganizationId,
		"client_id":       application.ClientId,
		"client_secret":   application.ClientSecret,
		"redirect_uris":   pq.Array(application.RedirectUris),
	})
	if err != nil {
		return err
	}
	// Query to get grant_type ids from names
	grantTypeIDs, err := r.getGrantIdsByNames(application.GrantTypes)
	if err != nil {
		return err
	}
	// Prepare the batch insert for client_grant_type table
	insertQuery := "INSERT INTO client_grant_type (application_id, grant_type_id) VALUES (:application_id, :grant_type_id)"
	var clientGrantTypes []map[string]interface{}
	for _, grantTypeID := range grantTypeIDs {
		clientGrantTypes = append(clientGrantTypes, map[string]interface{}{
			"application_id": application.Id,
			"grant_type_id":  grantTypeID,
		})
	}

	// Execute the batch insert
	_, err = r.db.NamedExec(insertQuery, clientGrantTypes)
	if err != nil {
		return err
	}
	return nil
}

func (r *applicationRepository) UpdateApplication(ctx context.Context, id string, updateApplication models.Application) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Prepare the update query
	updateQuery := "UPDATE application SET "
	updateFields := []string{}
	updateValues := []interface{}{}
	paramCount := 1

	if updateApplication.Name != "" {
		updateFields = append(updateFields, fmt.Sprintf("name = $%d", paramCount))
		updateValues = append(updateValues, updateApplication.Name)
		paramCount++
	}

	if updateApplication.RedirectUris != nil {
		updateFields = append(updateFields, fmt.Sprintf("redirect_uris = $%d", paramCount))
		updateValues = append(updateValues, pq.Array(updateApplication.RedirectUris))
		paramCount++
	}

	// If there are fields to update
	if len(updateFields) > 0 {
		updateQuery += strings.Join(updateFields, ", ") + fmt.Sprintf(" WHERE id = $%d", paramCount)
		updateValues = append(updateValues, id)

		_, err = tx.ExecContext(ctx, updateQuery, updateValues...)
		if err != nil {
			return err
		}
	}

	// Update grant types if provided
	if updateApplication.GrantTypes != nil {
		// Delete existing grant types
		_, err = tx.ExecContext(ctx, "DELETE FROM client_grant_type WHERE application_id = $1", id)
		if err != nil {
			return err
		}

		// Get new grant type IDs
		grantTypeIDs, err := r.getGrantIdsByNames(updateApplication.GrantTypes)
		if err != nil {
			return err
		}

		// Insert new grant types
		insertQuery := "INSERT INTO client_grant_type (application_id, grant_type_id) VALUES ($1, $2)"
		for _, grantTypeID := range grantTypeIDs {
			_, err = tx.ExecContext(ctx, insertQuery, id, grantTypeID)
			if err != nil {
				return err
			}
		}
	}

	// Commit the transaction
	return tx.Commit()
}

func (r *applicationRepository) DeleteApplication(ctx context.Context, id, orgId string) error {
	_, err := r.db.Exec("DELETE FROM application WHERE id=$1 AND organization_id=$2", id, orgId)
	return err
}

func (r *applicationRepository) ValidateClientId(ctx context.Context, clientId, orgId string) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM application WHERE client_id=$1 AND organization_id=$2", clientId, orgId)
	return count == 1, err
}

func (r *applicationRepository) ValidateClientSecret(ctx context.Context, clientId, clientSecret, orgId string) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM application WHERE client_id=$1 AND client_secret=$2 AND organization_id=$3", clientId, clientSecret, orgId)
	return count == 1, err
}

func (r *applicationRepository) ValidateRedirectUri(ctx context.Context, clientId, redirectUri, orgId string) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM application WHERE client_id=$1 AND $2 = ANY(redirect_uris) AND organization_id=$3", clientId, redirectUri, orgId)
	return count == 1, err
}

func (r *applicationRepository) getGrantIdsByNames(grantTypeNames []string) ([]string, error) {
	query := "SELECT id FROM grant_type WHERE name = ANY(:grant_type_names)"
	var grantTypeIDs []string

	// Using sqlx.Named to bind named parameters
	q, args, err := sqlx.Named(query, map[string]interface{}{
		"grant_type_names": pq.Array(grantTypeNames),
	})
	if err != nil {
		return []string{}, err
	}

	// Rebinding the query for the pq driver
	q, args, err = sqlx.In(q, args...)
	if err != nil {
		return []string{}, err
	}

	q = r.db.Rebind(q)
	err = r.db.Select(&grantTypeIDs, q, args...)
	if err != nil {
		return []string{}, err
	}
	return grantTypeIDs, nil
}
