package models

type OAuth2AuthorizeRequest struct {
	ResponseType   string
	ClientId       string
	RedirectUri    string
	Scope          string
	State          string
	SessionDataKey string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type OAuth2AuthorizeContext struct {
	OAuth2AuthorizeRequest OAuth2AuthorizeRequest `json:"oauth2_authorize_request"`
	AuthenticatedUser      AuthenticatedUser      `json:"authenticated_user"`
}

type OAuth2TokenRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (or OAuth2AuthorizeRequest) IsInitialRequestFromClient() bool {
	return or.SessionDataKey == ""
}

func (or OAuth2AuthorizeRequest) IsValidRequest() bool {
	if or.ResponseType == "" || or.ClientId == "" || or.RedirectUri == "" {
		return false
	}
	return true
}
