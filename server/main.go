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
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Cannot load env variables: %s", err)
		}

	}
	port := os.Getenv("PORT")
	clientURL := os.Getenv("CLIENT_URL")

	cfg := config.Config{
		Port:      port,
		ClientURL: clientURL,
	}

	fmt.Printf("CLIENT_URL: %v", cfg.ClientURL)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.ClientURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, http.StatusOK, "Hello")
	})

	log.Printf("Listening on port %v", cfg.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), r)
	if err != nil {
		log.Fatal("error starting server at port 3000")
	}

}
