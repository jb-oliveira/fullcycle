package main

import (
	"fmt"

	"github.com/jb-oliveira/fullcycle/tree/main/11-EVENTOS/pkg/rabbitmq"
	ampq "github.com/rabbitmq/amqp091-go"
)

func main() {
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msgs := make(chan ampq.Delivery)
	go rabbitmq.Consume(ch, "minha-fila", msgs)
	for msg := range msgs {
		fmt.Println(string(msg.Body))
		msg.Ack(false) // Requeu = false, pra não colocar ela de volta na fila
	}
}
