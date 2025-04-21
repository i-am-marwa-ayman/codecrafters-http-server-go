package main

import (
	"strings"
)

type handlerFunc func(*Request) *Respond

type Router struct {
	routes map[string]map[string]handlerFunc
}

func NewRouter() *Router {
	return &Router{make(map[string]map[string]handlerFunc)}
}
func (r *Router) AddRoute(method string, path string, hunhandler handlerFunc) {
	if _, ok := r.routes[method]; !ok {
		r.routes[method] = make(map[string]handlerFunc)
	}
	r.routes[method][path] = hunhandler
}

func (r *Router) Handle(req *Request) *Respond {
	methodRoute, ok := r.routes[req.method]
	if !ok {
		return NewRespond()
	}
	if handler, ok := methodRoute[req.path]; ok {
		return handler(req)
	}
	for path, handler := range methodRoute {
		if strings.HasPrefix(req.path, path) && path != "/" {
			return handler(req)
		}
	}
	return NewRespond()
}

func RootHandler(req *Request) *Respond {
	res := NewRespond()
	res.SetStatusLine(200, "OK")
	return res
}
func EchoHandler(req *Request) *Respond {
	res := NewRespond()
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
	return res
}
func UserAgentHandler(req *Request) *Respond {
	res := NewRespond()
	res.SetStatusLine(200, "OK")
	res.AddHeader("Content-Type", "text/plain")
	res.body = []byte(req.header["user-agent"])
	return res
}
func FilesPostHandler(req *Request) *Respond {
	res := NewRespond()
	fileName := req.path[7:]
	ok := AddFile(fileName, req.body)
	if ok {
		res.SetStatusLine(201, "Created")
	} else {
		res.SetStatusLine(500, "Internal Server Error")
	}
	return res
}
func FilesGetHandler(req *Request) *Respond {
	res := NewRespond()
	fileName := req.path[7:]
	str, err := GetFileContent(fileName)
	if err == nil {
		res.SetStatusLine(200, "OK")
		res.AddHeader("Content-Type", "application/octet-stream")
		res.body = str
	}
	return res
}
