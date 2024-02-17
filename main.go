package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

type key int

const (
	retry key = iota
)

// Server holds information about a backend server
type Server struct {
	URL          *url.URL
	Alive        bool
	mu           sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func (s *Server) SetAlive(alive bool) {
	s.mu.Lock()
	s.Alive = alive
	s.mu.Unlock()
}

func (s *Server) IsAlive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Alive
}

// ServerPool holds information about reachable servers
type ServerPool struct {
	servers []*Server
	current int
}

// NextIndex returns the index of the next server to use
func (s *ServerPool) NextIndex() int {
	return (s.current + 1) % len(s.servers)
}

// GetNextAlive returns the next alive server to use and sets the current index
func (s *ServerPool) GetNextAlive() *Server {
	next := s.NextIndex()
	l := len(s.servers) + next
	for i := next; i < l; i++ {
		server := s.servers[i%len(s.servers)]
		if server.IsAlive() {
			s.current = i % len(s.servers)
			return server
		}
	}
	return nil
}

func (s *ServerPool) SetServerStatus(url *url.URL, alive bool) {
	for _, server := range s.servers {
		if server.URL.String() == url.String() {
			server.SetAlive(alive)
			break
		}
	}
}

// AddServer adds a server to the server pool
func (s *ServerPool) AddServer(server *Server) {
	s.servers = append(s.servers, server)
}

func loadbalance(w http.ResponseWriter, r *http.Request) {
	server := serverPool.GetNextAlive()
	if server != nil {
		// proxy request
		server.ReverseProxy.ServeHTTP(w, r)
	}
	http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
}

func isServerAlive(u *url.URL) bool {
	conn, err := net.DialTimeout("tcp", u.Host, 2*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func healthCheck() {
	t := time.NewTicker(time.Minute)
	for range t.C {
		for _, s := range serverPool.servers {
			alive := isServerAlive(s.URL)
			s.SetAlive(alive)
		}
	}
}

// GetRetryFromContext returns the number of retries from the request context
func getRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(retry).(int); ok {
		return retry
	}
	return 0
}


var serverPool ServerPool

func main() {
	var serverList string
	var port int
	flag.StringVar(&serverList, "backends", "http://localhost:8080,http://localhost:8081", "Load balanced backends")
	flag.IntVar(&port, "port", 9090, "Port to serve on")
	flag.Parse()

	servers := strings.Split(serverList, ",")
	for _, se := range servers {
		serverURL, err := url.Parse(se)
		if err != nil {
			panic(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(serverURL)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
			retryCount := getRetryFromContext(r)
			if retryCount < 3 {
				select {
				case <-time.After(1 * time.Second):
					ctx := context.WithValue(r.Context(), retry, retryCount+1)
					proxy.ServeHTTP(w, r.WithContext(ctx))
				}
				return
			}

			// after 3 retries, mark this server as down
			serverPool.SetServerStatus(serverURL, false)
		}

		server := &Server{
			URL:   serverURL,
			Alive: true,
		}
		serverPool.AddServer(server)
	}

	go healthCheck()

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(loadbalance),
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
