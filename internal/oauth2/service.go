package oauth2

import (
	"context"
	"errors"

	"github.com/shashimalcse/tiny-is/internal/application"
	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/oauth2/grant_handlers"
	"github.com/shashimalcse/tiny-is/internal/oauth2/models"
	"github.com/shashimalcse/tiny-is/internal/oauth2/token"
)

type OAuth2Service interface {
	ValidateAuthroizeRequest(ctx context.Context, authroizeContext models.OAuth2AuthorizeContext) error
	AddOAuth2AuthorizeContextToCacheBySessionDataKey(ctx context.Context, sessionDataKey string, authroizeContext models.OAuth2AuthorizeContext)
	GetOAuth2AuthorizeContextFromCacheBySessionDataKey(ctx context.Context, sessionDataKey string) (models.OAuth2AuthorizeContext, error)
	AddOAuth2AuthorizeContextToCacheByAuthCode(ctx context.Context, code string, authroizeContext models.OAuth2AuthorizeContext)
	GetOAuth2AuthorizeContextFromCacheByAuthCode(ctx context.Context, code string) (models.OAuth2AuthorizeContext, error)
	ValidateTokenRequest(ctx context.Context, tokenContext models.OAuth2TokenContext) error
	GetGrantHandler(grantType string) (grant_handlers.GrantHandler, error)
}

type oauth2Service struct {
	cacheService       *cache.CacheService
	tokenService       token.TokenService
	applicationService application.ApplicationService
	grantHandlers      map[string]grant_handlers.GrantHandler
}

func NewOAuth2Service(cacheService *cache.CacheService, applicationService application.ApplicationService) OAuth2Service {
	service := &oauth2Service{
		cacheService:       cacheService,
		applicationService: applicationService,
		grantHandlers:      make(map[string]grant_handlers.GrantHandler),
	}
	service.registerGrantHandlers()
	service.tokenService = token.NewTokenService([]byte("secret"))
	return service
}

func (s *oauth2Service) registerGrantHandlers() {
	s.grantHandlers["authorization_code"] = grant_handlers.NewAuthorizationCodeGrantHandler(s.cacheService, s.tokenService)
	s.grantHandlers["refresh_token"] = grant_handlers.NewRefreshTokenGrantHandler(s.cacheService, s.tokenService)
}

func (s *oauth2Service) GetGrantHandler(grantType string) (grant_handlers.GrantHandler, error) {
	grantHandler := s.grantHandlers[grantType]
	if grantHandler == nil {
		return nil, errors.New("unsupported grant type")
	}
	return grantHandler, nil
}

func (s *oauth2Service) ValidateAuthroizeRequest(ctx context.Context, authroizeContext models.OAuth2AuthorizeContext) error {
	validClientId, err := s.applicationService.ValidateClientId(ctx, authroizeContext.OAuth2AuthorizeRequest.ClientId, authroizeContext.OAuth2AuthorizeRequest.OrganizationId)
	if err != nil {
		return err
	}
	if !validClientId {
		return errors.New("invalid client id")
	}
	validRedirectUri, err := s.applicationService.ValidateRedirectUri(ctx, authroizeContext.OAuth2AuthorizeRequest.ClientId, authroizeContext.OAuth2AuthorizeRequest.RedirectUri, authroizeContext.OAuth2AuthorizeRequest.OrganizationId)
	if err != nil {
		return err
	}
	if !validRedirectUri {
		return errors.New("invalid redirect uri")
	}
	return nil
}

func (s *oauth2Service) ValidateTokenRequest(ctx context.Context, tokenContext models.OAuth2TokenContext) error {
	validClientId, err := s.applicationService.ValidateClientId(ctx, tokenContext.OAuth2TokenRequest.ClientId, tokenContext.OAuth2TokenRequest.OrganizationId)
	if err != nil {
		return err
	}
	if !validClientId {
		return errors.New("invalid client id")
	}
	ValidClientSecret, err := s.applicationService.ValidateClientSecret(ctx, tokenContext.OAuth2TokenRequest.ClientId, tokenContext.OAuth2TokenRequest.ClientSecret, tokenContext.OAuth2TokenRequest.OrganizationId)
	if err != nil {
		return err
	}
	if !ValidClientSecret {
		return errors.New("invalid client secret")
	}
	return nil
}

func (s *oauth2Service) AddOAuth2AuthorizeContextToCacheBySessionDataKey(ctx context.Context, sessionDataKey string, authroizeContext models.OAuth2AuthorizeContext) {
	s.cacheService.AddOAuth2AuthorizeContextToCacheBySessionDataKey(sessionDataKey, authroizeContext)
}

func (s *oauth2Service) GetOAuth2AuthorizeContextFromCacheBySessionDataKey(ctx context.Context, sessionDataKey string) (models.OAuth2AuthorizeContext, error) {
	oauth2AuthorizeContext, found := s.cacheService.GetOAuth2AuthorizeContextFromCacheBySessionDataKey(sessionDataKey)
	if !found {
		return models.OAuth2AuthorizeContext{}, errors.New("invalid session_data_key")
	}
	return oauth2AuthorizeContext, nil
}

func (s *oauth2Service) AddOAuth2AuthorizeContextToCacheByAuthCode(ctx context.Context, code string, authroizeContext models.OAuth2AuthorizeContext) {
	s.cacheService.AddOAuth2AuthorizeContextToCacheByAuthCode(code, authroizeContext)
}

func (s *oauth2Service) GetOAuth2AuthorizeContextFromCacheByAuthCode(ctx context.Context, code string) (models.OAuth2AuthorizeContext, error) {
	oauth2AuthorizeContext, found := s.cacheService.GetOAuth2AuthorizeContextFromCacheByAuthCode(code)
	if !found {
		return models.OAuth2AuthorizeContext{}, errors.New("invalid code")
	}
	return oauth2AuthorizeContext, nil
}
