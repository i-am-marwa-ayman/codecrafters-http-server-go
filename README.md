# Simple HTTP Server

A lightweight, concurrent HTTP server written in Go. This project was completed as part of the CodeCrafters "Build Your Own HTTP Server" challenge.

## Features

- **HTTP Request Parsing**: Cleanly parses HTTP requests including headers and body
- **Echo Endpoint**: Returns any string following `/echo/` in the request path
- **User-Agent Endpoint**: Returns the client's User-Agent header
- **File Operations**:
  - GET `/files/{filename}`: Serves files from a specified directory
  - POST `/files/{filename}`: Creates or updates files
- **Concurrent Connections**: Uses Go routines to handle multiple connections simultaneously

## Getting Started

### Prerequisites

- Go 1.16 or higher

### Installation

1. Clone this repository
```bash
git clone https://github.com/i-am-marwa-ayman/codecrafters-http-server-go.git
cd codecrafters-http-server-go
```

2. Build the server
```bash
go build main.go
```

### Usage

Run the server with default settings:
```bash
./main
```

Specify a directory for file operations:
```bash
./main --directory /path/to/files
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Returns 200 OK |
| GET | `/echo/{string}` | Returns the provided string |
| GET | `/user-agent` | Returns the User-Agent header from the request |
| GET | `/files/{filename}` | Serves the specified file |
| POST | `/files/{filename}` | Creates or updates the specified file |

## How It Works

1. The server listens on port 4221
2. When a connection is received, it's handled in a separate goroutine
3. The request is parsed and routed to the appropriate handler
4. The response is generated and sent back to the client

## Project Structure

- `main.go`: Contains all the code for this simple server
  - `Request` struct: Represents an HTTP request
  - `NewRequest()`: Parses raw request data
  - `GetRespond()`: Handles GET requests
  - `PostRespond()`: Handles POST requests
  - `HandleConnection()`: Manages incoming connections
  - `main()`: Sets up the server and listens for connections

