package webrtc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imssyang/gweb/internal/log"
	webrtc_ "github.com/imssyang/gweb/internal/webrtc"
	"github.com/pion/webrtc/v4"
)

type RTCOfferREQ struct {
	BaseREQ
	PreferNetwork string   `json:"preferNetwork"`
	PeerBindPort  bool     `json:"peerBindPort"`
	Description   string   `json:"description"`
	ICEServerURLs []string `json:"iceServerURLs"`
}

func (r *RTCOfferREQ) GetNetType() webrtc_.NetType {
	var nt webrtc_.NetType
	nt = webrtc_.NetTypeUDP
	if r.PreferNetwork == "tcp" {
		nt = webrtc_.NetTypeTCP
	}
	return nt
}

func (r *RTCOfferREQ) GetDescription() (webrtc.SessionDescription, error) {
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

type RTCOfferRSP struct {
	BaseRSP
	PreferNetwork string `json:"preferNetwork"`
	Description   string `json:"description"`
}

func (r *RTCOfferRSP) SetDescription(desc webrtc.SessionDescription) error {
	descJson, err := json.Marshal(desc)
	if err != nil {
		return err
	}

	r.Description = base64.StdEncoding.EncodeToString(descJson)
	return nil
}

func (r *Router) offer() {
	r.Engine.POST("/"+r.Name+"/offer", func(c *gin.Context) {
		connIDParam := c.DefaultQuery("connid", "")
		rsp := RTCOfferRSP{
			BaseRSP: BaseRSP{
				ConnID: connIDParam,
			},
		}
		req := RTCOfferREQ{
			BaseREQ: BaseREQ{
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
