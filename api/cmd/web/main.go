package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	"github.com/pl3lee/restjson/internal/payment"
	"github.com/pl3lee/restjson/internal/ratelimit"
	"github.com/pl3lee/restjson/internal/utils"
	"github.com/redis/go-redis/v9"
)

type appConfig struct {
	port                string
	clientURL           string
	dbUrl               string
	baseURL             string
	googleClientID      string
	googleClientSecret  string
	db                  *database.Queries
	s3Bucket            string
	s3Region            string
	s3Client            *s3.Client
	rdb                 *redis.Client
	freeFileLimit       int
	proFileLimit        int
	stripeSecretKey     string
	stripeWebhookSecret string
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
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		log.Fatal("REDIS_URL not set")
	}
	redisOpts, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatal("invalid redis url")
	}
	freeFileLimitStr := os.Getenv("FREE_FILE_LIMIT")
	if freeFileLimitStr == "" {
		log.Fatal("invalid file limit")
	}
	freeFileLimit, err := strconv.Atoi(freeFileLimitStr)
	if err != nil {
		log.Fatal("file limit should be an integer")
	}
	proFileLimitStr := os.Getenv("PRO_FILE_LIMIT")
	if proFileLimitStr == "" {
		log.Fatal("invalid file limit")
	}
	proFileLimit, err := strconv.Atoi(proFileLimitStr)
	if err != nil {
		log.Fatal("file limit should be an integer")
	}
	stripeSecretKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeSecretKey == "" {
		log.Fatal("STRIPE_SECRET_KEY not set")
	}
	stripeWebhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if stripeWebhookSecret == "" {
		log.Fatal("STRIPE_WEBHOOK_SECRET not set")
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

	rdb := redis.NewClient(redisOpts)

	cfg := &appConfig{
		port:                port,
		clientURL:           clientURL,
		dbUrl:               dbUrl,
		baseURL:             baseUrl,
		googleClientID:      googleClientID,
		googleClientSecret:  googleClientSecret,
		db:                  dbQueries,
		s3Bucket:            s3Bucket,
		s3Region:            s3Region,
		s3Client:            client,
		rdb:                 rdb,
		freeFileLimit:       freeFileLimit,
		proFileLimit:        proFileLimit,
		stripeSecretKey:     stripeSecretKey,
		stripeWebhookSecret: stripeWebhookSecret,
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
		Rdb:                cfg.rdb,
		S3Bucket:           cfg.s3Bucket,
		S3Region:           cfg.s3Region,
		S3Client:           cfg.s3Client,
	}
	return authConfig
}

func loadJsonConfig(cfg *appConfig) *jsonfile.JsonConfig {
	jsonConfig := &jsonfile.JsonConfig{
		Db:            cfg.db,
		BaseURL:       cfg.baseURL,
		ClientURL:     cfg.clientURL,
		S3Bucket:      cfg.s3Bucket,
		S3Region:      cfg.s3Region,
		S3Client:      cfg.s3Client,
		Rdb:           cfg.rdb,
		FreeFileLimit: cfg.freeFileLimit,
		ProFileLimit:  cfg.proFileLimit,
	}
	return jsonConfig
}

func loadPaymentConfig(cfg *appConfig) *payment.PaymentConfig {
	paymentConfig := &payment.PaymentConfig{
		Db:                  cfg.db,
		BaseURL:             cfg.baseURL,
		ClientURL:           cfg.clientURL,
		Rdb:                 cfg.rdb,
		StripeSecretKey:     cfg.stripeSecretKey,
		StripeWebhookSecret: cfg.stripeWebhookSecret,
	}
	return paymentConfig
}

func main() {
	appConfig := loadAppConfig()
	authConfig := loadAuthConfig(appConfig)
	jsonConfig := loadJsonConfig(appConfig)
	paymentConfig := loadPaymentConfig(appConfig)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/", webRouter(appConfig, authConfig, jsonConfig, paymentConfig))
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

func webRouter(appConfig *appConfig, authConfig *auth.AuthConfig, jsonConfig *jsonfile.JsonConfig, paymentConfig *payment.PaymentConfig) http.Handler {
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
	r.Post("/webhooks/stripe", paymentConfig.HandlerStripeWebhook)

	r.Group(func(r chi.Router) {
		// middleware, capacity 10, refill rate 1, expiration 60 seconds
		r.Use(ratelimit.TokenBucketRateLimiter(appConfig.rdb, 10, 1, 60))
		r.Use(authConfig.SessionMiddleware)

		r.Get("/me", authConfig.HandlerGetMe)
		r.Put("/logout", authConfig.HandlerLogout)
		r.Delete("/users", authConfig.HandlerDeleteAccount)
		r.Post("/apikeys", authConfig.HandlerCreateApiKey)
		r.Get("/apikeys", authConfig.HandlerGetAllApiKeys)
		r.Delete("/apikeys/{keyHash}", authConfig.HandlerDeleteApiKey)

		r.Post("/jsonfiles", jsonConfig.HandlerCreateJson)
		r.Get("/jsonfiles", jsonConfig.HandlerGetJsonFiles)

		r.Post("/subscriptions/checkout", paymentConfig.HandlerCheckout)
		r.Post("/subscriptions/success", paymentConfig.HandlerSuccess)
		r.Get("/subscriptions", paymentConfig.HandlerGetSubscriptionStatus)
		r.Get("/subscriptions/manage", paymentConfig.HandlerCustomerPortal)

		r.Group(func(r chi.Router) {
			r.Use(jsonConfig.JsonFileMiddleware)
			r.Use(jsonConfig.JsonFileContentMiddleware)

			r.Get("/jsonfiles/{fileId}", jsonConfig.HandlerGetJson)
			r.Get("/jsonfiles/{fileId}/metadata", jsonConfig.HandlerGetJsonMetadata)
			r.Patch("/jsonfiles/{fileId}", jsonConfig.HandlerRenameJsonFile)
			r.Put("/jsonfiles/{fileId}", jsonConfig.HandlerUpdateJson)
			r.Delete("/jsonfiles/{fileId}", jsonConfig.HandlerDeleteJsonFile)

			r.Get("/jsonfiles/{fileId}/routes", jsonConfig.HandlerGetDynamicRoutes)
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
		// middleware, capacity 5, refill rate 1, expiration 60 seconds
		r.Use(ratelimit.TokenBucketRateLimiter(jsonConfig.Rdb, 5, 1, 60))
		r.Use(authConfig.ApiKeyMiddleware)

		r.Group(func(r chi.Router) {
			r.Use(jsonConfig.JsonFileMiddleware)
			r.Use(jsonConfig.JsonFileContentMiddleware)

			r.Get("/{fileId}", jsonConfig.HandlerGetJson)

			r.Group(func(r chi.Router) {
				r.Use(jsonConfig.ResourceMiddleware)

				r.Get("/{fileId}/{resource}", jsonConfig.HandlerGetResource)
				r.Put("/{fileId}/{resource}", jsonConfig.HandlerUpdateResource)
				r.Patch("/{fileId}/{resource}", jsonConfig.HandlerPartialUpdateResource)

				r.Group(func(r chi.Router) {
					r.Use(jsonConfig.ResourceArrayMiddleware)

					r.Get("/{fileId}/{resource}/{id}", jsonConfig.HandlerGetResourceItem)
					r.Post("/{fileId}/{resource}", jsonConfig.HandlerCreateResourceItem)
					r.Put("/{fileId}/{resource}/{id}", jsonConfig.HandlerUpdateResourceItem)
					r.Patch("/{fileId}/{resource}/{id}", jsonConfig.HandlerPartialUpdateResourceItem)
					r.Delete("/{fileId}/{resource}/{id}", jsonConfig.HandlerDeleteResourceItem)
				})
			})
		})
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, http.StatusOK, "Hello world from public api!")
	})
	return r
}
