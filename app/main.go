package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	reader := make([]byte, 1024)
	conn.Read(reader)

	respond := "HTTP/1.1 200 OK\r\n\r\n"
	if strings.HasPrefix(string(reader), "GET /echo/") {
		str := strings.Split(string(reader), " ")[1]
		str = str[6:]
		respond = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(str), str)
	} else if !strings.HasPrefix(string(reader), "GET / ") {
		respond = "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	conn.Write([]byte(respond))
	conn.Close()
}
