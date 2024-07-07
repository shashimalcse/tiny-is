package models

type AuthenticatedUser struct {
	Id             string `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	OrganizationId string `json:"organization_id"`
}

type AuthenticateResult struct {
	Authenticated     bool
	AuthenticatedUser AuthenticatedUser
}
