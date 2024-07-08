package cache

import (
	"testing"

	"github.com/shashimalcse/tiny-is/internal/oauth2/models"
	org_models "github.com/shashimalcse/tiny-is/internal/organization/models"
	server_models "github.com/shashimalcse/tiny-is/internal/server/models"
)

func TestOAuth2AuthorizeContextToCacheBySessionDataKey(t *testing.T) {
	cacheService := NewCacheService()
	testSessionDataKey := "test-session-data-key"
	testAuthorizeContext := models.OAuth2AuthorizeContext{
		OAuth2AuthorizeRequest: server_models.OAuth2AuthorizeRequest{
			ClientId:       "test-client-id",
			OrganizationId: "test-organization-id",
		},
	}
	cacheService.AddOAuth2AuthorizeContextToCacheBySessionDataKey(testSessionDataKey, testAuthorizeContext)

	_, found := cacheService.GetOAuth2AuthorizeContextFromCacheBySessionDataKey(testSessionDataKey)
	if !found {
		t.Errorf("Expected to find the authorize context from cache")
	}
	cacheService.DeleteOAuth2AuthorizeContextFromCacheBySessionDataKey(testSessionDataKey)
	_, found = cacheService.GetOAuth2AuthorizeContextFromCacheBySessionDataKey(testSessionDataKey)
	if found {
		t.Errorf("Expected not to find the authorize context from cache")
	}
}

func TestOAuth2AuthorizeContextToCacheByAuthCode(t *testing.T) {
	cacheService := NewCacheService()
	testCode := "test-code"
	testAuthorizeContext := models.OAuth2AuthorizeContext{
		OAuth2AuthorizeRequest: server_models.OAuth2AuthorizeRequest{
			ClientId:       "test-client-id",
			OrganizationId: "test-organization-id",
		},
	}
	cacheService.AddOAuth2AuthorizeContextToCacheByAuthCode(testCode, testAuthorizeContext)

	_, found := cacheService.GetOAuth2AuthorizeContextFromCacheByAuthCode(testCode)
	if !found {
		t.Errorf("Expected to find the authorize context from cache")
	}
	cacheService.DeleteOAuth2AuthorizeContextFromCacheByAuthCode(testCode)
	_, found = cacheService.GetOAuth2AuthorizeContextFromCacheByAuthCode(testCode)
	if found {
		t.Errorf("Expected not to find the authorize context from cache")
	}
}

func TestOrganizationToCache(t *testing.T) {
	cacheService := NewCacheService()
	testOrganization := org_models.Organization{
		Id:   "test-id",
		Name: "test-name",
	}
	cacheService.SetOrganization(testOrganization)

	_, found := cacheService.GetOrganizationById(testOrganization.Id)
	if !found {
		t.Errorf("Expected to find the organization from cache")
	}
	_, found = cacheService.GetOrganizationByName(testOrganization.Name)
	if !found {
		t.Errorf("Expected to find the organization from cache")
	}

	cacheService.DeleteOrganizationById(testOrganization.Id)
	_, found = cacheService.GetOrganizationById(testOrganization.Id)
	if found {
		t.Errorf("Expected not to find the organization from cache")
	}
	cacheService.DeleteOrganizationByName(testOrganization.Name)
	_, found = cacheService.GetOrganizationByName(testOrganization.Name)
	if found {
		t.Errorf("Expected not to find the organization from cache")
	}
}
