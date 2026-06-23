package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	t.Parallel()

	s := NewMemStorage()
	assert.NoError(t, s.Ping(t.Context()))
}
