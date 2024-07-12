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

type UserAttributeUpdateRequest struct {
	Attributes []models.UserAttribute `json:"attributes"`
}

func GetUserResponse(user models.User) UserResponse {
	return UserResponse{
		Id:       user.Id,
		Username: user.Username,
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
