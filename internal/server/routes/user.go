package routes

import (
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/server/handlers"
	"github.com/shashimalcse/tiny-is/internal/user"
)

func RegisterUserRoutes(mux *http.ServeMux, userService *user.UserService) {
	handler := handlers.NewUserHandler(userService)

	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		handler.CreateUser(w, r)
	})
	mux.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		handler.GetUsers(w, r)
	})
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetUserByID(w, r)
	})
}
