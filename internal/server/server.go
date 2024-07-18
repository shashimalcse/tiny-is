package server

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shashimalcse/tiny-is/internal/application"
	cs "github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/config"
	"github.com/shashimalcse/tiny-is/internal/oauth2/token"
	"github.com/shashimalcse/tiny-is/internal/organization"
	"github.com/shashimalcse/tiny-is/internal/server/routes"
	"github.com/shashimalcse/tiny-is/internal/server/utils"
	"github.com/shashimalcse/tiny-is/internal/session"
	"github.com/shashimalcse/tiny-is/internal/user"
)

func StartServer(cfg *config.Config) {

	cacheService := cs.NewCacheService()
	sessionStore := session.NewInMemorySessionStore()

	db, err := sqlx.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal(err)
	}
	organizationService := organization.NewOrganizationService(cacheService, organization.NewOrganizationRepository(db))
	applicationService := application.NewApplicationService(cacheService, application.NewApplicationRepository(db))
	userService := user.NewUserService(cacheService, user.NewUserRepository(db))
	tokenService := token.NewTokenService(cacheService, token.NewTokenRepository(db), []byte("secret"))
	err = utils.InitServer(cfg, db, organizationService, applicationService, userService)
	if err != nil {
		log.Fatal(err)
	}
	router := routes.NewRouter(cfg, cacheService, sessionStore, organizationService, applicationService, userService, tokenService)
	loggedRouter := LoggingMiddleware(router)
	if err := http.ListenAndServe(":9444", loggedRouter); err != nil {
		panic(err)
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method: %s, URL: %s, Origin: %s", r.Method, r.URL.String(), r.Header.Get("Origin"))
		next.ServeHTTP(w, r)
		log.Printf("Request completed with Method: %s, URL: %s", r.Method, r.URL.String())
	})
}
