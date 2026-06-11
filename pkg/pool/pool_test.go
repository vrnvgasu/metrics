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

	p, err := pool.New(func() *testItem { return &testItem{} })
	require.NoError(t, err)
	obj := p.Get()
	require.NotNil(t, obj)
}

func TestPool_New_NilReturnsError(t *testing.T) {
	t.Parallel()

	_, err := pool.New[*testItem](nil)
	require.Error(t, err)
}

func TestPool_PutResetsState(t *testing.T) {
	t.Parallel()

	p, err := pool.New(func() *testItem {
		return &testItem{}
	})
	require.NoError(t, err)

	obj := p.Get()
	obj.Value = "dirty"
	obj.Count = 42
	p.Put(obj)

	got := p.Get()
	require.Equal(t, testItem{}, *got)
}
