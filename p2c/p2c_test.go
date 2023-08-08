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
	p2c := p2c.New(p2c.NewServer("127.0.0.1", &fakeLoadFetcher{load: 1}))
	assert.Equal(t, p2c.Next().Address(), "127.0.0.1")
}

func TestP2CWithTwoServers(t *testing.T) {
	p2c := p2c.New(
		p2c.NewServer("127.0.0.1", &fakeLoadFetcher{load: 1}),
		p2c.NewServer("192.168.0.1", &fakeLoadFetcher{load: 2}),
	)
	assert.Equal(t, p2c.Next().Address(), "127.0.0.1")
}
