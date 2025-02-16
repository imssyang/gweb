package webrtc

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	webrtc_ "github.com/imssyang/gweb/internal/webrtc"
)

type RTCCloseREQ struct {
	BaseREQ
}

type RTCCloseRSP struct {
	BaseRSP
}

type RTCCloseMSG struct {
	REQ RTCCloseREQ
	RSP RTCCloseRSP
}

func NewRTCCloseMSG(connID string) *RTCCloseMSG {
	return &RTCCloseMSG{
		REQ: RTCCloseREQ{
			BaseREQ: BaseREQ{
				ConnID: connID,
			},
		},
		RSP: RTCCloseRSP{
			BaseRSP: BaseRSP{
				ConnID: connID,
			},
		},
	}
}

func (m *RTCCloseMSG) Close(c *gin.Context) {
	connData := webrtc_.Connection(m.REQ.ConnID)
	if connData == nil {
		m.RSP.Err = fmt.Sprintf("No connection: %v", m.REQ.ConnID)
		c.JSON(http.StatusBadRequest, m.RSP)
		return
	}

	err := connData.Pool.CloseConnection(m.REQ.ConnID)
	if err != nil {
		m.RSP.Err = fmt.Sprintf("Failed to close connection(%v): %v", m.REQ.ConnID, err)
		c.JSON(http.StatusBadRequest, m.RSP)
		return
	}

	c.JSON(http.StatusOK, m.RSP)
}

func (r *Router) close() {
	r.Engine.DELETE("/"+r.Name+"/connection", func(c *gin.Context) {
		connID := c.DefaultQuery("connid", "")
		msg := NewRTCCloseMSG(connID)

		if len(connID) == 0 {
			msg.RSP.Err = fmt.Sprintf("Failed to parse param: %v", c.Params)
			c.JSON(http.StatusBadRequest, msg.RSP)
			return
		}

		msg.Close(c)
	})
}