package routes

import (
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/config"
	"github.com/shashimalcse/tiny-is/internal/security"
	"github.com/shashimalcse/tiny-is/internal/server/handlers"
	tinyhttp "github.com/shashimalcse/tiny-is/internal/server/http"
	"github.com/shashimalcse/tiny-is/internal/server/middlewares"
	"github.com/shashimalcse/tiny-is/internal/user"
)

func RegisterUserRoutes(mux *tinyhttp.TinyServeMux, cfg *config.Config, keyManager *security.KeyManager, userService user.UserService) {
	handler := handlers.NewUserHandler(userService)
	getUsersHandler := middlewares.ChainMiddleware(handler.GetUsers, middlewares.ErrorMiddleware(), middlewares.JWTMiddleware(cfg, keyManager))
	getUserByIDHandler := middlewares.ChainMiddleware(handler.GetUserByID, middlewares.ErrorMiddleware(), middlewares.JWTMiddleware(cfg, keyManager))
	createUserHandler := middlewares.ChainMiddleware(handler.CreateUser, middlewares.ErrorMiddleware(), middlewares.JWTMiddleware(cfg, keyManager))
	addUserAttributesHandler := middlewares.ChainMiddleware(handler.AddUserAttributes, middlewares.ErrorMiddleware(), middlewares.JWTMiddleware(cfg, keyManager))
	patchUserAttributesHandler := middlewares.ChainMiddleware(handler.PatchUserAttributes, middlewares.ErrorMiddleware(), middlewares.JWTMiddleware(cfg, keyManager))
	createAttributeHandler := middlewares.ChainMiddleware(handler.CreateAttribute, middlewares.ErrorMiddleware(), middlewares.JWTMiddleware(cfg, keyManager))
	getAttributesHandler := middlewares.ChainMiddleware(handler.GetAttributes, middlewares.ErrorMiddleware(), middlewares.JWTMiddleware(cfg, keyManager))
	mux.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) { getUsersHandler(w, r) })
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) { getUserByIDHandler(w, r) })
	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) { createUserHandler(w, r) })

	mux.HandleFunc("POST /attributes", func(w http.ResponseWriter, r *http.Request) { createAttributeHandler(w, r) })
	mux.HandleFunc("GET /attributes", func(w http.ResponseWriter, r *http.Request) { getAttributesHandler(w, r) })
	// user attributes
	mux.HandleFunc("POST /users/{id}/attributes", func(w http.ResponseWriter, r *http.Request) { addUserAttributesHandler(w, r) })
	mux.HandleFunc("PATCH /users/{id}/attributes", func(w http.ResponseWriter, r *http.Request) { patchUserAttributesHandler(w, r) })
}
