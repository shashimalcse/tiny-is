package organization

import (
	"context"

	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/organization/models"
)

type OrganizationService interface {
	CreateOrganization(ctx context.Context, organization models.Organization) error
	IsOrganizationExistByName(ctx context.Context, name string) (bool, error)
	GetOrganizationByName(ctx context.Context, name string) (models.Organization, error)
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

func (s *organizationService) CreateOrganization(ctx context.Context, organization models.Organization) error {
	
	err := s.repo.CreateOrganization(ctx, organization)
	if err != nil {
		return err
	}
	s.cacheService.SetOrganization(organization)
	return nil
}

func (s *organizationService) IsOrganizationExistByName(ctx context.Context, name string) (bool, error) {

	organization, err := s.GetOrganizationByName(ctx, name)
	if err != nil {
		return false, err
	}
	if organization.Name != name {
		return false, nil
	}
	return true, nil
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
