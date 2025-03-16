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
	"github.com/pl3lee/restjson/internal/auth"
	"github.com/pl3lee/restjson/internal/database"
	"github.com/pl3lee/restjson/internal/jsonfile"
	"github.com/pl3lee/restjson/internal/utils"
)

type appConfig struct {
	port               string
	clientURL          string
	dbUrl              string
	baseURL            string
	googleClientID     string
	googleClientSecret string
	db                 *database.Queries
	s3Bucket           string
	s3Region           string
	s3Client           *s3.Client
}

func loadAppConfig() *appConfig {
	// env variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Cannot read .env file, this is normal when running in a docker container: %v\n", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not set")
	}
	clientURL := os.Getenv("CLIENT_URL")
	if clientURL == "" {
		log.Fatal("CLIENT_URL not set")
	}
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL not set")
	}
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		log.Fatal("BASE_URL not set")
	}
	googleClientID, googleClientSecret := os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET")
	if googleClientID == "" {
		log.Fatal("GOOGLE_CLIENT_ID not set")
	}
	if googleClientSecret == "" {
		log.Fatal("GOOGLE_CLIENT_SECRET not set")
	}
	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		log.Fatal("S3_BUCKET not set")
	}
	s3Region := os.Getenv("S3_REGION")
	if s3Region == "" {
		log.Fatal("S3_REGION not set")
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

	cfg := &appConfig{
		port:               port,
		clientURL:          clientURL,
		dbUrl:              dbUrl,
		baseURL:            baseUrl,
		googleClientID:     googleClientID,
		googleClientSecret: googleClientSecret,
		db:                 dbQueries,
		s3Bucket:           s3Bucket,
		s3Region:           s3Region,
		s3Client:           client,
	}
	return cfg
}

func loadAuthConfig(cfg *appConfig) *auth.AuthConfig {
	authConfig := &auth.AuthConfig{
		Db:                 cfg.db,
		BaseURL:            cfg.baseURL,
		GoogleClientID:     cfg.googleClientID,
		GoogleClientSecret: cfg.googleClientSecret,
		ClientURL:          cfg.clientURL,
	}
	return authConfig
}

func loadJsonConfig(cfg *appConfig) *jsonfile.JsonConfig {
	jsonConfig := &jsonfile.JsonConfig{
		Db:        cfg.db,
		BaseURL:   cfg.baseURL,
		ClientURL: cfg.clientURL,
		S3Bucket:  cfg.s3Bucket,
		S3Region:  cfg.s3Region,
		S3Client:  cfg.s3Client,
	}
	return jsonConfig
}

func main() {
	appConfig := loadAppConfig()
	authConfig := loadAuthConfig(appConfig)
	jsonConfig := loadJsonConfig(appConfig)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/", webRouter(appConfig, authConfig, jsonConfig))
	r.Mount("/public", publicRouter(authConfig, jsonConfig))

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%v", appConfig.port),
		Handler:           r,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Listening on port %v", appConfig.port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("error starting server at port %v", appConfig.port)
	}

}

func webRouter(appConfig *appConfig, authConfig *auth.AuthConfig, jsonConfig *jsonfile.JsonConfig) http.Handler {
	r := chi.NewRouter()

	corsWeb := cors.Handler(cors.Options{
		AllowedOrigins:   []string{appConfig.clientURL},
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
		r.Use(authConfig.SessionMiddleware)

		r.Get("/me", authConfig.HandlerGetMe)
		r.Put("/logout", authConfig.HandlerLogout)
		r.Post("/apikeys", authConfig.HandlerCreateApiKey)
		r.Get("/apikeys", authConfig.HandlerGetAllApiKeys)
		r.Delete("/apikeys/{keyHash}", authConfig.HandlerDeleteApiKey)

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

	return r
}

func publicRouter(authConfig *auth.AuthConfig, jsonConfig *jsonfile.JsonConfig) http.Handler {
	r := chi.NewRouter()

	corsPublic := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	})
	r.Use(corsPublic)
	r.Group(func(r chi.Router) {
		r.Use(authConfig.ApiKeyMiddleware)

		r.Group(func(r chi.Router) {
			r.Use(jsonConfig.JsonFileMiddleware)

			r.Get("/{fileId}", jsonConfig.HandlerGetJson)
			// TODO: auto generate routes, say
			// GET /posts
			r.Get("/{fileId}/{resource}", jsonConfig.HandlerGetResource)
			// GET /posts/:id
			r.Get("/{fileId}/{resource}/{id}", jsonConfig.HandlerGetResourceById)
			// POST /posts
			// PUT /posts/:id
			// DELETE /posts/:id
		})
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, http.StatusOK, "Hello world from public api!")
	})
	return r
}
