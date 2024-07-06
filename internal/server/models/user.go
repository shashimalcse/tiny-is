package models

import "github.com/shashimalcse/tiny-is/internal/user/models"

type UserResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type UserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
