# GPTO

## Description

GPTO is an AI-powered assistant tool for e-commerce platforms. It has two main functionalities: indexing product data into a database, and acting as a chatbot to assist customers on e-commerce websites.

## Project Structure

.
├── cmd
│ ├── gpt
│ │ └── main.go # Main file for chatbot functionality
│ └── indexer
│ └── main.go # Main file for data indexing functionality
├── internal # Internal packages
├── migrations # Database migrations
├── pkg # Public packages
├── Dockerfile # Dockerfile for building the application
├── docker-compose.yml # Docker Compose file for running services
├── getProduct.bash # Bash script for getting product data
├── go.mod # Go module file
├── go.sum # Go module checksums
├── output.csv # Output file for indexed data

## Prerequisites

- Go (version 1.18.x or higher)
- Docker
- Docker Compose

## Setup & Installation

1. Clone the repository
   ```sh
   git clone https://github.com/dynamicdevs/ai-chat-store-api
   cd ai-chat-store-api
   ```

### Running in Docker

1. Build the Docker image

```sh
docker build -t gpto .
```

2. Start the services using Docker Compose

```sh
docker-compose up
```

### Running without Docker

1. Running the Indexer
   The indexer is responsible for indexing product data into the database. To run the indexer:

```sh
go run cmd/indexer/main.go
```

This will index the product data into the database, ready for the chatbot to use.

2. Running the Chatbot
   The chatbot can be used to assist customers on e-commerce websites. To run the chatbot:

```sh
go run cmd/gpt/main.go
```
