package main

func somaInt(m map[string]int) int {
	var soma int
	for _, v := range m {
		soma += v
	}
	return soma
}

func somaFloat(m map[string]float64) float64 {
	var soma float64
	for _, v := range m {
		soma += v
	}
	return soma
}

func soma[T int | float64](m map[string]T) T {
	var soma T
	for _, v := range m {
		soma += v
	}
	return soma
}

type MyNumber int

type Number interface {
	~int | ~float64
	// the tilde serves to identify any type of integer as in the case of myNumber
}

func somaConstraint[T Number](m map[string]T) T {
	var soma T
	for _, v := range m {
		soma += v
	}
	return soma
}

func compara[T comparable](a T, b T) bool {
	if a == b {
		return true
	}
	return false
}

func main() {

	m := map[string]int{"A": 100, "B": 220, "C": 450}
	mF := map[string]float64{"A": 100.0, "B": 220.0, "C": 450.0}
	mM := map[string]MyNumber{"A": 100, "B": 220, "C": 450}
	println(somaInt(m))
	println(somaFloat(mF))
	println(soma(m))
	println(soma(mF))
	println(somaConstraint(m))
	println(somaConstraint(mF))
	println(somaConstraint(mM))

	//
	println(compara(1, 1.0))
}
