package application

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/shashimalcse/tiny-is/internal/application/models"
	"github.com/shashimalcse/tiny-is/internal/cache"
)

type ApplicationService interface {
	GetApplications(ctx context.Context, orgId string) ([]models.Application, error)
	GetApplicationByID(ctx context.Context, id, orgId string) (models.Application, error)
	CreateApplication(ctx context.Context, application models.Application) error
	UpdateApplication(ctx context.Context, id, orgId string, application models.Application) error
	DeleteApplication(ctx context.Context, id, orgId string) error
	ValidateClientId(ctx context.Context, clientId, orgId string) (bool, error)
	ValidateClientSecret(ctx context.Context, clientId, clientSecret, orgId string) (bool, error)
	ValidateRedirectUri(ctx context.Context, clientId, redirectUri, orgId string) (bool, error)
}

type applicationService struct {
	cacheService *cache.CacheService
	repo         ApplicationRepository
}

func NewApplicationService(cacheService *cache.CacheService, repo ApplicationRepository) ApplicationService {
	return &applicationService{
		cacheService: cacheService,
		repo:         repo,
	}
}

func (s *applicationService) GetApplications(ctx context.Context, orgId string) ([]models.Application, error) {
	return s.repo.GetApplications(ctx, orgId)
}

func (s *applicationService) GetApplicationByID(ctx context.Context, id, orgId string) (models.Application, error) {
	return s.repo.GetApplicationByID(ctx, id, orgId)
}

func (s *applicationService) CreateApplication(ctx context.Context, application models.Application) error {
	appId := uuid.New().String()
	clientId, err := GenerateClientId()
	if err != nil {
		return err
	}
	clientSecret, err := GenerateClientSecreat()
	if err != nil {
		return err
	}
	application.Id = appId
	application.ClientId = clientId
	application.ClientSecret = clientSecret
	return s.repo.CreateApplication(ctx, application)
}

func (s *applicationService) UpdateApplication(ctx context.Context, id, orgId string, application models.Application) error {
	_, err := s.GetApplicationByID(ctx, id, orgId)
	if err != nil {
		return err
	}
	return s.repo.UpdateApplication(ctx, id, application)
}

func (s *applicationService) DeleteApplication(ctx context.Context, id, orgId string) error {
	_, err := s.GetApplicationByID(ctx, id, orgId)
	if err != nil {
		return err
	}
	return s.repo.DeleteApplication(ctx, id, orgId)
}

func (s *applicationService) ValidateClientId(ctx context.Context, clientId, orgId string) (bool, error) {
	return s.repo.ValidateClientId(ctx, clientId, orgId)
}

func (s *applicationService) ValidateClientSecret(ctx context.Context, clientId, clientSecret, orgId string) (bool, error) {
	return s.repo.ValidateClientSecret(ctx, clientId, clientSecret, orgId)
}

func (s *applicationService) ValidateRedirectUri(ctx context.Context, clientId, redirectUri, orgId string) (bool, error) {
	return s.repo.ValidateRedirectUri(ctx, clientId, redirectUri, orgId)
}

func GenerateClientId() (string, error) {
	bytes := make([]byte, 10)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func GenerateClientSecreat() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
