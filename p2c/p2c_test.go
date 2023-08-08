package p2c_test

import (
	"testing"

	"github.com/mauricioabreu/load-balancingo/p2c"
	"github.com/stretchr/testify/assert"
)

type fakeLoadFetcher struct {
	load int
}

func (f *fakeLoadFetcher) FetchLoad() int {
	return f.load
}

func TestP2CWithOneServerOnly(t *testing.T) {
	b := p2c.New(p2c.NewServer("127.0.0.1", &fakeLoadFetcher{load: 1}))
	assert.Equal(t, b.Next().Address(), "127.0.0.1")
}

func TestP2CWithTwoServers(t *testing.T) {
	b := p2c.New(
		p2c.NewServer("127.0.0.1", &fakeLoadFetcher{load: 1}),
		p2c.NewServer("192.168.0.1", &fakeLoadFetcher{load: 2}),
	)
	assert.Equal(t, b.Next().Address(), "127.0.0.1")
}

func TestRandomDistribution(t *testing.T) {
	b := p2c.New(
		p2c.NewServer("127.0.0.1", &fakeLoadFetcher{load: 1}),
		p2c.NewServer("192.168.0.1", &fakeLoadFetcher{load: 1}),
		p2c.NewServer("192.170.0.1", &fakeLoadFetcher{load: 1}),
	)

	counter := make(map[string]int)
	iterations := 100

	for i := 0; i < iterations; i++ {
		srv := b.Next()
		counter[srv.Address()]++
	}

	expectedCount := iterations / len(b.Servers())
	tolerance := expectedCount / 2

	for addr, count := range counter {
		if count < (expectedCount-tolerance) || count > (expectedCount+tolerance) {
			t.Fatalf("server %s has an unexpected selection count: %d", addr, count)
		}
	}
}
