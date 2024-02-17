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
