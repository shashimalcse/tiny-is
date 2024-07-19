package models

type Application struct {
	Id             string   `db:"id" json:"id"`
	Name           string   `db:"name" json:"name"`
	OrganizationId string   `db:"organization_id" json:"organization_id"`
	ClientId       string   `db:"client_id" json:"client_id,omitempty"`
	ClientSecret   string   `db:"client_secret" json:"client_secret,omitempty"`
	RedirectUris   []string `db:"redirect_uris" json:"redirect_uris,omitempty"`
	GrantTypes     []string `json:"grant_types,omitempty"`
}
