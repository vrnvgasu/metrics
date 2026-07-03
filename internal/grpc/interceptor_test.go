package grpc

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func mustCIDR(t *testing.T, cidr string) *net.IPNet {
	t.Helper()
	_, subnet, err := net.ParseCIDR(cidr)
	require.NoError(t, err)
	return subnet
}

func ctxWithRealIP(ip string) context.Context {
	if ip == "" {
		return context.Background()
	}
	md := metadata.New(map[string]string{RealIPMetaKey: ip})
	return metadata.NewIncomingContext(context.Background(), md)
}

func TestTrustedSubnetInterceptor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		subnet      *net.IPNet
		ip          string
		wantCalled  bool
		wantErrCode codes.Code
	}{
		{
			name:       "nil subnet passes through",
			subnet:     nil,
			ip:         "",
			wantCalled: true,
		},
		{
			name:       "ip in subnet allowed",
			subnet:     mustCIDR(t, "10.0.0.0/8"),
			ip:         "10.1.2.3",
			wantCalled: true,
		},
		{
			name:        "ip out of subnet denied",
			subnet:      mustCIDR(t, "10.0.0.0/8"),
			ip:          "192.168.0.1",
			wantCalled:  false,
			wantErrCode: codes.PermissionDenied,
		},
		{
			name:        "missing ip denied",
			subnet:      mustCIDR(t, "10.0.0.0/8"),
			ip:          "",
			wantCalled:  false,
			wantErrCode: codes.PermissionDenied,
		},
		{
			name:        "malformed ip denied",
			subnet:      mustCIDR(t, "10.0.0.0/8"),
			ip:          "not-an-ip",
			wantCalled:  false,
			wantErrCode: codes.PermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			called := false
			handler := func(_ context.Context, _ any) (any, error) {
				called = true
				return "ok", nil
			}

			interceptor := TrustedSubnetInterceptor(tt.subnet)
			resp, err := interceptor(ctxWithRealIP(tt.ip), nil, &grpc.UnaryServerInfo{}, handler)

			assert.Equal(t, tt.wantCalled, called)

			if tt.wantErrCode != codes.OK {
				require.Error(t, err)
				assert.Equal(t, tt.wantErrCode, status.Code(err))
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "ok", resp)
		})
	}
}

func TestRealIPFromContext(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "1.2.3.4", realIPFromContext(ctxWithRealIP("1.2.3.4")))
	assert.Empty(t, realIPFromContext(context.Background()))
}
