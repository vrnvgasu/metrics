package retry

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fastCfg = &RetryConfig{MaxRetries: 3, StartRetryInterval: 0, AddRetryPeriod: 0}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()
	require.NotNil(t, cfg)
	assert.Equal(t, 3, cfg.MaxRetries)
	assert.Greater(t, cfg.StartRetryInterval.Nanoseconds(), int64(0))
}

func TestRetryableError(t *testing.T) {
	t.Parallel()

	orig := errors.New("original")
	re := NewRetryableError(orig)

	assert.Equal(t, orig.Error(), re.Error())
	assert.ErrorIs(t, re, orig)
}

func TestRetryWithSettings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		f           func() error
		cfg         *RetryConfig
		expectedErr require.ErrorAssertionFunc
	}{
		{
			name:        "success on first try",
			f:           func() error { return nil },
			cfg:         fastCfg,
			expectedErr: require.NoError,
		},
		{
			name:        "non-retryable error stops immediately",
			f:           func() error { return errors.New("permanent") },
			cfg:         fastCfg,
			expectedErr: require.Error,
		},
		{
			name: "retryable error eventually succeeds",
			f: func() func() error {
				calls := 0
				return func() error {
					calls++
					if calls < 3 {
						return NewRetryableError(errors.New("temp"))
					}
					return nil
				}
			}(),
			cfg:         fastCfg,
			expectedErr: require.NoError,
		},
		{
			name:        "retryable error exhausts retries",
			f:           func() error { return NewRetryableError(errors.New("always fails")) },
			cfg:         fastCfg,
			expectedErr: require.Error,
		},
		{
			name:        "nil config uses default",
			f:           func() error { return nil },
			cfg:         nil,
			expectedErr: require.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := RetryWithSettings(tt.f, tt.cfg)
			tt.expectedErr(t, err)
		})
	}
}

func TestRetry(t *testing.T) {
	t.Parallel()
	err := Retry(func() error { return nil })
	require.NoError(t, err)
}
