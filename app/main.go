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

var router *Router

func initRouter() {
	router = NewRouter()
	router.AddRoute("GET", "/", RootHandler)
	router.AddRoute("GET", "/echo", EchoHandler)
	router.AddRoute("GET", "/user-agent", UserAgentHandler)
	router.AddRoute("GET", "/files", FilesGetHandler)
	router.AddRoute("POST", "/files", FilesPostHandler)
}
func CompressData(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}
	gz.Close()

	return b.Bytes(), nil
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

func HandleRequest(req *Request) ([]byte, bool) {
	res := router.Handle(req)
	cancel := false
	if req.header["connection"] == "close" {
		cancel = true
		res.AddHeader("Connection", "close")
	}
	return res.Serialize(), cancel
}
func HandleConnection(conn net.Conn) {
	defer conn.Close()
	cancel := false
	for !cancel {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		req := NewRequest(buffer[:n])
		if req == nil {
			return
		}
		var respond []byte
		respond, cancel = HandleRequest(req)
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
	initRouter()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go HandleConnection(conn)
	}
}
