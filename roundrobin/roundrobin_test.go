package roundrobin_test

import (
	"testing"

	"github.com/mauricioabreu/load-balancingo/roundrobin"
	"github.com/stretchr/testify/assert"
)

func TestRoundRobin(t *testing.T) {
	b := roundrobin.New([]string{"127.0.0.1", "192.168.0.1", "192.170.0.1"})
	assert.Equal(t, b.Next(), "127.0.0.1")
	assert.Equal(t, b.Next(), "192.168.0.1")
	assert.Equal(t, b.Next(), "192.170.0.1")
	assert.Equal(t, b.Next(), "127.0.0.1")
}
