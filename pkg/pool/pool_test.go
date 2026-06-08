package pool_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vrnvgasu/metrics/pkg/pool"
)

type testItem struct {
	Value string
	Count int
}

func (t *testItem) Reset() {
	t.Value = ""
	t.Count = 0
}

func TestPool_GetReturnsNewObject(t *testing.T) {
	t.Parallel()

	p := pool.New(func() *testItem { return &testItem{} })
	obj := p.Get()
	require.NotNil(t, obj)
}

func TestPool_PutResetsState(t *testing.T) {
	t.Parallel()

	p := pool.New(func() *testItem {
		return &testItem{}
	})

	obj := p.Get()
	obj.Value = "dirty"
	obj.Count = 42
	p.Put(obj)

	got := p.Get()
	require.Equal(t, testItem{}, *got)
}
