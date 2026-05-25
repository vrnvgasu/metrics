package compress

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// реалистичный батч из 30 gauge-метрик
var benchPayload = []byte(`[
{"id":"Alloc","type":"gauge","value":1234567.0},
{"id":"BuckHashSys","type":"gauge","value":2345.0},
{"id":"Frees","type":"gauge","value":98765.0},
{"id":"GCCPUFraction","type":"gauge","value":0.001},
{"id":"GCSys","type":"gauge","value":4096.0},
{"id":"HeapAlloc","type":"gauge","value":1234567.0},
{"id":"HeapIdle","type":"gauge","value":5678901.0},
{"id":"HeapInuse","type":"gauge","value":3456789.0},
{"id":"HeapObjects","type":"gauge","value":12345.0},
{"id":"HeapReleased","type":"gauge","value":2345678.0},
{"id":"HeapSys","type":"gauge","value":9012345.0},
{"id":"LastGC","type":"gauge","value":1700000000000000000.0},
{"id":"Lookups","type":"gauge","value":0.0},
{"id":"MCacheInuse","type":"gauge","value":4800.0},
{"id":"MCacheSys","type":"gauge","value":15600.0},
{"id":"MSpanInuse","type":"gauge","value":123456.0},
{"id":"MSpanSys","type":"gauge","value":130000.0},
{"id":"Mallocs","type":"gauge","value":111111.0},
{"id":"NextGC","type":"gauge","value":2469134.0},
{"id":"NumForcedGC","type":"gauge","value":0.0},
{"id":"NumGC","type":"gauge","value":42.0},
{"id":"OtherSys","type":"gauge","value":987654.0},
{"id":"PauseTotalNs","type":"gauge","value":1234567.0},
{"id":"StackInuse","type":"gauge","value":327680.0},
{"id":"StackSys","type":"gauge","value":327680.0},
{"id":"Sys","type":"gauge","value":12345678.0},
{"id":"TotalAlloc","type":"gauge","value":99999999.0},
{"id":"TotalMemory","type":"gauge","value":17179869184.0},
{"id":"FreeMemory","type":"gauge","value":4294967296.0},
{"id":"PollCount","type":"counter","delta":42}
]`)

func BenchmarkGzipCompress(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		result, err := GzipCompress(benchPayload)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}

func TestGzipCompress(t *testing.T) {
	result, err := GzipCompress(benchPayload)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Less(t, len(result), len(benchPayload))
}
