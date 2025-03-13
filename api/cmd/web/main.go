package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pl3lee/webjson/internal/auth"
	"github.com/pl3lee/webjson/internal/database"
	"github.com/pl3lee/webjson/internal/jsonfile"
	"github.com/pl3lee/webjson/internal/utils"
)

type appConfig struct {
	port               string
	clientURL          string
	dbUrl              string
	webBaseURL         string
	googleClientID     string
	googleClientSecret string
	db                 *database.Queries
	s3Bucket           string
	s3Region           string
	s3Client           *s3.Client
}

func main() {
	// env variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Cannot read .env file, this is normal when running in a docker container: %v\n", err)
	}
	port := os.Getenv("WEB_PORT")
	if port == "" {
		log.Fatal("WEB_PORT not set")
	}
	clientURL := os.Getenv("SHARED_CLIENT_URL")
	if clientURL == "" {
		log.Fatal("SHARED_CLIENT_URL not set")
	}
	dbUrl := os.Getenv("SHARED_DB_URL")
	if dbUrl == "" {
		log.Fatal("SHARED_DB_URL not set")
	}
	baseUrl := os.Getenv("WEB_BASE_URL")
	if baseUrl == "" {
		log.Fatal("WEB_BASE_URL not set")
	}
	googleClientID, googleClientSecret := os.Getenv("WEB_GOOGLE_CLIENT_ID"), os.Getenv("WEB_GOOGLE_CLIENT_SECRET")
	if googleClientID == "" {
		log.Fatal("WEB_GOOGLE_CLIENT_ID not set")
	}
	if googleClientSecret == "" {
		log.Fatal("WEB_GOOGLE_CLIENT_SECRET not set")
	}
	s3Bucket := os.Getenv("SHARED_S3_BUCKET")
	if s3Bucket == "" {
		log.Fatal("SHARED_S3_BUCKET not set")
	}
	s3Region := os.Getenv("SHARED_S3_REGION")
	if s3Region == "" {
		log.Fatal("SHARED_S3_REGION not set")
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("cannot load aws config: %v", err)
	}
	client := s3.NewFromConfig(awsCfg)

	pgDb, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("cannot open database: %v", err)
	}
	dbQueries := database.New(pgDb)

	cfg := appConfig{
		port:               port,
		clientURL:          clientURL,
		dbUrl:              dbUrl,
		webBaseURL:         baseUrl,
		googleClientID:     googleClientID,
		googleClientSecret: googleClientSecret,
		db:                 dbQueries,
		s3Bucket:           s3Bucket,
		s3Region:           s3Region,
		s3Client:           client,
	}

	authConfig := auth.AuthConfig{
		Db:                 cfg.db,
		WebBaseURL:         cfg.webBaseURL,
		GoogleClientID:     cfg.googleClientID,
		GoogleClientSecret: cfg.googleClientSecret,
		ClientURL:          cfg.clientURL,
	}

	jsonConfig := jsonfile.JsonConfig{
		Db:         cfg.db,
		WebBaseURL: cfg.webBaseURL,
		ClientURL:  cfg.clientURL,
		S3Bucket:   cfg.s3Bucket,
		S3Region:   cfg.s3Region,
		S3Client:   cfg.s3Client,
	}

	log.Printf("env vars: %v\n", cfg)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	corsWeb := cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.clientURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
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
		r.Put("/logout", authConfig.HandlerLogout)

		r.Post("/jsonfiles", jsonConfig.HandlerCreateJson)
		r.Get("/jsonfiles", jsonConfig.HandlerGetJsonFiles)

		r.Group(func(r chi.Router) {
			r.Use(jsonConfig.JsonFileMiddleware)

			r.Get("/jsonfiles/{fileId}", jsonConfig.HandlerGetJson)
			r.Get("/jsonfiles/{fileId}/metadata", jsonConfig.HandlerGetJsonMetadata)
			r.Patch("/jsonfiles/{fileId}", jsonConfig.HandlerRenameJsonFile)
			r.Put("/jsonfiles/{fileId}", jsonConfig.HandlerUpdateJson)
			r.Delete("/jsonfiles/{fileId}", jsonConfig.HandlerDeleteJsonFile)
		})
	})
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%v", cfg.port),
		Handler:           r,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Listening on port %v", cfg.port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("error starting server at port %v", cfg.port)
	}

}
