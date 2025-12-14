package main

import "fmt"

func main() {

	//ch := make(chan int, 1)
	//go subscribe(ch)
	//publish(ch)
	//time.Sleep(time.Second)

	ch := make(chan string)
	go receiveOnlyChannel("Hello", ch)
	sendOnlyChannel(ch)
}

func receiveOnlyChannel(nome string, hello chan<- string) {
	hello <- nome
}

func sendOnlyChannel(data <-chan string) {
	fmt.Println(<-data)
}
