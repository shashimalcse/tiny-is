package models

import "github.com/shashimalcse/tiny-is/internal/user/models"

type UserResponse struct {
	Id             string `json:"id"`
	OrganizationId string `json:"organization_id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
}

type UserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type AttributeCreateRequest struct {
	Name string `json:"name"`
}

type AttributePatchRequest struct {
	AddedAttributes   []models.Attribute `json:"added_attributes"`
	RemovedAttributes []models.Attribute `json:"removed_attributes"`
}

type UserAttributeUpdateRequest struct {
	Attributes []models.UserAttribute `json:"attributes"`
}

type UserAttributePatchRequest struct {
	AddedAttributes   []models.UserAttribute `json:"added_attributes"`
	RemovedAttributes []models.UserAttribute `json:"removed_attributes"`
}

func GetUserResponse(user models.User) UserResponse {
	return UserResponse{
		Id:             user.Id,
		Username:       user.Username,
		OrganizationId: user.OrganizationId,
		Email:          user.Email,
	}
}

func GetUsersResponse(users []models.User) []UserResponse {
	if users == nil {
		return []UserResponse{}
	}
	var usersResponse []UserResponse
	for _, user := range users {
		usersResponse = append(usersResponse, GetUserResponse(user))
	}
	return usersResponse
}
