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
		outputWeight float64
	}{
		{
			name:   "Simple test 1",
			input:  New([]*WeightedServer{}),
			outputWeight: 0.0,
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
			outputWeight: 0.3,
		},
		{
			name: "Simple test 3",
			input: New([]*WeightedServer{
				{
					URL:    &url.URL{Host: "localhost:3000"},
					Alive:  true,
					Weight: 0.3,
				},
				{
					URL:    &url.URL{Host: "localhost:3001"},
					Alive:  true,
					Weight: 1.2,
				},
			}),
			outputWeight: 1.5,
		},
	}

	for _, tt := range cases {
		tt := tt
		wsp := &WeightedServerPool{}
		wsp.Set(tt.input.servers)
		assert.Equal(t, tt.outputWeight, wsp.Weight, tt.name)
		assert.Equal(t, len(tt.input.servers), len(wsp.servers), tt.name)
	}
}
