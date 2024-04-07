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
	router.ping()
	router.json()
	router.index()
	router.posts()
	router.users()
	router.raw()
}

func (r *Router) ping() {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
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

func (r *Router) index() {
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

}

func (r *Router) posts() {
	r.GET("/posts/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
		})
	})
}

func (r *Router) users() {
	r.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "Users",
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
