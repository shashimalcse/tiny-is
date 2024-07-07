package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/shashimalcse/tiny-is/internal/user/models"
)

type UserRepository interface {
	GetUsers(ctx context.Context, orgId string) ([]models.User, error)
	GetUserByID(ctx context.Context, id, orgId string) (models.User, error)
	GetUserByUsername(ctx context.Context, username, orgId string) (models.User, error)
	CreateUser(ctx context.Context, User models.User) error
	GetHashedPasswordByUsername(ctx context.Context, username, orgId string) (string, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUsers(ctx context.Context, orgId string) ([]models.User, error) {
	var Users []models.User
	err := r.db.Select(&Users, "SELECT id, organization_id, username, email FROM org_user WHERE organization_id=$1", orgId)
	if err != nil {
		return nil, err
	}
	return Users, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id, orgId string) (models.User, error) {
	var User models.User
	err := r.db.Get(&User, "SELECT id, organization_id, username, email FROM org_user WHERE id=$1 AND organization_id=$2", id, orgId)
	if err != nil {
		return models.User{}, err
	}
	return User, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username, orgId string) (models.User, error) {
	var User models.User
	err := r.db.Get(&User, "SELECT id, organization_id, username, email FROM org_user WHERE username=$1 AND organization_id=$2", username, orgId)
	if err != nil {
		return models.User{}, err
	}
	return User, nil
}

func (r *userRepository) CreateUser(ctx context.Context, User models.User) error {
	_, err := r.db.Exec("INSERT INTO org_user (id, organization_id, username, email, password_hash) VALUES ($1, $2, $3, $4, $5)", User.Id, User.OrganizationId, User.Username, User.Email, User.PasswordHash)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetHashedPasswordByUsername(ctx context.Context, username, orgId string) (string, error) {
	var password string
	err := r.db.Get(&password, "SELECT password_hash FROM org_user WHERE username=$1 AND organization_id=$2", username, orgId)
	if err != nil {
		return "", err
	}
	return password, nil
}
