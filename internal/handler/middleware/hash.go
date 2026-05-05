package middleware

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/config"
	"github.com/vrnvgasu/metrics/internal/handler/response"
	serviceerrors "github.com/vrnvgasu/metrics/internal/service/errors"
	"github.com/vrnvgasu/metrics/pkg/hash"
)

const (
	hashHeader = "HashSHA256"
)

func Hash(cnf *config.ServerCnf) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cnf.Key == "" {
			c.Next()

			return
		}

		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		signature, err := hash.PrepareHeaderHashSHA256(cnf.Key, bodyBytes)
		if err != nil {
			response.ResponseError(c, err)

			return
		}

		hh := c.Request.Header.Get(hashHeader)
		_ = hh
		if c.Request.Header.Get(hashHeader) != signature {
			response.ResponseError(c, serviceerrors.BadRequestError())

			return
		}

		c.Next()
	}
}
