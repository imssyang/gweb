package formatify

import (
	"fmt"
	"io"
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
	router.command()
}

func (r *Router) index() {
	r.Engine.GET("/"+Name, func(c *gin.Context) {
		c.HTML(http.StatusOK, Name+"/index", gin.H{})
	})
	//log.Println(PycmdDumps("/sss/test -a 1 -b 2 -c 3 adsf", 2))
}

func (r *Router) command() {
	const Prefix = "/" + Name + "/command"

	r.Engine.POST(Prefix+"/contract", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(500, "Internal Server Error")
			return
		}

		formated, err := PycmdDumps(string(body), 0)
		if err != nil {
			c.String(501, "PycmdDumps Error %v", err)
			return
		}
		fmt.Println(3333, body, "ddd", string(body), "abcd", formated)
		c.String(200, formated)
	})

	r.Engine.POST(Prefix+"/expand", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(500, "Internal Server Error")
			return
		}

		formated, err := PycmdDumps(string(body), 2)
		if err != nil {
			c.String(501, "PycmdDumps Error %v", err)
			return
		}
		fmt.Println(1234, body, "ddd", string(body), "abcd", formated)
		c.String(200, formated)
	})
}
