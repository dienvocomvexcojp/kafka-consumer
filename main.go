package main

import (
	"fmt"
	"os"

	"github.com/comvex-jp/backend-service-go-framework/v8/app"
	"github.com/comvex-jp/backend-service-go-framework/v8/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	core := logger.NewConsoleCore(logger.InfoLevel)
	logger.Init(logger.Options{}, core)

	// Create the main application
	application := app.New()

	// Initialize the console kernel with a root command
	application.ConsoleKernel.NewCommand("kafka-cli")

	// Register commands
	application.ConsoleKernel.RegisterCommand(NewProducerCommand())
	application.ConsoleKernel.RegisterCommand(NewConsumerCommand())

	// Boot the application
	application.Boot()

	// Run the console kernel
	if err := application.ConsoleKernel.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
