package handler

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed `templates`
var tmplFS embed.FS

func NewRouter(handler *Handler) *gin.Engine {
	r := gin.New()

	sub, _ := fs.Sub(tmplFS, "templates")
	r.LoadHTMLFS(http.FS(sub), "*")

	r.POST("/update/:mtype/:name/:value", handler.Update)
	r.GET("/value/:mtype/:name", handler.Find)
	r.GET("/", handler.List)

	return r
}
