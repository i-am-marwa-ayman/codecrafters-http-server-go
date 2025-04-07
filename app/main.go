package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

type Request struct {
	method  string
	path    string
	version string
	header  map[string]string
	body    string
}

func NewRequest(data []byte) *Request {
	lines := strings.Split(string(data), "\r\n")
	if len(lines) == 0 {
		return nil
	}
	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) != 3 {
		return nil
	}
	req := &Request{
		method:  requestLine[0],
		path:    requestLine[1],
		version: requestLine[2],
		header:  make(map[string]string),
	}
	for _, i := range lines {
		if i == "" {
			break // empty line end of header
		}
		mapEle := strings.Split(i, ": ")
		if len(mapEle) == 2 {
			key := strings.ToLower(mapEle[0])
			value := mapEle[1]
			req.header[key] = value
		}
	}
	return req
}
func GetRespond(req *Request) string {
	respond := "HTTP/1.1 200 OK\r\n\r\n"
	if strings.HasPrefix(req.path, "/echo") {
		str := req.path[6:]
		respond = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(str), str)
	} else if strings.HasPrefix(req.path, "/user-agent") {
		str := req.header["user-agent"]
		respond = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(str), str)
	} else if req.path != "/" {
		respond = "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	return respond
}
func HandleConnction(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}
	req := NewRequest(buffer[:n])
	if req == nil {
		return
	}
	respond := GetRespond(req)
	conn.Write([]byte(respond))
}
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	var dir = flag.String("directory", ".", "")
	flag.Parse()
	err := os.Chdir(*dir)
	if err != nil {
		fmt.Println("Failed to change the cur directory")
		os.Exit(1)
	}

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go HandleConnction(conn)
	}
}
