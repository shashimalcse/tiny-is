package handlers

import (
	"net/http"
	"net/url"

	"github.com/shashimalcse/tiny-is/internal/authn"
	"github.com/shashimalcse/tiny-is/internal/server/models"
)

type AuthnHandler struct {
	authn *authn.Authn
}

func NewAuthnHandler(authn *authn.Authn) *AuthnHandler {
	return &AuthnHandler{
		authn: authn,
	}
}

func (handler AuthnHandler) Login(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	sessionDataKey := r.Form.Get("session_data_key")
	if sessionDataKey == "" {
		http.Error(w, "session_data_key is required", http.StatusBadRequest)
		return
	}
	isUserExists, err := handler.authn.ValidateUser(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if isUserExists {
		oauth2AuthorizeContext, found := handler.authn.CacheService.GetOAuth2AuthorizeContextFromCacheBySessionDataKey(sessionDataKey)
		if !found {
			http.Error(w, "invalid session_data_key", http.StatusBadRequest)
			return
		}
		userId, err := handler.authn.GetUserIdByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		authroizedUser := models.AuthenticatedUser{
			Id:       userId,
			Username: username,
		}
		oauth2AuthorizeContext.AuthenticatedUser = authroizedUser
		handler.authn.CacheService.AddOAuth2AuthorizeContextToCacheBySessionDataKey(sessionDataKey, oauth2AuthorizeContext)
		redirectURL, err := url.Parse("/authorize")
		if err != nil {
			http.Error(w, "invalid redirect_uri", http.StatusBadRequest)
			return
		}
		query := redirectURL.Query()
		query.Set("session_data_key", sessionDataKey)
		redirectURL.RawQuery = query.Encode()
		http.Redirect(w, r, redirectURL.String(), http.StatusFound)
	} else {
		w.Write([]byte("login failed"))
	}
}

func (handler AuthnHandler) LoginForm(w http.ResponseWriter, r *http.Request) {

	sessionDataKey := r.URL.Query().Get("session_data_key")
	if sessionDataKey == "" {
		http.Error(w, "session_data_key is required", http.StatusBadRequest)
		return
	}
	handler.authn.GetLoginPage(sessionDataKey).Render(r.Context(), w)
}
