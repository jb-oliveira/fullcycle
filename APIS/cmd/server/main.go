package main

import (
	"log"
	"net/http"

	"github.com/jb-oliveira/fullcycle/tree/main/APIS/configs"
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/infra/database"
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
	http.ListenAndServe(":8000", nil)
}

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {

}

func NewProductHandler(db database.ProductInterface) *ProductHandler {

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
	log.Println("Conexão com banco estabelecida")
}
