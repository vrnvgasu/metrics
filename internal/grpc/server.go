package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
	pb "github.com/vrnvgasu/metrics/internal/proto"
)

//go:generate mockgen -destination=./mocks/metric.go . MetricService
type MetricService interface {
	CreateOrUpdate(context.Context, models.MetricsList) error
}

//go:generate mockgen -destination=./mocks/publisher.go . Publisher
type Publisher interface {
	Notify(ctx context.Context, metricIDList []string, ip string)
}

type MetricsServer struct {
	pb.UnimplementedMetricsServer

	service   MetricService
	publisher Publisher
}

// UpdateMetrics принимает батч метрик и сохраняет их в хранилище.
func (s *MetricsServer) UpdateMetrics(
	ctx context.Context, in *pb.UpdateMetricsRequest,
) (*pb.UpdateMetricsResponse, error) {
	list := ToModels(in.GetMetrics())

	if err := s.service.CreateOrUpdate(ctx, list); err != nil {
		return nil, status.Errorf(codes.Internal, "could not save metrics: %v", err)
	}

	if s.publisher != nil {
		s.publisher.Notify(ctx, list.IDList(), realIPFromContext(ctx))
	}

	return pb.UpdateMetricsResponse_builder{}.Build(), nil
}

type Server struct {
	srv  *grpc.Server
	addr string
}

func NewServer(service MetricService, publisher Publisher, cnf *config.ServerCnf) (*Server, error) {
	var opts []grpc.ServerOption

	if cnf.TrustedSubnet != "" {
		_, subnet, err := net.ParseCIDR(cnf.TrustedSubnet)
		if err != nil {
			return nil, fmt.Errorf("grpc.NewServer ParseCIDR: %w", err)
		}
		opts = append(opts, grpc.UnaryInterceptor(TrustedSubnetInterceptor(subnet)))
	}

	srv := grpc.NewServer(opts...)
	pb.RegisterMetricsServer(srv, &MetricsServer{service: service, publisher: publisher})

	return &Server{srv: srv, addr: cnf.GRPCAddress}, nil
}

func (s *Server) Run() error {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("grpc.Server Run Listen: %w", err)
	}

	if err = s.srv.Serve(listen); err != nil {
		return fmt.Errorf("grpc.Server Run Serve: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
