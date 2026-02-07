package handler

import (
	"encoding/json"
	"sync"

	"github.com/jb-oliveira/fullcycle/CleanArch/internal/event"
	"github.com/streadway/amqp"
)

type orderCreatedHandlerRabbitMQ struct {
	RabbitMQChannel *amqp.Channel
}

func NewOrderCreatedHandlerRabbitMQ(rabbitMQChannel *amqp.Channel) event.EventHandler {
	return &orderCreatedHandlerRabbitMQ{RabbitMQChannel: rabbitMQChannel}
}

func (h *orderCreatedHandlerRabbitMQ) Handle(event event.Event, wg *sync.WaitGroup) error {
	defer wg.Done()
	orderJSON, err := json.Marshal(event.GetPayload())
	println("Order created: ", string(orderJSON))
	if err != nil {
		return err
	}
	msgRabbitMQ := amqp.Publishing{
		ContentType: "application/json",
		Body:        orderJSON,
	}
	return h.RabbitMQChannel.Publish(
		"amq.direct", // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		msgRabbitMQ,
	)
}
