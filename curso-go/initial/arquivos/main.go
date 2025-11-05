package main

import (
	"bufio"
	"fmt"
	"os"
)

const fileName = "Arquivo.txt"

func main() {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	tamanho, err := f.WriteString("Hello World!")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Tamanho do arquivo: %v\n", tamanho)
	f.Close()

	arquivo, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(arquivo))

	arquivo2, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer arquivo2.Close()

	reader := bufio.NewReader(arquivo2)
	buffer := make([]byte, 3)
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			break
		}
		fmt.Println(string(buffer[:n]))
	}

	err = os.Remove(fileName)
	if err != nil {
		panic(err)
	}
}
