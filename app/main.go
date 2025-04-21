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
	res := NewRespond()
	cancel := false
	if strings.HasPrefix(req.path, "/echo") {
		str := req.path[6:]
		if HasValidEncodingScheme(req.header["accept-encoding"]) {
			data, err := CompressData([]byte(str))
			if err == nil {
				res.SetStatusLine(200, "OK")
				res.AddHeader("Content-Encoding", "gzip")
				res.AddHeader("Content-Type", "text/plain")
				res.body = data
			}
		} else {
			res.SetStatusLine(200, "OK")
			res.AddHeader("Content-Type", "text/plain")
			res.body = []byte(str)
		}
	} else if strings.HasPrefix(req.path, "/user-agent") {
		res.SetStatusLine(200, "OK")
		res.AddHeader("Content-Type", "text/plain")
		res.body = []byte(req.header["user-agent"])
	} else if strings.HasPrefix(req.path, "/files") {
		fileName := req.path[7:]
		if req.method == "GET" {
			str, err := GetFileContent(fileName)
			if err == nil {
				res.SetStatusLine(200, "OK")
				res.AddHeader("Content-Type", "application/octet-stream")
				res.body = str
			}
		} else if req.method == "POST" {
			ok := AddFile(fileName, req.body)
			if ok {
				res.SetStatusLine(201, "Created")
			} else {
				res.SetStatusLine(500, "Internal Server Error")
			}
		}
	} else if req.path == "/" {
		res.SetStatusLine(200, "OK")
	}
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

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go HandleConnection(conn)
	}
}
