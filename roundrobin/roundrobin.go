package roundrobin

import "sync"

type Balancer struct {
	servers []string
	index   int
	m       sync.Mutex
}

func New(servers []string) *Balancer {
	return &Balancer{
		servers: servers,
	}
}

func (b *Balancer) Next() string {
	server := b.servers[b.index]
	b.index = (b.index + 1) % len(b.servers)
	return server
}

func (b *Balancer) Add(host string) {
	b.m.Lock()
	defer b.m.Unlock()

	b.servers = append(b.servers, host)
}

func (b *Balancer) Servers() []string {
	return b.servers
}
