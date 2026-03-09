package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/logger"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		gin.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		c.Writer = &loggingResponseWriter{
			ResponseWriter: c.Writer, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		c.Next()

		duration := time.Since(start)

		logger.Log.Infow("http request",
			"uri", c.Request.RequestURI,
			"method", c.Request.Method,
			"duration", duration,
			"status", responseData.status,
			"size", responseData.size,
		)
	}
}
