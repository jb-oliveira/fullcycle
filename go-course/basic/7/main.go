package main

func sumInt(m map[string]int) int {
	var sum int
	for _, v := range m {
		sum += v
	}
	return sum
}

func sumFloat(m map[string]float64) float64 {
	var sum float64
	for _, v := range m {
		sum += v
	}
	return sum
}

func sum[T int | float64](m map[string]T) T {
	var sum T
	for _, v := range m {
		sum += v
	}
	return sum
}

type MyNumber int

type Number interface {
	~int | ~float64
	// the tilde serves to identify any type of integer as in the case of myNumber
}

func sumConstraint[T Number](m map[string]T) T {
	var sum T
	for _, v := range m {
		sum += v
	}
	return sum
}

func compare[T comparable](a T, b T) bool {
	if a == b {
		return true
	}
	return false
}

func main() {

	m := map[string]int{"A": 100, "B": 220, "C": 450}
	mF := map[string]float64{"A": 100.0, "B": 220.0, "C": 450.0}
	mM := map[string]MyNumber{"A": 100, "B": 220, "C": 450}
	println(sumInt(m))
	println(sumFloat(mF))
	println(sum(m))
	println(sum(mF))
	println(sumConstraint(m))
	println(sumConstraint(mF))
	println(sumConstraint(mM))

	//
	println(compare(1, 1.0))
}
