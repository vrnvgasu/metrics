package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mockagent "github.com/vrnvgasu/metrics/internal/agent/mocks"
	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/pkg/helper"
)

// okResponse — пустой успешный HTTP-ответ для мока.
func okResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(nil)),
	}
}

func gunzip(t *testing.T, data []byte) []byte {
	t.Helper()
	r, err := gzip.NewReader(bytes.NewReader(data))
	require.NoError(t, err)
	out, err := io.ReadAll(r)
	require.NoError(t, err)
	return out
}

func TestAgent_sendBatch(t *testing.T) {
	t.Parallel()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	metrics := []*models.Metrics{
		{ID: "Alloc", MType: models.Gauge, Value: helper.NewRefFloat64(1.5)},
	}

	tests := []struct {
		name          string
		key           string
		realIP        string
		publicKey     *rsa.PublicKey
		wantHash      bool
		wantRealIP    bool
		wantEncrypted bool
	}{
		{
			name:       "plain batch with real ip",
			realIP:     "10.0.0.5",
			wantRealIP: true,
		},
		{
			name:     "with hash key",
			key:      "secret",
			wantHash: true,
		},
		{
			name:       "empty real ip; no header",
			realIP:     "",
			wantRealIP: false,
		},
		{
			name:          "encrypted batch",
			publicKey:     &priv.PublicKey,
			wantEncrypted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			client := mockagent.NewMockClient(ctrl)

			var gotReq *http.Request
			var gotBody []byte
			client.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
				gotReq = req
				gotBody, _ = io.ReadAll(req.Body)
				return okResponse(), nil
			})

			a := NewAgent(client, 1)
			a.SetRealIP(tt.realIP)
			if tt.publicKey != nil {
				a.SetPublicKey(tt.publicKey)
			}

			cnf := &config.AgentCnf{Address: "localhost:8080", Key: tt.key}
			require.NoError(t, a.sendBatch(metrics, cnf))

			assert.Equal(t, "http://localhost:8080/updates", gotReq.URL.String())
			assert.Equal(t, "application/json", gotReq.Header.Get("Content-Type"))
			assert.Equal(t, "gzip", gotReq.Header.Get("Content-Encoding"))

			if tt.wantHash {
				assert.NotEmpty(t, gotReq.Header.Get(hashHeader))
			} else {
				assert.Empty(t, gotReq.Header.Get(hashHeader))
			}

			if tt.wantRealIP {
				assert.Equal(t, tt.realIP, gotReq.Header.Get(realIPHeader))
			} else {
				assert.Empty(t, gotReq.Header.Get(realIPHeader))
			}

			if tt.wantEncrypted {
				assert.Equal(t, "true", gotReq.Header.Get(encryptedHeader))
				return
			}

			assert.Empty(t, gotReq.Header.Get(encryptedHeader))

			var got []*models.Metrics
			require.NoError(t, json.Unmarshal(gunzip(t, gotBody), &got))
			require.Len(t, got, 1)
			assert.Equal(t, "Alloc", got[0].ID)
		})
	}
}

func TestAgent_sendBatch_ClientError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	client := mockagent.NewMockClient(ctrl)
	client.EXPECT().Do(gomock.Any()).Return(nil, errors.New("network down"))

	a := NewAgent(client, 1)
	err := a.sendBatch(
		[]*models.Metrics{{ID: "Alloc", MType: models.Gauge, Value: helper.NewRefFloat64(1.0)}},
		&config.AgentCnf{Address: "localhost:8080"},
	)
	require.Error(t, err)
}

func TestAgent_SendMetrics(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	client := mockagent.NewMockClient(ctrl)
	// остаток (< batchCount) отправляется ровно одним запросом при закрытии канала
	client.EXPECT().Do(gomock.Any()).DoAndReturn(func(*http.Request) (*http.Response, error) {
		return okResponse(), nil
	}).Times(1)

	a := NewAgent(client, 10)
	for i := range 3 {
		a.pushMetric(models.Metrics{ID: "Alloc", MType: models.Gauge, Value: helper.NewRefFloat64(float64(i))})
	}
	close(a.Metrics)

	require.NoError(t, a.SendMetrics(t.Context(), &config.AgentCnf{Address: "localhost:8080"}))
}

func TestAgent_SendMetrics_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	client := mockagent.NewMockClient(ctrl)
	client.EXPECT().Do(gomock.Any()).Return(nil, errors.New("network down")).AnyTimes()

	a := NewAgent(client, 10)
	a.pushMetric(models.Metrics{ID: "Alloc", MType: models.Gauge, Value: helper.NewRefFloat64(1)})
	close(a.Metrics)

	err := a.SendMetrics(t.Context(), &config.AgentCnf{Address: "localhost:8080"})
	require.Error(t, err)
}

func TestAgent_SetHeaderHashSHA256(t *testing.T) {
	t.Parallel()

	a := NewAgent(nil, 1)
	body := []byte(`{"id":"Alloc","type":"gauge","value":1.0}`)

	tests := []struct {
		name         string
		secret       string
		expectHeader bool
		expectedErr  require.ErrorAssertionFunc
	}{
		{
			name:         "empty secret; skip",
			secret:       "",
			expectHeader: false,
			expectedErr:  require.NoError,
		},
		{
			name:         "valid secret; sets header",
			secret:       "my-secret",
			expectHeader: true,
			expectedErr:  require.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(http.MethodPost, "/", nil)
			require.NoError(t, err)

			err = a.SetHeaderHashSHA256(tt.secret, body, req)
			tt.expectedErr(t, err)

			if tt.expectHeader {
				assert.NotEmpty(t, req.Header.Get("HashSHA256"))
			} else {
				assert.Empty(t, req.Header.Get("HashSHA256"))
			}
		})
	}
}
