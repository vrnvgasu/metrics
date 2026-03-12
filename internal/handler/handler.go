package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	models "github.com/vrnvgasu/metrics/internal/model"
	seriveerrors "github.com/vrnvgasu/metrics/internal/service/errors"
)

type MetricService interface {
	CreateOrUpdate(context.Context, *models.Metrics) error
	FindByTypeAndID(ctx context.Context, mtype, id string) (*models.Metrics, error)
	AllMetrics(context.Context) models.MetricsList
}

type HealthService interface {
	CheckPing(ctx context.Context) error
}

type ResponseError struct {
	Code     string `json:"code"`
	HTTPCode int    `json:"httpCode"`
	Title    string `json:"title"`
	Message  string `json:"message"`
	Error    string `json:"error,omitempty"`
}

type Handler struct {
	MetricService MetricService
	HealthService HealthService
}

func NewHandler(m MetricService, h HealthService) *Handler {
	return &Handler{
		MetricService: m,
		HealthService: h,
	}
}

func (h *Handler) responseError(c *gin.Context, err error) {
	if err != nil {
		_ = c.Error(err)
	}

	var serviceError *seriveerrors.ServiceError

	switch {
	case errors.As(err, &serviceError):
		h.parseServiceError(c, serviceError)
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseError{
			Message:  string(seriveerrors.ErrInternal),
			HTTPCode: http.StatusInternalServerError,
			Title:    "Unhandled error",
		})
	}
}

func (h *Handler) parseServiceError(c *gin.Context, err *seriveerrors.ServiceError) {
	sourceError := ""
	if err.SourceError != nil {
		sourceError = err.SourceError.Error()
	}

	c.AbortWithStatusJSON(err.HTTPCode, ResponseError{
		Code:     string(err.Type),
		HTTPCode: err.HTTPCode,
		Title:    err.Title,
		Message:  err.Message,
		Error:    sourceError,
	})
}
