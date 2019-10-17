package share

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNode(t *testing.T) {
	n, _ := NewNode("hello", 0)
	assert.Equal(t, DEFAULT_PORT, n.Port)

	n1, _ := NewNode("hello:30", 0)
	assert.Equal(t, int16(30), n1.Port)

}

func TestParse(t *testing.T) {
	nodes, _ := Parse("hello:1,hello:2")

	assert.Len(t, nodes, 2)
	assert.Equal(t, int16(2), nodes[1].Port)
}
