# Kafka Consumer/Producer CLI

A command-line tool for producing and consuming messages from Kafka topics using the `github.com/comvex-jp/backend-service-go-framework/v8` framework.

## Features

- **Producer Command**: Send messages to Kafka topics using the framework's Kafka dispatcher
- **Consumer Command**: Consume messages from Kafka topics using the framework's event router
- **Environment Configuration**: Kafka broker configuration via environment variables
- **Framework Integration**: Both producer and consumer use the framework's built-in components
- **Event Router Integration**: Uses the framework's event router for robust message processing
- **Graceful Shutdown**: Proper signal handling for clean shutdown

## Prerequisites

- Go 1.25.2 or later
- Kafka server running and accessible

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Build the application:
   ```bash
   go build -o kafka-cli .
   ```

## Configuration

Set the following environment variables for Kafka configuration:

```bash
export KAFKA_BROKER=localhost:9092
export KAFKA_CONSUMER_GROUP=kafka-consumer-group
export KAFKA_PROCESS_COUNT=1  # Optional, defaults to 1
```

Or create a `.env` file in the project root:

```bash
cp env.example .env
# Edit .env with your Kafka configuration
```

### Environment Variables

- `KAFKA_BROKER`: Kafka broker address (required)
- `KAFKA_CONSUMER_GROUP`: Consumer group ID (required for consumer)
- `KAFKA_PROCESS_COUNT`: Number of concurrent processors (optional, default: 1)

## Usage

### Producer Command

Send a message to a Kafka topic:

```bash
./kafka-cli producer <topic> <message>
```

**Example:**

```bash
./kafka-cli producer my-topic "Hello, Kafka!"
```

### Consumer Command

Consume messages from a Kafka topic:

```bash
./kafka-cli consumer <topic>
```

**Example:**

```bash
./kafka-cli consumer my-topic
```

The consumer will run continuously using the framework's event router and print received messages to the console. It supports:

- Multiple concurrent processors (configurable via `KAFKA_PROCESS_COUNT`)
- Consumer group management
- Message headers display
- Graceful shutdown with `Ctrl+C`

## Commands

- `producer <topic> <message>` - Send a message to the specified topic
- `consumer <topic>` - Consume messages from the specified topic
- `--help` - Show help information

## Dependencies

- `github.com/comvex-jp/backend-service-go-framework/v8` - Backend service framework with event router and Kafka dispatcher
- `github.com/segmentio/kafka-go` - Kafka client library (used by both event router and dispatcher)
- `github.com/joho/godotenv` - Environment variable loading from .env files
