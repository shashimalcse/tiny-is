package cache

import (
	"github.com/patrickmn/go-cache"
	"github.com/shashimalcse/tiny-is/internal/organization/models"
)

var (
	organization_name_cache_prefix = "organization_name_"
	organization_id_cache_prefix   = "organization_id_"
)

func (cacheService *CacheService) GetOrganizationByName(name string) (models.Organization, bool) {
	authorizeContext, found := cacheService.c.Get(organization_name_cache_prefix + name)
	if !found {
		return models.Organization{}, false
	}
	return authorizeContext.(models.Organization), true
}

func (cacheService *CacheService) GetOrganizationById(id string) (models.Organization, bool) {
	authorizeContext, found := cacheService.c.Get(organization_id_cache_prefix + id)
	if !found {
		return models.Organization{}, false
	}
	return authorizeContext.(models.Organization), true
}

func (cacheService *CacheService) SetOrganizationByName(organization models.Organization) {
	cacheService.c.Set(organization_name_cache_prefix+organization.Name, organization, cache.DefaultExpiration)
}
