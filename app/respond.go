package main

import (
	"bytes"
	"fmt"
)

type Respond struct {
	version string
	code    int
	msg     string
	header  map[string]string
	body    []byte
}

func NewRespond() *Respond {
	return &Respond{
		version: "HTTP/1.1",
		code:    404,
		msg:     "Not Found",
		header:  make(map[string]string),
		body:    nil,
	}
}
func (res *Respond) Serialize() []byte {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%s %d %s\r\n", res.version, res.code, res.msg))
	if res.body != nil && res.header["Content-Length"] == "" {
		res.header["Content-Length"] = fmt.Sprintf("%d", len(res.body))
	}
	for k, v := range res.header {
		b.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	b.WriteString("\r\n")
	if res.body != nil {
		b.Write(res.body)
	}
	return b.Bytes()
}
func (res *Respond) AddHeader(key string, val string) {
	res.header[key] = val
}
func (res *Respond) SetStatusLine(code int, msg string) {
	res.code = code
	res.msg = msg
}
