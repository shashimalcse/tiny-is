package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/shashimalcse/tiny-is/internal/oauth2"
	"github.com/shashimalcse/tiny-is/internal/server/models"
)

type OAuth2Handler struct {
	oauth2 *oauth2.OAuth2
}

func NewOAuth2Handler(oauth2 *oauth2.OAuth2) *OAuth2Handler {
	return &OAuth2Handler{
		oauth2: oauth2,
	}
}

func (handler OAuth2Handler) GetOAuth2AuthorizeRequest(w http.ResponseWriter, r *http.Request) models.OAuth2AuthorizeRequest {

	oauth2AuthorizeRequest := models.OAuth2AuthorizeRequest{
		ResponseType:   r.URL.Query().Get("response_type"),
		ClientId:       r.URL.Query().Get("client_id"),
		RedirectUri:    r.URL.Query().Get("redirect_uri"),
		Scope:          r.URL.Query().Get("scope"),
		State:          r.URL.Query().Get("state"),
		SessionDataKey: r.URL.Query().Get("session_data_key"),
	}

	return oauth2AuthorizeRequest
}

func (handler OAuth2Handler) GetOAuth2TokenRequest(w http.ResponseWriter, r *http.Request) models.OAuth2TokenRequest {

	oauth2TokenRequest := models.OAuth2TokenRequest{
		GrantType:    r.Form.Get("grant_type"),
		Code:         r.Form.Get("code"),
		RedirectUri:  r.Form.Get("redirect_uri"),
		ClientId:     r.Form.Get("client_id"),
		ClientSecret: r.Form.Get("client_secret"),
	}
	return oauth2TokenRequest
}

func (handler OAuth2Handler) Authorize(w http.ResponseWriter, r *http.Request) {

	oauth2AuthorizeRequest := handler.GetOAuth2AuthorizeRequest(w, r)

	if oauth2AuthorizeRequest.IsInitialRequestFromClient() {
		if !oauth2AuthorizeRequest.IsValidRequest() {
			http.Error(w, "invalid request!", http.StatusBadRequest)
			return
		}
		err := handler.validateClientForAuthorize(oauth2AuthorizeRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sessionDataKey := uuid.New().String()
		oauth2AuthorizeContext := models.OAuth2AuthorizeContext{
			OAuth2AuthorizeRequest: oauth2AuthorizeRequest,
		}
		handler.oauth2.CacheService.AddOAuth2AuthorizeContextToCacheBySessionDataKey(sessionDataKey, oauth2AuthorizeContext)
		http.Redirect(w, r, "/login?session_data_key="+sessionDataKey, http.StatusFound)
		return
	}

	sessionDataKey := r.URL.Query().Get("session_data_key")
	if sessionDataKey == "" {
		http.Error(w, "session_data_key is required", http.StatusBadRequest)
		return
	}
	oauth2AuthorizeContext, found := handler.oauth2.CacheService.GetOAuth2AuthorizeContextFromCacheBySessionDataKey(sessionDataKey)
	if !found {
		http.Error(w, "invalid session_data_key", http.StatusBadRequest)
		return
	}

	code := uuid.New().String()
	handler.oauth2.CacheService.AddCodeToCache(code, oauth2AuthorizeContext)
	state := oauth2AuthorizeContext.OAuth2AuthorizeRequest.State
	redirectURI := oauth2AuthorizeContext.OAuth2AuthorizeRequest.RedirectUri

	redirectURL, err := url.ParseRequestURI(redirectURI)
	if err != nil {
		http.Error(w, "invalid redirect_uri", http.StatusBadRequest)
		return
	}

	query := redirectURL.Query()
	query.Set("code", code)
	if state != "" {
		query.Set("state", state)
	}
	redirectURL.RawQuery = query.Encode()
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
        <!DOCTYPE html>
        <html>
        <head>
        </head>
        <body>
            <script>
                window.location.replace("%s");
            </script>
        </body>
        </html>
    `, redirectURL.String())
}

func (handler OAuth2Handler) Token(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	oauth2TokenRequest := handler.GetOAuth2TokenRequest(w, r)

	err := handler.validateClientForToken(oauth2TokenRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	grantHandler := handler.oauth2.GrantTypes[oauth2TokenRequest.GrantType]
	if grantHandler == nil {
		http.Error(w, "unsupported grant type", http.StatusBadRequest)
		return

	}
	tokenResponse, err := grantHandler(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenResponse)
}

func (handler OAuth2Handler) validateClientForAuthorize(oauth2AuthorizeRequest models.OAuth2AuthorizeRequest) error {
	isClientIdExist, err := handler.isValidateClientId(oauth2AuthorizeRequest.ClientId)
	if err != nil {
		return fmt.Errorf("internal server error: %w", err)
	}
	if !isClientIdExist {
		return fmt.Errorf("invalid client_id")
	}
	isRedirectUriCorrect, err := handler.isValidateRedirectURI(oauth2AuthorizeRequest.ClientId, oauth2AuthorizeRequest.RedirectUri)
	if err != nil {
		return fmt.Errorf("internal server error: %w", err)
	}
	if !isRedirectUriCorrect {
		return fmt.Errorf("invalid redirect_uri")
	}
	return nil
}

func (handler OAuth2Handler) validateClientForToken(oauth2TokenRequest models.OAuth2TokenRequest) error {
	isClientIdExist, err := handler.isValidateClientId(oauth2TokenRequest.ClientId)
	if err != nil {
		return fmt.Errorf("internal server error: %w", err)
	}
	if !isClientIdExist {
		return fmt.Errorf("invalid client_id")
	}
	isClientSecretCorrect, err := handler.isValidateClientSecret(oauth2TokenRequest.ClientId, oauth2TokenRequest.ClientSecret)
	if err != nil {
		return fmt.Errorf("internal server error: %w", err)
	}
	if !isClientSecretCorrect {
		return fmt.Errorf("invalid client_secret")
	}
	return nil
}

func (handler OAuth2Handler) isValidateClientId(clientId string) (bool, error) {
	return handler.oauth2.ApplicationService.ValidateClientId(clientId)
}

func (handler OAuth2Handler) isValidateClientSecret(clientId, clientSecret string) (bool, error) {
	return handler.oauth2.ApplicationService.ValidateClientSecret(clientId, clientSecret)
}

func (handler OAuth2Handler) isValidateRedirectURI(clientId, redirectURI string) (bool, error) {
	return handler.oauth2.ApplicationService.ValidateRedirectUri(clientId, redirectURI)
}

