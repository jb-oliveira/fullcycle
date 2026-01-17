package main

import (
	"fmt"
	"os"
)

func main() {
	i := 0

	for i < 10 {
		f, err := os.Create(fmt.Sprint("./tmp/example", i, ".txt"))
		if err != nil {
			panic(err)
		}
		f.WriteString("Hello World")
		f.Close()
		i++
	}
}
