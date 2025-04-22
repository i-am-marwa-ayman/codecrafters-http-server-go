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
- **Content Compression**: Supports gzip compression for optimized data transfer
- **Persistent Connections**: Implements HTTP keep-alive for improved performance

## Getting Started

### Installation

1. Clone this repository
```bash
git clone https://github.com/i-am-marwa-ayman/codecrafters-http-server-go.git
cd codecrafters-http-server-go
```

2. Build the server
```bash
./your_program.sh
```

### Usage

Run the server with default settings:
```bash
./app
```

Specify a directory for file operations:
```bash
./app --directory /path/to/files
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Returns 200 OK |
| GET | `/echo/{string}` | Returns the provided string |
| GET | `/user-agent` | Returns the User-Agent header from the request |
| GET | `/files/{filename}` | Serves the specified file |
| POST | `/files/{filename}` | Creates or updates the specified file |

## Content Compression

The server automatically detects if a client supports gzip compression by checking the Accept-Encoding header. If supported, responses are compressed to reduce bandwidth usage and improve performance.

## Connection Management

The server supports HTTP persistent connections (keep-alive) by default, which allows multiple requests to be sent over a single TCP connection. This reduces latency and improves performance by eliminating the need to establish a new connection for each request.

- Connections remain open for subsequent requests by default
- Clients can request connection closure by sending the "Connection: close" header
- The server properly handles connection state for each client

## How It Works

1. The server listens on port 4221
2. When a connection is received, it's handled in a separate goroutine
3. The request is parsed and routed to the appropriate handler
4. The response is generated and sent back to the client
5. The connection remains open for additional requests unless explicitly closed

## Example Usage

### Echo Service

```http
GET /echo/hello-world HTTP/1.1
Host: localhost:4221
Accept-Encoding: gzip

HTTP/1.1 200 OK
Content-Encoding: gzip
Content-Type: text/plain
Content-Length: 30

[compressed content]
```

Without compression:
```http
GET /echo/hello-world HTTP/1.1
Host: localhost:4221

HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 11

hello-world
```

### User-Agent Information

```http
GET /user-agent HTTP/1.1
Host: localhost:4221
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)

HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 41

Mozilla/5.0 (Windows NT 10.0; Win64; x64)
```

### File Operations

Creating a file:
```http
POST /files/example.txt HTTP/1.1
Host: localhost:4221
Content-Length: 18

This is a test file

HTTP/1.1 201 Created
```

Retrieving a file:
```http
GET /files/example.txt HTTP/1.1
Host: localhost:4221

HTTP/1.1 200 OK
Content-Type: application/octet-stream
Content-Length: 18

This is a test file
```

### Persistent Connection Example

Multiple requests on the same connection:
```http
GET /echo/first-request HTTP/1.1
Host: localhost:4221

HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 13

first-request

GET /echo/second-request HTTP/1.1
Host: localhost:4221

HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 14

second-request
```

Closing a connection:
```http
GET /echo/final-request HTTP/1.1
Host: localhost:4221
Connection: close

HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 13
Connection: close

final-request
```

## Performance

The server is designed to handle multiple concurrent connections efficiently using Go's goroutines. Each connection is processed independently, allowing the server to scale with available system resources. The implementation of persistent connections further improves performance by reducing the overhead of establishing new TCP connections for each request.

## Extending the Server

The modular design makes it easy to extend the server with additional functionality:

1. Add new handlerfunc in `router.go`
2. Add new route in init router function in `main.go`
3. Implement new response types in `respond.go`
4. Add additional file operations in `fileHelper.go`

### Adding a New Route

```go
// In initRouter() function in main.go
router.AddRoute("GET", "/newpath", NewPathHandler)

// Then define your handler function in router.go
func NewPathHandler(req *Request) *Respond {
    res := NewRespond()
    res.SetStatusLine(200, "OK")
    res.AddHeader("Content-Type", "text/plain")
    res.body = []byte("Hello from new path!")
    return res
}
```

## Future Improvements

- Implement better request/respond structure
- Implement request timeouts and connection limits

## Feedback

Feedback is more than welcome, whether it's a suggestion or a roast.