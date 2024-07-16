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
	CreateAttribute(ctx context.Context, name, orgId string) error
	GetAttributes(ctx context.Context, orgId string) ([]models.Attribute, error)
	PatchAttributes(ctx context.Context, orgId string, addedAttributes []models.Attribute, removedAttributes []models.Attribute) error
	AddUserAttributes(ctx context.Context, userId string, attributes []models.UserAttribute) error
	PatchUserAttributes(ctx context.Context, userId string, addedAttributes []models.UserAttribute, removedAttributes []models.UserAttribute) error
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
	user, err := s.repo.GetUserByID(ctx, id, orgId)
	if err != nil {
		return models.User{}, err
	}
	attributes, err := s.repo.GetUserAttributes(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	user.Attributes = attributes
	return user, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username, orgId string) (models.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username, orgId)
	if err != nil {
		return models.User{}, err
	}
	attributes, err := s.repo.GetUserAttributes(ctx, user.Id)
	if err != nil {
		return models.User{}, err
	}
	user.Attributes = attributes
	return user, nil
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

func (s *userService) CreateAttribute(ctx context.Context, name, orgId string) error {
	attributeId := uuid.New().String()
	return s.repo.CreateAttribute(ctx, attributeId, name, orgId)
}

func (s *userService) GetAttributes(ctx context.Context, orgId string) ([]models.Attribute, error) {
	return s.repo.GetAttributes(ctx, orgId)
}

func (s *userService) PatchAttributes(ctx context.Context, orgId string, addedAttributes []models.Attribute, removedAttributes []models.Attribute) error {
	return s.repo.PatchAttributes(ctx, orgId, addedAttributes, removedAttributes)
}

func (s *userService) AddUserAttributes(ctx context.Context, userId string, attributes []models.UserAttribute) error {
	return s.repo.AddUserAttributes(ctx, userId, attributes)
}

func (s *userService) PatchUserAttributes(ctx context.Context, userId string, addedAttributes []models.UserAttribute, removedAttributes []models.UserAttribute) error {
	return s.repo.PatchUserAttributes(ctx, userId, addedAttributes, removedAttributes)
}

func getPasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
