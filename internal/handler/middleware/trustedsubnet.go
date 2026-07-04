package middleware

import (
	"net"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/handler/response"
	serviceerrors "github.com/vrnvgasu/metrics/internal/service/errors"
)

const realIPHeader = "X-Real-IP"

// TrustedSubnet — middleware, ограничивающее доступ агентам из доверенной подсети.
func TrustedSubnet(subnet *net.IPNet) gin.HandlerFunc {
	return func(c *gin.Context) {
		if subnet == nil {
			c.Next()
			return
		}

		ip := net.ParseIP(c.Request.Header.Get(realIPHeader))
		if ip == nil || !subnet.Contains(ip) {
			response.ResponseError(c, serviceerrors.ForbiddenError())
			return
		}

		c.Next()
	}
}
