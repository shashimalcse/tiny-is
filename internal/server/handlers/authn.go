package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/shashimalcse/tiny-is/internal/authn"
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
		return models.LoginRequest{}, fmt.Errorf("organization not found")
	}
	orgName := r.Header.Get("org_name")
	if orgName == "" {
		return models.LoginRequest{}, fmt.Errorf("organization not found")
	}
	loginRequest := models.LoginRequest{
		Username:         r.Form.Get("username"),
		Password:         r.Form.Get("password"),
		OrganizationId:   orgId,
		OrganizationName: orgName,
	}
	return loginRequest, nil
}

func (handler AuthnHandler) Login(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	sessionDataKey := r.Form.Get("session_data_key")
	if sessionDataKey == "" {
		http.Error(w, "session_data_key is required", http.StatusBadRequest)
		return
	}
	loginRequest, err := handler.GetLoginRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	authenticateResult, err := handler.authnService.AuthenticateUser(ctx, loginRequest.Username, loginRequest.Password, loginRequest.OrganizationId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if authenticateResult.Authenticated {
		oauth2AuthorizeContext, err := handler.authnService.GetOAuth2AuthorizeContextFromCacheBySessionDataKey(ctx, sessionDataKey)
		oauth2AuthorizeContext.AuthenticatedUser = authenticateResult.AuthenticatedUser
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
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
		w.Write([]byte("login failed"))
	}
}

func (handler AuthnHandler) LoginForm(w http.ResponseWriter, r *http.Request) {

	orgName := r.Header.Get("org_name")
	if orgName == "" {
		http.Error(w, "Organization not found!", http.StatusNotFound)
		return
	}
	sessionDataKey := r.URL.Query().Get("session_data_key")
	if sessionDataKey == "" {
		http.Error(w, "session_data_key is required", http.StatusBadRequest)
		return
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
		http.Error(w, "session_data_key is required", http.StatusBadRequest)
		return
	}
	handler.authnService.GetLoginPage(ctx, sessionDataKey, orgName).Render(r.Context(), w)
}
