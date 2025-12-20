package main

import (
	"fmt"
	"time"

	"github.com/jb-oliveira/fullcycle/tree/main/11-EVENTOS/pkg/rabbitmq"
)

func main() {
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	for i := 1; i <= 100; i++ {
		rabbitmq.Publish(ch, fmt.Sprintf("Mensagem: %d", i))
		time.Sleep(200 * time.Millisecond)
	}

}
