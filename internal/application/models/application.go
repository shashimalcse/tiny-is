package models

import "strings"

type Application struct {
	Id           string `db:"id" json:"id"`
	ClientId     string `db:"client_id" json:"client_id,omitempty"`
	ClientSecret string `db:"client_secret" json:"client_secret,omitempty"`
	RedirectUri  string `db:"redirect_uri" json:"redirect_uri,omitempty"`
	GrantTypes   string `db:"grant_types" json:"grant_types,omitempty"`
}

func (app Application) GetGrantTypes() []string {
	return strings.Split(app.GrantTypes, " ")
}
