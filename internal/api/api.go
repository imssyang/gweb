package api

import (
	"log"
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
	router.json()
	router.posts()
	router.raw()
}

func (r *Router) detect() {
	r.HEAD("/", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})
}

func (r *Router) index() {
	r.GET("/index", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/formatify")
	})
}

func (r *Router) json() {
	r.GET("/json", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO",
			"tag":  "<br>",
		}

		log.Println(data)
		c.AsciiJSON(http.StatusOK, data)
	})
}

func (r *Router) posts() {
	r.GET("/posts/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
		})
	})
}

func (r *Router) raw() {
	r.GET("/raw", func(c *gin.Context) {
		c.HTML(http.StatusOK, "raw.tmpl", map[string]interface{}{
			"now": time.Date(2017, 07, 01, 0, 0, 0, 0, time.UTC),
		})
	})
}
