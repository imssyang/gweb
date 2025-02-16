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
	StreamURLs    []string `json:"streamURLs"`
	PreferNetwork string   `json:"preferNetwork"`
	PeerBindPort  bool     `json:"peerBindPort"`
	Description   string   `json:"description"`
	ICEServerURLs []string `json:"iceServerURLs"`
}

func (m *RTCOfferREQ) GetNetwork() webrtc_.NetworkType {
	var t webrtc_.NetworkType
	t = webrtc_.NetworkUDP
	if m.PreferNetwork == "tcp" {
		t = webrtc_.NetworkTCP
	}
	return t
}

func (m *RTCOfferREQ) GetDescription() (webrtc.SessionDescription, error) {
	descJson, err := base64.StdEncoding.DecodeString(m.Description)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	var descObj webrtc.SessionDescription
	if err = json.Unmarshal(descJson, &descObj); err != nil {
		return webrtc.SessionDescription{}, err
	}

	if descObj.Type != webrtc.SDPTypeOffer {
		return webrtc.SessionDescription{}, fmt.Errorf("Invalid SDP type: %s", descObj.Type.String())
	}

	return descObj, nil
}

type RTCOfferRSP struct {
	BaseRSP
	PreferNetwork string `json:"preferNetwork"`
	Description   string `json:"description"`
}

func (m *RTCOfferRSP) SetDescription(desc webrtc.SessionDescription) error {
	descJson, err := json.Marshal(desc)
	if err != nil {
		return err
	}

	m.Description = base64.StdEncoding.EncodeToString(descJson)
	return nil
}

type RTCOfferMSG struct {
	REQ RTCOfferREQ
	RSP RTCOfferRSP
}

func NewRTCOfferMSG(connID string) *RTCOfferMSG {
	return &RTCOfferMSG{
		REQ: RTCOfferREQ{
			BaseREQ: BaseREQ{
				ConnID: connID,
			},
		},
		RSP: RTCOfferRSP{
			BaseRSP: BaseRSP{
				ConnID: connID,
			},
		},
	}
}

func (m *RTCOfferMSG) connection(c *gin.Context) (*webrtc_.ConnectionData, error) {
	connID := m.REQ.ConnID
	connData := webrtc_.Connection(connID)
	if connData != nil {
		m.RSP.Err = fmt.Sprintf("Repead connected: %v", connID)
		c.JSON(http.StatusBadRequest, m.RSP)
		return nil, fmt.Errorf("RepeadConnection: %v", connID)
	}

	network := m.REQ.GetNetwork()
	poolData, err := webrtc_.Pool(network, m.REQ.PeerBindPort)
	if err != nil {
		m.RSP.Err = fmt.Sprintf("Failed to find webrtc pool: %v", err)
		c.JSON(http.StatusServiceUnavailable, m.RSP)
		return nil, err
	}

	m.RSP.PreferNetwork = network.String()
	connData, err = poolData.CreateConnection(connID, m.REQ.ICEServerURLs)
	if err != nil {
		m.RSP.Err = fmt.Sprintf("Failed to create webrtc connection: %v", err)
		c.JSON(http.StatusServiceUnavailable, m.RSP)
		return nil, err
	}

	if len(m.REQ.StreamURLs) > 0 {
		err = connData.SetStreamURLs(m.REQ.StreamURLs)
		if err != nil {
			m.RSP.Err = fmt.Sprintf("Invalid URLs: %v", m.REQ.StreamURLs)
			c.JSON(http.StatusBadRequest, m.RSP)
			return nil, err
		}
	}

	return connData, nil
}

func (m *RTCOfferMSG) Offer(c *gin.Context) {
	if len(m.REQ.StreamURLs) == 0 {
		m.RSP.Err = fmt.Sprintf("%v no any URLs", m.REQ.ConnID)
		c.JSON(http.StatusBadRequest, m.RSP)
		return
	}

	log.Zap.Debugf("ConnID[%v] StreamURLs: %+v", m.REQ.ConnID, m.REQ.StreamURLs)

	connData, err := m.connection(c)
	if err != nil {
		return
	}

	localDesc, err := connData.SetLocalDescription(webrtc.SDPTypeOffer, false)
	if err != nil {
		m.RSP.Err = fmt.Sprintf("Fail create PeerOfferSDP: %v", err)
		c.JSON(http.StatusServiceUnavailable, m.RSP)
		return
	}

	log.Zap.Debugf("ConnID[%v] LocalDescription: %+v", m.REQ.ConnID, localDesc.Type)

	m.RSP.SetDescription(localDesc)
	c.JSON(http.StatusOK, m.RSP)
}

func (m *RTCOfferMSG) Answer(c *gin.Context) {
	remoteDesc, err := m.REQ.GetDescription()
	if err != nil {
		m.RSP.Err = fmt.Sprintf("Failed to parse description: %v", err)
		c.JSON(http.StatusBadRequest, m.RSP)
		return
	}

	log.Zap.Debugf("ConnID[%v] RemoteDescription: %+v", m.REQ.ConnID, remoteDesc.Type)

	connData, err := m.connection(c)
	if err != nil {
		return
	}

	err = connData.SetRemoteDescription(remoteDesc)
	if err != nil {
		m.RSP.Err = fmt.Sprintf("Failed to set remote description: %v", err)
		c.JSON(http.StatusServiceUnavailable, m.RSP)
		return
	}

	localDesc, err := connData.SetLocalDescription(webrtc.SDPTypeAnswer, false)
	if err != nil {
		m.RSP.Err = fmt.Sprintf("Failed to set local description: %v", err)
		c.JSON(http.StatusServiceUnavailable, m.RSP)
		return
	}

	log.Zap.Debugf("ConnID[%v] LocalDescription: %+v", m.REQ.ConnID, localDesc.Type)

	m.RSP.SetDescription(localDesc)
	c.JSON(http.StatusOK, m.RSP)
}

func (r *Router) offer() {
	r.Engine.POST("/"+r.Name+"/offer", func(c *gin.Context) {
		connID := c.DefaultQuery("connid", "")
		msg := NewRTCOfferMSG(connID)

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

		if len(msg.REQ.Description) == 0 {
			msg.Offer(c)
		} else {
			msg.Answer(c)
		}
	})
}
