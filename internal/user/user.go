package user

import (
	"github.com/jmoiron/sqlx"
	"github.com/shashimalcse/tiny-is/internal/cache"
)

type UserService struct {
	cacheService *cache.CacheService
	db           *sqlx.DB
}

func NewUserService(cacheService *cache.CacheService, db *sqlx.DB) *UserService {
	return &UserService{
		cacheService: cacheService,
		db:           db,
	}
}

func (u UserService) GetUserIdByUsername(username string) (string, error) {
	var userId string
	err := u.db.Get(&userId, "SELECT id FROM org_user WHERE username=$1", username)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (u UserService) CreateUser(username, password string) error {
	_, err := u.db.Exec("INSERT INTO org_user (username, password) VALUES ($1, $2)", username, password)
	if err != nil {
		return err
	}
	return nil
}

func (u UserService) GetHashedPasswordByUsername(username string) (string, error) {
	var password string
	err := u.db.Get(&password, "SELECT password FROM org_user WHERE username=$1", username)
	if err != nil {
		return "", err
	}
	return password, nil
}
