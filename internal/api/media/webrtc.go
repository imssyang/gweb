package media

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

const WebRTCName = Name + "/webrtc"

type WebRTCRouter struct {
	*gin.Engine
	*gin.RouterGroup
}

func WebRTCRegister(engine *gin.Engine) {
	router := &WebRTCRouter{
		Engine:      engine,
		RouterGroup: engine.Group(WebRTCName),
	}
	router.offer()
	router.candidate()
}

type WebRTCBaseREQ struct {
	ConnID string `json:"connID"`
}

type WebRTCBaseRSP struct {
	Err    string `json:"error"`
	ConnID string `json:"connID"`
}

type WebRTCOfferREQ struct {
	WebRTCBaseREQ
	PreferNetwork string   `json:"preferNetwork"`
	PeerBindPort  bool     `json:"peerBindPort"`
	Description   string   `json:"description"`
	ICEServerURLs []string `json:"iceServerURLs"`
}

func (r *WebRTCOfferREQ) GetNetType() webrtc_.NetType {
	var nt webrtc_.NetType
	nt = webrtc_.NetTypeUDP
	if r.PreferNetwork == "tcp" {
		nt = webrtc_.NetTypeTCP
	}
	return nt
}

func (r *WebRTCOfferREQ) GetDescription() (webrtc.SessionDescription, error) {
	descJson, err := base64.StdEncoding.DecodeString(r.Description)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	var descObj webrtc.SessionDescription
	if err = json.Unmarshal(descJson, &descObj); err != nil {
		return webrtc.SessionDescription{}, err
	}

	return descObj, nil
}

type WebRTCOfferRSP struct {
	WebRTCBaseRSP
	PreferNetwork string `json:"preferNetwork"`
	Description   string `json:"description"`
}

func (r *WebRTCOfferRSP) SetDescription(desc webrtc.SessionDescription) error {
	descJson, err := json.Marshal(desc)
	if err != nil {
		return err
	}

	r.Description = base64.StdEncoding.EncodeToString(descJson)
	return nil
}

func (r *WebRTCRouter) offer() {
	r.Engine.POST("/"+WebRTCName+"/offer", func(c *gin.Context) {
		connIDParam := c.DefaultQuery("connid", "")
		rsp := WebRTCOfferRSP{
			WebRTCBaseRSP: WebRTCBaseRSP{
				ConnID: connIDParam,
			},
		}
		req := WebRTCOfferREQ{
			WebRTCBaseREQ: WebRTCBaseREQ{
				ConnID: connIDParam,
			},
		}
		if err := c.BindJSON(&req); err != nil {
			rsp.Err = fmt.Sprintf("Failed to parse JSON: %v", err)
			c.JSON(http.StatusBadRequest, rsp)
			return
		}

		remoteDesc, err := req.GetDescription()
		if err != nil {
			rsp.Err = fmt.Sprintf("Failed to parse description: %v", err)
			c.JSON(http.StatusBadRequest, rsp)
			return
		}

		log.Zap.Debugf("ConnID[%v] RemoteDescription: %+v", req.ConnID, remoteDesc.Type)

		netType := req.GetNetType()
		webrtcPool, err := webrtc_.Pool(netType, req.PeerBindPort)
		if err != nil {
			rsp.Err = fmt.Sprintf("Failed to find webrtc pool: %v", err)
			c.JSON(http.StatusServiceUnavailable, rsp)
			return
		}

		rsp.PreferNetwork = netType.String()
		connID := webrtc_.ConnectionID(req.ConnID)
		connData, err := webrtcPool.CreateConnection(connID, req.ICEServerURLs)
		if err != nil {
			rsp.Err = fmt.Sprintf("Failed to create webrtc connection: %v", err)
			c.JSON(http.StatusServiceUnavailable, rsp)
			return
		}

		err = connData.SetRemoteDescription(remoteDesc)
		if err != nil {
			rsp.Err = fmt.Sprintf("Failed to set remote description: %v", err)
			c.JSON(http.StatusServiceUnavailable, rsp)
			return
		}

		localDesc, err := connData.SetLocalDescription(webrtc.SDPTypeAnswer, false)
		if err != nil {
			rsp.Err = fmt.Sprintf("Failed to set local description: %v", err)
			c.JSON(http.StatusServiceUnavailable, rsp)
			return
		}

		log.Zap.Debugf("ConnID[%v] LocalDescription: %+v", req.ConnID, localDesc.Type)

		rsp.SetDescription(localDesc)
		c.JSON(http.StatusOK, rsp)
	})
}

type WebRTCCandidateREQ struct {
	WebRTCBaseREQ
	Candidates []webrtc.ICECandidateInit `json:"iceCandidates"`
}

type WebRTCCandidateRSP struct {
	WebRTCBaseRSP
	Candidates     []webrtc.ICECandidateInit `json:"iceCandidates"`
	GatheringState string                    `json:"iceGatheringState"`
}

func (r *WebRTCCandidateRSP) SetCandidates(candidates ...webrtc.ICECandidateInit) {
	for _, candidate := range candidates {
		r.Candidates = append(r.Candidates, candidate)
	}
}

func (r *WebRTCRouter) candidate() {
	r.Engine.POST("/"+WebRTCName+"/candidate", func(c *gin.Context) {
		connIDParam := c.DefaultQuery("connid", "")
		rsp := WebRTCCandidateRSP{
			WebRTCBaseRSP: WebRTCBaseRSP{
				ConnID: connIDParam,
			},
		}
		req := WebRTCCandidateREQ{
			WebRTCBaseREQ: WebRTCBaseREQ{
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
		c.JSON(http.StatusOK, rsp)
	})

	r.Engine.GET("/"+WebRTCName+"/candidate", func(c *gin.Context) {
		connIDParam := c.DefaultQuery("connid", "")
		rsp := WebRTCCandidateRSP{
			WebRTCBaseRSP: WebRTCBaseRSP{
				ConnID: connIDParam,
			},
			Candidates: make([]webrtc.ICECandidateInit, 0),
		}
		req := WebRTCCandidateREQ{
			WebRTCBaseREQ: WebRTCBaseREQ{
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
