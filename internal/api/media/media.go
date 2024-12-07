package media

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	webrtc_ "github.com/imssyang/gweb/internal/webrtc"
	"github.com/pion/webrtc/v4"
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
	router.offer()
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

func (r *Router) offer() {
	r.Engine.POST("/"+Name+"/webrtc/offer", func(c *gin.Context) {
		//body, err := io.ReadAll(c.Request.Body)
		//if err != nil {
		//	c.String(http.StatusInternalServerError, "Internal Server Error")
		//	return
		//}
		//fmt.Printf("bodyLen %v\n", len(body))

		var session webrtc.SessionDescription
		if err := c.BindJSON(&session); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Failed to parse JSON: %v", err),
			})
			return
		}
		//fmt.Printf("Parsed SessionDescription: %+v\n", session)

		connData, err := webrtc_.Pool.CreateConnection("mediaui")
		if err != nil {
			c.String(http.StatusBadRequest+1, "CreateConnection error %v", err)
			return
		}

		connData.SetRemoteDescription(session)
		answerDesc, err := connData.SetLocalDescription(webrtc.SDPTypeAnswer)
		if err != nil {
			c.String(http.StatusBadRequest+2, "SetLocalDescription error %v", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "SessionDescription parsed successfully",
			"type":    answerDesc.Type.String(),
			"sdp":     answerDesc.SDP,
		})
	})
}
