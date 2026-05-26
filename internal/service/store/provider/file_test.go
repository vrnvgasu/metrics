package provider

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/pkg/helper"
)

func TestNewProducer(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		f, err := os.CreateTemp("", "producer_*.json")
		require.NoError(t, err)
		name := f.Name()
		f.Close()
		defer os.Remove(name)

		p, err := NewProducer(name)
		require.NoError(t, err)
		require.NoError(t, p.Close())
	})

	t.Run("invalid path", func(t *testing.T) {
		t.Parallel()
		_, err := NewProducer("/nonexistent/dir/file.json")
		require.Error(t, err)
	})
}

func TestNewConsumer(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		f, err := os.CreateTemp("", "consumer_*.json")
		require.NoError(t, err)
		name := f.Name()
		f.Close()
		defer os.Remove(name)

		c, err := NewConsumer(name)
		require.NoError(t, err)
		require.NoError(t, c.Close())
	})

	t.Run("invalid path", func(t *testing.T) {
		t.Parallel()
		_, err := NewConsumer("/nonexistent/dir/file.json")
		require.Error(t, err)
	})
}

func TestProducerConsumer_WriteRead(t *testing.T) {
	t.Parallel()

	f, err := os.CreateTemp("", "metrics_*.json")
	require.NoError(t, err)
	name := f.Name()
	f.Close()
	defer os.Remove(name)

	list := models.MetricsList{
		{ID: "gauge1", MType: models.Gauge, Value: helper.NewRefFloat64(1.1)},
		{ID: "counter1", MType: models.Counter, Delta: helper.NewRefInt64(5)},
	}

	p, err := NewProducer(name)
	require.NoError(t, err)
	require.NoError(t, p.WriteMetrics(&list))
	require.NoError(t, p.Close())

	c, err := NewConsumer(name)
	require.NoError(t, err)
	result, err := c.ReadMetrics()
	require.NoError(t, err)
	require.NoError(t, c.Close())

	require.Len(t, result, 2)
	assert.Equal(t, "gauge1", result[0].ID)
	assert.Equal(t, 1.1, *result[0].Value)
	assert.Equal(t, "counter1", result[1].ID)
	assert.Equal(t, int64(5), *result[1].Delta)
}

func TestConsumer_ReadEmpty(t *testing.T) {
	t.Parallel()

	f, err := os.CreateTemp("", "empty_*.json")
	require.NoError(t, err)
	name := f.Name()
	f.Close()
	defer os.Remove(name)

	c, err := NewConsumer(name)
	require.NoError(t, err)
	defer c.Close()

	result, err := c.ReadMetrics()
	require.NoError(t, err)
	assert.Empty(t, result)
}
