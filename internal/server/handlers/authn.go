package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/shashimalcse/tiny-is/internal/authn"
	"github.com/shashimalcse/tiny-is/internal/server/middlewares"
	"github.com/shashimalcse/tiny-is/internal/server/models"
)

type AuthnHandler struct {
	authnService authn.AuthnService
}

func NewAuthnHandler(authnService authn.AuthnService) *AuthnHandler {
	return &AuthnHandler{
		authnService: authnService,
	}
}

func (handler AuthnHandler) GetLoginRequest(w http.ResponseWriter, r *http.Request) (models.LoginRequest, error) {

	orgId := r.Header.Get("org_id")
	if orgId == "" {
		return models.LoginRequest{}, fmt.Errorf("Organization not found!")
	}
	orgName := r.Header.Get("org_name")
	if orgName == "" {
		return models.LoginRequest{}, fmt.Errorf("Organization not found!")
	}
	loginRequest := models.LoginRequest{
		Username:         r.Form.Get("username"),
		Password:         r.Form.Get("password"),
		OrganizationId:   orgId,
		OrganizationName: orgName,
	}
	return loginRequest, nil
}

func (handler AuthnHandler) Login(w http.ResponseWriter, r *http.Request) error {

	err := r.ParseForm()
	if err != nil {
		return middlewares.NewAPIError(http.StatusBadRequest, "invalid request")
	}

	sessionDataKey := r.Form.Get("session_data_key")
	if sessionDataKey == "" {
		return middlewares.NewAPIError(http.StatusBadRequest, "session_data_key is required")
	}
	loginRequest, err := handler.GetLoginRequest(w, r)
	if err != nil {
		return middlewares.NewAPIError(http.StatusBadRequest, err.Error())
	}
	ctx := r.Context()
	authenticateResult, err := handler.authnService.AuthenticateUser(ctx, loginRequest.Username, loginRequest.Password, loginRequest.OrganizationId)
	if err != nil {
		return middlewares.NewAPIError(http.StatusInternalServerError, err.Error())
	}
	if authenticateResult.Authenticated {
		oauth2AuthorizeContext, err := handler.authnService.GetOAuth2AuthorizeContextFromCacheBySessionDataKey(ctx, sessionDataKey)
		oauth2AuthorizeContext.AuthenticatedUser = authenticateResult.AuthenticatedUser
		if err != nil {
			return middlewares.NewAPIError(http.StatusBadRequest, err.Error())
		}
		sessionDuration := 30 * time.Minute
		sessionID := handler.authnService.CreateSession(ctx, oauth2AuthorizeContext, sessionDuration)
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Expires:  time.Now().Add(sessionDuration),
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, cookie)
		handler.authnService.AddOAuth2AuthorizeContextToCacheBySessionDataKey(ctx, sessionDataKey, oauth2AuthorizeContext)
		u := &url.URL{
			Path:     fmt.Sprintf("/o/%s/authorize", loginRequest.OrganizationName),
			RawQuery: "session_data_key=" + url.QueryEscape(sessionDataKey),
		}
		http.Redirect(w, r, u.String(), http.StatusFound)
	} else {
		return middlewares.NewAPIError(http.StatusUnauthorized, "Invalid credentials")
	}
	return nil
}

func (handler AuthnHandler) GetLoginForm(w http.ResponseWriter, r *http.Request) error {

	orgName := r.Header.Get("org_name")
	if orgName == "" {
		return middlewares.NewAPIError(http.StatusNotFound, "Organization not found!")
	}
	sessionDataKey := r.URL.Query().Get("session_data_key")
	if sessionDataKey == "" {
		return middlewares.NewAPIError(http.StatusBadRequest, "session_data_key is required")
	}
	ctx := r.Context()
	cookie, err := r.Cookie("session_id")
	if err == nil {
		if _, found := handler.authnService.GetSession(ctx, cookie.Value); found {
			u := &url.URL{
				Path:     fmt.Sprintf("/o/%s/authorize", orgName),
				RawQuery: "session_data_key=" + url.QueryEscape(sessionDataKey),
			}
			http.Redirect(w, r, u.String(), http.StatusFound)
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		})
	}
	if sessionDataKey == "" {
		return middlewares.NewAPIError(http.StatusBadRequest, "session_data_key is required")
	}
	handler.authnService.GetLoginPage(ctx, sessionDataKey, orgName).Render(r.Context(), w)
	return nil
}
