package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jb-oliveira/fullcycle/tree/main/APIS/configs"
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/dto"
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/entity"
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

	productDB := database.NewProductDB(configs.GetDB())
	productHandler := ProductHandler{productDB: productDB}

	http.HandleFunc("/products", productHandler.CreateProduct)
	http.ListenAndServe(":8000", nil)
}

type ProductHandler struct {
	productDB database.ProductInterface
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var productDTO dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&productDTO)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Deveria ser pelo Caso de Uso, mas por enquanto ta indo direto mesmo
	p, err := entity.NewProduct(productDTO.Name, productDTO.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.productDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{productDB: db}
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
