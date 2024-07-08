package authn

import (
	"context"
	"errors"
	"time"

	"github.com/a-h/templ"
	"github.com/shashimalcse/tiny-is/internal/authn/models"
	"github.com/shashimalcse/tiny-is/internal/authn/screens"
	"github.com/shashimalcse/tiny-is/internal/cache"
	oauth2_models "github.com/shashimalcse/tiny-is/internal/oauth2/models"
	"github.com/shashimalcse/tiny-is/internal/session"
	"github.com/shashimalcse/tiny-is/internal/user"
)

type AuthnService interface {
	GetLoginPage(ctx context.Context, sessionDataKey, organizationName string) templ.Component
	AuthenticateUser(ctx context.Context, username, password, orgId string) (models.AuthenticateResult, error)
	GetOAuth2AuthorizeContextFromCacheBySessionDataKey(ctx context.Context, sessionDataKey string) (oauth2_models.OAuth2AuthorizeContext, error)
	AddOAuth2AuthorizeContextToCacheBySessionDataKey(ctx context.Context, sessionDataKey string, authroizeContext oauth2_models.OAuth2AuthorizeContext)
	CreateSession(ctx context.Context, oauth2AuthorizeContext oauth2_models.OAuth2AuthorizeContext, sessionDuration time.Duration) string
	GetSession(ctx context.Context, sessionID string) (session.SessionInfo, bool)
}

type authnService struct {
	cacheService cache.CacheService
	SessionStore session.SessionStore
	userService  user.UserService
}

func NewAuthnService(cacheService cache.CacheService, sessionStore session.SessionStore, userService user.UserService) AuthnService {
	service := &authnService{
		cacheService: cacheService,
		SessionStore: sessionStore,
		userService:  userService,
	}
	return service
}

func (s *authnService) GetLoginPage(ctx context.Context, sessionDataKey, organizationName string) templ.Component {
	return screens.LoginPage(sessionDataKey, organizationName)
}

func (s *authnService) AuthenticateUser(ctx context.Context, username, password, orgId string) (models.AuthenticateResult, error) {
	authenticated, err := s.userService.AuthenticateUser(ctx, username, password, orgId)
	if err != nil {
		return models.AuthenticateResult{}, err
	}
	user, err := s.userService.GetUserByUsername(ctx, username, orgId)
	if err != nil {
		return models.AuthenticateResult{}, err
	}
	authenticatedUser := models.AuthenticatedUser{
		Id:             user.Id,
		Username:       user.Username,
		Email:          user.Email,
		OrganizationId: user.OrganizationId,
	}
	return models.AuthenticateResult{Authenticated: authenticated, AuthenticatedUser: authenticatedUser}, nil
}

func (s *authnService) GetOAuth2AuthorizeContextFromCacheBySessionDataKey(ctx context.Context, sessionDataKey string) (oauth2_models.OAuth2AuthorizeContext, error) {
	oauth2AuthorizeContext, found := s.cacheService.GetOAuth2AuthorizeContextFromCacheBySessionDataKey(sessionDataKey)
	if !found {
		return oauth2_models.OAuth2AuthorizeContext{}, errors.New("invalid session_data_key")
	}
	return oauth2AuthorizeContext, nil
}

func (s *authnService) AddOAuth2AuthorizeContextToCacheBySessionDataKey(ctx context.Context, sessionDataKey string, authroizeContext oauth2_models.OAuth2AuthorizeContext) {
	s.cacheService.AddOAuth2AuthorizeContextToCacheBySessionDataKey(sessionDataKey, authroizeContext)
}

func (s *authnService) CreateSession(ctx context.Context, oauth2AuthorizeContext oauth2_models.OAuth2AuthorizeContext, sessionDuration time.Duration) string {
	sessionID := s.SessionStore.CreateSession(oauth2AuthorizeContext.AuthenticatedUser.Id, oauth2AuthorizeContext.OAuth2AuthorizeRequest.OrganizationId, oauth2AuthorizeContext.OAuth2AuthorizeRequest.ClientId, sessionDuration)
	return sessionID
}

func (s *authnService) GetSession(ctx context.Context, sessionID string) (session.SessionInfo, bool) {
	return s.SessionStore.GetSession(sessionID)
}
