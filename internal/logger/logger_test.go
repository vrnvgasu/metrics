package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name        string
		level       string
		expectedErr require.ErrorAssertionFunc
	}{
		{name: "valid level info", level: "info", expectedErr: require.NoError},
		{name: "valid level debug", level: "debug", expectedErr: require.NoError},
		{name: "valid level warn", level: "warn", expectedErr: require.NoError},
		{name: "invalid level", level: "dummy", expectedErr: require.Error},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := Initialize(tt.level)
			tt.expectedErr(t, err)
		})
	}
}
