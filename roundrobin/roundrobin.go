package roundrobin

import "sync"

type server struct {
	Address       string
	Weigth        int
	CurrentWeigth int
}

func NewServer(address string) *server {
	return &server{
		Address:       address,
		Weigth:        1,
		CurrentWeigth: 0,
	}
}

func (s *server) WithWeigth(weigth int) *server {
	if weigth <= 0 {
		panic("weight must be greater than 0")
	}
	s.Weigth = weigth
	return s
}

type Balancer struct {
	servers []*server
	m       sync.Mutex
}

func New(servers ...*server) *Balancer {
	return &Balancer{
		servers: servers,
	}
}

func (b *Balancer) Next() *server {
	b.m.Lock()
	defer b.m.Unlock()

	maxIndex := -1
	totalWeight := 0

	for idx, server := range b.servers {
		b.servers[idx].CurrentWeigth += server.Weigth

		totalWeight += server.Weigth

		if maxIndex == -1 || b.servers[maxIndex].CurrentWeigth < b.servers[idx].CurrentWeigth {
			maxIndex = idx
		}
	}

	b.servers[maxIndex].CurrentWeigth -= totalWeight

	return b.servers[maxIndex]
}

func (b *Balancer) Add(s *server) {
	b.m.Lock()
	defer b.m.Unlock()

	b.servers = append(b.servers, s)
}

func (b *Balancer) Servers() []*server {
	return b.servers
}
