package main

import "fmt"

type Respond struct {
	version string
	code    int
	msg     string
	header  string
	body    string
}

func NewRespond() *Respond {
	return &Respond{
		version: "HTTP/1.1",
		code:    404,
		msg:     "Not Found",
		header:  "",
		body:    "",
	}
}
func (res *Respond) ToString() string {
	return fmt.Sprintf("%s %d %s\r\n%s\r\n%s", res.version, res.code, res.msg, res.header, res.body)
}
func (res *Respond) OkRespond(header string, body string) {
	res.code = 200
	res.msg = "OK"
	res.header = header
	res.body = body
}
