package webrtc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imssyang/gweb/internal/log"
	webrtc_ "github.com/imssyang/gweb/internal/webrtc"
	"github.com/pion/webrtc/v4"
)

type RTCAnswerREQ struct {
	BaseREQ
	Description string `json:"description"`
}

func (r *RTCAnswerREQ) GetDescription() (webrtc.SessionDescription, error) {
	descJson, err := base64.StdEncoding.DecodeString(m.Description)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	var descObj webrtc.SessionDescription
	if err = json.Unmarshal(descJson, &descObj); err != nil {
		return webrtc.SessionDescription{}, err
	}

	if descObj.Type != webrtc.SDPTypeAnswer {
		return webrtc.SessionDescription{}, fmt.Errorf("Invalid SDP type: %s", descObj.Type.String())
	}

	return descObj, nil
}

type RTCAnswerRSP struct {
	BaseRSP
}

type RTCAnswerMSG struct {
	REQ RTCAnswerREQ
	RSP RTCAnswerRSP
}

func NewRTCAnswerMSG(connID string) *RTCAnswerMSG {
	return &RTCAnswerMSG{
		REQ: RTCAnswerREQ{
			BaseREQ: BaseREQ{
				ConnID: connID,
			},
		},
		RSP: RTCAnswerRSP{
			BaseRSP: BaseRSP{
				ConnID: connID,
			},
		},
	}
}

func (m *RTCAnswerMSG) Answer(c *gin.Context) {
	connID := webrtc_.ConnectionID(m.REQ.ConnID)
	connData, err := webrtc_.GetConnection(connID, 3*time.Second)
	if err != nil {
		m.RSP.Err = fmt.Sprintf("Failed to find webrtc connection: %v", err)
		c.JSON(http.StatusBadRequest, m.RSP)
		return
	}

	remoteDesc, err := m.REQ.GetDescription()
	if err != nil {
		m.RSP.Err = fmt.Sprintf("Failed to parse description: %v", err)
		c.JSON(http.StatusBadRequest, m.RSP)
		return
	}

	log.Zap.Debugf("ConnID[%v] RemoteDescription: %+v", m.REQ.ConnID, remoteDesc.Type)

	err = connData.SetRemoteDescription(remoteDesc)
	if err != nil {
		m.RSP.Err = fmt.Sprintf("Failed to set remote description: %v", err)
		c.JSON(http.StatusServiceUnavailable, m.RSP)
		return
	}

	c.JSON(http.StatusOK, m.RSP)
}

func (r *Router) answer() {
	r.Engine.POST("/"+r.Name+"/answer", func(c *gin.Context) {
		connID := c.DefaultQuery("connid", "")
		msg := NewRTCAnswerMSG(connID)

		if len(connID) == 0 {
			msg.RSP.Err = fmt.Sprintf("Failed to parse param: %v", c.Params)
			c.JSON(http.StatusBadRequest, msg.RSP)
			return
		}

		if err := c.BindJSON(&msg.REQ); err != nil {
			msg.RSP.Err = fmt.Sprintf("Failed to parse JSON: %v", err)
			c.JSON(http.StatusBadRequest, msg.RSP)
			return
		}

		msg.Answer(c)
	})
}
