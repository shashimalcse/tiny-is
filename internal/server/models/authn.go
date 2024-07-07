package models

type LoginRequest struct {
	OrganizationId   string `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	Username         string `json:"username"`
	Password         string `json:"password"`
}
