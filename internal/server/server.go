package server

import (
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
	"github.com/shashimalcse/tiny-is/internal/application"
	cs "github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/server/routes"
	"github.com/shashimalcse/tiny-is/internal/user"
)

func StartServer() {
	c := cache.New(5*time.Minute, 10*time.Minute)

	cacheService := cs.NewCacheService(c)

	db, err := sqlx.Connect("postgres", "user=postgres dbname=tinydb sslmode=disable password=tinydb host=localhost")
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}

	applicationService := application.NewApplicationService(cacheService, db)
	userService := user.NewUserService(cacheService, db)
	router := routes.NewRouter(cacheService, applicationService, userService)
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