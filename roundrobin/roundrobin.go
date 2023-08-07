// Package roundrobin provides a simple and weighted implementation of the round robin algorithm.
package roundrobin

import (
	"errors"
	"sync"
)

var (
	ErrServerAlreadyExists = errors.New("server already exists")
	ErrServerDoesNotExist  = errors.New("server does not exist")
)

type server struct {
	Address       string
	Weight        int
	currentWeight int
}

// NewServer creates a new server with a default weight of 1.
func NewServer(address string) *server {
	return &server{
		Address:       address,
		Weight:        1,
		currentWeight: 1,
	}
}

// WithWeight sets the server's weight. If provided weight is 0 or negative, it defaults to 1.
func (s *server) WithWeight(weight int) *server {
	if weight <= 0 {
		weight = 1
	}

	s.Weight = weight

	return s
}

// Balancer is a weighted round robin balancer for servers.
type Balancer struct {
	servers     []*server
	m           sync.Mutex
	totalWeight int
}

// New creates a new balancer with the provided servers.
func New(servers ...*server) *Balancer {
	totalWeight := 0
	for _, server := range servers {
		totalWeight += server.Weight
		server.currentWeight = server.Weight
	}

	return &Balancer{
		servers:     servers,
		totalWeight: totalWeight,
	}
}

// Next selects the next server based on the weighted round robin algorithm.
func (b *Balancer) Next() *server {
	b.m.Lock()
	defer b.m.Unlock()

	var next *server

	for _, s := range b.servers {
		if next == nil || s.currentWeight > next.currentWeight {
			next = s
		}
	}

	if next != nil {
		next.currentWeight -= b.totalWeight

		for i := range b.servers {
			b.servers[i].currentWeight += b.servers[i].Weight
		}
	}

	return next
}

// Add adds a server to the balancer.
func (b *Balancer) Add(s *server) error {
	b.m.Lock()
	defer b.m.Unlock()

	if b.exists(s.Address) {
		return ErrServerAlreadyExists
	}

	b.servers = append(b.servers, s)
	b.totalWeight += s.Weight

	return nil
}

// Remove removes a server from the balancer by its address.
func (b *Balancer) Remove(address string) error {
	b.m.Lock()
	defer b.m.Unlock()

	if !b.exists(address) {
		return ErrServerDoesNotExist
	}

	for i := range b.servers {
		if b.servers[i].Address == address {
			b.totalWeight -= b.servers[i].Weight
			b.servers = append(b.servers[:i], b.servers[i+1:]...)

			break
		}
	}

	return nil
}

func (b *Balancer) exists(address string) bool {
	for _, s := range b.servers {
		if s.Address == address {
			return true
		}
	}

	return false
}

// Servers returns a list of servers currently in the balancer
func (b *Balancer) Servers() []*server {
	return b.servers
}
