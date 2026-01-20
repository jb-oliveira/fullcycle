package main

import "fmt"

func main() {
	var slice = []int{1, 2, 3, 4, 5, 6}
	fmt.Printf("len=%d cap=%d %v\n", len(slice), cap(slice), slice)
	fmt.Printf("len=%d cap=%d %v\n", len(slice[:0]), cap(slice[:0]), slice[:0])
	newSlice := slice[:2]
	fmt.Printf("len=%d cap=%d %v\n", len(newSlice), cap(newSlice), newSlice)
	newSlice = slice[2:]
	fmt.Printf("len=%d cap=%d %v\n", len(newSlice), cap(newSlice), newSlice)

	newSlice = append(slice, 7)
	fmt.Printf("len=%d cap=%d %v\n", len(newSlice), cap(newSlice), newSlice)
	// for i, v := range slice {
	// 	fmt.Printf("Index %d Value %v\n", i, v)
	// }

	for v := range slice {
		fmt.Printf("O Valor é: %d\n", v)
	}
}
