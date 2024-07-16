package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/server/middlewares"
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

func (handler UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	orgId := r.Header.Get("org_id")
	if orgId == "" {
		return middlewares.NewAPIError(http.StatusNotFound, "Organization not found!")
	}
	userId := r.URL.Query().Get("id")
	user, err := handler.userService.GetUserByID(ctx, userId, orgId)
	if err != nil {
		return middlewares.NewAPIError(http.StatusInternalServerError, err.Error())
	}
	if user.Id == "" {
		return middlewares.NewAPIError(http.StatusNotFound, "User not found!")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.GetUserResponse(user))
	w.WriteHeader(http.StatusOK)
	return nil
}

func (handler UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	orgId := r.Header.Get("org_id")
	if orgId == "" {
		return middlewares.NewAPIError(http.StatusNotFound, "Organization not found!")
	}
	users, err := handler.userService.GetUsers(ctx, orgId)
	if err != nil {
		return middlewares.NewAPIError(http.StatusInternalServerError, err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.GetUsersResponse(users))
	return nil
}

func (handler UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	orgId := r.Header.Get("org_id")
	if orgId == "" {
		return middlewares.NewAPIError(http.StatusNotFound, "Organization not found!")
	}
	var userCreateRequest models.UserCreateRequest
	err := json.NewDecoder(r.Body).Decode(&userCreateRequest)
	if err != nil {
		return middlewares.NewAPIError(http.StatusBadRequest, "Invalid request payload")
	}
	user := user_models.User{
		OrganizationId: orgId,
		Username:       userCreateRequest.Username,
		Password:       userCreateRequest.Password,
		Email:          userCreateRequest.Email,
	}
	err = handler.userService.CreateUser(ctx, user)
	if err != nil {
		return middlewares.NewAPIError(http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (handler UserHandler) CreateAttribute(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	orgId := r.Header.Get("org_id")
	if orgId == "" {
		return middlewares.NewAPIError(http.StatusNotFound, "Organization not found!")
	}
	var attribute models.AttributeCreateRequest
	err := json.NewDecoder(r.Body).Decode(&attribute)
	if err != nil {
		return middlewares.NewAPIError(http.StatusBadRequest, "Invalid request payload")
	}
	err = handler.userService.CreateAttribute(ctx, attribute.Name, orgId)
	if err != nil {
		return middlewares.NewAPIError(http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (handler UserHandler) GetAttributes(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	orgId := r.Header.Get("org_id")
	if orgId == "" {
		return middlewares.NewAPIError(http.StatusNotFound, "Organization not found!")
	}
	attributes, err := handler.userService.GetAttributes(ctx, orgId)
	if err != nil {
		return middlewares.NewAPIError(http.StatusInternalServerError, err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(attributes)
	return nil
}

func (handler UserHandler) PatchAttributes(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	orgId := r.Header.Get("org_id")
	if orgId == "" {
		return middlewares.NewAPIError(http.StatusNotFound, "Organization not found!")
	}
	var attribute models.AttributePatchRequest
	err := json.NewDecoder(r.Body).Decode(&attribute)
	if err != nil {
		return middlewares.NewAPIError(http.StatusBadRequest, "Invalid request payload")
	}
	err = handler.userService.PatchAttributes(ctx, orgId, attribute.AddedAttributes, attribute.RemovedAttributes)
	if err != nil {
		return middlewares.NewAPIError(http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (handler UserHandler) AddUserAttributes(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	orgId := r.Header.Get("org_id")
	if orgId == "" {
		return middlewares.NewAPIError(http.StatusNotFound, "Organization not found!")
	}
	userId := r.URL.Query().Get("id")
	var attributes models.UserAttributeUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&attributes)
	if err != nil {
		return middlewares.NewAPIError(http.StatusBadRequest, "Invalid request payload")
	}
	err = handler.userService.AddUserAttributes(ctx, userId, attributes.Attributes)
	if err != nil {
		return middlewares.NewAPIError(http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (handler UserHandler) PatchUserAttributes(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	orgId := r.Header.Get("org_id")
	if orgId == "" {
		return middlewares.NewAPIError(http.StatusNotFound, "Organization not found!")
	}
	userId := r.URL.Query().Get("id")
	var attributes models.UserAttributePatchRequest
	err := json.NewDecoder(r.Body).Decode(&attributes)
	if err != nil {
		return middlewares.NewAPIError(http.StatusBadRequest, "Invalid request payload")
	}
	err = handler.userService.PatchUserAttributes(ctx, userId, attributes.AddedAttributes, attributes.RemovedAttributes)
	if err != nil {
		return middlewares.NewAPIError(http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
