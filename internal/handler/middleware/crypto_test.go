package middleware

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vrnvgasu/metrics/pkg/crypto"
)

func testDecryptRouter(priv *rsa.PrivateKey) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Decrypt(priv))
	r.POST("/", func(c *gin.Context) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(c.Request.Body)
		c.String(http.StatusOK, buf.String())
	})
	return r
}

func TestDecrypt_NilKey(t *testing.T) {
	t.Parallel()

	r := testDecryptRouter(nil)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("plain"))
	req.Header.Set(cryptoHeader, "true")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "plain", w.Body.String())
}

func TestDecrypt_NoHeader(t *testing.T) {
	t.Parallel()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	r := testDecryptRouter(priv)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("plain"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "plain", w.Body.String())
}

func TestDecrypt_ValidBody(t *testing.T) {
	t.Parallel()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	plaintext := []byte("secret metrics data")
	encrypted, err := crypto.Encrypt(&priv.PublicKey, plaintext)
	require.NoError(t, err)

	r := testDecryptRouter(priv)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(encrypted))
	req.Header.Set(cryptoHeader, "true")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(plaintext), w.Body.String())
}

func TestDecrypt_InvalidBody(t *testing.T) {
	t.Parallel()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	r := testDecryptRouter(priv)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("not encrypted"))
	req.Header.Set(cryptoHeader, "true")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}
