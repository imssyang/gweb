package webrtc

import (
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine, name string) {
	router := NewRouter(engine, name+"/webrtc")
	router.offer()
	router.answer()
	router.candidate()
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

type BaseREQ struct {
	ConnID string `json:"connID"`
}

type BaseRSP struct {
	Err    string `json:"error"`
	ConnID string `json:"connID"`
}
