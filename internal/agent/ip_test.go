package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutboundIP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		address   string
		expectErr require.ErrorAssertionFunc
	}{
		{
			name:      "localhost",
			address:   "localhost:8080",
			expectErr: require.NoError,
		},
		{
			name:      "remote address",
			address:   "8.8.8.8:53",
			expectErr: require.NoError,
		},
		{
			name:      "missing port",
			address:   "localhost",
			expectErr: require.Error,
		},
		{
			name:      "empty address",
			address:   "",
			expectErr: require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ip, err := OutboundIP(tt.address)
			tt.expectErr(t, err)

			if err == nil {
				require.NotNil(t, ip)
				assert.NotEmpty(t, ip.String())
			}
		})
	}
}
