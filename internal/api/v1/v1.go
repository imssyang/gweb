package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const Name = "v1"

type Router struct {
	*gin.RouterGroup
}

func Register(engine *gin.Engine) {
	router := &Router{
		RouterGroup: engine.Group(Name),
	}
	router.json()
}

func teste(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"sucess": "so beautiful",
	})
}

func (r *Router) json() {
	r.GET("/", teste)
}
