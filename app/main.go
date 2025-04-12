package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
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
	lines := strings.Split(string(data), "\r\n\r\n")
	if len(lines) == 0 {
		return nil
	}

	req := &Request{}
	if len(lines) == 2 {
		req.body = lines[1]
		fmt.Println(req.body)
	}

	lines = strings.Split(lines[0], "\r\n")

	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) != 3 {
		return nil
	}
	req.method = requestLine[0]
	req.path = requestLine[1]
	req.version = requestLine[2]
	req.header = make(map[string]string)

	for i := 1; i < len(lines); i++ {
		parts := strings.SplitN(lines[i], ": ", 2)
		if len(parts) == 2 {
			key := strings.ToLower(parts[0])
			value := parts[1]
			req.header[key] = value
		}
	}

	return req
}
func GetFileContent(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	return string(content), nil
}
func CompressData(data []byte) (string, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(data)
	if err != nil {
		return "", err
	}
	gz.Close()

	return b.String(), nil
}
func HasValidEncodingScheme(Schemes string) bool {
	options := strings.Split(Schemes, ", ")
	for _, op := range options {
		if op == "gzip" {
			return true
		}
	}
	return false
}
func GetRespond(req *Request) string {
	respond := "HTTP/1.1 404 Not Found\r\n\r\n"
	if strings.HasPrefix(req.path, "/echo") {
		str := req.path[6:]
		if HasValidEncodingScheme(req.header["accept-encoding"]) {
			data, err := CompressData([]byte(str))
			if err == nil {
				respond = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Encoding: gzip\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(data), data)
			}
		} else {
			respond = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(str), str)
		}
	} else if strings.HasPrefix(req.path, "/user-agent") {
		str := req.header["user-agent"]
		respond = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(str), str)
	} else if strings.HasPrefix(req.path, "/files") {
		fileName := req.path[7:]
		str, err := GetFileContent(fileName)
		if err != nil {
			respond = "HTTP/1.1 404 Not Found\r\n\r\n"
		} else {
			respond = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(str), str)
		}
	} else if req.path == "/" {
		respond = "HTTP/1.1 200 OK\r\n\r\n"
	}
	return respond
}
func AddFile(fileName string, content string) bool {
	file, err := os.Create(fileName)
	defer file.Close()

	if err != nil {
		return false
	}
	_, err = file.WriteString(content)
	if err != nil {
		return false
	}
	return true
}
func PostRespond(req *Request) string {
	respond := "HTTP/1.1 200 OK\r\n\r\n"
	if strings.HasPrefix(req.path, "/files") {
		fileName := req.path[7:]
		done := AddFile(fileName, req.body)
		if done {
			respond = "HTTP/1.1 201 Created\r\n\r\n"
		} else {
			respond = "HTTP/1.1 500 Internal Server Error\r\n\r\n"
		}
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
	var respond string
	if req.method == "GET" {
		respond = GetRespond(req)
	} else {
		respond = PostRespond(req)
	}
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
