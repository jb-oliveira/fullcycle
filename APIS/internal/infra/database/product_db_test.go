package database

import (
	"testing"

	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupProductTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err)

	err = db.AutoMigrate(&entity.Product{})
	require.NoError(t, err)

	return db
}

func TestNewProduct(t *testing.T) {
	db := setupProductTestDB(t)
	productDB := NewProductDB(db)

	assert.NotNil(t, productDB)
	assert.NotNil(t, productDB.db)
}

func TestProduct_Create(t *testing.T) {
	t.Run("should create product successfully", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		product, err := entity.NewProduct("Laptop", 999.99)
		require.NoError(t, err)

		err = productDB.Create(product)
		assert.NoError(t, err)

		var count int64
		db.Model(&entity.Product{}).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("should create multiple products", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		product1, err := entity.NewProduct("Laptop", 999.99)
		require.NoError(t, err)
		err = productDB.Create(product1)
		require.NoError(t, err)

		product2, err := entity.NewProduct("Mouse", 29.99)
		require.NoError(t, err)
		err = productDB.Create(product2)
		require.NoError(t, err)

		var count int64
		db.Model(&entity.Product{}).Count(&count)
		assert.Equal(t, int64(2), count)
	})
}

func TestProduct_FindByID(t *testing.T) {
	t.Run("should find product by ID successfully", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		product, err := entity.NewProduct("Keyboard", 79.99)
		require.NoError(t, err)

		err = productDB.Create(product)
		require.NoError(t, err)

		foundProduct, err := productDB.FindByID(product.ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, foundProduct)
		assert.Equal(t, product.ID, foundProduct.ID)
		assert.Equal(t, product.Name, foundProduct.Name)
		assert.Equal(t, product.Price, foundProduct.Price)
	})

	t.Run("should return error when product not found", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		foundProduct, err := productDB.FindByID("019ab24a-dc97-72a4-9056-cc09f4c13bef")
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NotNil(t, foundProduct)
	})

	t.Run("should return error for invalid ID format", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		foundProduct, err := productDB.FindByID("invalid-uuid")
		assert.Error(t, err)
		assert.Nil(t, foundProduct)
	})
}

func TestProduct_FindAll(t *testing.T) {
	t.Run("should find all products with pagination", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		// Create test products
		products := []struct {
			name  string
			price float64
		}{
			{"Product A", 10.00},
			{"Product C", 30.00},
			{"Product E", 50.00},
			{"Product D", 40.00},
			{"Product B", 20.00},
		}

		for _, p := range products {
			product, err := entity.NewProduct(p.name, p.price)
			require.NoError(t, err)
			err = productDB.Create(product)
			require.NoError(t, err)
		}

		// Test first page
		result, err := productDB.FindAll(1, 2, "prd_name ASC")
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Product A", result[0].Name)
		assert.Equal(t, "Product B", result[1].Name)

		// Test second page
		result, err = productDB.FindAll(2, 2, "prd_name ASC")
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Product C", result[0].Name)
		assert.Equal(t, "Product D", result[1].Name)
	})

	t.Run("should return empty list when no products exist", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		result, err := productDB.FindAll(1, 10, "prd_name ASC")
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("should sort by price descending", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		product1, _ := entity.NewProduct("Cheap", 10.00)
		product2, _ := entity.NewProduct("Expensive", 100.00)
		product3, _ := entity.NewProduct("Medium", 50.00)

		productDB.Create(product1)
		productDB.Create(product2)
		productDB.Create(product3)

		result, err := productDB.FindAll(1, 10, "prd_price DESC")
		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, "Expensive", result[0].Name)
		assert.Equal(t, "Medium", result[1].Name)
		assert.Equal(t, "Cheap", result[2].Name)
	})
}

func TestProduct_Update(t *testing.T) {
	t.Run("should update product successfully", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		product, err := entity.NewProduct("Old Name", 99.99)
		require.NoError(t, err)

		err = productDB.Create(product)
		require.NoError(t, err)

		// Update product
		product.Name = "New Name"
		product.Price = 149.99

		err = productDB.Update(product)
		assert.NoError(t, err)

		// Verify update
		foundProduct, err := productDB.FindByID(product.ID.String())
		assert.NoError(t, err)
		assert.Equal(t, "New Name", foundProduct.Name)
		assert.Equal(t, 149.99, foundProduct.Price)
	})

	t.Run("should update only price", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		product, err := entity.NewProduct("Product", 50.00)
		require.NoError(t, err)

		err = productDB.Create(product)
		require.NoError(t, err)

		// Update only price
		product.Price = 75.00

		err = productDB.Update(product)
		assert.NoError(t, err)

		// Verify update
		foundProduct, err := productDB.FindByID(product.ID.String())
		assert.NoError(t, err)
		assert.Equal(t, "Product", foundProduct.Name)
		assert.Equal(t, 75.00, foundProduct.Price)
	})
}

func TestProduct_Delete(t *testing.T) {
	t.Run("should delete product successfully", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		product, err := entity.NewProduct("To Delete", 99.99)
		require.NoError(t, err)

		err = productDB.Create(product)
		require.NoError(t, err)

		// Delete product
		err = productDB.Delete(product.ID.String())
		assert.NoError(t, err)

		// Verify deletion (soft delete)
		var count int64
		db.Model(&entity.Product{}).Unscoped().Where("id = ?", product.ID).Count(&count)
		assert.Equal(t, int64(1), count)

		// Verify not found in normal query
		foundProduct, err := productDB.FindByID(product.ID.String())
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NotNil(t, foundProduct)
	})

	t.Run("should return error for invalid ID format", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		err := productDB.Delete("invalid-uuid")
		assert.Error(t, err)
	})

	t.Run("should not error when deleting non-existent product", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		err := productDB.Delete("019ab24a-dc97-72a4-9056-cc09f4c13bef")
		assert.NoError(t, err)
	})
}

func TestProduct_Count(t *testing.T) {
	t.Run("should return zero when no products exist", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		count, err := productDB.Count()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("should return correct count with products", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		// Create test products
		for i := 1; i <= 5; i++ {
			product, err := entity.NewProduct("Product", 10.00)
			require.NoError(t, err)
			err = productDB.Create(product)
			require.NoError(t, err)
		}

		count, err := productDB.Count()
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})

	t.Run("should not count deleted products", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		// Create products
		product1, _ := entity.NewProduct("Product 1", 10.00)
		product2, _ := entity.NewProduct("Product 2", 20.00)
		product3, _ := entity.NewProduct("Product 3", 30.00)

		productDB.Create(product1)
		productDB.Create(product2)
		productDB.Create(product3)

		// Delete one product
		err := productDB.Delete(product2.ID.String())
		require.NoError(t, err)

		// Count should be 2 (soft delete)
		count, err := productDB.Count()
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})

	t.Run("should return count after updates", func(t *testing.T) {
		db := setupProductTestDB(t)
		productDB := NewProductDB(db)

		// Create product
		product, err := entity.NewProduct("Original", 50.00)
		require.NoError(t, err)
		err = productDB.Create(product)
		require.NoError(t, err)

		// Update product
		product.Name = "Updated"
		product.Price = 75.00
		err = productDB.Update(product)
		require.NoError(t, err)

		// Count should still be 1
		count, err := productDB.Count()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})
}
