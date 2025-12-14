package main

import "time"

func worker(workerId int, data chan int) {
	for x := range data {
		println("Worker", workerId, "received", x)
		time.Sleep(time.Second)
	}
}

func main() {
	data := make(chan int)

	numWorkers := 100000
	for i := 0; i < numWorkers; i++ {
		go worker(i, data)
	}

	for i := 0; i < 1000000; i++ {
		data <- i
	}
}
