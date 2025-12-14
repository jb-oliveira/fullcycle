package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Message struct {
	Id      int64
	Payload string
}

func main() {
	ch1 := make(chan Message)
	ch2 := make(chan Message)
	var id int64 = 0

	go func() {
		for {
			atomic.AddInt64(&id, 1)
			time.Sleep(time.Second)
			ch1 <- Message{
				Id:      int64(id),
				Payload: "Hello",
			}
		}

	}()
	go func() {
		for {
			atomic.AddInt64(&id, 1)
			time.Sleep(2 * time.Second)
			ch2 <- Message{
				Id:      id,
				Payload: "World",
			}
		}

	}()

	for {
		select {
		case msg := <-ch1:
			fmt.Printf("Received From RabbitMq ID: %d Message: %s\n", msg.Id, msg.Payload)
		case msg := <-ch2:
			fmt.Printf("Received From Kafka ID: %d Message: %s\n", msg.Id, msg.Payload)
		case <-time.After(1 * time.Second):
			println("timeout")
		}
	}

}
