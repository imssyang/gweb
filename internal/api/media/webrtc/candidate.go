package webrtc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imssyang/gweb/internal/log"
	webrtc_ "github.com/imssyang/gweb/internal/webrtc"
	"github.com/pion/webrtc/v4"
)

type ICECandidatesREQ struct {
	BaseREQ
	Candidates []webrtc.ICECandidateInit `json:"iceCandidates"`
}

type ICECandidatesRSP struct {
	BaseRSP
	Candidates     []webrtc.ICECandidateInit `json:"iceCandidates"`
	GatheringState string                    `json:"iceGatheringState"`
}

type ICECandidatesMSG struct {
	REQ ICECandidatesREQ
	RSP ICECandidatesRSP
}

func NewICECandidatesMSG(connID string) *ICECandidatesMSG {
	return &ICECandidatesMSG{
		REQ: ICECandidatesREQ{
			BaseREQ: BaseREQ{
				ConnID: connID,
			},
			Candidates: make([]webrtc.ICECandidateInit, 0),
		},
		RSP: ICECandidatesRSP{
			BaseRSP: BaseRSP{
				ConnID: connID,
			},
			Candidates: make([]webrtc.ICECandidateInit, 0),
		},
	}
}

func (m *ICECandidatesMSG) GetCandidates(c *gin.Context) {
	connID := webrtc_.ConnectionID(m.REQ.ConnID)
	connData, err := webrtc_.ConnectionWithTimeout(connID, 3*time.Second)
	if err != nil {
		m.RSP.Err = fmt.Sprintf("Failed to find webrtc connection: %v", err)
		c.JSON(http.StatusBadRequest, m.RSP)
		return
	}

	candidates, gatheringState := connData.FetchLocalCandidates()
	for _, candidate := range candidates {
		m.RSP.Candidates = append(m.RSP.Candidates, candidate.ToJSON())
	}
	m.RSP.GatheringState = gatheringState.String()

	log.Zap.Debugf("ConnID[%v] LocalCandidates: %v", m.REQ.ConnID, m.REQ.Candidates)
	c.JSON(http.StatusOK, m.RSP)
}

func (m *ICECandidatesMSG) SetCandidates(c *gin.Context) {
	connID := webrtc_.ConnectionID(m.REQ.ConnID)
	connData, err := webrtc_.ConnectionWithTimeout(connID, 3*time.Second)
	if err != nil {
		m.RSP.Err = fmt.Sprintf("Failed to find webrtc connection: %v", err)
		c.JSON(http.StatusBadRequest, m.RSP)
		return
	}

	err = connData.AddRemoteCandidates(m.REQ.Candidates...)
	if err != nil {
		m.RSP.Err = fmt.Sprintf("Failed to add webrtc candidate: %v", err)
		c.JSON(http.StatusBadRequest, m.RSP)
		return
	}

	log.Zap.Debugf("ConnID[%v] add RemoteCandidates: %v", m.REQ.ConnID, m.REQ.Candidates)

	candidates, gatheringState := connData.FetchLocalCandidates()
	for _, candidate := range candidates {
		m.RSP.Candidates = append(m.RSP.Candidates, candidate.ToJSON())
	}
	m.RSP.GatheringState = gatheringState.String()

	c.JSON(http.StatusOK, m.RSP)
}

func (r *Router) candidate() {
	r.Engine.GET("/"+r.Name+"/candidate", func(c *gin.Context) {
		connID := c.DefaultQuery("connid", "")
		msg := NewICECandidatesMSG(connID)

		if len(connID) == 0 {
			msg.RSP.Err = fmt.Sprintf("Failed to parse param: %v", c.Params)
			c.JSON(http.StatusBadRequest, msg.RSP)
			return
		}

		msg.GetCandidates(c)
	})

	r.Engine.POST("/"+r.Name+"/candidate", func(c *gin.Context) {
		connID := c.DefaultQuery("connid", "")
		msg := NewICECandidatesMSG(connID)

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

		msg.SetCandidates(c)
	})
}
