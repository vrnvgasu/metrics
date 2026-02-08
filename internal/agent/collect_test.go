package agent

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
)

func Test_addPollCount(t *testing.T) {
	t.Parallel()

	a := NewAgent(nil)
	require.Equal(t, a.Metrics.Len(), 0)

	a.addPollCount()
	require.Equal(t, a.Metrics.Len(), 1)
	element := a.Metrics.Front()
	require.NotNil(t, element)

	var delta int64 = 1
	m := element.Value.(models.Metrics)
	require.Equal(t, models.Metrics{
		ID:    pollCount,
		MType: models.Counter,
		Delta: &delta,
	}, m)
	require.Nil(t, m.Value)

	require.Equal(t, int64(1), a.pollCount.Load())
}

func Test_addRandomValue(t *testing.T) {
	t.Parallel()

	a := NewAgent(nil)
	require.Equal(t, a.Metrics.Len(), 0)

	a.addRandomValue()
	require.Equal(t, a.Metrics.Len(), 1)
	element := a.Metrics.Front()
	require.NotNil(t, element)

	mActual := element.Value.(models.Metrics)
	require.Equal(t, randomValue, mActual.ID)
	require.Equal(t, models.Gauge, mActual.MType)
	require.NotNil(t, mActual.Value)
	require.Nil(t, mActual.Delta)
}

func Test_addStatsMetric(t *testing.T) {
	t.Parallel()

	var memStats runtime.MemStats
	a := NewAgent(nil)
	require.Equal(t, a.Metrics.Len(), 0)

	a.addStatsMetric(memStats)
	require.Equal(t, len(gaugesMemStatNames), a.Metrics.Len())

	actualNames := make(map[string]struct{}, a.Metrics.Len())
	for {
		el := a.Metrics.Front()
		if el == nil {
			break
		}
		m := el.Value.(models.Metrics)
		var existInList bool
		for _, name := range gaugesMemStatNames {
			if name == m.ID {
				existInList = true

				_, ok := actualNames[name]
				require.False(t, ok, fmt.Sprintf("name %s the name appeared twice", name))

				actualNames[name] = struct{}{}

				break
			}
		}

		a.Metrics.Remove(el)

		require.True(t, existInList)
	}
}

func TestCollect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	a := NewAgent(nil)
	require.Equal(t, 0, a.Metrics.Len())

	err := a.Collect(ctx, &config.AgentCnf{PollInterval: 1})
	require.NoError(t, err)
	require.Equal(t, len(gaugesMemStatNames)+1+1, a.Metrics.Len())
}
