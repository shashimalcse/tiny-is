package organization

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/organization/models"
)

type OrganizationService interface {
	CreateOrganization(ctx context.Context, organization models.Organization) (models.Organization, error)
	DeleteOrganization(ctx context.Context, orgId string) error
	IsOrganizationExistByName(ctx context.Context, name string) (bool, error)
	GetOrganizationByName(ctx context.Context, name string) (models.Organization, error)
	GetOrganizationById(ctx context.Context, orgId string) (models.Organization, error)
}

type organizationService struct {
	cacheService cache.CacheService
	repo         OrganizationRepository
}

func NewOrganizationService(cacheService cache.CacheService, repo OrganizationRepository) OrganizationService {
	return &organizationService{
		cacheService: cacheService,
		repo:         repo,
	}
}

func (s *organizationService) CreateOrganization(ctx context.Context, organization models.Organization) (models.Organization, error) {
	if organization.Name == "" {
		return models.Organization{}, errors.New("organization name is required")
	}
	orgId := uuid.New().String()
	organization.Id = orgId
	err := s.repo.CreateOrganization(ctx, organization)
	if err != nil {
		return models.Organization{}, err
	}
	s.cacheService.SetOrganization(organization)
	return organization, nil
}

func (s *organizationService) IsOrganizationExistByName(ctx context.Context, name string) (bool, error) {
	exists, err := s.repo.IsOrganizationExistByName(ctx, name)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	return false, nil
}

func (s *organizationService) GetOrganizationByName(ctx context.Context, name string) (models.Organization, error) {
	organization, found := s.cacheService.GetOrganizationByName(name)
	if found {
		return organization, nil
	}
	organization, err := s.repo.GetOrganizationByName(ctx, name)
	if err != nil {
		return models.Organization{}, err
	}
	s.cacheService.SetOrganization(organization)
	return organization, nil
}

func (s *organizationService) GetOrganizationById(ctx context.Context, orgId string) (models.Organization, error) {
	organization, found := s.cacheService.GetOrganizationById(orgId)
	if found {
		return organization, nil
	}
	organization, err := s.repo.GetOrganizationByID(ctx, orgId)
	if err != nil {
		return models.Organization{}, err
	}
	s.cacheService.SetOrganization(organization)
	return organization, nil

}

func (s *organizationService) DeleteOrganization(ctx context.Context, orgId string) error {
	org, err := s.GetOrganizationById(ctx, orgId)
	if err != nil {
		return err
	}
	err = s.repo.DeleteOrganization(ctx, orgId)
	if err != nil {
		return err
	}
	s.cacheService.DeleteOrganizationById(orgId)
	s.cacheService.DeleteOrganizationByName(org.Name)
	return nil
}
