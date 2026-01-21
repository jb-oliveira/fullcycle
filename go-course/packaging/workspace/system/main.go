package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jb-oliveira/fullcycle/tree/main/curso-go/packaging/workspace/math"
)

func main() {
	m := math.Math{A: 1, B: 2}
	fmt.Println(m.Add())

	fmt.Println(uuid.NewV7())

}
