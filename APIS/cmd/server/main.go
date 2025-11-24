package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/configs"
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/entity"
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/infra/database"
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/infra/webserver/handlers"
)

func main() {
	// Inicializa o banco
	initDB()
	// inicializa a web
	_, err := configs.LoadWebConfig(".")
	if err != nil {
		log.Fatalf("falha ao carregar configuração web: %v", err)
	}
	log.Println("Configuração carregada com sucesso")

	productDB := database.NewProductDB(configs.GetDB())
	productHandler := handlers.NewProductHandler(productDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/products", productHandler.CreateProduct)
	r.Get("/products/{id}", productHandler.GetProduct)
	r.Put("/products/{id}", productHandler.UpdateProduct)
	r.Delete("/products/{id}", productHandler.DeleteProduct)
	r.Get("/products", productHandler.GetProducts)

	userDB := database.NewUserDB(configs.GetDB())
	userHandler := handlers.NewUserHandler(userDB)

	r.Post("/users", userHandler.CreateUser)

	http.ListenAndServe(":8000", r)
}

func initDB() {
	_, err := configs.LoadDbConfig(".")
	if err != nil {
		log.Fatalf("falha ao carregar configuração do banco: %v", err)
	}
	err = configs.InitGorm()
	if err != nil {
		log.Fatalf("falha ao conectar ao banco: %v", err)
	}

	db := configs.GetDB()
	if db == nil {
		log.Fatal("instância do banco é nula")
	}

	// Remover o auto migrate e depois ver qual melhor migration pra GO
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	log.Println("Conexão com banco estabelecida")
}
