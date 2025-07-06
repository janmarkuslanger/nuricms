package server

import (
	"fmt"
	"net/http"
)

func NewServer() *Server {
	mux := http.NewServeMux()
	return &Server{Mux: mux}
}

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

type HandlerFunc func(ctx Context)

type Server struct {
	Mux *http.ServeMux
}

func (s Server) AddHandler(pattern string, handler HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(Context{
			Writer:  w,
			Request: r,
		})
	})

	finalHandler := Chain(baseHandler, middlewares...)
	s.Mux.Handle(pattern, finalHandler)
}

func (s Server) Handle(pattern string, handler HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	fmt.Println("Serving route(s) with pattern: " + pattern)

	s.AddHandler(pattern, handler, middlewares...)

	if pattern[len(pattern)-1:] != "/" {
		s.AddHandler(pattern+"/", handler, middlewares...)
	}
}
