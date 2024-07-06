package models

type AuthorizationServer struct {
	Id     string  `db:"id" json:"id"`
	Name   string  `db:"name" json:"name,omitempty"`
	Scopes []Scope `json:"scopes,omitempty"`
}

type Scope struct {
	Id                    string `db:"id" json:"id"`
	AuthorizationServerId string `db:"authorization_server_id" json:"authorization_server_id,omitempty"`
	Name                  string `db:"name" json:"name,omitempty"`
	Description           string `db:"description" json:"description,omitempty"`
}

type Policy struct {
	Id                    string   `db:"id" json:"id"`
	AuthorizationServerId string   `db:"authorization_server_id" json:"authorization_server_id,omitempty"`
	Name                  string   `db:"name" json:"name,omitempty"`
	Applications          []string `json:"application,omitempty"`
}

type Rule struct {
	Id         string   `db:"id" json:"id"`
	PolicyId   string   `db:"policy_id" json:"policy_id,omitempty"`
	Name       string   `db:"name" json:"name,omitempty"`
	GrantTypes []string `json:"grant_types,omitempty"`
	Users      []string `json:"users,omitempty"`
	Scopes     []string `json:"scopes,omitempty"`
}
