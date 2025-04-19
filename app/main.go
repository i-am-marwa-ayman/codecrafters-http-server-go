package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

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

func HandleGetRequest(req *Request) string {
	res := NewRespond()
	if strings.HasPrefix(req.path, "/echo") {
		str := req.path[6:]
		if HasValidEncodingScheme(req.header["accept-encoding"]) {
			data, err := CompressData([]byte(str))
			if err == nil {
				header := fmt.Sprintf("Content-Encoding: gzip\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n", len(data))
				res.OkRespond(header, data)
			}
		} else {
			header := fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n", len(str))
			res.OkRespond(header, str)
		}
	} else if strings.HasPrefix(req.path, "/user-agent") {
		str := req.header["user-agent"]
		header := fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n", len(str))
		res.OkRespond(header, str)
	} else if strings.HasPrefix(req.path, "/files") {
		fileName := req.path[7:]
		str, err := GetFileContent(fileName)
		if err == nil {
			header := fmt.Sprintf("Content-Type: application/octet-stream\r\nContent-Length: %d\r\n", len(str))
			res.OkRespond(header, str)
		}
	} else if req.path == "/" {
		res.OkRespond("", "")
	}
	return res.ToString()
}
func HandlePostRequest(req *Request) string {
	res := NewRespond()
	if strings.HasPrefix(req.path, "/files") {
		fileName := req.path[7:]
		ok := AddFile(fileName, req.body)
		if ok {
			res.code = 201
			res.msg = "Created"
		} else {
			res.code = 500
			res.msg = "Internal Server Error"
		}
	}
	return res.ToString()
}
func HandleConnection(conn net.Conn) {
	defer conn.Close()
	ok := false
	for !ok {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		req := NewRequest(buffer[:n])
		if req == nil {
			return
		}
		_, ok = req.header["connection"]
		var respond string
		if req.method == "GET" {
			respond = HandleGetRequest(req)
		} else if req.method == "POST" {
			respond = HandlePostRequest(req)
		}
		conn.Write([]byte(respond))
	}
	fmt.Println("connection closed")
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
		go HandleConnection(conn)
	}
}
