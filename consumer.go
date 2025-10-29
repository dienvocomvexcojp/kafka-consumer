package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/comvex-jp/backend-service-go-framework/v8/console"
	"github.com/comvex-jp/backend-service-go-framework/v8/logger"
	"github.com/comvex-jp/backend-service-go-framework/v8/router/event"
)

// ConsumerCommand implements the console.Command interface for consuming messages from Kafka
type ConsumerCommand struct {
	console.Base
}

// NewConsumerCommand creates a new ConsumerCommand instance
func NewConsumerCommand() *ConsumerCommand {
	return &ConsumerCommand{
		Base: console.Base{
			Name:        "consumer",
			Description: "Consume messages from a Kafka topic",
		},
	}
}

// Handle processes the consumer command
func (c *ConsumerCommand) Handle(args console.Arguments, options console.Options) error {
	argItems := args.GetItems()
	if len(argItems) < 1 {
		return fmt.Errorf("usage: consumer <topic>")
	}

	topic := argItems[0]

	// Get Kafka broker from environment variable
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		return fmt.Errorf("KAFKA_BROKER environment variable is not set")
	}

	// Get Kafka consumer group from environment variable
	consumerGroup := os.Getenv("KAFKA_CONSUMER_GROUP")
	if consumerGroup == "" {
		return fmt.Errorf("KAFKA_CONSUMER_GROUP environment variable is not set")
	}

	// Get process count from environment variable (optional, default to 1)
	processCount := 1
	if processCountStr := os.Getenv("KAFKA_PROCESS_COUNT"); processCountStr != "" {
		if count, err := strconv.Atoi(processCountStr); err == nil {
			processCount = count
		}
	}

	// Create event server
	server := event.NewServer("kafka-consumer")
	server.RegisterHandlerMiddlewares([]event.MiddlewareHandlerFunc{})

	// Setup processor configuration
	server.SetupProcessorConfig(event.ProcessorConfig{
		Brokers:      []string{broker},
		GroupId:      consumerGroup,
		ProcessCount: processCount,
	})

	// Register event handler for the topic
	server.Route(topic, func(ctx context.Context) {
		// Get the event from context
		evt, err := event.FromContext(ctx)
		if err != nil {
			logger.Logw(ctx, logger.ErrorLevel, "Error getting event from context", "error", err)
			return
		}

		// Print the message details
		logger.Logw(ctx, logger.InfoLevel, "Received message from topic", "topic", evt.Topic(), "message", string(evt.Value()))

		// Print headers if any
		if headers := evt.Headers(); len(headers) > 0 {
			logger.Logw(ctx, logger.InfoLevel, "Headers", "headers", headers)
		}
	})

	fmt.Printf("Starting consumer for topic: %s\n", topic)
	fmt.Printf("Consumer group: %s\n", consumerGroup)
	fmt.Printf("Process count: %d\n", processCount)
	fmt.Println("Press Ctrl+C to stop...")

	// Set up signal handling for graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	serverDone := make(chan error, 1)
	go func() {
		serverDone <- server.Serve()
	}()

	// Wait for signal or server completion
	select {
	case sig := <-sigchan:
		fmt.Printf("\nCaught signal %v: terminating\n", sig)
		server.Stop()
		<-serverDone // Wait for server to stop gracefully
		return nil
	case err := <-serverDone:
		if err != nil {
			return fmt.Errorf("server error: %v", err)
		}
		return nil
	}
}
