package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetServerStatus(t *testing.T) {
	servers := []*Server{
		{
			URL:   &url.URL{Host: "localhost:3000"},
			Alive: true,
		},
		{
			URL:   &url.URL{Host: "localhost:3001"},
			Alive: true,
		},
	}
	serverPool := &ServerPool{
		servers: servers,
		current: 0,
	}

	addr := url.URL{Host: "localhost:3000"}
	serverPool.SetServerStatus(&addr, false)

	assert.Equal(t, false, serverPool.servers[0].Alive)
}

func TestGetNextAlive(t *testing.T) {
	servers := []*Server{
		{
			URL:   &url.URL{Host: "localhost:3000"},
			Alive: true,
		},
		{
			URL:   &url.URL{Host: "localhost:3001"},
			Alive: false,
		},
		{
			URL:   &url.URL{Host: "localhost:3002"},
			Alive: false,
		},
		{
			URL:   &url.URL{Host: "localhost:3003"},
			Alive: true,
		},
	}
	serverPool := &ServerPool{
		servers: servers,
		current: 0,
	}

	s1 := serverPool.GetNextAlive()
	assert.Equal(t, "localhost:3003", s1.URL.Host)

	s2 := serverPool.GetNextAlive()
	assert.Equal(t, "localhost:3000", s2.URL.Host)

	u := url.URL{Host: "localhost:3002"}
	serverPool.SetServerStatus(&u, true)
	s3 := serverPool.GetNextAlive()
	assert.Equal(t, "localhost:3002", s3.URL.Host)
}
