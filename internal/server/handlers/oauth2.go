package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/shashimalcse/tiny-is/internal/oauth2"
	oauth2_models "github.com/shashimalcse/tiny-is/internal/oauth2/models"
	"github.com/shashimalcse/tiny-is/internal/server/middlewares"
	"github.com/shashimalcse/tiny-is/internal/server/models"
)

type OAuth2Handler struct {
	oauth2Service oauth2.OAuth2Service
}

func NewOAuth2Handler(oauth2Service oauth2.OAuth2Service) *OAuth2Handler {
	return &OAuth2Handler{
		oauth2Service: oauth2Service,
	}
}

func (handler OAuth2Handler) GetOAuth2AuthorizeRequest(w http.ResponseWriter, r *http.Request) (models.OAuth2AuthorizeRequest, error) {

	orgId := r.Header.Get("org_id")
	if orgId == "" {
		return models.OAuth2AuthorizeRequest{}, fmt.Errorf("Organization not found!")
	}
	orgName := r.Header.Get("org_name")
	if orgName == "" {
		return models.OAuth2AuthorizeRequest{}, fmt.Errorf("Organization not found!")
	}
	oauth2AuthorizeRequest := models.OAuth2AuthorizeRequest{
		ResponseType:     r.URL.Query().Get("response_type"),
		ClientId:         r.URL.Query().Get("client_id"),
		RedirectUri:      r.URL.Query().Get("redirect_uri"),
		Scope:            r.URL.Query().Get("scope"),
		State:            r.URL.Query().Get("state"),
		SessionDataKey:   r.URL.Query().Get("session_data_key"),
		OrganizationId:   orgId,
		OrganizationName: orgName,
	}

	return oauth2AuthorizeRequest, nil
}

func (handler OAuth2Handler) GetOAuth2TokenRequest(w http.ResponseWriter, r *http.Request) (models.OAuth2TokenRequest, error) {

	orgId := r.Header.Get("org_id")
	if orgId == "" {
		return models.OAuth2TokenRequest{}, fmt.Errorf("Organization not found!")
	}
	orgName := r.Header.Get("org_name")
	if orgName == "" {
		return models.OAuth2TokenRequest{}, fmt.Errorf("Organization not found!")
	}
	oauth2TokenRequest := models.OAuth2TokenRequest{
		GrantType:        r.Form.Get("grant_type"),
		Code:             r.Form.Get("code"),
		RefreshToken:     r.Form.Get("refresh_token"),
		ClientId:         r.Form.Get("client_id"),
		ClientSecret:     r.Form.Get("client_secret"),
		OrganizationId:   orgId,
		OrganizationName: orgName,
	}
	return oauth2TokenRequest, nil
}

func (handler OAuth2Handler) Authorize(w http.ResponseWriter, r *http.Request) error {

	oauth2AuthorizeRequest, err := handler.GetOAuth2AuthorizeRequest(w, r)
	if err != nil {
		return middlewares.NewAPIError(http.StatusBadRequest, err.Error())
	}
	ctx := r.Context()
	if oauth2AuthorizeRequest.IsInitialRequestFromClient() {
		if !oauth2AuthorizeRequest.IsValidRequest() {
			return middlewares.NewAPIError(http.StatusBadRequest, "Invalid request")
		}
		oauth2AuthorizeContext := oauth2_models.OAuth2AuthorizeContext{
			OAuth2AuthorizeRequest: oauth2AuthorizeRequest,
		}
		err := handler.oauth2Service.ValidateAuthroizeRequest(ctx, oauth2AuthorizeContext)
		if err != nil {
			return middlewares.NewAPIError(http.StatusBadRequest, err.Error())
		}
		sessionDataKey := uuid.New().String()
		oauth2AuthorizeContext.OAuth2AuthorizeRequest.SessionDataKey = sessionDataKey
		handler.oauth2Service.AddOAuth2AuthorizeContextToCacheBySessionDataKey(ctx, sessionDataKey, oauth2AuthorizeContext)
		u := &url.URL{
			Path:     fmt.Sprintf("/o/%s/login", oauth2AuthorizeRequest.OrganizationName),
			RawQuery: "session_data_key=" + url.QueryEscape(sessionDataKey),
		}
		http.Redirect(w, r, u.String(), http.StatusFound)
		return nil
	}

	sessionDataKey := r.URL.Query().Get("session_data_key")
	if sessionDataKey == "" {
		return middlewares.NewAPIError(http.StatusBadRequest, "session_data_key is required")
	}
	oauth2AuthorizeContext, err := handler.oauth2Service.GetOAuth2AuthorizeContextFromCacheBySessionDataKey(ctx, sessionDataKey)
	if err != nil {
		return middlewares.NewAPIError(http.StatusBadRequest, err.Error())
	}

	code := uuid.New().String()
	handler.oauth2Service.AddOAuth2AuthorizeContextToCacheByAuthCode(ctx, code, oauth2AuthorizeContext)
	state := oauth2AuthorizeContext.OAuth2AuthorizeRequest.State
	redirectURI := oauth2AuthorizeContext.OAuth2AuthorizeRequest.RedirectUri

	redirectURL, err := url.ParseRequestURI(redirectURI)
	if err != nil {
		return middlewares.NewAPIError(http.StatusBadRequest, "Invalid redirect uri")
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
	return nil
}

func (handler OAuth2Handler) Token(w http.ResponseWriter, r *http.Request) error {

	if err := r.ParseForm(); err != nil {
		return middlewares.NewAPIError(http.StatusBadRequest, "Invalid request payload")
	}

	oauth2TokenRequest, err := handler.GetOAuth2TokenRequest(w, r)
	if err != nil {
		return middlewares.NewAPIError(http.StatusBadRequest, err.Error())
	}
	ctx := r.Context()
	oauth2TokenContext := oauth2_models.OAuth2TokenContext{
		OAuth2TokenRequest: oauth2TokenRequest,
	}
	err = handler.oauth2Service.ValidateTokenRequest(ctx, oauth2TokenContext)
	if err != nil {
		return middlewares.NewAPIError(http.StatusBadRequest, err.Error())
	}

	grantHandler, err := handler.oauth2Service.GetGrantHandler(oauth2TokenRequest.GrantType)
	if err != nil {
		return middlewares.NewAPIError(http.StatusBadRequest, err.Error())
	}
	tokenResponse, err := grantHandler.HandleGrant(r.Context(), oauth2TokenContext)
	if err != nil {
		if err.Error() == "invalid_code" {
			return middlewares.NewAPIError(http.StatusBadRequest, "Invalid code")
		} else if err.Error() == "invalid_refresh_token" {
			return middlewares.NewAPIError(http.StatusBadRequest, "Invalid refresh token")
		}
		return middlewares.NewAPIError(http.StatusInternalServerError, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenResponse)
	return nil
}
