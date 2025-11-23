package dto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProductInput_JSONMarshaling(t *testing.T) {
	t.Run("should marshal to JSON correctly", func(t *testing.T) {
		input := CreateProductInput{
			Name:  "Laptop",
			Price: 999.99,
		}

		jsonData, err := json.Marshal(input)
		require.NoError(t, err)

		expected := `{"name":"Laptop","price":999.99}`
		assert.JSONEq(t, expected, string(jsonData))
	})

	t.Run("should unmarshal from JSON correctly", func(t *testing.T) {
		jsonData := `{"name":"Mouse","price":29.99}`

		var input CreateProductInput
		err := json.Unmarshal([]byte(jsonData), &input)
		require.NoError(t, err)

		assert.Equal(t, "Mouse", input.Name)
		assert.Equal(t, 29.99, input.Price)
	})

	t.Run("should handle special characters in name", func(t *testing.T) {
		input := CreateProductInput{
			Name:  "Product @#$% Special!",
			Price: 49.99,
		}

		jsonData, err := json.Marshal(input)
		require.NoError(t, err)

		var decoded CreateProductInput
		err = json.Unmarshal(jsonData, &decoded)
		require.NoError(t, err)

		assert.Equal(t, input.Name, decoded.Name)
		assert.Equal(t, input.Price, decoded.Price)
	})
}

func TestUpdateProductInput_JSONMarshaling(t *testing.T) {
	t.Run("should marshal to JSON correctly", func(t *testing.T) {
		input := UpdateProductInput{
			Name:  "Updated Product",
			Price: 149.99,
		}

		jsonData, err := json.Marshal(input)
		require.NoError(t, err)

		expected := `{"name":"Updated Product","price":149.99}`
		assert.JSONEq(t, expected, string(jsonData))
	})

	t.Run("should unmarshal from JSON correctly", func(t *testing.T) {
		jsonData := `{"name":"Keyboard","price":79.99}`

		var input UpdateProductInput
		err := json.Unmarshal([]byte(jsonData), &input)
		require.NoError(t, err)

		assert.Equal(t, "Keyboard", input.Name)
		assert.Equal(t, 79.99, input.Price)
	})
}

func TestProductOutput_JSONMarshaling(t *testing.T) {
	t.Run("should marshal to JSON correctly", func(t *testing.T) {
		output := ProductOutput{
			ID:    "019ab24a-dc97-72a4-9056-cc09f4c13bef",
			Name:  "Monitor",
			Price: 299.99,
		}

		jsonData, err := json.Marshal(output)
		require.NoError(t, err)

		expected := `{"id":"019ab24a-dc97-72a4-9056-cc09f4c13bef","name":"Monitor","price":299.99}`
		assert.JSONEq(t, expected, string(jsonData))
	})

	t.Run("should unmarshal from JSON correctly", func(t *testing.T) {
		jsonData := `{"id":"019ab24a-dc97-72a4-9056-cc09f4c13bef","name":"Headphones","price":89.99}`

		var output ProductOutput
		err := json.Unmarshal([]byte(jsonData), &output)
		require.NoError(t, err)

		assert.Equal(t, "019ab24a-dc97-72a4-9056-cc09f4c13bef", output.ID)
		assert.Equal(t, "Headphones", output.Name)
		assert.Equal(t, 89.99, output.Price)
	})
}

func TestProductListOutput_JSONMarshaling(t *testing.T) {
	t.Run("should marshal to JSON correctly", func(t *testing.T) {
		output := ProductListOutput{
			Products: []ProductOutput{
				{
					ID:    "019ab24a-dc97-72a4-9056-cc09f4c13bef",
					Name:  "Product 1",
					Price: 10.00,
				},
				{
					ID:    "019ab24a-dc97-72a4-9056-cc09f4c13bf0",
					Name:  "Product 2",
					Price: 20.00,
				},
			},
			Page:  1,
			Limit: 10,
			Total: 2,
		}

		jsonData, err := json.Marshal(output)
		require.NoError(t, err)

		var decoded ProductListOutput
		err = json.Unmarshal(jsonData, &decoded)
		require.NoError(t, err)

		assert.Equal(t, 2, len(decoded.Products))
		assert.Equal(t, 1, decoded.Page)
		assert.Equal(t, 10, decoded.Limit)
		assert.Equal(t, 2, decoded.Total)
	})

	t.Run("should handle empty product list", func(t *testing.T) {
		output := ProductListOutput{
			Products: []ProductOutput{},
			Page:     1,
			Limit:    10,
			Total:    0,
		}

		jsonData, err := json.Marshal(output)
		require.NoError(t, err)

		var decoded ProductListOutput
		err = json.Unmarshal(jsonData, &decoded)
		require.NoError(t, err)

		assert.Equal(t, 0, len(decoded.Products))
		assert.Equal(t, 1, decoded.Page)
		assert.Equal(t, 10, decoded.Limit)
		assert.Equal(t, 0, decoded.Total)
	})
}

func TestCreateProductInput_ZeroValues(t *testing.T) {
	t.Run("should handle zero values", func(t *testing.T) {
		input := CreateProductInput{}

		assert.Equal(t, "", input.Name)
		assert.Equal(t, 0.0, input.Price)
	})
}

func TestUpdateProductInput_ZeroValues(t *testing.T) {
	t.Run("should handle zero values", func(t *testing.T) {
		input := UpdateProductInput{}

		assert.Equal(t, "", input.Name)
		assert.Equal(t, 0.0, input.Price)
	})
}

func TestProductOutput_ZeroValues(t *testing.T) {
	t.Run("should handle zero values", func(t *testing.T) {
		output := ProductOutput{}

		assert.Equal(t, "", output.ID)
		assert.Equal(t, "", output.Name)
		assert.Equal(t, 0.0, output.Price)
	})
}

func TestCreateProductInput_LargePrice(t *testing.T) {
	t.Run("should handle large price values", func(t *testing.T) {
		input := CreateProductInput{
			Name:  "Expensive Item",
			Price: 99999999.99,
		}

		jsonData, err := json.Marshal(input)
		require.NoError(t, err)

		var decoded CreateProductInput
		err = json.Unmarshal(jsonData, &decoded)
		require.NoError(t, err)

		assert.Equal(t, input.Name, decoded.Name)
		assert.Equal(t, input.Price, decoded.Price)
	})
}

func TestCreateProductInput_SmallPrice(t *testing.T) {
	t.Run("should handle small price values", func(t *testing.T) {
		input := CreateProductInput{
			Name:  "Cheap Item",
			Price: 0.01,
		}

		jsonData, err := json.Marshal(input)
		require.NoError(t, err)

		var decoded CreateProductInput
		err = json.Unmarshal(jsonData, &decoded)
		require.NoError(t, err)

		assert.Equal(t, input.Name, decoded.Name)
		assert.Equal(t, input.Price, decoded.Price)
	})
}

func TestProductListOutput_Pagination(t *testing.T) {
	t.Run("should handle different pagination values", func(t *testing.T) {
		output := ProductListOutput{
			Products: []ProductOutput{},
			Page:     5,
			Limit:    20,
			Total:    100,
		}

		jsonData, err := json.Marshal(output)
		require.NoError(t, err)

		var decoded ProductListOutput
		err = json.Unmarshal(jsonData, &decoded)
		require.NoError(t, err)

		assert.Equal(t, 5, decoded.Page)
		assert.Equal(t, 20, decoded.Limit)
		assert.Equal(t, 100, decoded.Total)
	})
}
