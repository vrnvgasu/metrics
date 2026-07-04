package agent

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/vrnvgasu/metrics/internal/config"
	grpcserver "github.com/vrnvgasu/metrics/internal/grpc"
	mockgrpc "github.com/vrnvgasu/metrics/internal/grpc/mocks"
	models "github.com/vrnvgasu/metrics/internal/model"
	pb "github.com/vrnvgasu/metrics/internal/proto"
	"github.com/vrnvgasu/metrics/pkg/helper"
)

func freeAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := l.Addr().String()
	require.NoError(t, l.Close())
	return addr
}

func TestAgent_sendBatchGRPC(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mockgrpc.NewMockMetricService(ctrl)
	publisher := mockgrpc.NewMockPublisher(ctrl)

	var got models.MetricsList
	done := make(chan struct{})
	service.EXPECT().
		CreateOrUpdate(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, list models.MetricsList) error {
			got = list
			return nil
		})
	// сервер должен увидеть IP агента из метаданных
	publisher.EXPECT().
		Notify(gomock.Any(), gomock.Any(), "10.0.0.5").
		Do(func(_ context.Context, _ []string, _ string) { close(done) })

	addr := freeAddr(t)
	srv, err := grpcserver.NewServer(service, publisher, &config.ServerCnf{GRPCAddress: addr})
	require.NoError(t, err)

	go func() { _ = srv.Run() }()
	t.Cleanup(srv.Stop)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	a := NewAgent(nil, 1)
	a.SetRealIP("10.0.0.5")
	a.SetGRPCClient(pb.NewMetricsClient(conn))

	metrics := []*models.Metrics{
		{ID: "Alloc", MType: models.Gauge, Value: helper.NewRefFloat64(3.0)},
		{ID: "PollCount", MType: models.Counter, Delta: helper.NewRefInt64(5)},
	}
	require.NoError(t, a.sendBatchGRPC(context.Background(), metrics))

	<-done
	require.Len(t, got, 2)
	assert.Equal(t, "Alloc", got[0].ID)
	assert.Equal(t, "PollCount", got[1].ID)
}

func TestAgent_sendBatchGRPC_CanceledCtx(t *testing.T) {
	t.Parallel()

	addr := freeAddr(t)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	a := NewAgent(nil, 1)
	a.SetGRPCClient(pb.NewMetricsClient(conn))

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // контекст отменён до вызова

	metrics := []*models.Metrics{
		{ID: "Alloc", MType: models.Gauge, Value: helper.NewRefFloat64(1.0)},
	}
	assert.NoError(t, a.sendBatchGRPC(ctx, metrics))
}
