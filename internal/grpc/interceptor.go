package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const RealIPMetaKey = "x-real-ip"

func realIPFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(RealIPMetaKey)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func TrustedSubnetInterceptor(subnet *net.IPNet) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if subnet == nil {
			return handler(ctx, req)
		}

		ip := net.ParseIP(realIPFromContext(ctx))
		if ip == nil || !subnet.Contains(ip) {
			return nil, status.Error(codes.PermissionDenied, "ip address is not in trusted subnet")
		}

		return handler(ctx, req)
	}
}
