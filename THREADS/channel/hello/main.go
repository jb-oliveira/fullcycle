package main

func main() {

	canal := make(chan string)

	go func() {
		canal <- "Hello World!"
	}()

	msg := <-canal
	println(msg)
}
