package agent

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
