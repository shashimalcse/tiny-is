package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/shashimalcse/tiny-is/internal/server/models"
	"github.com/shashimalcse/tiny-is/internal/user"
	user_models "github.com/shashimalcse/tiny-is/internal/user/models"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userService *user.UserService
}

func NewUserHandler(userService *user.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (handler UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")
	user, err := handler.userService.GetUserById(userId)
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
	users, err := handler.userService.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.GetUsersResponse(users))
}

func (handler UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var userCreateRequest models.UserCreateRequest
	err := json.NewDecoder(r.Body).Decode(&userCreateRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId := uuid.New().String()
	hashedPassword, err := hashPassword(userCreateRequest.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := user_models.User{
		Id:       userId,
		Username: userCreateRequest.Username,
		Password: hashedPassword,
	}
	err = handler.userService.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
