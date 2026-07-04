package handler

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vrnvgasu/metrics/internal/config"
	"github.com/vrnvgasu/metrics/internal/repository/mem"
	"github.com/vrnvgasu/metrics/internal/service/metric"
)

func TestNewRouter_TrustedSubnet(t *testing.T) {
	t.Parallel()

	s := metric.NewService(mem.NewMemStorage())

	t.Run("valid CIDR", func(t *testing.T) {
		t.Parallel()
		h := NewHandler(s, nil, newTestPublisher(t), &config.ServerCnf{TrustedSubnet: "192.168.0.0/24"})
		_, err := NewRouter(h)
		require.NoError(t, err)
	})

	t.Run("invalid CIDR", func(t *testing.T) {
		t.Parallel()
		h := NewHandler(s, nil, newTestPublisher(t), &config.ServerCnf{TrustedSubnet: "not-a-cidr"})
		_, err := NewRouter(h)
		require.Error(t, err)
	})
}

func TestNewRouter_CryptoKey(t *testing.T) {
	t.Parallel()

	s := metric.NewService(mem.NewMemStorage())

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	keyPath := filepath.Join(t.TempDir(), "private.pem")
	f, err := os.Create(keyPath)
	require.NoError(t, err)
	require.NoError(t, pem.Encode(f, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	}))
	require.NoError(t, f.Close())

	t.Run("valid key", func(t *testing.T) {
		t.Parallel()
		h := NewHandler(s, nil, newTestPublisher(t), &config.ServerCnf{CryptoKey: keyPath})
		_, err := NewRouter(h)
		require.NoError(t, err)
	})

	t.Run("missing key file", func(t *testing.T) {
		t.Parallel()
		h := NewHandler(s, nil, newTestPublisher(t), &config.ServerCnf{CryptoKey: "/nonexistent.pem"})
		_, err := NewRouter(h)
		require.Error(t, err)
	})
}
