package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	data := []byte("hello metrics")
	encrypted, err := Encrypt(&priv.PublicKey, data)
	require.NoError(t, err)

	decrypted, err := Decrypt(priv, encrypted)
	require.NoError(t, err)

	assert.Equal(t, data, decrypted)
}
