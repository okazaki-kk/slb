package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	cases := []struct {
		name   string
		input  *WeightedServerPool
		output *WeightedServerPool
	}{
		{
			name:   "Simple test 1",
			input:  New([]*WeightedServer{}),
			output: &WeightedServerPool{servers: []*WeightedServer{}, Weight: 0},
		},
		{
			name: "Simple test 2",
			input: New([]*WeightedServer{
				{
					URL:    &url.URL{Host: "localhost:3000"},
					Alive:  true,
					Weight: 0.3,
				},
			}),
			output: &WeightedServerPool{servers: []*WeightedServer{}, Weight: 0.3},
		},
	}

	for _, tt := range cases {
		tt := tt
		wsp := &WeightedServerPool{}
		wsp.Set(tt.input.servers)
		assert.Equal(t, tt.output.Weight, wsp.Weight, tt.name)
		assert.Equal(t, len(tt.input.servers), len(wsp.servers), tt.name)
	}
}
