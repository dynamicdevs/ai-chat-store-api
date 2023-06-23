# GPTO

## Description

GPTO is an AI-powered assistant tool for e-commerce platforms. It has two main functionalities: indexing product data into a database, and acting as a chatbot to assist customers on e-commerce websites.

## Swagger

<url>/docs

## Generating Code Documentation

To generate and view code documentation locally, use the `godoc` tool. Follow these steps:

1. Ensure that you have Go installed on your machine.

```sh
 go install golang.org/x/tools/cmd/godoc@latest
```

2. Run the `godoc` server:

   ```sh
   godoc -http=:6060
   ```

3. Open your web browser and go to: [http://localhost:6060](http://localhost:6060).

4. To find the documentation for your project, navigate to the pkg tab.
   Here, you will find a list of all the packages available in your $GOPATH. Locate your project's package and click on it.
   For example, if your package is at github.com/Abraxas-365/commerce-chat, you can directly go to:
   [http://localhost:6060/pkg/github.com/Abraxas-365/commerce-chat/]

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
