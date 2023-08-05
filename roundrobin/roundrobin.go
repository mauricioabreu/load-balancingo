package roundrobin

type Balancer struct {
	servers []string
	index   int
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
