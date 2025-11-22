package tax

import "testing"

func TestCalculateTax(t *testing.T) {
	amount := 500.0
	expected := 5.0

	result := CalculateTax(amount)

	if result != expected {
		t.Errorf("Expected %f but got %f", expected, result)
	}
}

func TestCalculateTax_InBatch(t *testing.T) {
	type CalcTex struct {
		amount, expected float64
	}
	table := []CalcTex{
		{amount: 500, expected: 5},
		{amount: 1000, expected: 10},
		{amount: 1500, expected: 10},
	}

	for _, item := range table {
		result := CalculateTax(item.amount)
		if result != item.expected {
			t.Errorf("Expected %f but got %f", item.expected, result)
		}
	}
}

func BenchmarkCalculateTax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalculateTax(500)
	}
}

func BenchmarkCalculateTax2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalculateTax2(500)
	}
}

func FuzzCalculateTaxErro(f *testing.F) {
	seed := []float64{-1, -2, 500, 1000, 1501}
	for _, amount := range seed {
		f.Add(amount)
	}
	f.Fuzz(func(t *testing.T, amount float64) {
		result := CalculateTaxErrado(amount)
		if amount <= 0 && result != 0 {
			t.Errorf("Received %f but expected 0.0", result)
		}
		if amount >= 20000 && result != 20 {
			t.Errorf("Received %f but expected 20.0", result)
		}
	})
}
