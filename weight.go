package main

import (
	"math/rand"
	"net/http/httputil"
	"net/url"
	"sort"
	"sync"
)

const (
	DefaultWeight = 1.0
	BTreeBorder   = 10
)

type WeightedServer struct {
	URL          *url.URL
	Alive        bool
	mu           sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
	Weight       float64
}

type WeightedServerPool struct {
	servers []*WeightedServer
	current int
	Weight  float64
	Rand    *rand.Rand
}

type Option struct {
	DefaultWeight int
	BTreeBorder   int
}

func (s *WeightedServerPool) Len() int {
	return len(s.servers)
}

func (s *WeightedServerPool) Swap(i, j int) {
	s.servers[i], s.servers[j] = s.servers[j], s.servers[i]
}

func (s *WeightedServerPool) Less(i, j int) bool {
	return s.servers[i].Weight < s.servers[j].Weight
}

func (s *WeightedServerPool) Set(list []*WeightedServer) {
	sortedPool := WeightedServerPool{
		servers: list,
	}
	sort.Sort(&sortedPool)

	weightSum := 0.0
	for _, server := range list {
		weightSum += server.Weight
	}
	s.Weight = weightSum
	s.servers = sortedPool.servers
}

func New(list []*WeightedServer) *WeightedServerPool {
	wsp := &WeightedServerPool{
		servers: list,
		Weight:  0,
		Rand:    rand.New(rand.NewSource(0)),
	}

	return wsp
}
