package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/config"
	"github.com/vrnvgasu/metrics/internal/logger"
	models "github.com/vrnvgasu/metrics/internal/model"
	serviceerrors "github.com/vrnvgasu/metrics/internal/service/errors"
	"github.com/vrnvgasu/metrics/pkg/hash"
)

const (
	hashHeader = "HashSHA256"
)

type MetricService interface {
	CreateOrUpdate(context.Context, []*models.Metrics) error
	FindByTypeAndID(ctx context.Context, mtype, id string) (*models.Metrics, error)
	AllMetrics(context.Context) (models.MetricsList, error)
}

type HealthService interface {
	CheckPing(ctx context.Context) error
}

type ResponseError struct {
	Code        string `json:"code"`
	HTTPCode    int    `json:"http_code"`
	UserMessage string `json:"user_message"`
	Error       string `json:"error,omitempty"`
}

type Handler struct {
	MetricService MetricService
	HealthService HealthService

	cfg *config.ServerCnf
}

func NewHandler(m MetricService, h HealthService, cfg *config.ServerCnf) *Handler {
	return &Handler{
		MetricService: m,
		HealthService: h,
		cfg:           cfg,
	}
}

func (h *Handler) responseError(c *gin.Context, err error) {
	if err != nil {
		_ = c.Error(err)
	}

	var serviceError *serviceerrors.ServiceError

	switch {
	case errors.As(err, &serviceError):
		h.parseServiceError(c, serviceError)
	default:
		logger.Log.Errorw("http request",
			"uri", c.Request.RequestURI,
			"method", c.Request.Method,
			"error", err,
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseError{
			HTTPCode:    http.StatusInternalServerError,
			UserMessage: "Unhandled error",
		})
	}
}

func (h *Handler) parseServiceError(c *gin.Context, err *serviceerrors.ServiceError) {
	sourceError := ""
	if err.SourceError != nil {
		sourceError = err.SourceError.Error()
	}

	c.AbortWithStatusJSON(err.HTTPCode, ResponseError{
		Code:        string(err.Type),
		HTTPCode:    err.HTTPCode,
		UserMessage: err.Message,
		Error:       sourceError,
	})
}

func (h *Handler) validateHeaderHashSHA256(c *gin.Context) bool {
	if h.cfg.Key == "" {
		return true
	}

	bodyBytes, err := c.GetRawData()
	if err != nil {
		h.responseError(c, err)

		return false
	}

	if len(bodyBytes) == 0 {
		return true
	}

	signature, err := hash.PrepareHeaderHashSHA256(h.cfg.Key, bodyBytes)
	if err != nil {
		h.responseError(c, err)

		return false
	}

	if c.Request.Header.Get(hashHeader) != signature {
		h.responseError(c, serviceerrors.BadRequestError())

		return false
	}

	return true
}
