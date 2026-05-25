package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrepareHeaderHashSHA256(t *testing.T) {
	t.Parallel()

	data := []byte(`[{"id":"Alloc","type":"gauge","value":1234567.0}]`)

	h1, err := PrepareHeaderHashSHA256("secret", data)
	require.NoError(t, err)
	require.NotEmpty(t, h1)

	h2, err := PrepareHeaderHashSHA256("secret", data)
	require.NoError(t, err)
	require.Equal(t, h1, h2)

	h3, err := PrepareHeaderHashSHA256("other-secret", data)
	require.NoError(t, err)
	require.NotEqual(t, h1, h3)
}
