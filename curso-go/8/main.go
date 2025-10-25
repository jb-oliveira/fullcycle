package main

import (
	"fmt"

	"github.com/google/uuid"
	mat "github.com/jb-oliveira/fullcycle/tree/main/curso-go/8/matematica"
)

func main() {

	s := mat.Soma(1, 2)
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
	for k, v := range slc {
		println(k, v)
	}

	i := 0
	for i < 10 {
		println(i)
		i++
	}

	// println(val)
}
