package media

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const Name = "media"

type Router struct {
	*gin.Engine
	*gin.RouterGroup
}

func Register(engine *gin.Engine) {
	router := &Router{
		Engine:      engine,
		RouterGroup: engine.Group(Name),
	}
	router.index()
	WebRTCRegister(engine)
}

func (r *Router) index() {
	r.Engine.GET("/"+Name, func(c *gin.Context) {
		c.HTML(http.StatusOK, Name+"/index", gin.H{
			"title":  "Media",
			"icon":   "img/media.svg",
			"style":  "css/media.min.css",
			"main":   "/js/media.min.js",
			"prefix": "/media",
		})
	})
}
