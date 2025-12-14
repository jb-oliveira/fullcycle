package main

import (
	"fmt"
	"sync"
	"time"
)

func task(name string, times int, waitGroup *sync.WaitGroup) {
	for i := 0; i < times; i++ {
		fmt.Printf("%s %d\n", name, i)
		time.Sleep(time.Millisecond * 250)
		waitGroup.Done()
	}
}

func main() {
	// go task("task 1", 10)
	// go task("task 2", 10)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(30)

	go task("task 1", 10, &waitGroup)
	go task("task 2", 10, &waitGroup)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("task 3 %d\n", i)
			time.Sleep(time.Millisecond * 250)
			waitGroup.Done()
		}
	}()

	waitGroup.Wait()
}
