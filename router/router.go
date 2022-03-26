package router

import (
	"context"
	"net/http"
	"replace-url-gin/config"
	"replace-url-gin/handler"

	"github.com/gin-gonic/gin"
)

type router struct {
	ctx     context.Context
	config  *config.Config
	route   *gin.Engine
	handler handler.Handler
}

func Register(ctx context.Context, conf *config.Config, route *gin.Engine) Router {
	return &router{
		ctx:     ctx,
		config:  conf,
		route:   route,
		handler: handler.NewHandler(ctx, conf),
	}
}

func (r router) All() {
	r.BaseRouter()
	r.NotFound()
}

func (r router) NotFound() {
	r.route.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Page not found",
		})
	})
}
