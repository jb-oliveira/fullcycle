package main

import "fmt"

func main() {

	fmt.Println("First Line")
	defer fmt.Println("Second Line")
	defer fmt.Println("Third Line")
	fmt.Println("Quarta Linha")
}
