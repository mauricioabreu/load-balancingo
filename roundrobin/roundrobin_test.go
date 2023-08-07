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
		roundrobin.NewServer("127.0.0.1").WithWeight(3),
		roundrobin.NewServer("192.168.0.1").WithWeight(2),
		roundrobin.NewServer("192.170.0.1").WithWeight(1),
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

	err := b.Add(roundrobin.NewServer("192.170.0.1"))
	assert.NoError(t, err)
	assert.Equal(t, len(b.Servers()), 3)
}

func TestAddServerWhenAlreadyExists(t *testing.T) {
	b := roundrobin.New(
		roundrobin.NewServer("127.0.0.1"),
	)
	assert.Equal(t, len(b.Servers()), 1)

	err := b.Add(roundrobin.NewServer("127.0.0.1"))
	assert.Equal(t, len(b.Servers()), 1)
	assert.ErrorIs(t, err, roundrobin.ErrServerAlreadyExists)
}

func TestRemoveServer(t *testing.T) {
	b := roundrobin.New(
		roundrobin.NewServer("127.0.0.1"),
		roundrobin.NewServer("192.168.0.1"),
	)
	assert.Equal(t, len(b.Servers()), 2)

	b.Remove("192.168.0.1")
	assert.Equal(t, len(b.Servers()), 1)
}

func TestRemoveServerWhenItDoesNotExist(t *testing.T) {
	b := roundrobin.New(
		roundrobin.NewServer("127.0.0.1"),
	)
	assert.Equal(t, len(b.Servers()), 1)

	err := b.Remove("192.168.0.1")
	assert.ErrorIs(t, err, roundrobin.ErrServerDoesNotExist)
	assert.Equal(t, len(b.Servers()), 1)
}
