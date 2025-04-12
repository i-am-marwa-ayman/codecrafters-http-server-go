package main

import (
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
