package formatify

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const Name = "formatify"

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
}

func (r *Router) index() {
	r.Engine.GET("/"+Name, func(c *gin.Context) {
		c.HTML(http.StatusOK, Name+"/index", gin.H{})
	})
	log.Println(PycmdDumps("/sss/test -a 1 -b 2 -c 3 adsf", 2))
}
