package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
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

func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		return false
	}

	defer conn.Close()
	return true
}

// HealthCheck pings the backends and update the status
func (s *ServerPool) HealthCheck() {
	for _, b := range s.Servers {
		status := "up"
		alive := isBackendAlive(&b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", &b.URL, status)
	}
}

var serverPool *ServerPool

func main() {
	var serverList string
	var port int
	flag.StringVar(&serverList, "backends", "", "Load balanced backends, use commas to separate")
	flag.IntVar(&port, "port", 3030, "Port to serve")
	flag.Parse()

	if len(serverList) == 0 {
		log.Println("no backends provided, please provide one or more backends to load balance")
		return
	}

	servers := strings.Split(serverList, ",")
	for _, s := range servers {
		serverURL, err := url.Parse(s)
		if err != nil {
			log.Fatal(err)
		}
		serverPool.Servers = append(serverPool.Servers, &Server{
			URL:          *serverURL,
			Alive:        true,
			ReverseProxy: httputil.NewSingleHostReverseProxy(serverURL),
		})
	}

	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
