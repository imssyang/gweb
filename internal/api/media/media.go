package media

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imssyang/gweb/internal/api/media/webrtc"
	"github.com/imssyang/gweb/internal/conf"
)

func Register(engine *gin.Engine) {
	router := NewRouter(engine, "media")
	router.index()
	webrtc.Register(engine, router.Name)
}

type Router struct {
	Name string
	*gin.Engine
	*gin.RouterGroup
}

func NewRouter(engine *gin.Engine, name string) *Router {
	return &Router{
		Name:        name,
		Engine:      engine,
		RouterGroup: engine.Group(name),
	}
}

func (r *Router) index() {
	r.Engine.GET("/"+r.Name, func(c *gin.Context) {
		c.HTML(http.StatusOK, r.Name+"/index", gin.H{
			"title":         "Media",
			"icon":          "img/media.svg",
			"style":         "css/media.min.css",
			"main":          "/js/media.min.js",
			"urlGroup":      r.Name,
			"iceServerURLs": conf.App.WebRTC.ICEServers,
		})
	})
}
