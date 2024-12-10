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

func (r *ICECandidatesRSP) SetCandidates(candidates ...webrtc.ICECandidateInit) {
	for _, candidate := range candidates {
		r.Candidates = append(r.Candidates, candidate)
	}
}

func (r *Router) candidate() {
	r.Engine.POST("/"+r.Name+"/candidate", func(c *gin.Context) {
		connIDParam := c.DefaultQuery("connid", "")
		rsp := ICECandidatesRSP{
			BaseRSP: BaseRSP{
				ConnID: connIDParam,
			},
		}
		req := ICECandidatesREQ{
			BaseREQ: BaseREQ{
				ConnID: connIDParam,
			},
			Candidates: make([]webrtc.ICECandidateInit, 0),
		}
		if err := c.BindJSON(&req); err != nil {
			rsp.Err = fmt.Sprintf("Failed to parse JSON: %v", err)
			c.JSON(http.StatusBadRequest, rsp)
			return
		}

		connID := webrtc_.ConnectionID(req.ConnID)
		connData, err := webrtc_.GetConnection(connID, 3*time.Second)
		if err != nil {
			rsp.Err = fmt.Sprintf("Failed to find webrtc connection: %v", err)
			c.JSON(http.StatusBadRequest, rsp)
			return
		}

		err = connData.AddRemoteCandidates(req.Candidates...)
		if err != nil {
			rsp.Err = fmt.Sprintf("Failed to add webrtc candidate: %v", err)
			c.JSON(http.StatusBadRequest, rsp)
			return
		}

		log.Zap.Debugf("ConnID[%v] add RemoteCandidates: %v", req.ConnID, req.Candidates)

		candidates, gatheringState := connData.FetchLocalCandidates()
		for _, candidate := range candidates {
			rsp.Candidates = append(rsp.Candidates, candidate.ToJSON())
		}
		rsp.GatheringState = gatheringState.String()

		c.JSON(http.StatusOK, rsp)
	})

	r.Engine.GET("/"+r.Name+"/candidate", func(c *gin.Context) {
		connIDParam := c.DefaultQuery("connid", "")
		rsp := ICECandidatesRSP{
			BaseRSP: BaseRSP{
				ConnID: connIDParam,
			},
			Candidates: make([]webrtc.ICECandidateInit, 0),
		}
		req := ICECandidatesREQ{
			BaseREQ: BaseREQ{
				ConnID: connIDParam,
			},
		}

		connID := webrtc_.ConnectionID(req.ConnID)
		connData, err := webrtc_.GetConnection(connID, 3*time.Second)
		if err != nil {
			rsp.Err = fmt.Sprintf("Failed to find webrtc connection: %v", err)
			c.JSON(http.StatusBadRequest, rsp)
			return
		}

		candidates, gatheringState := connData.FetchLocalCandidates()
		for _, candidate := range candidates {
			rsp.Candidates = append(rsp.Candidates, candidate.ToJSON())
		}
		rsp.GatheringState = gatheringState.String()

		log.Zap.Debugf("ConnID[%v] LocalCandidates: %v", req.ConnID, req.Candidates)
		c.JSON(http.StatusOK, rsp)
	})
}
