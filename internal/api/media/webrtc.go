package media

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

const WebRTCName = Name + "/webrtc"

func WebRTCRegister(engine *gin.Engine) {
	router := &WebRTCRouter{
		Engine:      engine,
		RouterGroup: engine.Group(WebRTCName),
	}
	router.offer()
}

type WebRTCOfferREQ struct {
	ConnID        string   `json:"connID"`
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
	ConnID        string `json:"connID"`
	PreferNetwork string `json:"preferNetwork"`
	Description   string `json:"description"`
}

func NewWebRTCOfferRSP(connID string, preferNetwork webrtc_.NetType, desc webrtc.SessionDescription) (*WebRTCOfferRSP, error) {
	r := &WebRTCOfferRSP{
		ConnID:        connID,
		PreferNetwork: string(preferNetwork),
	}
	err := r.SetDescription(desc)
	return r, err
}

func (r *WebRTCOfferRSP) SetDescription(desc webrtc.SessionDescription) error {
	descJson, err := json.Marshal(desc)
	if err != nil {
		return err
	}

	r.Description = base64.StdEncoding.EncodeToString(descJson)
	return nil
}

type WebRTCRouter struct {
	*gin.Engine
	*gin.RouterGroup
}

func (r *WebRTCRouter) offer() {
	r.Engine.POST("/"+WebRTCName+"/:connType", func(c *gin.Context) {
		var offerREQ WebRTCOfferREQ
		if err := c.BindJSON(&offerREQ); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Failed to parse JSON: %v", err),
			})
			return
		}

		remoteDesc, err := offerREQ.GetDescription()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Failed to parse description: %v", err),
			})
			return
		}

		log.Zap.Debugf("RemoteDescription: %+v\n", remoteDesc.Type)

		netType := offerREQ.GetNetType()
		webrtcPool, err := webrtc_.Pool(netType, offerREQ.PeerBindPort)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": fmt.Sprintf("Failed to find webrtc pool: %v", err),
			})
			return
		}

		connID := webrtc_.ConnectionID(offerREQ.ConnID)
		connData, err := webrtcPool.CreateConnection(connID, offerREQ.ICEServerURLs)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": fmt.Sprintf("Failed to create webrtc connection: %v", err),
			})
			return
		}

		err = connData.SetRemoteDescription(remoteDesc)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": fmt.Sprintf("Failed to set remote description: %v", err),
			})
			return
		}

		localDesc, err := connData.SetLocalDescription(webrtc.SDPTypeAnswer, true)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": fmt.Sprintf("Failed to set local description: %v", err),
			})
			return
		}

		log.Zap.Debugf("LocalDescription: %+v\n", localDesc.Type)

		offerRsp, _ := NewWebRTCOfferRSP(offerREQ.ConnID, netType, localDesc)
		c.JSON(http.StatusOK, gin.H{
			"connID":        offerRsp.ConnID,
			"preferNetwork": offerRsp.PreferNetwork,
			"description":   offerRsp.Description,
		})
	})
}
