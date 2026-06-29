package handler

import (
	"crypto/rsa"
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/handler/middleware"
	"github.com/vrnvgasu/metrics/pkg/crypto"
)

//go:embed `templates`
var tmplFS embed.FS

// NewRouter создает gin.Engine с маршрутами и middleware (logger, decrypt, gzip, recovery).
func NewRouter(handler *Handler) (*gin.Engine, error) {
	var privateKey *rsa.PrivateKey
	if handler.cfg.CryptoKey != "" {
		var err error
		privateKey, err = crypto.LoadPrivateKey(handler.cfg.CryptoKey)
		if err != nil {
			return nil, fmt.Errorf("handler.NewRouter LoadPrivateKey: %w", err)
		}
	}

	var trustedSubnet *net.IPNet
	if handler.cfg.TrustedSubnet != "" {
		var err error
		if _, trustedSubnet, err = net.ParseCIDR(handler.cfg.TrustedSubnet); err != nil {
			return nil, fmt.Errorf("handler.NewRouter ParseCIDR: %w", err)
		}
	}

	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(middleware.TrustedSubnet(trustedSubnet))
	r.Use(middleware.Decrypt(privateKey))
	r.Use(middleware.Gzip())
	r.Use(gin.Recovery())

	sub, _ := fs.Sub(tmplFS, "templates")
	r.LoadHTMLFS(http.FS(sub), "*")

	r.GET("/", handler.List)

	valueGroup := r.Group("/value")
	{
		valueGroup.POST("/", handler.Value)
		valueGroup.GET("/:mtype/:name", handler.Find)
	}

	r.POST("/updates", middleware.Hash(handler.cfg), handler.UpdateJSONList)

	updateGroup := r.Group("/update")
	{
		updateGroup.POST("/", handler.UpdateJSON)
		updateGroup.POST("/:mtype/:name/:value", handler.Update)
	}

	r.GET("/ping", handler.Ping)

	return r, nil
}
