# RabbitMQ Connection Manager with Automatic Reconnection

This package provides a robust RabbitMQ connection manager that automatically handles connection failures and reconnections with exponential backoff.

## Features

- **Automatic Reconnection**: Automatically reconnects when connection is lost
- **Exponential Backoff**: Uses exponential backoff with jitter to avoid overwhelming the server
- **Thread-Safe**: All operations are thread-safe using mutexes
- **Connection Monitoring**: Continuously monitors connection health
- **Graceful Error Handling**: Proper error handling and logging
- **Backward Compatibility**: Maintains compatibility with existing code

## Usage

### Basic Usage with Automatic Reconnection

```go
package main

import (
    "log"
    "your-project/pkg/rabbitmq"
    amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
    // Create connection manager
    cm := rabbitmq.NewConnectionManager(
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

    // Create channel for messages
    messageChannel := make(chan amqp.Delivery)

    // Start consuming with automatic reconnection
    go func() {
        err := cm.ConsumeWithReconnect(messageChannel)
        if err != nil {
            log.Printf("Consumer error: %v", err)
        }
    }()

    // Process messages
    for msg := range messageChannel {
        log.Printf("Received: %s", string(msg.Body))
        
        // Process your message here
        // ...
        
        // Acknowledge the message
        msg.Ack(false)
    }
}
```

### Using Default Values

```go
// Uses default values:
// URL: "amqp://admin:password@localhost:5672/"
// Queue: "minha-fila"
// Consumer: "go-consumer"
cm := rabbitmq.NewConnectionManager("", "", "")
```

### Manual Connection Management

```go
cm := rabbitmq.NewConnectionManager("", "", "")

// Check if connected
if cm.IsConnected() {
    log.Println("Connected to RabbitMQ")
}

// Get channel for manual operations
channel, err := cm.GetChannel()
if err != nil {
    log.Printf("Failed to get channel: %v", err)
}

// Close connection
err = cm.Close()
if err != nil {
    log.Printf("Failed to close: %v", err)
}
```

## Connection Manager Methods

### `NewConnectionManager(url, queueName, consumerTag string) *ConnectionManager`
Creates a new connection manager. Empty strings use default values.

### `NewConnectionManagerWithExchange(url, queueName, consumerTag, exchangeName, exchangeType string) *ConnectionManager`
Creates a new connection manager with exchange configuration for advanced messaging patterns.

### `Connect() error`
Establishes initial connection to RabbitMQ, starts monitoring, and declares exchange if configured.

### `ConsumeWithReconnect(out chan<- amqp.Delivery) error`
Starts consuming messages with automatic reconnection on connection loss.

### `Publish(body string) error`
Publishes a simple text message to the default queue with automatic reconnection.

### `PublishToQueue(queueName, body string) error`
Publishes a simple text message to a specific queue with automatic reconnection.

### `PublishWithOptions(queueName, routingKey string, msg amqp.Publishing) error`
Publishes a message with custom publishing options and automatic reconnection.

### `PublishToExchange(exchangeName, routingKey, body string) error`
Publishes a simple text message to an exchange with automatic reconnection.

### `PublishToExchangeWithOptions(exchangeName, routingKey string, msg amqp.Publishing) error`
Publishes a message to an exchange with custom options and automatic reconnection.

### `ConsumeFromExchange(exchangeName, exchangeType, queueName, routingKey string, out chan<- amqp.Delivery) error`
Consumes messages from a queue bound to an exchange with automatic reconnection.

### `DeclareExchange(name, exchangeType string, durable, autoDelete, internal, noWait bool, args amqp.Table) error`
Declares an exchange with specified parameters.

### `SetExchange(exchangeName, exchangeType string)`
Updates the default exchange configuration.

### `GetExchange() (string, string)`
Returns the current exchange name and type.

### `GetChannel() (*amqp.Channel, error)`
Returns the current channel (thread-safe).

### `IsConnected() bool`
Checks if the connection is currently active.

### `Close() error`
Closes the connection and channel gracefully.

## Reconnection Behavior

- **Initial Backoff**: 1 second
- **Maximum Backoff**: 30 seconds
- **Backoff Strategy**: Exponential with doubling
- **Connection Monitoring**: Automatic detection of connection loss
- **Recovery**: Automatic reconnection with full consumer restart

## Error Handling

The connection manager handles various error scenarios:

- Network connectivity issues
- RabbitMQ server restarts
- Authentication failures (with retry)
- Channel errors
- Consumer cancellation

## Backward Compatibility

The original functions are still available for existing code:

```go
// Legacy functions (no reconnection)
channel, err := rabbitmq.OpenChannel()
err = rabbitmq.Consume(channel, messageChannel)
err = rabbitmq.Publish(channel, "Hello World!")
```

## Publishing Messages

The ConnectionManager provides several methods for publishing messages:

### Simple Publishing
```go
cm := rabbitmq.NewConnectionManager("", "", "")
err := cm.Connect()
defer cm.Close()

// Publish to default queue
err = cm.Publish("Hello World!")
```

### Publishing to Specific Queue
```go
// Publish to a specific queue
err = cm.PublishToQueue("my-queue", "Hello specific queue!")
```

### Advanced Publishing with Custom Options
```go
// Publish with custom message properties
customMessage := amqp.Publishing{
    ContentType:  "application/json",
    Body:         []byte(`{"message": "Hello JSON!"}`),
    DeliveryMode: amqp.Persistent, // Make message persistent
    Priority:     5,               // Set priority
    Headers: amqp.Table{
        "source": "my-app",
        "type":   "notification",
    },
}

err = cm.PublishWithOptions("json-queue", "json-queue", customMessage)
```

## Exchange-Based Messaging

For advanced messaging patterns, use exchanges instead of direct queue publishing:

### Creating Connection Manager with Exchange
```go
// Create with exchange configuration
cm := rabbitmq.NewConnectionManagerWithExchange(
    "amqp://admin:password@localhost:5672/",
    "",                // queue name (optional for exchange publishing)
    "my-publisher",    // consumer tag
    "events-exchange", // exchange name
    "topic",          // exchange type (direct, topic, fanout, headers)
)
```

### Publishing to Exchanges
```go
// Simple exchange publishing
err = cm.PublishToExchange("events-exchange", "user.created", "User created event")

// Advanced exchange publishing
eventMessage := amqp.Publishing{
    ContentType: "application/json",
    Body:        []byte(`{"userId": 123, "action": "created"}`),
    Headers: amqp.Table{
        "source": "user-service",
    },
}
err = cm.PublishToExchangeWithOptions("events-exchange", "user.created", eventMessage)
```

### Consuming from Exchanges
```go
messageChannel := make(chan amqp.Delivery)

// Consume from exchange with routing pattern
go cm.ConsumeFromExchange(
    "events-exchange", // exchange name
    "topic",          // exchange type
    "user-events",    // queue name
    "user.*",         // routing key pattern
    messageChannel,
)

// Process messages
for msg := range messageChannel {
    log.Printf("Event: %s - %s", msg.RoutingKey, string(msg.Body))
    msg.Ack(false)
}
```

### Exchange Types and Patterns

**Topic Exchange** - Pattern-based routing:
```go
// Publisher
cm.PublishToExchange("events", "user.profile.updated", "message")
cm.PublishToExchange("events", "order.payment.completed", "message")

// Consumer patterns
"user.*"           // All user events
"*.completed"      // All completion events  
"user.profile.#"   // All user profile events
"#"               // All events
```

**Direct Exchange** - Exact routing key match:
```go
cm.PublishToExchange("direct-exchange", "orders", "Order message")
// Only queues bound with "orders" routing key receive this
```

**Fanout Exchange** - Broadcast to all bound queues:
```go
cm.PublishToExchange("broadcast-exchange", "", "Broadcast message")
// All bound queues receive this message
```

See [EXCHANGE_GUIDE.md](EXCHANGE_GUIDE.md) for detailed exchange usage examples.

## Testing

Run tests with:

```bash
# Unit tests only
go test ./pkg/rabbitmq -short

# Integration tests (requires running RabbitMQ)
go test ./pkg/rabbitmq
```

## Dependencies

- `github.com/rabbitmq/amqp091-go` - RabbitMQ client library

## Configuration

The connection manager uses these default values:
- **URL**: `amqp://admin:password@localhost:5672/`
- **Queue**: `minha-fila`
- **Consumer Tag**: `go-consumer`

All values can be customized when creating the connection manager.