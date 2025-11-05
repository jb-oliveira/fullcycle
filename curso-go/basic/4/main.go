package main

import (
	"errors"
	"fmt"
)

func main() {

	fmt.Println(sum(1, 2))
	value, err := sum2(1, 2)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(value)
	value, err = sum2(51, 2)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(value)
	}

	total := func() int {
		return sumVariadicas(2345, 345354, 345, 435534, 34534, 53, 54435, 34) * 2
	}()
	fmt.Println(total)

	arr := []int{
		2, 3, 4, 5}
	fmt.Println(sumVariadicas(arr...))

}

func sumVariadicas(numeros ...int) int {
	total := 0
	for _, num := range numeros {
		total += num
	}
	return total
}

func sum(a, b int) (int, bool) {
	if a+b > 50 {
		return a + b, true
	}
	return a + b, false
}

func sum2(a, b int) (int, error) {
	if a+b > 50 {
		return 0, errors.New("Soma maior que 50")
	}
	return a + b, nil
}
