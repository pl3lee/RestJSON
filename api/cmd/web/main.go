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
	"github.com/pl3lee/webjson/internal/auth"
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
	baseUrl := os.Getenv("WEB_BASE_URL")
	googleClientID, googleClientSecret := os.Getenv("WEB_GOOGLE_CLIENT_ID"), os.Getenv("WEB_GOOGLE_CLIENT_SECRET")

	pgDb, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("cannot open database: %v", err)
	}
	dbQueries := database.New(pgDb)

	cfg := config.Config{
		Port:               port,
		ClientURL:          clientURL,
		DbUrl:              dbUrl,
		Secret:             sec,
		WebBaseURL:         baseUrl,
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		Db:                 dbQueries,
	}

	authConfig := auth.AuthConfig{
		Db:                 cfg.Db,
		Secret:             cfg.Secret,
		WebBaseURL:         baseUrl,
		GoogleClientID:     cfg.GoogleClientID,
		GoogleClientSecret: cfg.GoogleClientSecret,
		ClientURL:          cfg.ClientURL,
	}

	log.Printf("env vars: %v\n", cfg)
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

	// public routes
	r.Get("/auth/google/login", authConfig.HandlerGoogleLogin)
	r.Get("/auth/google/callback", authConfig.HandlerGoogleCallback)

	r.Group(func(r chi.Router) {
		r.Use(authConfig.AuthMiddleware)

		r.Get("/me", authConfig.HandlerGetMe)
		r.Post("/logout", authConfig.HandlerLogout)
	})
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%v", cfg.Port),
		Handler:           r,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Listening on port %v", cfg.Port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("error starting server at port %v", cfg.Port)
	}

}
