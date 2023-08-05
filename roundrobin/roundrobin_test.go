package roundrobin_test

import (
	"testing"

	"github.com/mauricioabreu/load-balancingo/roundrobin"
	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	b := roundrobin.New([]string{"127.0.0.1", "192.168.0.1", "192.170.0.1"})
	assert.Equal(t, b.Next(), "127.0.0.1")
	assert.Equal(t, b.Next(), "192.168.0.1")
	assert.Equal(t, b.Next(), "192.170.0.1")
	assert.Equal(t, b.Next(), "127.0.0.1")
}

func TestAddServer(t *testing.T) {
	b := roundrobin.New([]string{"127.0.0.1", "192.168.0.1"})
	assert.Equal(t, len(b.Servers()), 2)

	b.Add("192.170.0.1")
	assert.Equal(t, len(b.Servers()), 3)
}
