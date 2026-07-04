// Package grpc реализует gRPC-транспорт сервера метрик: сервис Metrics,
// interceptor проверки доверенной подсети и конвертеры proto <-> доменная модель.
package grpc

import (
	models "github.com/vrnvgasu/metrics/internal/model"
	pb "github.com/vrnvgasu/metrics/internal/proto"
)

// ToModels конвертирует protobuf-метрики в доменную модель.
func ToModels(in []*pb.Metric) models.MetricsList {
	list := make(models.MetricsList, 0, len(in))
	for _, m := range in {
		metric := models.Metrics{ID: m.GetId()}

		switch m.GetType() {
		case pb.Metric_COUNTER:
			metric.MType = models.Counter
			delta := m.GetDelta()
			metric.Delta = &delta
		case pb.Metric_GAUGE:
			metric.MType = models.Gauge
			value := m.GetValue()
			metric.Value = &value
		}

		list = append(list, metric)
	}

	return list
}

// ToProto конвертирует доменные метрики в protobuf.
func ToProto(list []*models.Metrics) []*pb.Metric {
	result := make([]*pb.Metric, 0, len(list))
	for _, m := range list {
		b := pb.Metric_builder{Id: m.ID}

		switch m.MType {
		case models.Counter:
			b.Type = pb.Metric_COUNTER
			if m.Delta != nil {
				b.Delta = *m.Delta
			}
		case models.Gauge:
			b.Type = pb.Metric_GAUGE
			if m.Value != nil {
				b.Value = *m.Value
			}
		}

		result = append(result, b.Build())
	}

	return result
}
