package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
}

func Register(engine *gin.Engine) {
	router := &Router{
		Engine: engine,
	}
	router.detect()
	router.index()
	router.health()
}

func (r *Router) detect() {
	r.HEAD("/", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})
}

func (r *Router) index() {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/formatify")
	})
}

func (r *Router) health() {
	r.GET("/health", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", map[string]interface{}{
			"now": time.Now(),
			"abc": time.Date(2024, 05, 01, 0, 0, 0, 0, time.UTC),
		})
	})
}
