package models

type Attribute struct {
	ID             string `db:"id"`
	OrganizationID string `db:"organization_id"`
	Name           string `db:"name"`
}

type UserAttribute struct {
	ID    string `db:"id"`
	Name  string `db:"name"`
	Value string `db:"value"`
}
