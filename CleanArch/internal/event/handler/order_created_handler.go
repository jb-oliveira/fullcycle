package handler

import (
	"encoding/json"
	"sync"

	"github.com/jb-oliveira/fullcycle/CleanArch/pkg/events"
	"github.com/streadway/amqp"
)

type orderCreatedHandlerRabbitMQ struct {
	RabbitMQChannel *amqp.Channel
}

func NewOrderCreatedHandlerRabbitMQ(rabbitMQChannel *amqp.Channel) events.EventHandlerInterface {
	return &orderCreatedHandlerRabbitMQ{RabbitMQChannel: rabbitMQChannel}
}

func (h *orderCreatedHandlerRabbitMQ) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	orderJSON, err := json.Marshal(event.GetPayload())
	if err != nil {
		println("Error marshaling order:", err.Error())
		return
	}
	println("Order created: ", string(orderJSON))
	msgRabbitMQ := amqp.Publishing{
		ContentType: "application/json",
		Body:        orderJSON,
	}
	err = h.RabbitMQChannel.Publish(
		"amq.direct", // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		msgRabbitMQ,
	)
	if err != nil {
		println("Error publishing to RabbitMQ:", err.Error())
	}
}
