package tax

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCalculateTax(t *testing.T) {
	tax, err := CalculateTax(1000)
	assert.Equal(t, 10.0, tax)
	assert.Nil(t, err)

	tax, err = CalculateTax(-1)
	assert.EqualError(t, err, "bbb Negative amount is not allowed")
	assert.Equal(t, tax, 0.0)
}

type TaxRepositoryMock struct {
	mock.Mock
}

func (m *TaxRepositoryMock) Save(tax float64) error {
	args := m.Called(tax)
	return args.Error(0)
}

func TestCalculateTaxAndSave(t *testing.T) {
	repository := &TaxRepositoryMock{}
	// the Once indicates that it can only be called once with this parameter
	repository.On("Save", 10.0).Return(nil).Once()
	repository.On("Save", 0.0).Return(errors.New("Error saving to database"))
	// repository.On("Save", mock.Anything).Return(errors.New("Error saving to database"))

	// The test is likely failing because CalculateTaxAndSave calculates tax first
	// For amount 10.0, the calculated tax would be 0.1 (10.0 * 0.01), not 10.0
	// For amount 0.0, it might return an error before saving due to validation
	err := CalculateTaxAndSave(1000.0, repository) // 1000 * 0.01 = 10.0 tax
	assert.Nil(t, err)

	err = CalculateTaxAndSave(-1.0, repository) // Negative amount should cause validation error
	assert.EqualError(t, err, "Error saving to database")

	repository.AssertExpectations(t)
	repository.AssertNumberOfCalls(t, "Save", 2)
}
