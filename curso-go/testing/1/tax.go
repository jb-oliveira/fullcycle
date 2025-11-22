package tax

import "time"

func CalculateTax(amount float64) float64 {
	if amount <= 0 {
		return 0
	}
	if amount >= 1000 {
		return 10
	}
	return 5
}

func CalculateTaxErrado(amount float64) float64 {
	if amount <= 0 {
		return 0
	}
	if amount >= 1000 {
		return 10
	}
	if amount >= 20000 {
		return 20
	}
	return 5
}

func CalculateTax2(amount float64) float64 {
	time.Sleep(time.Millisecond)
	if amount == 0 {
		return 0
	}
	if amount >= 1000 {
		return 10
	}
	return 5
}
