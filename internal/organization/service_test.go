package organization

import (
	"context"
	"testing"

	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/organization/models"
)

func NewMockOrganizationService() OrganizationService {
	repo := NewMockOrganizationRepository()
	cache := cache.NewCacheService()
	return &organizationService{repo: repo, cacheService: cache}
}

func TestServiceCreateOrganization(t *testing.T) {
	organizationService := NewMockOrganizationService()
	organization := models.Organization{Name: "org-1"}
	org, err := organizationService.CreateOrganization(context.Background(), organization)
	if err != nil {
		t.Errorf("failed to create organization: %v", err)
	}
	err = organizationService.DeleteOrganization(context.Background(), org.Id)
	if err != nil {
		t.Errorf("failed to delete organization: %v", err)
	}
}

func TestServiceCreateOrganizationWithNullName(t *testing.T) {
	organizationService := NewMockOrganizationService()
	organization := models.Organization{}
	_, err := organizationService.CreateOrganization(context.Background(), organization)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestServiceCreateOrganizationWithEmptyName(t *testing.T) {
	organizationService := NewMockOrganizationService()
	organization := models.Organization{Name: ""}
	_, err := organizationService.CreateOrganization(context.Background(), organization)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestServiceIsOrganizationExistByName(t *testing.T) {
	organizationService := NewMockOrganizationService()
	organization := models.Organization{Name: "org-3"}
	org, err := organizationService.CreateOrganization(context.Background(), organization)
	if err != nil {
		t.Errorf("failed to create organization: %v", err)
	}
	exists, err := organizationService.IsOrganizationExistByName(context.Background(), organization.Name)
	if err != nil {
		t.Errorf("failed to check organization existence: %v", err)
	}
	if !exists {
		t.Errorf("expected organization to exist")
	}
	err = organizationService.DeleteOrganization(context.Background(), org.Id)
	if err != nil {
		t.Errorf("failed to delete organization: %v", err)
	}
}

func TestServiceIsOrganizationExistByNameWithNonExistingOrganization(t *testing.T) {
	organizationService := NewMockOrganizationService()
	exists, err := organizationService.IsOrganizationExistByName(context.Background(), "non-existing")
	if err != nil {
		t.Errorf("failed to check organization existence: %v", err)
	}
	if exists {
		t.Errorf("expected organization not to exist")
	}
}

func TestServiceGetOrganizationByName(t *testing.T) {
	organizationService := NewMockOrganizationService()
	organization := models.Organization{Name: "org-4"}
	_, err := organizationService.CreateOrganization(context.Background(), organization)
	if err != nil {
		t.Errorf("failed to create organization: %v", err)
	}
	org, err := organizationService.GetOrganizationByName(context.Background(), organization.Name)
	if err != nil {
		t.Errorf("failed to get organization: %v", err)
	}
	if org.Name != organization.Name {
		t.Errorf("expected organization name to be %s, got %s", organization.Name, org.Name)
	}
	err = organizationService.DeleteOrganization(context.Background(), org.Id)
	if err != nil {
		t.Errorf("failed to delete organization: %v", err)
	}
}

func TestServiceGetOrganizationByNameWithNonExistingOrganization(t *testing.T) {
	organizationService := NewMockOrganizationService()
	_, err := organizationService.GetOrganizationByName(context.Background(), "non-existing")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestServiceGetOrganizationById(t *testing.T) {
	organizationService := NewMockOrganizationService()
	organization := models.Organization{Name: "org-5"}
	org, err := organizationService.CreateOrganization(context.Background(), organization)
	if err != nil {
		t.Errorf("failed to create organization: %v", err)
	}
	_, err = organizationService.GetOrganizationById(context.Background(), org.Id)
	if err != nil {
		t.Errorf("failed to get organization: %v", err)
	}
	err = organizationService.DeleteOrganization(context.Background(), org.Id)
	if err != nil {
		t.Errorf("failed to delete organization: %v", err)
	}
}

func TestServiceGetOrganizationByIdWithNonExistingOrganization(t *testing.T) {
	organizationService := NewMockOrganizationService()
	_, err := organizationService.GetOrganizationById(context.Background(), "non-existing")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestServiceDeleteOrganization(t *testing.T) {
	organizationService := NewMockOrganizationService()
	organization := models.Organization{Name: "org-6"}
	org, err := organizationService.CreateOrganization(context.Background(), organization)
	if err != nil {
		t.Errorf("failed to create organization: %v", err)
	}
	err = organizationService.DeleteOrganization(context.Background(), org.Id)
	if err != nil {
		t.Errorf("failed to delete organization: %v", err)
	}
}

func TestServiceDeleteOrganizationWithNonExistingOrganization(t *testing.T) {
	organizationService := NewMockOrganizationService()
	err := organizationService.DeleteOrganization(context.Background(), "non-existing")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
