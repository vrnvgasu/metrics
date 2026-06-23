package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTempPEM(t *testing.T, blockType string, data []byte) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.pem")
	require.NoError(t, err)
	require.NoError(t, pem.Encode(f, &pem.Block{Type: blockType, Bytes: data}))
	require.NoError(t, f.Close())
	return f.Name()
}

func genKeyPair(t *testing.T) (pubPath, keyPath string) {
	t.Helper()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	require.NoError(t, err)

	dir := t.TempDir()
	pubPath = filepath.Join(dir, "public.pem")
	keyPath = filepath.Join(dir, "private.pem")

	pf, err := os.Create(pubPath)
	require.NoError(t, err)
	require.NoError(t, pem.Encode(pf, &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}))
	pf.Close()

	kf, err := os.Create(keyPath)
	require.NoError(t, err)
	require.NoError(t, pem.Encode(kf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}))
	kf.Close()

	return pubPath, keyPath
}

func TestLoadPublicKey(t *testing.T) {
	t.Parallel()

	pubPath, _ := genKeyPair(t)

	pub, err := LoadPublicKey(pubPath)
	require.NoError(t, err)
	assert.NotNil(t, pub)
}

func TestLoadPublicKey_Errors(t *testing.T) {
	t.Parallel()

	t.Run("file not found", func(t *testing.T) {
		_, err := LoadPublicKey("/nonexistent/path.pem")
		assert.Error(t, err)
	})

	t.Run("no PEM block", func(t *testing.T) {
		f, err := os.CreateTemp(t.TempDir(), "*.pem")
		require.NoError(t, err)
		f.WriteString("not a pem")
		f.Close()
		_, err = LoadPublicKey(f.Name())
		assert.Error(t, err)
	})

	t.Run("wrong PEM type", func(t *testing.T) {
		path := writeTempPEM(t, "RSA PRIVATE KEY", []byte("data"))
		_, err := LoadPublicKey(path)
		assert.Error(t, err)
	})
}

func TestLoadPrivateKey(t *testing.T) {
	t.Parallel()

	_, keyPath := genKeyPair(t)

	priv, err := LoadPrivateKey(keyPath)
	require.NoError(t, err)
	assert.NotNil(t, priv)
}

func TestLoadPrivateKey_Errors(t *testing.T) {
	t.Parallel()

	t.Run("file not found", func(t *testing.T) {
		_, err := LoadPrivateKey("/nonexistent/path.pem")
		assert.Error(t, err)
	})

	t.Run("no PEM block", func(t *testing.T) {
		f, err := os.CreateTemp(t.TempDir(), "*.pem")
		require.NoError(t, err)
		f.WriteString("not a pem")
		f.Close()
		_, err = LoadPrivateKey(f.Name())
		assert.Error(t, err)
	})

	t.Run("invalid key data", func(t *testing.T) {
		path := writeTempPEM(t, "RSA PRIVATE KEY", []byte("invalid"))
		_, err := LoadPrivateKey(path)
		assert.Error(t, err)
	})
}

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	data := []byte("hello metrics")
	encrypted, err := Encrypt(&priv.PublicKey, data)
	require.NoError(t, err)

	decrypted, err := Decrypt(priv, encrypted)
	require.NoError(t, err)

	assert.Equal(t, data, decrypted)
}
