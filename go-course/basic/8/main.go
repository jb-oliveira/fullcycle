package main

import (
	"fmt"

	"github.com/google/uuid"
	mat "github.com/jb-oliveira/fullcycle/go-course/basic/8/math"
)

func main() {

	s := mat.Sum(1, 2)
	println(s)
	val, err := uuid.NewV7()
	if err != nil {
		println(err)

	} else {
		fmt.Printf("Valor = %v", val)
	}

	for i := 0; i < 10; i++ {
		println(i)
	}

	slc := []string{"A", "B", "C"}
	newFunction(slc...)

	i := 0
	for i < 10 {
		println(i)
		i++
	}

	// println(val)
}

func newFunction(slc ...string) {
	for k, v := range slc {
		println(k, v)
	}
}
