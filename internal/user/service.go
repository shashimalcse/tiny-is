package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/user/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUsers(ctx context.Context, orgId string) ([]models.User, error)
	GetUserByID(ctx context.Context, id, orgId string) (models.User, error)
	GetUserByUsername(ctx context.Context, username, orgId string) (models.User, error)
	CreateUser(ctx context.Context, User models.User) error
	AuthenticateUser(ctx context.Context, username, password, orgId string) (bool, error)
}

type userService struct {
	cacheService cache.CacheService
	repo         UserRepository
}

func NewUserService(cacheService cache.CacheService, repo UserRepository) UserService {
	return &userService{
		cacheService: cacheService,
		repo:         repo,
	}
}

func (s *userService) GetUsers(ctx context.Context, orgId string) ([]models.User, error) {
	return s.repo.GetUsers(ctx, orgId)
}

func (s *userService) GetUserByID(ctx context.Context, id, orgId string) (models.User, error) {
	return s.repo.GetUserByID(ctx, id, orgId)
}

func (s *userService) GetUserByUsername(ctx context.Context, username, orgId string) (models.User, error) {
	return s.repo.GetUserByUsername(ctx, username, orgId)
}

func (s *userService) CreateUser(ctx context.Context, user models.User) error {
	userId := uuid.New().String()
	passwordHash, err := getPasswordHash(user.Password)
	if err != nil {
		return err
	}
	user.Id = userId
	user.PasswordHash = passwordHash
	return s.repo.CreateUser(ctx, user)
}

func (s *userService) AuthenticateUser(ctx context.Context, username, password, orgId string) (bool, error) {
	hashedPassword, err := s.repo.GetHashedPasswordByUsername(ctx, username, orgId)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}

func getPasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
