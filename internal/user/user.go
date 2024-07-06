package user

import (
	"github.com/jmoiron/sqlx"
	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/user/models"
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

func (u UserService) GetUserById(userId string) (models.User, error) {
	var user models.User
	err := u.db.Get(&user, "SELECT id, username FROM org_user WHERE id=$1", userId)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u UserService) GetUsers() ([]models.User, error) {
	var users []models.User
	err := u.db.Select(&users, "SELECT id, username FROM org_user")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u UserService) CreateUser(user models.User) error {
	_, err := u.db.Exec("INSERT INTO org_user (id, username, password_hash) VALUES ($1, $2, $3)", user.Id, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u UserService) GetHashedPasswordByUsername(username string) (string, error) {
	var password string
	err := u.db.Get(&password, "SELECT password_hash FROM org_user WHERE username=$1", username)
	if err != nil {
		return "", err
	}
	return password, nil
}
