package grpc

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/vrnvgasu/metrics/internal/config"
	mockgrpc "github.com/vrnvgasu/metrics/internal/grpc/mocks"
	models "github.com/vrnvgasu/metrics/internal/model"
	pb "github.com/vrnvgasu/metrics/internal/proto"
)

func freeAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := l.Addr().String()
	require.NoError(t, l.Close())
	return addr
}

func TestNewServer(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	service := mockgrpc.NewMockMetricService(ctrl)

	t.Run("no trusted subnet", func(t *testing.T) {
		t.Parallel()
		srv, err := NewServer(service, nil, &config.ServerCnf{GRPCAddress: ":3200"})
		require.NoError(t, err)
		assert.NotNil(t, srv)
	})

	t.Run("with trusted subnet", func(t *testing.T) {
		t.Parallel()
		srv, err := NewServer(service, nil, &config.ServerCnf{
			GRPCAddress:   ":3200",
			TrustedSubnet: "10.0.0.0/8",
		})
		require.NoError(t, err)
		assert.NotNil(t, srv)
	})

	t.Run("invalid CIDR", func(t *testing.T) {
		t.Parallel()
		srv, err := NewServer(service, nil, &config.ServerCnf{
			GRPCAddress:   ":3200",
			TrustedSubnet: "not-a-cidr",
		})
		require.Error(t, err)
		assert.Nil(t, srv)
	})
}

func TestServer_RunStop(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	service := mockgrpc.NewMockMetricService(ctrl)
	service.EXPECT().CreateOrUpdate(gomock.Any(), gomock.Any()).Return(nil)

	addr := freeAddr(t)
	srv, err := NewServer(service, nil, &config.ServerCnf{GRPCAddress: addr})
	require.NoError(t, err)

	runErr := make(chan error, 1)
	go func() { runErr <- srv.Run() }()
	t.Cleanup(srv.Stop)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	client := pb.NewMetricsClient(conn)
	req := pb.UpdateMetricsRequest_builder{
		Metrics: []*pb.Metric{
			pb.Metric_builder{Id: "Alloc", Type: pb.Metric_GAUGE, Value: 1.0}.Build(),
		},
	}.Build()

	_, err = client.UpdateMetrics(context.Background(), req)
	require.NoError(t, err)

	srv.Stop()
	assert.NoError(t, <-runErr)
}

func TestMetricsServer_UpdateMetrics(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	service := mockgrpc.NewMockMetricService(ctrl)
	publisher := mockgrpc.NewMockPublisher(ctrl)

	var savedList models.MetricsList
	service.EXPECT().
		CreateOrUpdate(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, list models.MetricsList) error {
			savedList = list
			return nil
		})
	publisher.EXPECT().Notify(gomock.Any(), []string{"PollCount"}, "10.0.0.1")

	srv := &MetricsServer{service: service, publisher: publisher}

	ctx := ctxWithRealIP("10.0.0.1")
	req := pb.UpdateMetricsRequest_builder{
		Metrics: []*pb.Metric{
			pb.Metric_builder{Id: "PollCount", Type: pb.Metric_COUNTER, Delta: 3}.Build(),
		},
	}.Build()

	resp, err := srv.UpdateMetrics(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	require.Len(t, savedList, 1)
	assert.Equal(t, "PollCount", savedList[0].ID)
	assert.Equal(t, int64(3), *savedList[0].Delta)
}

func TestMetricsServer_UpdateMetrics_ServiceError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	service := mockgrpc.NewMockMetricService(ctrl)
	service.EXPECT().CreateOrUpdate(gomock.Any(), gomock.Any()).Return(errors.New("db down"))

	// publisher не должен вызываться при ошибке сохранения
	srv := &MetricsServer{service: service, publisher: nil}

	req := pb.UpdateMetricsRequest_builder{
		Metrics: []*pb.Metric{
			pb.Metric_builder{Id: "Alloc", Type: pb.Metric_GAUGE, Value: 1.0}.Build(),
		},
	}.Build()

	resp, err := srv.UpdateMetrics(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.Internal, status.Code(err))
}
