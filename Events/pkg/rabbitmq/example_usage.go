package rabbitmq

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ExampleUsage demonstrates how to use the ConnectionManager with reconnection
func ExampleUsage() {
	// Create a connection manager
	cm := NewConnectionManager(
		"amqp://admin:password@localhost:5672/", // RabbitMQ URL
		"minha-fila",                            // Queue name
		"go-consumer",                           // Consumer tag
	)

	// Connect initially
	err := cm.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cm.Close()

	// Create a channel to receive messages
	messageChannel := make(chan amqp.Delivery)

	// Start consuming with automatic reconnection
	go func() {
		err := cm.ConsumeWithReconnect(messageChannel)
		if err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()

	// Process messages
	for {
		select {
		case msg := <-messageChannel:
			log.Printf("Received message: %s", string(msg.Body))

			// Process the message here
			// ...

			// Acknowledge the message
			err := msg.Ack(false)
			if err != nil {
				log.Printf("Failed to ack message: %v", err)
			}

		case <-time.After(30 * time.Second):
			log.Println("No messages received in 30 seconds")
		}
	}
}

// ExampleWithManualReconnection shows how to handle reconnection manually
func ExampleWithManualReconnection() {
	cm := NewConnectionManager("", "", "") // Use defaults

	for {
		err := cm.Connect()
		if err != nil {
			log.Printf("Connection failed: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Connection successful, start consuming
		messageChannel := make(chan amqp.Delivery)

		go func() {
			err := cm.ConsumeWithReconnect(messageChannel)
			if err != nil {
				log.Printf("Consumer stopped: %v", err)
			}
		}()

		// Process messages until connection is lost
		for {
			if !cm.IsConnected() {
				log.Println("Connection lost, will reconnect...")
				break
			}

			select {
			case msg := <-messageChannel:
				log.Printf("Processing: %s", string(msg.Body))
				msg.Ack(false)
			case <-time.After(1 * time.Second):
				// Check connection status periodically
			}
		}

		cm.Close()
		time.Sleep(2 * time.Second) // Wait before reconnecting
	}
}

// ExamplePublisher demonstrates how to publish messages with reconnection
func ExamplePublisher() {
	// Create a connection manager
	cm := NewConnectionManager(
		"amqp://admin:password@localhost:5672/", // RabbitMQ URL
		"minha-fila",                            // Queue name
		"go-publisher",                          // Consumer tag (not used for publishing)
	)

	// Connect initially
	err := cm.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cm.Close()

	// Publish messages with automatic reconnection
	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("Message %d - %s", i, time.Now().Format(time.RFC3339))

		// Simple publish to default queue
		err := cm.Publish(message)
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
		} else {
			log.Printf("Published: %s", message)
		}

		time.Sleep(2 * time.Second)
	}
}

// ExampleAdvancedPublisher demonstrates advanced publishing options
func ExampleAdvancedPublisher() {
	cm := NewConnectionManager("", "", "")

	err := cm.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cm.Close()

	// Publish to specific queue
	err = cm.PublishToQueue("custom-queue", "Hello Custom Queue!")
	if err != nil {
		log.Printf("Failed to publish to custom queue: %v", err)
	}

	// Publish with custom options
	customMessage := amqp.Publishing{
		ContentType:  "application/json",
		Body:         []byte(`{"message": "Hello JSON!", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`),
		DeliveryMode: amqp.Persistent, // Make message persistent
		Priority:     5,               // Set priority
		Headers: amqp.Table{
			"source": "go-publisher",
			"type":   "example",
		},
	}

	err = cm.PublishWithOptions("json-queue", "json-queue", customMessage)
	if err != nil {
		log.Printf("Failed to publish JSON message: %v", err)
	}
}

// ExampleExchangePublisher demonstrates publishing to exchanges
func ExampleExchangePublisher() {
	// Create connection manager with exchange configuration
	cm := NewConnectionManagerWithExchange(
		"amqp://admin:password@localhost:5672/",
		"", // queue name not needed for exchange publishing
		"exchange-publisher",
		"events-exchange", // exchange name
		"topic",           // exchange type
	)

	err := cm.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cm.Close()

	// Publish to different routing keys
	routingKeys := []string{
		"user.created",
		"user.updated",
		"order.placed",
		"order.shipped",
	}

	for i, routingKey := range routingKeys {
		message := fmt.Sprintf("Event %d: %s at %s", i+1, routingKey, time.Now().Format(time.RFC3339))

		err := cm.PublishToExchange("events-exchange", routingKey, message)
		if err != nil {
			log.Printf("Failed to publish to exchange: %v", err)
		} else {
			log.Printf("Published to exchange with routing key %s: %s", routingKey, message)
		}

		time.Sleep(1 * time.Second)
	}
}

// ExampleExchangeConsumer demonstrates consuming from exchanges
func ExampleExchangeConsumer() {
	cm := NewConnectionManager("", "", "exchange-consumer")

	err := cm.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cm.Close()

	// Create channel for messages
	messageChannel := make(chan amqp.Delivery)

	// Start consuming from exchange with specific routing pattern
	go func() {
		err := cm.ConsumeFromExchange(
			"events-exchange", // exchange name
			"topic",           // exchange type
			"user-events",     // queue name
			"user.*",          // routing key pattern (topic exchange)
			messageChannel,
		)
		if err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()

	// Process messages
	for msg := range messageChannel {
		log.Printf("Received from exchange - Routing Key: %s, Message: %s",
			msg.RoutingKey, string(msg.Body))

		// Process the message
		// ...

		// Acknowledge the message
		err := msg.Ack(false)
		if err != nil {
			log.Printf("Failed to ack message: %v", err)
		}
	}
}

// ExampleMultipleExchangeTypes demonstrates different exchange types
func ExampleMultipleExchangeTypes() {
	cm := NewConnectionManager("", "", "multi-exchange-demo")

	err := cm.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cm.Close()

	// 1. Direct Exchange
	err = cm.DeclareExchange("direct-exchange", "direct", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare direct exchange: %v", err)
	}

	// 2. Topic Exchange
	err = cm.DeclareExchange("topic-exchange", "topic", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare topic exchange: %v", err)
	}

	// 3. Fanout Exchange
	err = cm.DeclareExchange("fanout-exchange", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare fanout exchange: %v", err)
	}

	// 4. Headers Exchange
	err = cm.DeclareExchange("headers-exchange", "headers", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare headers exchange: %v", err)
	}

	// Publish to different exchange types

	// Direct exchange - exact routing key match
	err = cm.PublishToExchange("direct-exchange", "orders", "Direct message to orders")
	if err != nil {
		log.Printf("Failed to publish to direct exchange: %v", err)
	}

	// Topic exchange - pattern matching
	err = cm.PublishToExchange("topic-exchange", "user.profile.updated", "Topic message")
	if err != nil {
		log.Printf("Failed to publish to topic exchange: %v", err)
	}

	// Fanout exchange - broadcast to all bound queues
	err = cm.PublishToExchange("fanout-exchange", "", "Fanout broadcast message")
	if err != nil {
		log.Printf("Failed to publish to fanout exchange: %v", err)
	}

	// Headers exchange - route based on headers
	headersMessage := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Headers-based routing message"),
		Headers: amqp.Table{
			"type":     "notification",
			"priority": "high",
		},
	}
	err = cm.PublishToExchangeWithOptions("headers-exchange", "", headersMessage)
	if err != nil {
		log.Printf("Failed to publish to headers exchange: %v", err)
	}

	log.Println("Published messages to all exchange types")
}
