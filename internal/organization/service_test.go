package organization

import (
	"github.com/shashimalcse/tiny-is/internal/cache"
)

func NewMockOrganizationService() OrganizationService {
	repo := NewMockOrganizationRepository()
	cache := cache.NewCacheService()
	return &organizationService{repo: repo, cacheService: cache}
}
