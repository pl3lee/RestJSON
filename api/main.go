package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/pl3lee/webjson/internal/config"
	"github.com/pl3lee/webjson/internal/utils"
)

func main() {
	// env variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Cannot read .env file, this is normal when running in a docker container: %v\n", err)
	}
	port := os.Getenv("PORT")
	clientURL := os.Getenv("CLIENT_URL")
	dbUrl := os.Getenv("DB_URL")
	sec := os.Getenv("AUTH_SECRET")
	googleClientID, googleClientSecret := os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET")
	githubClientID, githubClientSecret := os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET")

	cfg := config.Config{
		Port:               port,
		ClientURL:          clientURL,
		DbUrl:              dbUrl,
		Secret:             sec,
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		GithubClientID:     githubClientID,
		GithubClientSecret: githubClientSecret,
	}

	log.Printf("env vars: %s\n", cfg)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	corsWeb := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.ClientURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	corsPublic := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	})
	r.Route("/web", func(r chi.Router) {
		r.Use(corsWeb.Handler)
		r.Options("/*", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			utils.RespondWithJSON(w, http.StatusOK, "Hello world from web router!")
		})
	})

	r.Route("/public", func(r chi.Router) {
		r.Use(corsPublic.Handler)
		r.Options("/*", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			utils.RespondWithJSON(w, http.StatusOK, "Hello world from public router!")
		})

	})

	log.Printf("Listening on port %v", cfg.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), r)
	if err != nil {
		log.Fatalf("error starting server at port %v", cfg.Port)
	}

}
