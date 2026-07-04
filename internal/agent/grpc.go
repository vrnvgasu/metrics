package agent

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/metadata"

	metricsgrpc "github.com/vrnvgasu/metrics/internal/grpc"
	models "github.com/vrnvgasu/metrics/internal/model"
	pb "github.com/vrnvgasu/metrics/internal/proto"
)

// sendBatchGRPC отправляет батч метрик на сервер по gRPC.
// IP-адрес агента передается в метаданных запроса (ключ x-real-ip).
func (a *Agent) sendBatchGRPC(ctx context.Context, m []*models.Metrics) error {
	req := pb.UpdateMetricsRequest_builder{
		Metrics: metricsgrpc.ToProto(m),
	}.Build()

	if a.realIP != "" {
		md := metadata.New(map[string]string{metricsgrpc.RealIPMetaKey: a.realIP})
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	if _, err := a.grpcClient.UpdateMetrics(ctx, req); err != nil {
		if errors.Is(ctx.Err(), context.Canceled) {
			return nil
		}

		return fmt.Errorf("agent.sendBatchGRPC UpdateMetrics: %w", err)
	}

	return nil
}
