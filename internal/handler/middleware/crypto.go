package middleware

import (
	"bytes"
	"crypto/rsa"
	"io"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/handler/response"
	serviceerrors "github.com/vrnvgasu/metrics/internal/service/errors"
	"github.com/vrnvgasu/metrics/pkg/crypto"
)

const cryptoHeader = "X-Encrypted"

// Decrypt — middleware для расшифровки тела запроса приватным RSA-ключом.
// Пропускает запрос без изменений, если заголовок X-Encrypted не установлен или ключ не задан.
func Decrypt(privateKey *rsa.PrivateKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		if privateKey == nil || c.Request.Header.Get(cryptoHeader) == "" {
			c.Next()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			response.ResponseError(c, serviceerrors.BadRequestError())
			return
		}

		decrypted, err := crypto.Decrypt(privateKey, body)
		if err != nil {
			response.ResponseError(c, err)
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(decrypted))
		c.Request.ContentLength = int64(len(decrypted))
		c.Next()
	}
}
