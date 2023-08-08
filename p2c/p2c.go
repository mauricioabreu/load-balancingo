package p2c

import (
	"math/rand"
	"sync"
	"time"
)

const (
	maxLoad       = 101
	intervalFetch = 10 * time.Second
)

type server struct {
	address   string
	fetchLoad LoadFetcher
	loadCh    chan int
}

// NewServer creates a new server.
func NewServer(address string, lf LoadFetcher) *server {
	return &server{
		address:   address,
		fetchLoad: lf,
		loadCh:    make(chan int, 1),
	}
}

func (s *server) Address() string {
	return s.address
}

type LoadFetcher interface {
	FetchLoad() int
}

type RandomLoadFetcher struct{}

func (r *RandomLoadFetcher) FetchLoad() int {
	return rand.Intn(maxLoad)
}

func (s *server) StartFetchLoader(interval time.Duration) {
	ticker := time.NewTicker(interval)
	s.loadCh <- s.fetchLoad.FetchLoad()

	go func() {
		for range ticker.C {
			s.loadCh <- s.fetchLoad.FetchLoad()
		}
	}()
}

type Balancer struct {
	servers []*server
	m       sync.Mutex
}

func New(servers ...*server) *Balancer {
	for _, s := range servers {
		s.StartFetchLoader(intervalFetch)
	}

	return &Balancer{servers: servers}
}

func (b *Balancer) Next() *server {
	b.m.Lock()
	defer b.m.Unlock()

	if len(b.servers) == 1 {
		return b.servers[0]
	}

	srv1, srv2 := b.getRandomServers()

	load1 := <-srv1.loadCh
	load2 := <-srv2.loadCh

	if load1 < load2 {
		return srv1
	}

	return srv2
}

func (b *Balancer) getRandomServers() (srv1, srv2 *server) {
	if len(b.servers) == 2 { //nolint:gomnd // 2 seems pretty obvious
		return b.servers[0], b.servers[1]
	}

	idx1 := rand.Intn(len(b.servers))
	idx2 := rand.Intn(len(b.servers))

	if idx1 == idx2 {
		idx2 = rand.Intn(len(b.servers))
	}

	return b.servers[idx1], b.servers[idx2]
}
