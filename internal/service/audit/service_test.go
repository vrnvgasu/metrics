package audit_test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/vrnvgasu/metrics/internal/config"
	"github.com/vrnvgasu/metrics/internal/service/audit"
	mockaudit "github.com/vrnvgasu/metrics/internal/service/audit/mocks"
)

func TestPublisher(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockObserver1 := mockaudit.NewMockObserver(ctrl)
	mockObserver2 := mockaudit.NewMockObserver(ctrl)

	mockObserver1.EXPECT().Observe(t.Context(), gomock.Any()).Return(nil)
	mockObserver1.EXPECT().Close()
	mockObserver2.EXPECT().Observe(t.Context(), gomock.Any()).Return(nil)
	mockObserver2.EXPECT().Close()

	publisher := audit.NewEvent()
	publisher.Register(mockObserver1)
	publisher.Register(mockObserver2)

	publisher.Notify(t.Context(), []string{"a"}, "0.0.0.0")
	time.Sleep(time.Millisecond * 100)
	require.NoError(t, publisher.Close())
}

func TestFileAudit(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	tests := []struct {
		name string
		cfg  *config.ServerCnf
	}{
		{
			name: "file audit not create",
			cfg:  &config.ServerCnf{},
		},
		{
			name: "file audit create",
			cfg:  &config.ServerCnf{AuditFile: filepath.Join(dir, "testFileAudit.log")},
		},
		{
			name: "url audit create",
			cfg:  &config.ServerCnf{AuditURL: srv.URL},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			publisher, err := audit.NewAudit(tt.cfg)
			require.NoError(t, err)

			publisher.Notify(t.Context(), []string{"a"}, "0.0.0.0")
			time.Sleep(time.Millisecond * 100)
			require.NoError(t, publisher.Close())

			if tt.cfg.AuditFile != "" {
				require.FileExists(t, tt.cfg.AuditFile)
			}
		})
	}
}

func TestURLAudit(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	a := audit.NewURLAudit(srv.URL)

	err := a.Observe(t.Context(), audit.EventMessage{
		Ts:        time.Now().Unix(),
		Metrics:   []string{"test"},
		IpAddress: "127.0.0.1",
	})
	require.NoError(t, err)
	require.NoError(t, a.Close())
}
