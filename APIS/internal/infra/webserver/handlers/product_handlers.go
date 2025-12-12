package handlers

import (
	"encoding/json"
	"errors"
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

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{productDB: db}
}

// Create Product Godoc
// @Summary Create a new product
// @Description Create a new product
// @Tags Products
// @Accept json
// @Produce json
// @Param product body dto.CreateProductInput true "Product to create"
// @Success 201 {object} dto.ProductOutput
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var productDTO dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&productDTO)
	if err != nil {
		ReturnHttpError(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}
	// Deveria ser pelo Caso de Uso, mas por enquanto ta indo direto mesmo
	p, err := entity.NewProduct(productDTO.Name, productDTO.Price)
	if err != nil {
		ReturnHttpError(w, err, http.StatusBadRequest)
		return
	}
	err = h.productDB.Create(p)
	if err != nil {
		ReturnHttpError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Get Product Godoc
// @Summary Get a product by ID
// @Description Get a product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductOutput
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		ReturnHttpError(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}
	product, err := h.productDB.FindByID(id)
	if err != nil {
		ReturnHttpError(w, errors.New("product not found"), http.StatusNotFound)
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

// Update Product Godoc
// @Summary Update a product
// @Description Update a product
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body dto.UpdateProductInput true "Product to update"
// @Success 200 {object} dto.ProductOutput
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// adquire o id
	id := chi.URLParam(r, "id")
	if id == "" {
		ReturnHttpError(w, errors.New("id is required"), http.StatusBadRequest)
		return
	}
	_, err := entityPkg.ParseID(id)
	if err != nil {
		ReturnHttpError(w, err, http.StatusBadRequest)
		return
	}
	// deserializa
	var productDTO dto.UpdateProductInput
	err = json.NewDecoder(r.Body).Decode(&productDTO)
	if err != nil {
		ReturnHttpError(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}
	// carrega o produto
	product, err := h.productDB.FindByID(id)
	if err != nil {
		ReturnHttpError(w, errors.New("product not found"), http.StatusNotFound)
		return
	}
	// salva o produto
	product.Name = productDTO.Name
	product.Price = productDTO.Price
	err = h.productDB.Update(product)
	if err != nil {
		ReturnHttpError(w, err, http.StatusInternalServerError)
		return
	}
	// Retorna o produto
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Delete Product Godoc
// @Summary Delete a product
// @Description Delete a product
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		ReturnHttpError(w, errors.New("id is required"), http.StatusBadRequest)
		return
	}
	_, err := entityPkg.ParseID(id)
	if err != nil {
		ReturnHttpError(w, err, http.StatusBadRequest)
		return
	}
	product, err := h.productDB.FindByID(id)
	if err != nil {
		ReturnHttpError(w, errors.New("product not found"), http.StatusNotFound)
		return
	}
	err = h.productDB.Delete(product.ID.String())
	if err != nil {
		ReturnHttpError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// Get Products Godoc
// @Summary Get all products
// @Description Get all products
// @Tags Products
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param sort query string false "Sort by"
// @Param sort_direction query string false "Sort direction"
// @Success 200 {object} dto.ProductOutput
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /products [get]
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
	sortDir := r.URL.Query().Get("sort_direction")
	sorts := make(map[string]string)
	sorts["id"] = "prd_id"
	sorts["name"] = "prd_name"
	sorts["price"] = "prd_price"
	sortComplete := sorts[sort]
	if sortComplete == "" {
		sortComplete = "prd_id"
	}
	if sortDir == "desc" {
		sortComplete = sortComplete + " desc"
	} else {
		sortDir = "asc"
		sortComplete = sortComplete + " " + sortDir
	}
	products, err := h.productDB.FindAll(pageInt, limitInt, sortComplete)
	if err != nil {
		ReturnHttpError(w, err, http.StatusInternalServerError)
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
		ReturnHttpError(w, err, http.StatusInternalServerError)
		return
	}
	result := entityPkg.NewPage(dtos, pageInt, limitInt, int(count), sort, sortDir)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
