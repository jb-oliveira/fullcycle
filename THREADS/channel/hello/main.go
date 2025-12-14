package main

func main() {

	canal := make(chan string)

	go func() {
		canal <- "Olá Mundo!"
	}()

	msg := <-canal
	println(msg)
}
