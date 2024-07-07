package models

type User struct {
	Id             string `db:"id" json:"id"`
	OrganizationId string `db:"organization_id" json:"organization_id"`
	Username       string `db:"username" json:"username"`
	Email          string `db:"email" json:"email"`
	PasswordHash   string `db:"password_hash"`
	Password       string `json:"password"`
}
