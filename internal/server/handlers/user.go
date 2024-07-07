package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/server/models"
	"github.com/shashimalcse/tiny-is/internal/user"
	user_models "github.com/shashimalcse/tiny-is/internal/user/models"
)

type UserHandler struct {
	userService user.UserService
}

func NewUserHandler(userService user.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (handler UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orgId := r.Header.Get("org_id")
	if orgId == "" {
		http.Error(w, "Organization not found!", http.StatusNotFound)
		return
	}
	userId := r.URL.Query().Get("id")
	user, err := handler.userService.GetUserByID(ctx, userId, orgId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user.Id == "" {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.GetUserResponse(user))
	w.WriteHeader(http.StatusOK)
}

func (handler UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orgId := r.Header.Get("org_id")
	if orgId == "" {
		http.Error(w, "Organization not found!", http.StatusNotFound)
		return
	}
	users, err := handler.userService.GetUsers(ctx, orgId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.GetUsersResponse(users))
}

func (handler UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	orgId := r.Header.Get("org_id")
	if orgId == "" {
		http.Error(w, "Organization not found!", http.StatusNotFound)
		return
	}
	var userCreateRequest models.UserCreateRequest
	err := json.NewDecoder(r.Body).Decode(&userCreateRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := user_models.User{
		OrganizationId: orgId,
		Username:       userCreateRequest.Username,
		Password:       userCreateRequest.Password,
		Email:          userCreateRequest.Email,
	}
	err = handler.userService.CreateUser(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
