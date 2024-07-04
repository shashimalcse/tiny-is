package application

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/jmoiron/sqlx"
	"github.com/shashimalcse/tiny-is/internal/application/models"
	"github.com/shashimalcse/tiny-is/internal/cache"
)

type ApplicationService struct {
	cacheService *cache.CacheService
	db    *sqlx.DB
}

func NewApplicationService(cacheService *cache.CacheService, db *sqlx.DB) *ApplicationService {
	return &ApplicationService{
		cacheService: cacheService,
		db:    db,
	}
}

func (a ApplicationService) GetApplications() ([]models.Application, error) {
	var applications []models.Application
	err := a.db.Select(&applications, "SELECT * FROM applications")
	if err != nil {
		return nil, err
	}
	return applications, nil
}

func (a ApplicationService) GetApplicationByID(id string) (models.Application, error) {
	var application models.Application
	err := a.db.Get(&application, "SELECT * FROM applications WHERE id=$1", id)
	if err != nil {
		return models.Application{}, err
	}
	return application, nil
}

func (a ApplicationService) CreateApplication(application models.Application) error {
	_, err := a.db.NamedExec("INSERT INTO applications (id, client_id, client_secret, redirect_uri, grant_types) VALUES (:id, :client_id, :client_secret, :redirect_uri, :grant_types)", application)
	if err != nil {
		return err
	}
	return nil
}

func (a ApplicationService) ValidateClientId(clientId string) (bool, error) {
	var count int
	err := a.db.Get(&count, "SELECT COUNT(*) FROM applications WHERE client_id=$1", clientId)
	return count == 1, err
}

func (a ApplicationService) ValidateClientSecret(clientId, clientSecret string) (bool, error) {
	var count int
	err := a.db.Get(&count, "SELECT COUNT(*) FROM applications WHERE client_id=$1 AND client_secret=$2", clientId, clientSecret)
	return count == 1, err
}

func (a ApplicationService) ValidateRedirectUri(clientId, redirectUri string) (bool, error) {
	var count int
	err := a.db.Get(&count, "SELECT COUNT(*) FROM applications WHERE client_id=$1 AND redirect_uri=$2", clientId, redirectUri)
	return count == 1, err
}

func (a ApplicationService) GenerateClientId() (string, error) {
	bytes := make([]byte, 10)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (a ApplicationService) GenerateClientSecreat() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
