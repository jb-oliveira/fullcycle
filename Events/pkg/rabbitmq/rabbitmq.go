package rabbitmq

import (
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// I asked Kiro to generate with reconnection, this is the connection manager to use see the example_usage.go file
// the course ones are down below look for CURSO...

// The Kiro one works like this...
// cm := rabbitmq.NewConnectionManager("", "", "")
// err := cm.Connect()
// messageChannel := make(chan amqp.Delivery)
// go cm.ConsumeWithReconnect(messageChannel)

// ConnectionManager manages RabbitMQ connection with automatic reconnection
type ConnectionManager struct {
	url           string
	conn          *amqp.Connection
	channel       *amqp.Channel
	mutex         sync.RWMutex
	reconnecting  bool
	reconnectChan chan bool
	closeChan     chan *amqp.Error
	queueName     string
	consumerTag   string
	exchangeName  string
	exchangeType  string
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager(url, queueName, consumerTag string) *ConnectionManager {
	if url == "" {
		url = "amqp://admin:password@localhost:5672/"
	}
	if queueName == "" {
		queueName = "minha-fila"
	}
	if consumerTag == "" {
		consumerTag = "go-consumer"
	}

	return &ConnectionManager{
		url:           url,
		reconnectChan: make(chan bool),
		queueName:     queueName,
		consumerTag:   consumerTag,
		exchangeName:  "", // Default to direct queue publishing
		exchangeType:  "direct",
	}
}

// NewConnectionManagerWithExchange creates a new connection manager with exchange configuration
func NewConnectionManagerWithExchange(url, queueName, consumerTag, exchangeName, exchangeType string) *ConnectionManager {
	cm := NewConnectionManager(url, queueName, consumerTag)
	cm.exchangeName = exchangeName
	if exchangeType != "" {
		cm.exchangeType = exchangeType
	}
	return cm
}

// Connect establishes connection to RabbitMQ
func (cm *ConnectionManager) Connect() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	var err error
	cm.conn, err = amqp.Dial(cm.url)
	if err != nil {
		return err
	}

	cm.channel, err = cm.conn.Channel()
	if err != nil {
		cm.conn.Close()
		return err
	}

	// Declare exchange if specified
	if cm.exchangeName != "" {
		err = cm.channel.ExchangeDeclare(
			cm.exchangeName, // name
			cm.exchangeType, // type
			true,            // durable
			false,           // auto-deleted
			false,           // internal
			false,           // no-wait
			nil,             // arguments
		)
		if err != nil {
			cm.channel.Close()
			cm.conn.Close()
			return err
		}
		log.Printf("Declared exchange: %s (type: %s)", cm.exchangeName, cm.exchangeType)
	}

	// Listen for connection close events
	cm.closeChan = make(chan *amqp.Error)
	cm.conn.NotifyClose(cm.closeChan)

	// Start monitoring connection in a separate goroutine
	go cm.monitorConnection()

	log.Println("Connected to RabbitMQ")
	return nil
}

// monitorConnection monitors the connection and triggers reconnection if needed
func (cm *ConnectionManager) monitorConnection() {
	for {
		select {
		case err := <-cm.closeChan:
			if err != nil {
				log.Printf("Connection lost: %v", err)
				cm.triggerReconnect()
			}
			return
		}
	}
}

// triggerReconnect initiates the reconnection process
func (cm *ConnectionManager) triggerReconnect() {
	cm.mutex.Lock()
	if cm.reconnecting {
		cm.mutex.Unlock()
		return
	}
	cm.reconnecting = true
	cm.mutex.Unlock()

	go cm.reconnect()
}

// reconnect attempts to reconnect with exponential backoff
func (cm *ConnectionManager) reconnect() {
	defer func() {
		cm.mutex.Lock()
		cm.reconnecting = false
		cm.mutex.Unlock()
	}()

	backoff := time.Second
	maxBackoff := 30 * time.Second

	for {
		log.Println("Attempting to reconnect to RabbitMQ...")

		err := cm.Connect()
		if err == nil {
			log.Println("Successfully reconnected to RabbitMQ")
			// Notify that reconnection is complete
			select {
			case cm.reconnectChan <- true:
			default:
			}
			return
		}

		log.Printf("Reconnection failed: %v. Retrying in %v", err, backoff)
		time.Sleep(backoff)

		// Exponential backoff with jitter
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}
}

// GetChannel returns the current channel (thread-safe)
func (cm *ConnectionManager) GetChannel() (*amqp.Channel, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	if cm.channel == nil {
		return nil, amqp.ErrClosed
	}
	return cm.channel, nil
}

// IsConnected checks if the connection is active
func (cm *ConnectionManager) IsConnected() bool {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	return cm.conn != nil && !cm.conn.IsClosed() && cm.channel != nil
}

// Close closes the connection and channel
func (cm *ConnectionManager) Close() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	var err error
	if cm.channel != nil {
		err = cm.channel.Close()
		cm.channel = nil
	}
	if cm.conn != nil {
		if connErr := cm.conn.Close(); connErr != nil && err == nil {
			err = connErr
		}
		cm.conn = nil
	}
	return err
}

// ConsumeWithReconnect consumes messages with automatic reconnection
func (cm *ConnectionManager) ConsumeWithReconnect(out chan<- amqp.Delivery) error {
	for {
		if !cm.IsConnected() {
			log.Println("Not connected, waiting for connection...")
			<-cm.reconnectChan
		}

		channel, err := cm.GetChannel()
		if err != nil {
			log.Printf("Failed to get channel: %v", err)
			time.Sleep(time.Second)
			continue
		}

		msgs, err := channel.Consume(
			cm.queueName,   // queue
			cm.consumerTag, // consumer
			false,          // auto-ack
			false,          // exclusive
			false,          // no-local
			false,          // no-wait
			nil,            // args
		)
		if err != nil {
			log.Printf("Failed to start consuming: %v", err)
			time.Sleep(time.Second)
			continue
		}

		log.Println("Started consuming messages")

		// Consume messages until connection is lost
		for msg := range msgs {
			select {
			case out <- msg:
			case <-cm.reconnectChan:
				// Connection was lost, break out of loop to reconnect
				log.Println("Connection lost during consumption, will reconnect...")
				goto reconnect
			}
		}

	reconnect:
		log.Println("Message channel closed, waiting for reconnection...")
	}
}

// Publish publishes a message with automatic reconnection
func (cm *ConnectionManager) Publish(body string) error {
	return cm.PublishToQueue(cm.queueName, body)
}

// PublishToQueue publishes a message to a specific queue with automatic reconnection
func (cm *ConnectionManager) PublishToQueue(queueName, body string) error {
	return cm.PublishWithOptions(queueName, body, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
}

// PublishWithOptions publishes a message with custom publishing options and automatic reconnection
func (cm *ConnectionManager) PublishWithOptions(queueName, routingKey string, msg amqp.Publishing) error {
	maxRetries := 3
	retryDelay := time.Second
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if !cm.IsConnected() {
			log.Println("Not connected for publishing, waiting for connection...")
			select {
			case <-cm.reconnectChan:
				// Connection restored, continue
			case <-time.After(5 * time.Second):
				// Timeout waiting for reconnection
				if attempt == maxRetries-1 {
					return amqp.ErrClosed
				}
				continue
			}
		}

		channel, err := cm.GetChannel()
		if err != nil {
			log.Printf("Failed to get channel for publishing (attempt %d/%d): %v", attempt+1, maxRetries, err)
			lastErr = err
			if attempt < maxRetries-1 {
				time.Sleep(retryDelay)
				retryDelay *= 2 // Exponential backoff
				continue
			}
			return err
		}

		err = channel.Publish(
			cm.exchangeName, // exchange (use configured exchange or empty for direct queue)
			routingKey,      // routing key
			false,           // mandatory
			false,           // immediate
			msg,             // message
		)

		if err == nil {
			log.Printf("Message published successfully to queue: %s", queueName)
			return nil
		}

		log.Printf("Failed to publish message (attempt %d/%d): %v", attempt+1, maxRetries, err)
		lastErr = err

		// If it's a connection error, trigger reconnection
		if err == amqp.ErrClosed || err.Error() == "channel/connection is not open" {
			cm.triggerReconnect()
		}

		if attempt < maxRetries-1 {
			time.Sleep(retryDelay)
			retryDelay *= 2 // Exponential backoff
		}
	}

	return lastErr
}

// DeclareExchange declares an exchange with the specified parameters
func (cm *ConnectionManager) DeclareExchange(name, exchangeType string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	channel, err := cm.GetChannel()
	if err != nil {
		return err
	}

	return channel.ExchangeDeclare(name, exchangeType, durable, autoDelete, internal, noWait, args)
}

// PublishToExchange publishes a message to an exchange with routing key
func (cm *ConnectionManager) PublishToExchange(exchangeName, routingKey, body string) error {
	return cm.PublishToExchangeWithOptions(exchangeName, routingKey, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
}

// PublishToExchangeWithOptions publishes a message to an exchange with custom options and automatic reconnection
func (cm *ConnectionManager) PublishToExchangeWithOptions(exchangeName, routingKey string, msg amqp.Publishing) error {
	maxRetries := 3
	retryDelay := time.Second
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if !cm.IsConnected() {
			log.Println("Not connected for publishing to exchange, waiting for connection...")
			select {
			case <-cm.reconnectChan:
				// Connection restored, continue
			case <-time.After(5 * time.Second):
				// Timeout waiting for reconnection
				if attempt == maxRetries-1 {
					return amqp.ErrClosed
				}
				continue
			}
		}

		channel, err := cm.GetChannel()
		if err != nil {
			log.Printf("Failed to get channel for exchange publishing (attempt %d/%d): %v", attempt+1, maxRetries, err)
			lastErr = err
			if attempt < maxRetries-1 {
				time.Sleep(retryDelay)
				retryDelay *= 2 // Exponential backoff
				continue
			}
			return err
		}

		err = channel.Publish(
			exchangeName, // exchange
			routingKey,   // routing key
			false,        // mandatory
			false,        // immediate
			msg,          // message
		)

		if err == nil {
			log.Printf("Message published successfully to exchange: %s with routing key: %s", exchangeName, routingKey)
			return nil
		}

		log.Printf("Failed to publish message to exchange (attempt %d/%d): %v", attempt+1, maxRetries, err)
		lastErr = err

		// If it's a connection error, trigger reconnection
		if err == amqp.ErrClosed || err.Error() == "channel/connection is not open" {
			cm.triggerReconnect()
		}

		if attempt < maxRetries-1 {
			time.Sleep(retryDelay)
			retryDelay *= 2 // Exponential backoff
		}
	}

	return lastErr
}

// ConsumeFromExchange consumes messages from a queue bound to an exchange with automatic reconnection
func (cm *ConnectionManager) ConsumeFromExchange(exchangeName, exchangeType, queueName, routingKey string, out chan<- amqp.Delivery) error {
	for {
		if !cm.IsConnected() {
			log.Println("Not connected, waiting for connection...")
			<-cm.reconnectChan
		}

		channel, err := cm.GetChannel()
		if err != nil {
			log.Printf("Failed to get channel: %v", err)
			time.Sleep(time.Second)
			continue
		}

		// Declare exchange
		err = channel.ExchangeDeclare(
			exchangeName, // name
			exchangeType, // type
			true,         // durable
			false,        // auto-deleted
			false,        // internal
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			log.Printf("Failed to declare exchange: %v", err)
			time.Sleep(time.Second)
			continue
		}

		// Declare queue
		queue, err := channel.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			log.Printf("Failed to declare queue: %v", err)
			time.Sleep(time.Second)
			continue
		}

		// Bind queue to exchange
		err = channel.QueueBind(
			queue.Name,   // queue name
			routingKey,   // routing key
			exchangeName, // exchange
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			log.Printf("Failed to bind queue to exchange: %v", err)
			time.Sleep(time.Second)
			continue
		}

		msgs, err := channel.Consume(
			queue.Name,     // queue
			cm.consumerTag, // consumer
			false,          // auto-ack
			false,          // exclusive
			false,          // no-local
			false,          // no-wait
			nil,            // args
		)
		if err != nil {
			log.Printf("Failed to start consuming from exchange: %v", err)
			time.Sleep(time.Second)
			continue
		}

		log.Printf("Started consuming from exchange: %s, queue: %s, routing key: %s", exchangeName, queue.Name, routingKey)

		// Consume messages until connection is lost
		for msg := range msgs {
			select {
			case out <- msg:
			case <-cm.reconnectChan:
				// Connection was lost, break out of loop to reconnect
				log.Println("Connection lost during exchange consumption, will reconnect...")
				goto reconnect
			}
		}

	reconnect:
		log.Println("Exchange message channel closed, waiting for reconnection...")
	}
}

// SetExchange updates the default exchange configuration
func (cm *ConnectionManager) SetExchange(exchangeName, exchangeType string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.exchangeName = exchangeName
	cm.exchangeType = exchangeType
}

// GetExchange returns the current exchange configuration
func (cm *ConnectionManager) GetExchange() (string, string) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.exchangeName, cm.exchangeType
}

// COURSE: Legacy functions for backward compatibility
func OpenChannel() (*amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://admin:password@localhost:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return ch, nil
}

func Consume(ch *amqp.Channel, queue string, out chan<- amqp.Delivery) error {
	msgs, err := ch.Consume(
		queue,         // queue (queue)
		"go-consumer", // consumer (application name)
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		return err
	}

	// will read the message from msgs in msg and throw it into the GO Channel
	for msg := range msgs {
		out <- msg
	}

	return nil
}

func Publish(ch *amqp.Channel, exchange string, body string) error {
	return ch.Publish(
		exchange, // exchange
		"",       // routing key (nome da fila) deixa em branco a excengae vai fazer o trabalho
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
}
