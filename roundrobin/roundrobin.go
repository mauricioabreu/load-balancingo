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
	address       string
	weight        int
	currentWeight int
}

// NewServer creates a new server with a default weight of 1.
func NewServer(address string) *server {
	return &server{
		address:       address,
		weight:        1,
		currentWeight: 1,
	}
}

func (s *server) Address() string {
	return s.address
}

func (s *server) Weight() int {
	return s.weight
}

// WithWeight sets the server's weight. If provided weight is 0 or negative, it defaults to 1.
func (s *server) WithWeight(weight int) *server {
	if weight <= 0 {
		weight = 1
	}

	s.weight = weight

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
		totalWeight += server.weight
		server.currentWeight = server.weight
	}

	return &Balancer{
		servers:     servers,
		totalWeight: totalWeight,
	}
}

func (b *Balancer) Algorithm() string {
	return "round_robin"
}

// Next returns the next server in the load balancing rotation using the round-robin algorithm.
// It locks the Balancer's mutex to ensure thread safety and updates the current weight of each server.
// Returns nil if there are no servers available.
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
			b.servers[i].currentWeight += b.servers[i].weight
		}
	}

	return next
}

func (b *Balancer) NextAddress() string {
	next := b.Next()

	if next == nil {
		return ""
	}

	return next.Address()
}

// Add adds a server to the balancer.
func (b *Balancer) Add(s *server) error {
	b.m.Lock()
	defer b.m.Unlock()

	if b.exists(s.address) {
		return ErrServerAlreadyExists
	}

	b.servers = append(b.servers, s)
	b.totalWeight += s.weight

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
		if b.servers[i].address == address {
			b.totalWeight -= b.servers[i].weight
			b.servers = append(b.servers[:i], b.servers[i+1:]...)

			break
		}
	}

	return nil
}

func (b *Balancer) exists(address string) bool {
	for _, s := range b.servers {
		if s.address == address {
			return true
		}
	}

	return false
}

// Servers returns a list of servers currently in the balancer
func (b *Balancer) Servers() []*server {
	return b.servers
}
