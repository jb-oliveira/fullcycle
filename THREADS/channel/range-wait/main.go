package main

import "sync"

func main() {

	ch := make(chan int, 1)
	wg := sync.WaitGroup{}
	wg.Add(10)
	go publish(ch)
	go subscribe(ch, &wg)
	wg.Wait()
}

func subscribe(ch chan int, wg *sync.WaitGroup) {
	for i := range ch {
		println(i)
		wg.Done()
	}
}

func publish(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}
	close(ch)
}
