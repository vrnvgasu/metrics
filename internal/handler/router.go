package handler

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/handler/middleware"
)

//go:embed `templates`
var tmplFS embed.FS

func NewRouter(handler *Handler) *gin.Engine {
	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(middleware.Gzip())

	sub, _ := fs.Sub(tmplFS, "templates")
	r.LoadHTMLFS(http.FS(sub), "*")

	r.GET("/", handler.List)

	valueGroup := r.Group("/value")
	{
		valueGroup.POST("/", handler.Value)
		valueGroup.GET("/:mtype/:name", handler.Find)
	}

	r.POST("/updates", handler.UpdateJSONList)

	updateGroup := r.Group("/update")
	{
		updateGroup.POST("/", handler.UpdateJSON)
		updateGroup.POST("/:mtype/:name/:value", handler.Update)
	}

	r.GET("/ping", handler.Ping)

	return r
}
