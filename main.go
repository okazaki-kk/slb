package main

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	URL          url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

type ServerPool struct {
	Servers []*Server
	counter uint64
}

func (s *Server) IsAlive() bool {
	s.mux.RLock()
	defer s.mux.Unlock()
	return s.Alive
}

func (s *Server) SetAlive(alive bool) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.Alive = alive
}

func main() {
}
