package authz

import (
	"github.com/jmoiron/sqlx"
	"github.com/shashimalcse/tiny-is/internal/cache"
)

type AuthorizationServerService struct {
	cacheService *cache.CacheService
	db           *sqlx.DB
}

func NewAuthorizationServerService(cacheService *cache.CacheService, db *sqlx.DB) *AuthorizationServerService {
	return &AuthorizationServerService{
		cacheService: cacheService,
		db:           db,
	}
}
