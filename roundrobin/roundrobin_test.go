package roundrobin_test

import (
	"testing"

	"github.com/mauricioabreu/load-balancingo/roundrobin"
	"github.com/stretchr/testify/assert"
)

func TestSimpleRoundRobin(t *testing.T) {
	b := roundrobin.New(
		roundrobin.NewServer("127.0.0.1"),
		roundrobin.NewServer("192.168.0.1"),
		roundrobin.NewServer("192.170.0.1"),
	)
	assert.Equal(t, b.Next().Address, "127.0.0.1")
	assert.Equal(t, b.Next().Address, "192.168.0.1")
	assert.Equal(t, b.Next().Address, "192.170.0.1")
	assert.Equal(t, b.Next().Address, "127.0.0.1")
}

func TestWeightedRoundRobin(t *testing.T) {
	b := roundrobin.New(
		roundrobin.NewServer("127.0.0.1").WithWeigth(3),
		roundrobin.NewServer("192.168.0.1").WithWeigth(2),
		roundrobin.NewServer("192.170.0.1").WithWeigth(1),
	)
	assert.Equal(t, b.Next().Address, "127.0.0.1")
	assert.Equal(t, b.Next().Address, "192.168.0.1")
	assert.Equal(t, b.Next().Address, "127.0.0.1")
	assert.Equal(t, b.Next().Address, "192.170.0.1")
}

func TestAddServer(t *testing.T) {
	b := roundrobin.New(
		roundrobin.NewServer("127.0.0.1"),
		roundrobin.NewServer("192.168.0.1"),
	)
	assert.Equal(t, len(b.Servers()), 2)

	b.Add(roundrobin.NewServer("192.170.0.1"))
	assert.Equal(t, len(b.Servers()), 3)
}
