# Publishing Messages with ConnectionManager

This guide shows how to use the new Publish methods in the ConnectionManager.

## Quick Start

```go
package main

import (
    "log"
    "your-project/pkg/rabbitmq"
)

func main() {
    // Create connection manager
    cm := rabbitmq.NewConnectionManager("", "", "")
    
    err := cm.Connect()
    if err != nil {
        log.Fatal(err)
    }
    defer cm.Close()
    
    // Publish a simple message
    err = cm.Publish("Hello World!")
    if err != nil {
        log.Printf("Failed to publish: %v", err)
    }
}
```

## Three Ways to Publish

### 1. Simple Publish (to default queue)

```go
// Publishes to the queue specified when creating ConnectionManager
err := cm.Publish("Simple message")
```

### 2. Publish to Specific Queue

```go
// Publishes to a specific queue
err := cm.PublishToQueue("my-custom-queue", "Message for custom queue")
```

### 3. Advanced Publish with Custom Options

```go
import amqp "github.com/rabbitmq/amqp091-go"

// Create custom message with headers, priority, etc.
customMessage := amqp.Publishing{
    ContentType:  "application/json",
    Body:         []byte(`{"event": "user.created", "userId": 123}`),
    DeliveryMode: amqp.Persistent, // Survive broker restart
    Priority:     5,               // Message priority (0-9)
    Timestamp:    time.Now(),
    Headers: amqp.Table{
        "source":      "user-service",
        "event-type":  "domain-event",
        "retry-count": 0,
    },
}

err := cm.PublishWithOptions("events-queue", "events-queue", customMessage)
```

## Features

### Automatic Reconnection
All publish methods automatically handle connection failures:
- Retries up to 3 times
- Exponential backoff between retries
- Waits for reconnection if connection is lost
- Triggers reconnection on connection errors

### Error Handling
```go
err := cm.Publish("message")
if err != nil {
    // Handle error
    // Common errors:
    // - amqp.ErrClosed: Connection closed
    // - Timeout errors: Connection couldn't be established
    log.Printf("Publish failed: %v", err)
}
```

## Complete Example: Publisher Service

```go
package main

import (
    "fmt"
    "log"
    "time"
    "your-project/pkg/rabbitmq"
)

func main() {
    // Create connection manager
    cm := rabbitmq.NewConnectionManager(
        "amqp://admin:password@localhost:5672/",
        "events-queue",
        "publisher-service",
    )
    
    // Connect with retry
    for {
        err := cm.Connect()
        if err == nil {
            break
        }
        log.Printf("Connection failed, retrying in 5s: %v", err)
        time.Sleep(5 * time.Second)
    }
    defer cm.Close()
    
    log.Println("Publisher started")
    
    // Publish messages periodically
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()
    
    counter := 0
    for range ticker.C {
        counter++
        message := fmt.Sprintf("Event #%d at %s", counter, time.Now().Format(time.RFC3339))
        
        err := cm.Publish(message)
        if err != nil {
            log.Printf("Failed to publish: %v", err)
            continue
        }
        
        log.Printf("Published: %s", message)
    }
}
```

## Message Properties

When using `PublishWithOptions`, you can set various message properties:

```go
amqp.Publishing{
    // Content
    ContentType:     "application/json",  // MIME type
    ContentEncoding: "utf-8",            // Encoding
    Body:            []byte("message"),   // Message body
    
    // Delivery
    DeliveryMode: amqp.Persistent,       // 1=non-persistent, 2=persistent
    Priority:     5,                     // 0-9 priority
    
    // Metadata
    CorrelationId: "request-123",        // For RPC patterns
    ReplyTo:       "response-queue",     // For RPC patterns
    MessageId:     "msg-456",            // Unique message ID
    Timestamp:     time.Now(),           // Message timestamp
    Type:          "user.created",       // Message type
    UserId:        "service-account",    // User ID
    AppId:         "user-service",       // Application ID
    
    // Custom headers
    Headers: amqp.Table{
        "x-custom-header": "value",
        "retry-count":     0,
    },
}
```

## Best Practices

1. **Always defer Close()**
   ```go
   cm := rabbitmq.NewConnectionManager("", "", "")
   err := cm.Connect()
   defer cm.Close() // Always close when done
   ```

2. **Check connection before publishing**
   ```go
   if !cm.IsConnected() {
       log.Println("Not connected, waiting...")
       // Wait or retry
   }
   ```

3. **Use persistent messages for important data**
   ```go
   msg := amqp.Publishing{
       DeliveryMode: amqp.Persistent, // Survives broker restart
       Body:         []byte("important data"),
   }
   ```

4. **Add retry logic for critical messages**
   ```go
   maxRetries := 5
   for i := 0; i < maxRetries; i++ {
       err := cm.Publish("critical message")
       if err == nil {
           break
       }
       log.Printf("Retry %d/%d: %v", i+1, maxRetries, err)
       time.Sleep(time.Second * time.Duration(i+1))
   }
   ```

5. **Use headers for metadata**
   ```go
   msg := amqp.Publishing{
       Body: []byte("data"),
       Headers: amqp.Table{
           "source":    "my-service",
           "timestamp": time.Now().Unix(),
           "version":   "1.0",
       },
   }
   ```

## Comparison with Legacy Function

### Old Way (No Reconnection)
```go
ch, err := rabbitmq.OpenChannel()
if err != nil {
    log.Fatal(err)
}
defer ch.Close()

err = rabbitmq.Publish(ch, "message")
// If connection lost, publish fails permanently
```

### New Way (With Reconnection)
```go
cm := rabbitmq.NewConnectionManager("", "", "")
err := cm.Connect()
defer cm.Close()

err = cm.Publish("message")
// Automatically retries and reconnects if needed
```

## Troubleshooting

### Message not being published
- Check if RabbitMQ is running
- Verify connection credentials
- Check if queue exists
- Look for error logs

### Connection keeps dropping
- Check network stability
- Verify RabbitMQ server health
- Check firewall settings
- Review RabbitMQ logs

### Messages lost
- Use `DeliveryMode: amqp.Persistent`
- Ensure queue is durable
- Use publisher confirms (advanced)
- Implement application-level acknowledgments