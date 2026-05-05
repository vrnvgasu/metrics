package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/logger"
	serviceerrors "github.com/vrnvgasu/metrics/internal/service/errors"
)

type Error struct {
	Code        string `json:"code"`
	HTTPCode    int    `json:"http_code"`
	UserMessage string `json:"user_message"`
	Error       string `json:"error,omitempty"`
}

func ResponseError(c *gin.Context, err error) {
	if err != nil {
		c.Error(err)
	}

	var serviceError *serviceerrors.ServiceError

	switch {
	case errors.As(err, &serviceError):
		parseServiceError(c, serviceError)
	default:
		logger.Log.Errorw("http request",
			"uri", c.Request.RequestURI,
			"method", c.Request.Method,
			"error", err,
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			HTTPCode:    http.StatusInternalServerError,
			UserMessage: "Unhandled error",
		})
	}
}

func parseServiceError(c *gin.Context, err *serviceerrors.ServiceError) {
	sourceError := ""
	if err.SourceError != nil {
		sourceError = err.SourceError.Error()
	}

	c.AbortWithStatusJSON(err.HTTPCode, Error{
		Code:        string(err.Type),
		HTTPCode:    err.HTTPCode,
		UserMessage: err.Message,
		Error:       sourceError,
	})
}
