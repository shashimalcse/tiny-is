package application

import (
	"context"
	"encoding/json"
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
	redirectURIsJSON, err := json.Marshal(application.RedirectUris)
	if err != nil {
		return err
	}
	_, err = r.db.NamedExec("INSERT INTO application (id, name, organization_id, client_id, client_secret, redirect_uris) VALUES (:id, :name, :organization_id, :client_id, :client_secret, :redirect_uris)", map[string]interface{}{
		"id":              application.Id,
		"name":            application.Name,
		"organization_id": application.OrganizationId,
		"client_id":       application.ClientId,
		"client_secret":   application.ClientSecret,
		"redirect_uris":   string(redirectURIsJSON),
	})
	if err != nil {
		return err
	}
	grantTypeIDs, err := r.getGrantIdsByNames(application.GrantTypes)
	if err != nil {
		return err
	}
	insertQuery := "INSERT INTO client_grant_type (application_id, grant_type_id) VALUES (:application_id, :grant_type_id)"
	var clientGrantTypes []map[string]interface{}
	for _, grantTypeID := range grantTypeIDs {
		clientGrantTypes = append(clientGrantTypes, map[string]interface{}{
			"application_id": application.Id,
			"grant_type_id":  grantTypeID,
		})
	}
	_, err = r.db.NamedExec(insertQuery, clientGrantTypes)
	if err != nil {
		return err
	}
	return nil
}

func (r *applicationRepository) UpdateApplication(ctx context.Context, id string, updateApplication models.Application) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
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
	if len(updateFields) > 0 {
		updateQuery += strings.Join(updateFields, ", ") + fmt.Sprintf(" WHERE id = $%d", paramCount)
		updateValues = append(updateValues, id)

		_, err = tx.ExecContext(ctx, updateQuery, updateValues...)
		if err != nil {
			return err
		}
	}
	if updateApplication.GrantTypes != nil {
		_, err = tx.ExecContext(ctx, "DELETE FROM client_grant_type WHERE application_id = $1", id)
		if err != nil {
			return err
		}
		grantTypeIDs, err := r.getGrantIdsByNames(updateApplication.GrantTypes)
		if err != nil {
			return err
		}
		insertQuery := "INSERT INTO client_grant_type (application_id, grant_type_id) VALUES ($1, $2)"
		for _, grantTypeID := range grantTypeIDs {
			_, err = tx.ExecContext(ctx, insertQuery, id, grantTypeID)
			if err != nil {
				return err
			}
		}
	}
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
	query := `SELECT COUNT(*) FROM application WHERE client_id = ? AND organization_id = ? AND json_array_length(json_extract(redirect_uris, '$')) > 0 
        AND ? IN (SELECT value FROM json_each(redirect_uris))
    `
	err := r.db.GetContext(ctx, &count, query, clientId, orgId, redirectUri)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *applicationRepository) getGrantIdsByNames(grantTypeNames []string) ([]string, error) {
	query := "SELECT id FROM grant_type WHERE name IN (?)"
	var grantTypeIDs []string
	query, args, err := sqlx.In(query, grantTypeNames)
	if err != nil {
		return []string{}, err
	}
	query = r.db.Rebind(query)
	err = r.db.Select(&grantTypeIDs, query, args...)
	if err != nil {
		return []string{}, err
	}
	return grantTypeIDs, nil
}
