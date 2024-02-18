package main

import (
	"math"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	cases := []struct {
		name         string
		input        *WeightedServerPool
		outputWeight float64
	}{
		{
			name:         "Simple test 1",
			input:        New([]*WeightedServer{}),
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

func TestAdd(t *testing.T) {
	cases := []struct {
		name         string
		input        *WeightedServerPool
		addServer    *WeightedServer
		outputWeight float64
	}{
		{
			name:  "Simple test 1",
			input: New([]*WeightedServer{}),
			addServer: &WeightedServer{
				URL:    &url.URL{Host: "localhost:3000"},
				Alive:  true,
				Weight: 0.3,
			},
			outputWeight: 0.3,
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
			addServer: &WeightedServer{
				URL:    &url.URL{Host: "localhost:3000"},
				Alive:  true,
				Weight: 0.3,
			},
			outputWeight: 0.6,
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
			addServer: &WeightedServer{
				URL:    &url.URL{Host: "localhost:3000"},
				Alive:  true,
				Weight: 0.3,
			},
			outputWeight: 1.8,
		},
	}

	for _, tt := range cases {
		tt := tt
		wsp := New(tt.input.servers)
		wsp.Add(tt.addServer)
		assert.True(t, almostEqual(tt.outputWeight, wsp.Weight), tt.name)
		assert.Equal(t, len(tt.input.servers)+1, len(wsp.servers), tt.name)
	}
}

func TestNew(t *testing.T) {
	testcases := []struct {
		name         string
		input        []*WeightedServer
		outputWeight float64
		outputLen    int
	}{
		{
			name:         "Simple test 1",
			input:        []*WeightedServer{},
			outputWeight: 0.0,
			outputLen:    0,
		},
		{
			name: "Simple test 2",
			input: []*WeightedServer{
				{
					URL:    &url.URL{Host: "localhost:3000"},
					Alive:  true,
					Weight: 0.3,
				},
			},
			outputWeight: 0.3,
			outputLen:    1,
		},
		{
			name: "Simple test 3",
			input: []*WeightedServer{
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
			},
			outputWeight: 1.5,
			outputLen:    2,
		},
	}

	for _, tt := range testcases {
		tt := tt
		wsp := New(tt.input)
		assert.Equal(t, tt.outputWeight, wsp.Weight, tt.name)
		assert.Equal(t, tt.outputLen, len(wsp.servers), tt.name)
	}
}

func TestRemove(t *testing.T) {
	testcases := []struct {
		name         string
		input        *WeightedServerPool
		removeServer *WeightedServer
		outputWeight float64
		outputLen    int
	}{
		{
			name: "Simple test 1",
			input: New([]*WeightedServer{
				{
					URL:    &url.URL{Host: "localhost:3000"},
					Alive:  true,
					Weight: 0.3,
				},
			}),
			removeServer: &WeightedServer{
				URL:    &url.URL{Host: "localhost:3000"},
				Alive:  true,
				Weight: 0.3,
			},
			outputWeight: 0.0,
			outputLen:    0,
		},
		{
			name: "Simple test 2",
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
			removeServer: &WeightedServer{
				URL:    &url.URL{Host: "localhost:3000"},
				Alive:  true,
				Weight: 0.3,
			},
			outputWeight: 1.2,
			outputLen:    1,
		},
	}

	for _, tt := range testcases {
		tt := tt
		wsp := New(tt.input.servers)
		wsp.Remove(tt.removeServer)
		assert.Equal(t, tt.outputWeight, wsp.Weight, tt.name)
		assert.Equal(t, tt.outputLen, len(wsp.servers), tt.name)
	}
}

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
    return math.Abs(a - b) <= float64EqualityThreshold
}
