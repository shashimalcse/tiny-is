package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/shashimalcse/tiny-is/internal/server/models"
)

type CacheService struct {
	c *cache.Cache
}

func NewCacheService() *CacheService {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &CacheService{
		c: c,
	}
}

func (cacheService CacheService) AddOAuth2AuthorizeContextToCacheBySessionDataKey(sessionDataKey string, authorizeContext models.OAuth2AuthorizeContext) {
	cacheService.c.Set(sessionDataKey, authorizeContext, cache.NoExpiration)
}

func (cacheService CacheService) GetOAuth2AuthorizeContextFromCacheBySessionDataKey(sessionDataKey string) (models.OAuth2AuthorizeContext, bool) {
	authorizeContext, found := cacheService.c.Get(sessionDataKey)
	if !found {
		return models.OAuth2AuthorizeContext{}, false
	}
	return authorizeContext.(models.OAuth2AuthorizeContext), true
}

func (cacheService CacheService) AddCodeToCache(code string, authorizeContext models.OAuth2AuthorizeContext) {
	cacheService.c.Set(code, authorizeContext, cache.DefaultExpiration)
}

func (cacheService CacheService) GetCodeFromCache(code string) (models.OAuth2AuthorizeContext, bool) {
	authorizeContext, found := cacheService.c.Get(code)
	if !found || authorizeContext == nil {
		return models.OAuth2AuthorizeContext{}, false
	}
	return authorizeContext.(models.OAuth2AuthorizeContext), true
}
