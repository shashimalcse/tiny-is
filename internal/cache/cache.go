package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/shashimalcse/tiny-is/internal/oauth2/models"
	org_models "github.com/shashimalcse/tiny-is/internal/organization/models"
)

var (
	organization_name_cache_prefix = "organization_name_"
	organization_id_cache_prefix   = "organization_id_"
)

type CacheService interface {
	AddOAuth2AuthorizeContextToCacheBySessionDataKey(sessionDataKey string, authorizeContext models.OAuth2AuthorizeContext)
	GetOAuth2AuthorizeContextFromCacheBySessionDataKey(sessionDataKey string) (models.OAuth2AuthorizeContext, bool)
	DeleteOAuth2AuthorizeContextFromCacheBySessionDataKey(sessionDataKey string)
	AddOAuth2AuthorizeContextToCacheByAuthCode(code string, authorizeContext models.OAuth2AuthorizeContext)
	GetOAuth2AuthorizeContextFromCacheByAuthCode(code string) (models.OAuth2AuthorizeContext, bool)
	DeleteOAuth2AuthorizeContextFromCacheByAuthCode(code string)
	GetOrganizationByName(name string) (org_models.Organization, bool)
	GetOrganizationById(id string) (org_models.Organization, bool)
	SetOrganization(organization org_models.Organization)
	DeleteOrganizationByName(name string)
	DeleteOrganizationById(id string)
}

type cacheService struct {
	c *cache.Cache
}

func NewCacheService() CacheService {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &cacheService{
		c: c,
	}
}

func (s *cacheService) AddOAuth2AuthorizeContextToCacheBySessionDataKey(sessionDataKey string, authorizeContext models.OAuth2AuthorizeContext) {
	s.c.Set(sessionDataKey, authorizeContext, cache.NoExpiration)
}

func (s *cacheService) GetOAuth2AuthorizeContextFromCacheBySessionDataKey(sessionDataKey string) (models.OAuth2AuthorizeContext, bool) {
	authorizeContext, found := s.c.Get(sessionDataKey)
	if !found {
		return models.OAuth2AuthorizeContext{}, false
	}
	return authorizeContext.(models.OAuth2AuthorizeContext), true
}

func (s *cacheService) DeleteOAuth2AuthorizeContextFromCacheBySessionDataKey(sessionDataKey string) {
	s.c.Delete(sessionDataKey)
}

func (s *cacheService) AddOAuth2AuthorizeContextToCacheByAuthCode(code string, authorizeContext models.OAuth2AuthorizeContext) {
	s.c.Set(code, authorizeContext, cache.DefaultExpiration)
}

func (s *cacheService) GetOAuth2AuthorizeContextFromCacheByAuthCode(code string) (models.OAuth2AuthorizeContext, bool) {
	authorizeContext, found := s.c.Get(code)
	if !found || authorizeContext == nil {
		return models.OAuth2AuthorizeContext{}, false
	}
	return authorizeContext.(models.OAuth2AuthorizeContext), true
}

func (s *cacheService) DeleteOAuth2AuthorizeContextFromCacheByAuthCode(code string) {
	s.c.Delete(code)
}

func (s *cacheService) GetOrganizationByName(name string) (org_models.Organization, bool) {
	authorizeContext, found := s.c.Get(organization_name_cache_prefix + name)
	if !found {
		return org_models.Organization{}, false
	}
	return authorizeContext.(org_models.Organization), true
}

func (s *cacheService) GetOrganizationById(id string) (org_models.Organization, bool) {
	authorizeContext, found := s.c.Get(organization_id_cache_prefix + id)
	if !found {
		return org_models.Organization{}, false
	}
	return authorizeContext.(org_models.Organization), true
}

func (s *cacheService) SetOrganization(organization org_models.Organization) {
	s.c.Set(organization_name_cache_prefix+organization.Name, organization, cache.DefaultExpiration)
	s.c.Set(organization_id_cache_prefix+organization.Id, organization, cache.DefaultExpiration)
}

func (s *cacheService) DeleteOrganizationByName(name string) {
	s.c.Delete(organization_name_cache_prefix + name)
}

func (s *cacheService) DeleteOrganizationById(id string) {
	s.c.Delete(organization_id_cache_prefix + id)
}
