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

func NewServer(address string) *server {
	return &server{
		Address:       address,
		Weight:        1,
		currentWeight: 1,
	}
}

func (s *server) WithWeight(weight int) *server {
	if weight <= 0 {
		weight = 1
	}

	s.Weight = weight

	return s
}

type Balancer struct {
	servers     []*server
	m           sync.Mutex
	totalWeight int
}

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

func (b *Balancer) Add(s *server) error {
	b.m.Lock()
	defer b.m.Unlock()

	if b.Exists(s.Address) {
		return ErrServerAlreadyExists
	}

	b.servers = append(b.servers, s)
	b.totalWeight += s.Weight
	return nil
}

func (b *Balancer) Remove(address string) error {
	b.m.Lock()
	defer b.m.Unlock()

	if !b.Exists(address) {
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

func (b *Balancer) Exists(address string) bool {
	for _, s := range b.servers {
		if s.Address == address {
			return true
		}
	}

	return false
}

func (b *Balancer) Servers() []*server {
	return b.servers
}
