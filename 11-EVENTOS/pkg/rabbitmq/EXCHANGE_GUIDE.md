# RabbitMQ Exchange Guide

This guide explains how to use exchanges with the ConnectionManager for more advanced messaging patterns.

## What are Exchanges?

Exchanges are message routing agents in RabbitMQ. Instead of publishing directly to queues, you publish to exchanges, which then route messages to queues based on routing rules.

## Exchange Types

### 1. Direct Exchange
Routes messages with exact routing key matches.

```go
cm := NewConnectionManagerWithExchange("", "", "", "direct-exchange", "direct")
cm.Connect()

// Publish with specific routing key
cm.PublishToExchange("direct-exchange", "orders", "Order message")
cm.PublishToExchange("direct-exchange", "users", "User message")
```

### 2. Topic Exchange
Routes messages based on routing key patterns using wildcards.

```go
cm := NewConnectionManagerWithExchange("", "", "", "topic-exchange", "topic")
cm.Connect()

// Publish with pattern-based routing keys
cm.PublishToExchange("topic-exchange", "user.created", "User created event")
cm.PublishToExchange("topic-exchange", "user.updated", "User updated event")
cm.PublishToExchange("topic-exchange", "order.placed", "Order placed event")

// Consumers can bind with patterns:
// "user.*" - matches all user events
// "*.created" - matches all creation events
// "#" - matches all messages
```

### 3. Fanout Exchange
Broadcasts messages to all bound queues (ignores routing key).

```go
cm := NewConnectionManagerWithExchange("", "", "", "fanout-exchange", "fanout")
cm.Connect()

// Routing key is ignored in fanout
cm.PublishToExchange("fanout-exchange", "", "Broadcast message")
```

### 4. Headers Exchange
Routes based on message headers instead of routing key.

```go
cm := NewConnectionManager("", "", "")
cm.Connect()

// Declare headers exchange
cm.DeclareExchange("headers-exchange", "headers", true, false, false, false, nil)

// Publish with headers
headersMsg := amqp.Publishing{
    Body: []byte("Headers message"),
    Headers: amqp.Table{
        "type": "notification",
        "priority": "high",
    },
}
cm.PublishToExchangeWithOptions("headers-exchange", "", headersMsg)
```

## Creating Connection Managers with Exchanges

### Method 1: With Exchange Configuration
```go
cm := NewConnectionManagerWithExchange(
    "amqp://admin:password@localhost:5672/", // URL
    "my-queue",                              // Queue name
    "my-consumer",                           // Consumer tag
    "my-exchange",                           // Exchange name
    "topic",                                 // Exchange type
)
```

### Method 2: Set Exchange Later
```go
cm := NewConnectionManager("", "", "")
cm.SetExchange("my-exchange", "direct")
```

## Publishing to Exchanges

### Simple Publishing
```go
// Publish simple text message
err := cm.PublishToExchange("events-exchange", "user.created", "User John created")
```

### Advanced Publishing
```go
// Publish with custom options
message := amqp.Publishing{
    ContentType:  "application/json",
    Body:         []byte(`{"userId": 123, "action": "created"}`),
    DeliveryMode: amqp.Persistent,
    Priority:     5,
    Headers: amqp.Table{
        "source": "user-service",
        "version": "1.0",
    },
}

err := cm.PublishToExchangeWithOptions("events-exchange", "user.created", message)
```

## Consuming from Exchanges

The `ConsumeFromExchange` method automatically:
1. Declares the exchange
2. Declares a queue
3. Binds the queue to the exchange
4. Starts consuming

```go
messageChannel := make(chan amqp.Delivery)

go func() {
    err := cm.ConsumeFromExchange(
        "events-exchange", // exchange name
        "topic",          // exchange type
        "user-events",    // queue name
        "user.*",         // routing key pattern
        messageChannel,
    )
    if err != nil {
        log.Printf("Consumer error: %v", err)
    }
}()

// Process messages
for msg := range messageChannel {
    log.Printf("Routing Key: %s, Message: %s", msg.RoutingKey, string(msg.Body))
    msg.Ack(false)
}
```

## Exchange Declaration

### Manual Exchange Declaration
```go
// Declare exchange with custom parameters
err := cm.DeclareExchange(
    "my-exchange", // name
    "topic",       // type
    true,          // durable
    false,         // auto-delete
    false,         // internal
    false,         // no-wait
    nil,           // arguments
)
```

### Exchange Parameters
- **durable**: Exchange survives broker restart
- **auto-delete**: Exchange deleted when no queues bound
- **internal**: Exchange can't be published to directly
- **no-wait**: Don't wait for server confirmation

## Complete Examples

### Event-Driven Architecture
```go
package main

import (
    "log"
    "time"
    "your-project/pkg/rabbitmq"
)

func main() {
    // Publisher
    publisher := rabbitmq.NewConnectionManagerWithExchange(
        "", "", "event-publisher", "events", "topic")
    publisher.Connect()
    defer publisher.Close()

    // Consumer
    consumer := rabbitmq.NewConnectionManager("", "", "event-consumer")
    consumer.Connect()
    defer consumer.Close()

    // Start consuming
    messageChannel := make(chan amqp.Delivery)
    go func() {
        consumer.ConsumeFromExchange("events", "topic", "user-events", "user.*", messageChannel)
    }()

    // Process messages
    go func() {
        for msg := range messageChannel {
            log.Printf("Event: %s - %s", msg.RoutingKey, string(msg.Body))
            msg.Ack(false)
        }
    }()

    // Publish events
    events := []string{"user.created", "user.updated", "user.deleted"}
    for i, event := range events {
        message := fmt.Sprintf("Event %d: %s", i+1, event)
        publisher.PublishToExchange("events", event, message)
        time.Sleep(2 * time.Second)
    }
}
```

### Microservices Communication
```go
// Service A - Publisher
func publishOrderEvent(cm *rabbitmq.ConnectionManager, orderID int, status string) {
    routingKey := fmt.Sprintf("order.%s", status)
    message := fmt.Sprintf(`{"orderId": %d, "status": "%s", "timestamp": "%s"}`, 
        orderID, status, time.Now().Format(time.RFC3339))
    
    err := cm.PublishToExchange("orders-exchange", routingKey, message)
    if err != nil {
        log.Printf("Failed to publish order event: %v", err)
    }
}

// Service B - Consumer
func consumeOrderEvents(cm *rabbitmq.ConnectionManager) {
    messageChannel := make(chan amqp.Delivery)
    
    go func() {
        // Listen to all order events
        cm.ConsumeFromExchange("orders-exchange", "topic", "inventory-service", "order.*", messageChannel)
    }()
    
    for msg := range messageChannel {
        // Process order event
        log.Printf("Processing order event: %s", string(msg.Body))
        // Update inventory, send notifications, etc.
        msg.Ack(false)
    }
}
```

## Routing Patterns

### Topic Exchange Patterns
- `user.created` - Exact match
- `user.*` - Matches `user.created`, `user.updated`, etc.
- `*.created` - Matches `user.created`, `order.created`, etc.
- `user.#` - Matches `user.created`, `user.profile.updated`, etc.
- `#` - Matches all messages

### Direct Exchange
- Exact routing key match only
- Perfect for simple routing scenarios

### Fanout Exchange
- Ignores routing key
- Broadcasts to all bound queues
- Great for notifications, logging

### Headers Exchange
- Routes based on message headers
- More flexible than routing keys
- Can match any/all headers

## Best Practices

1. **Use Topic Exchanges for Event-Driven Architecture**
   ```go
   // Good routing key structure
   "service.entity.action" // user-service.user.created
   "domain.subdomain.event" // orders.payment.completed
   ```

2. **Make Exchanges Durable**
   ```go
   cm.DeclareExchange("my-exchange", "topic", true, false, false, false, nil)
   //                                        ^^^^ durable = true
   ```

3. **Use Meaningful Routing Keys**
   ```go
   // Good
   cm.PublishToExchange("events", "user.profile.updated", message)
   
   // Bad
   cm.PublishToExchange("events", "event1", message)
   ```

4. **Handle Connection Failures**
   ```go
   for {
       err := cm.PublishToExchange("events", "user.created", message)
       if err == nil {
           break
       }
       log.Printf("Publish failed, retrying: %v", err)
       time.Sleep(time.Second)
   }
   ```

5. **Use Headers for Complex Routing**
   ```go
   message := amqp.Publishing{
       Body: []byte("Complex message"),
       Headers: amqp.Table{
           "region":   "us-east",
           "priority": "high",
           "type":     "alert",
       },
   }
   ```

## Migration from Direct Queue Publishing

### Before (Direct Queue)
```go
cm := NewConnectionManager("", "my-queue", "")
cm.Publish("message")
```

### After (Exchange-based)
```go
cm := NewConnectionManagerWithExchange("", "", "", "my-exchange", "direct")
cm.PublishToExchange("my-exchange", "my-queue", "message")
```

## Troubleshooting

### Exchange Not Found
- Ensure exchange is declared before publishing
- Check exchange name spelling
- Verify exchange type matches

### Messages Not Routed
- Check routing key patterns
- Verify queue binding
- Ensure exchange type is correct

### Connection Issues
- All exchange methods include automatic reconnection
- Check RabbitMQ server status
- Verify credentials and permissions