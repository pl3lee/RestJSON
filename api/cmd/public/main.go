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
	port          string
	dbUrl         string
	publicBaseUrl string
	db            *database.Queries
	s3Bucket      string
	s3Region      string
	s3Client      *s3.Client
}

func main() {
	// env variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Cannot read .env file, this is normal when running in a docker container: %v\n", err)
	}
	port := os.Getenv("PUBLIC_PORT")
	if port == "" {
		log.Fatal("PUBLIC_PORT not set")
	}
	dbUrl := os.Getenv("SHARED_DB_URL")
	if dbUrl == "" {
		log.Fatal("SHARED_DB_URL not set")
	}
	baseUrl := os.Getenv("PUBLIC_BASE_URL")
	if baseUrl == "" {
		log.Fatal("PUBLIC_BASE_URL not set")
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
		port:          port,
		dbUrl:         dbUrl,
		publicBaseUrl: baseUrl,
		db:            dbQueries,
		s3Bucket:      s3Bucket,
		s3Region:      s3Region,
		s3Client:      client,
	}

	authConfig := auth.AuthConfig{
		Db: cfg.db,
	}

	jsonConfig := jsonfile.JsonConfig{
		Db:       cfg.db,
		S3Bucket: cfg.s3Bucket,
		S3Region: cfg.s3Region,
		S3Client: cfg.s3Client,
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	corsPublic := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	})
	r.Use(corsPublic)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, http.StatusOK, "Hello world from public api!")
	})

	r.Group(func(r chi.Router) {
		r.Use(authConfig.ApiKeyMiddleware)

		r.Group(func(r chi.Router) {
			r.Use(jsonConfig.JsonFileMiddleware)

			r.Get("/{fileId}", jsonConfig.HandlerGetJson)
			// TODO: auto generate routes, say
			// GET /posts
			// GET /posts/:id
			// POST /posts
			// PUT /posts/:id
			// DELETE /posts/:id
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
