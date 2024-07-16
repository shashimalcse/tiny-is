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

	// attribues
	mux.HandleFunc("POST /attributes", middlewares.ErrorMiddleware(handler.CreateAttribute))
	mux.HandleFunc("GET /attributes", middlewares.ErrorMiddleware(handler.GetAttributes))

	// user attributes
	mux.HandleFunc("POST /users/{id}/attributes", middlewares.ErrorMiddleware(handler.AddUserAttributes))
	mux.HandleFunc("PATCH /users/{id}/attributes", middlewares.ErrorMiddleware(handler.PatchUserAttributes))
}
