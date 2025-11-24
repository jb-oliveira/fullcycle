package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/dto"
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/entity"
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/infra/database"
	entityPkg "github.com/jb-oliveira/fullcycle/tree/main/APIS/pkg/entity"
)

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

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := h.productDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	dto := dto.ProductOutput{
		ID:    product.ID.String(),
		Name:  product.Name,
		Price: product.Price,
	}
	json.NewEncoder(w).Encode(dto)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// adquire o id
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// deserializa
	var productDTO dto.UpdateProductInput
	err = json.NewDecoder(r.Body).Decode(&productDTO)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// carrega o produto
	product, err := h.productDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// salva o produto
	product.Name = productDTO.Name
	product.Price = productDTO.Price
	err = h.productDB.Update(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Retorna o produto
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := h.productDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = h.productDB.Delete(product.ID.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	limit := r.URL.Query().Get("limit")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}
	sort := r.URL.Query().Get("sort")
	sortDirection := r.URL.Query().Get("sort_direction")
	products, err := h.productDB.FindAll(pageInt, limitInt, sort, sortDirection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dtos := []dto.ProductOutput{}
	for _, product := range products {
		dtos = append(dtos, dto.ProductOutput{
			ID:    product.ID.String(),
			Name:  product.Name,
			Price: product.Price,
		})
	}
	count, err := h.productDB.Count()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := dto.ProductListOutput{
		Products: dtos,
		Page:     pageInt,
		Limit:    limitInt,
		Total:    int(count),
	}
	if sort != "" {
		result.Sort = sort
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
