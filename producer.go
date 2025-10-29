package main

import (
	"context"
	"fmt"
	"os"

	"github.com/comvex-jp/backend-service-go-framework/v8/console"
	"github.com/comvex-jp/backend-service-go-framework/v8/dispatcher/kafka"
)

// ProducerCommand implements the console.Command interface for producing messages to Kafka
type ProducerCommand struct {
	console.Base
}

// NewProducerCommand creates a new ProducerCommand instance
func NewProducerCommand() *ProducerCommand {
	return &ProducerCommand{
		Base: console.Base{
			Name:        "producer",
			Description: "Send a message to a Kafka topic",
		},
	}
}

// Handle processes the producer command
func (c *ProducerCommand) Handle(args console.Arguments, options console.Options) error {
	argItems := args.GetItems()
	if len(argItems) < 2 {
		return fmt.Errorf("usage: producer <topic> <message>")
	}

	topic := argItems[0]
	message := argItems[1]

	// Get Kafka broker from environment variable
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		return fmt.Errorf("KAFKA_BROKER environment variable is not set")
	}

	// Create Kafka producer using the framework's dispatcher
	producer := kafka.NewProducer(
		[]string{broker},
		kafka.WithAutoTopicCreation(), // Enable automatic topic creation
	)
	defer producer.Close()

	// Add optional headers if needed
	// kafkaMessage.SetHeader("content-type", "text/plain")

	// Send message
	ctx := context.Background()
	for i := range 100 {
		msg := fmt.Sprintf("%s %d", message, i)
		// Create message
		kafkaMessage := kafka.Message{
			Value: []byte(msg),
			Key:   []byte(""), // Empty key for now
		}

		err := producer.Produce(ctx, topic, kafkaMessage)
		if err != nil {
			return fmt.Errorf("error producing message: %v", err)
		}
		fmt.Printf("Message produced successfully to topic '%s': %s\n", topic, msg)
	}

	return nil
}
