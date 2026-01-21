package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/jwtauth"
	"github.com/jb-oliveira/fullcycle/APIS/configs"
	"github.com/jb-oliveira/fullcycle/APIS/internal/entity"
	"github.com/jb-oliveira/fullcycle/APIS/internal/infra/database"
	"github.com/jb-oliveira/fullcycle/APIS/internal/infra/webserver/handlers"

	_ "github.com/jb-oliveira/fullcycle/APIS/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           FullCycle API
// @version         1.0
// @description     This is a sample API for FullCycle course
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Initialize the database
	initDB()
	// initialize the web
	_, err := configs.LoadWebConfig(".")
	if err != nil {
		log.Fatalf("failed to load web configuration: %v", err)
	}
	log.Println("Configuration loaded successfully")

	productDB := database.NewProductDB(configs.GetDB())
	productHandler := handlers.NewProductHandler(productDB)

	r := chi.NewRouter()
	logger := httplog.NewLogger("fullcycle-api", httplog.Options{
		JSON:     true, // Structured JSON for prod
		LogLevel: slog.LevelInfo,
		Concise:  true, // Clean logs with fewer details
	})
	r.Use(httplog.RequestLogger(logger))
	// r.Use(middleware.Logger)
	r.Use(MiddlewareVazio)
	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.GetWebConfig().TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
		r.Get("/", productHandler.GetProducts)
	})

	userDB := database.NewUserDB(configs.GetDB())
	userHandler := handlers.NewUserHandler(userDB, configs.GetWebConfig().TokenAuth, configs.GetWebConfig().JWTExpiration)

	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/auth", userHandler.Auth)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	http.ListenAndServe(":8000", r)
}

func initDB() {
	_, err := configs.LoadDbConfig(".")
	if err != nil {
		log.Fatalf("failed to load database configuration: %v", err)
	}
	err = configs.InitGorm()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	db := configs.GetDB()
	if db == nil {
		log.Fatal("database instance is null")
	}

	// Remove auto migrate and later see which is the best migration for GO
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	log.Println("Database connection established")
}

func MiddlewareVazio(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
