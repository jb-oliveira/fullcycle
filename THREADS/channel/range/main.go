package main

import "time"

func main() {

	ch := make(chan int, 1)
	go publish(ch)
	subscribe(ch)

}

func subscribe(ch chan int) {
	for i := range ch {
		println(i)
	}
}

func publish(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- i
		time.Sleep(time.Second)
	}
	close(ch)
}
