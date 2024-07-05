package formatify

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

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
	router.mode()
}

func (r *Router) index() {
	r.Engine.GET("/"+Name, func(c *gin.Context) {
		c.HTML(http.StatusOK, Name+"/index", gin.H{
			"title":  "Formatify",
			"icon":   "img/formatify.svg",
			"style":  "css/formatify.min.css",
			"main":   "/js/formatify.min.js",
			"prefix": "/formatify",
		})
	})
}

func (r *Router) mode() {
	r.Engine.POST("/"+Name+"/:mode/:action", func(c *gin.Context) {
		escapeValue := c.DefaultQuery("escape", "false")
		hasEscape, err := strconv.ParseBool(escapeValue)
		if err != nil {
			hasEscape = false
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		mode := c.Param("mode")
		switch mode {
		case "json", "python", "command":
			action := c.Param("action")
			indent := 0
			if action == "contract" {
				indent = 0
			} else if action == "expand" {
				indent = map[string]int{
					"json":    4,
					"python":  1,
					"command": 2,
				}[mode]
			}

			formatted, err := PyfmtDumps(mode, string(body), indent, hasEscape)
			if err != nil {
				c.String(http.StatusBadRequest, "PyfmtDumps error %v", err)
				return
			}
			c.String(http.StatusOK, formatted)
		default:
			fmt.Printf("Unsupport %v mode!\n", mode)
		}
	})
}
