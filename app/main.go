package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type request struct {
	method  string
	path    string
	version string
	header  map[string]string
	body    string
}

func NewRequest(b []byte) *request {
	req := &request{}
	ele := strings.Split(string(b), "\r\n")
	startLine := strings.Split(ele[0], " ")

	req.method = startLine[0]
	req.path = startLine[1]
	req.version = startLine[2]
	header := make(map[string]string)
	for _, i := range ele {
		mapEle := strings.Split(i, ": ")
		if len(mapEle) == 2 {
			key := strings.ToLower(mapEle[0])
			value := mapEle[1]
			header[key] = value
		}
	}
	req.header = header
	return req
}

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
	req := NewRequest(reader)
	respond := "HTTP/1.1 200 OK\r\n\r\n"
	if strings.HasPrefix(req.path, "/echo/") {
		str := req.path[6:]
		respond = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(str), str)
	} else if strings.HasPrefix(req.path, "/user-agent") {
		str := req.header["user-agent"]
		respond = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(str), str)
	} else if req.path != "/" {
		respond = "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	conn.Write([]byte(respond))
	conn.Close()
}
