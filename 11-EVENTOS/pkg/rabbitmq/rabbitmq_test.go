package rabbitmq

import (
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestNewConnectionManager(t *testing.T) {
	cm := NewConnectionManager("", "", "")

	if cm.url != "amqp://admin:password@localhost:5672/" {
		t.Errorf("Expected default URL, got %s", cm.url)
	}

	if cm.queueName != "minha-fila" {
		t.Errorf("Expected default queue name, got %s", cm.queueName)
	}

	if cm.consumerTag != "go-consumer" {
		t.Errorf("Expected default consumer tag, got %s", cm.consumerTag)
	}
}

func TestConnectionManagerCustomValues(t *testing.T) {
	url := "amqp://test:test@localhost:5672/"
	queue := "test-queue"
	consumer := "test-consumer"

	cm := NewConnectionManager(url, queue, consumer)

	if cm.url != url {
		t.Errorf("Expected URL %s, got %s", url, cm.url)
	}

	if cm.queueName != queue {
		t.Errorf("Expected queue name %s, got %s", queue, cm.queueName)
	}

	if cm.consumerTag != consumer {
		t.Errorf("Expected consumer tag %s, got %s", consumer, cm.consumerTag)
	}
}

func TestIsConnectedInitially(t *testing.T) {
	cm := NewConnectionManager("", "", "")

	if cm.IsConnected() {
		t.Error("Expected connection to be false initially")
	}
}

func TestGetChannelWhenNotConnected(t *testing.T) {
	cm := NewConnectionManager("", "", "")

	_, err := cm.GetChannel()
	if err == nil {
		t.Error("Expected error when getting channel without connection")
	}
}

// Integration test - only runs if RabbitMQ is available
func TestConnectionManagerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cm := NewConnectionManager("", "", "")

	// Try to connect - this will fail if RabbitMQ is not running
	err := cm.Connect()
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer cm.Close()

	// Test connection status
	if !cm.IsConnected() {
		t.Error("Expected connection to be true after successful connect")
	}

	// Test getting channel
	ch, err := cm.GetChannel()
	if err != nil {
		t.Errorf("Failed to get channel: %v", err)
	}
	if ch == nil {
		t.Error("Expected non-nil channel")
	}

	// Test close
	err = cm.Close()
	if err != nil {
		t.Errorf("Failed to close connection: %v", err)
	}

	// Give it a moment to close
	time.Sleep(100 * time.Millisecond)

	if cm.IsConnected() {
		t.Error("Expected connection to be false after close")
	}
}
func TestConnectionManager_Publish(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cm := NewConnectionManager("", "", "")

	// Try to connect - this will fail if RabbitMQ is not running
	err := cm.Connect()
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer cm.Close()

	// Test simple publish
	err = cm.Publish("Test message")
	if err != nil {
		t.Errorf("Failed to publish message: %v", err)
	}
}

func TestConnectionManager_PublishToQueue(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cm := NewConnectionManager("", "", "")

	err := cm.Connect()
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer cm.Close()

	// Test publish to specific queue
	err = cm.PublishToQueue("test-queue", "Test message to specific queue")
	if err != nil {
		t.Errorf("Failed to publish to specific queue: %v", err)
	}
}

func TestConnectionManager_PublishWithOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cm := NewConnectionManager("", "", "")

	err := cm.Connect()
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer cm.Close()

	// Test publish with custom options
	customMsg := amqp.Publishing{
		ContentType:  "application/json",
		Body:         []byte(`{"test": "message"}`),
		DeliveryMode: amqp.Persistent,
		Headers: amqp.Table{
			"test": "header",
		},
	}

	err = cm.PublishWithOptions("test-queue", "test-queue", customMsg)
	if err != nil {
		t.Errorf("Failed to publish with custom options: %v", err)
	}
}

func TestConnectionManager_PublishWhenNotConnected(t *testing.T) {
	cm := NewConnectionManager("amqp://invalid:invalid@nonexistent:5672/", "", "")

	// Try to publish without connection
	err := cm.Publish("Test message")
	if err == nil {
		t.Error("Expected error when publishing without connection")
	}
}
func TestNewConnectionManagerWithExchange(t *testing.T) {
	url := "amqp://test:test@localhost:5672/"
	queue := "test-queue"
	consumer := "test-consumer"
	exchange := "test-exchange"
	exchangeType := "topic"

	cm := NewConnectionManagerWithExchange(url, queue, consumer, exchange, exchangeType)

	if cm.url != url {
		t.Errorf("Expected URL %s, got %s", url, cm.url)
	}

	if cm.queueName != queue {
		t.Errorf("Expected queue name %s, got %s", queue, cm.queueName)
	}

	if cm.consumerTag != consumer {
		t.Errorf("Expected consumer tag %s, got %s", consumer, cm.consumerTag)
	}

	if cm.exchangeName != exchange {
		t.Errorf("Expected exchange name %s, got %s", exchange, cm.exchangeName)
	}

	if cm.exchangeType != exchangeType {
		t.Errorf("Expected exchange type %s, got %s", exchangeType, cm.exchangeType)
	}
}

func TestConnectionManager_SetGetExchange(t *testing.T) {
	cm := NewConnectionManager("", "", "")

	// Test initial values
	exchangeName, exchangeType := cm.GetExchange()
	if exchangeName != "" {
		t.Errorf("Expected empty exchange name initially, got %s", exchangeName)
	}
	if exchangeType != "direct" {
		t.Errorf("Expected 'direct' exchange type initially, got %s", exchangeType)
	}

	// Test setting values
	cm.SetExchange("test-exchange", "topic")
	exchangeName, exchangeType = cm.GetExchange()
	if exchangeName != "test-exchange" {
		t.Errorf("Expected 'test-exchange', got %s", exchangeName)
	}
	if exchangeType != "topic" {
		t.Errorf("Expected 'topic', got %s", exchangeType)
	}
}

func TestConnectionManager_PublishToExchange(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cm := NewConnectionManagerWithExchange("", "", "", "test-exchange", "direct")

	err := cm.Connect()
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer cm.Close()

	// Test publish to exchange
	err = cm.PublishToExchange("test-exchange", "test.routing.key", "Test message to exchange")
	if err != nil {
		t.Errorf("Failed to publish to exchange: %v", err)
	}
}

func TestConnectionManager_PublishToExchangeWithOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cm := NewConnectionManager("", "", "")

	err := cm.Connect()
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer cm.Close()

	// Declare exchange first
	err = cm.DeclareExchange("test-exchange-options", "topic", true, false, false, false, nil)
	if err != nil {
		t.Errorf("Failed to declare exchange: %v", err)
	}

	// Test publish with custom options
	customMsg := amqp.Publishing{
		ContentType:  "application/json",
		Body:         []byte(`{"test": "exchange message"}`),
		DeliveryMode: amqp.Persistent,
		Headers: amqp.Table{
			"exchange-test": "header",
		},
	}

	err = cm.PublishToExchangeWithOptions("test-exchange-options", "test.topic.key", customMsg)
	if err != nil {
		t.Errorf("Failed to publish to exchange with options: %v", err)
	}
}

func TestConnectionManager_DeclareExchange(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cm := NewConnectionManager("", "", "")

	err := cm.Connect()
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer cm.Close()

	// Test exchange declaration
	err = cm.DeclareExchange("test-declare-exchange", "fanout", true, false, false, false, nil)
	if err != nil {
		t.Errorf("Failed to declare exchange: %v", err)
	}
}
