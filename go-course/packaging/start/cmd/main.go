package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jb-oliveira/fullcycle/go-course/packaging/start/math"
)

func main() {
	m := math.Math{A: 1, B: 2}
	fmt.Println(m.Add())

	fmt.Println(uuid.NewV7())

}
