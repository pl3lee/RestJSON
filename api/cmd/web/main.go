package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pl3lee/webjson/internal/config"
	"github.com/pl3lee/webjson/internal/database"
	"github.com/pl3lee/webjson/internal/utils"
)

func main() {
	// env variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Cannot read .env file, this is normal when running in a docker container: %v\n", err)
	}
	port := os.Getenv("WEB_PORT")
	clientURL := os.Getenv("SHARED_CLIENT_URL")
	dbUrl := os.Getenv("SHARED_DB_URL")
	sec := os.Getenv("WEB_AUTH_SECRET")
	googleClientID, googleClientSecret := os.Getenv("WEB_GOOGLE_CLIENT_ID"), os.Getenv("WEB_GOOGLE_CLIENT_SECRET")
	githubClientID, githubClientSecret := os.Getenv("WEB_GITHUB_CLIENT_ID"), os.Getenv("WEB_GITHUB_CLIENT_SECRET")

	pgDb, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("cannot open database: %w", err)
	}
	dbQueries := database.New(pgDb)

	cfg := config.Config{
		Port:               port,
		ClientURL:          clientURL,
		DbUrl:              dbUrl,
		Secret:             sec,
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		GithubClientID:     githubClientID,
		GithubClientSecret: githubClientSecret,
		Db:                 dbQueries,
	}

	log.Printf("env vars: %s\n", cfg)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	corsWeb := cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.ClientURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(corsWeb)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, http.StatusOK, "Hello world from web api!")
	})

	log.Printf("Listening on port %v", cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), r)
	if err != nil {
		log.Fatalf("error starting server at port %v", cfg.Port)
	}

}
