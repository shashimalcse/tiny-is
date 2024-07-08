package routes

import (
	"github.com/shashimalcse/tiny-is/internal/server/handlers"
	"github.com/shashimalcse/tiny-is/internal/server/middlewares"
	"github.com/shashimalcse/tiny-is/internal/server/utils"
	"github.com/shashimalcse/tiny-is/internal/user"
)

func RegisterUserRoutes(mux *utils.OrgServeMux, userService user.UserService) {
	handler := handlers.NewUserHandler(userService)

	mux.HandleFunc("GET /users", middlewares.ErrorMiddleware(handler.GetUsers))
	mux.HandleFunc("GET /users/{id}", middlewares.ErrorMiddleware(handler.GetUserByID))
	mux.HandleFunc("POST /users", middlewares.ErrorMiddleware(handler.CreateUser))
}
