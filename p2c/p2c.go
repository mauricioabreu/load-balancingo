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
	load      int
	lock      sync.RWMutex
}

// NewServer creates a new server.
func NewServer(address string, lf LoadFetcher) *server {
	return &server{
		address:   address,
		fetchLoad: lf,
		load:      0,
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
	s.load = s.fetchLoad.FetchLoad()

	go func() {
		for range ticker.C {
			newLoad := s.fetchLoad.FetchLoad()

			s.lock.Lock()
			s.load = newLoad
			s.lock.Unlock()
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

// Next returns the next server to be used for load balancing using the power of two choices algorithm.
// It locks the mutex to avoid race conditions and returns the server with the lowest load.
// If there is only one server available, it returns that server.
// If there are two or more servers available, it randomly selects two servers and returns the one with the lowest load.
func (b *Balancer) Next() *server {
	b.m.Lock()
	defer b.m.Unlock()

	if len(b.servers) == 1 {
		return b.servers[0]
	}

	srv1, srv2 := b.getRandomServers()

	if srv1.load < srv2.load {
		return srv1
	}

	return srv2
}

func (b *Balancer) NextAddress() string {
	return b.Next().address
}

func (b *Balancer) Servers() []*server {
	return b.servers
}

func (b *Balancer) getRandomServers() (srv1, srv2 *server) {
	if len(b.servers) == 2 { //nolint:gomnd // 2 seems pretty obvious
		return b.servers[0], b.servers[1]
	}

	idx1 := rand.Intn(len(b.servers))
	idx2 := rand.Intn(len(b.servers))

	for idx1 == idx2 {
		idx2 = rand.Intn(len(b.servers))
	}

	return b.servers[idx1], b.servers[idx2]
}
