package webrtc

import (
	"fmt"
	"net"
	"sync"

	"github.com/pion/ice/v4"
	"github.com/pion/webrtc/v4"
)

type NetType string

const (
	NetTypeUDP = "udp"
	NetTypeTCP = "tcp"
)

type WebRTCPool struct {
	mu           sync.RWMutex
	API          *webrtc.API
	NetType      NetType
	LocalIP      string
	LocalPort    int
	EnableDetach bool
	ICEServerURL string
	connections  map[ConnectionID]*ConnectionData
}

func NewWebRTCPool(netType NetType, ip string, port int, enableDetach bool) (*WebRTCPool, error) {
	settingEngine := webrtc.SettingEngine{}
	if enableDetach {
		settingEngine.DetachDataChannels()
	}
	netIP := net.ParseIP(ip)
	switch netType {
	case NetTypeUDP:
		mux, err := ice.NewMultiUDPMuxFromPort(port)
		if err != nil {
			return nil, err
		}
		settingEngine.SetICEUDPMux(mux)
	case NetTypeTCP:
		settingEngine.SetNetworkTypes([]webrtc.NetworkType{
			webrtc.NetworkTypeTCP4,
			webrtc.NetworkTypeTCP6,
		})
		tcpListener, err := net.ListenTCP("tcp", &net.TCPAddr{
			IP:   netIP,
			Port: port,
		})
		if err != nil {
			return nil, err
		}
		tcpMux := webrtc.NewICETCPMux(nil, tcpListener, 8)
		settingEngine.SetICETCPMux(tcpMux)
	}

	media := &webrtc.MediaEngine{}
	if err := media.RegisterDefaultCodecs(); err != nil {
		return nil, err
	}

	return &WebRTCPool{
		NetType:      netType,
		LocalIP:      ip,
		LocalPort:    port,
		EnableDetach: enableDetach,
		API: webrtc.NewAPI(
			webrtc.WithSettingEngine(settingEngine),
			webrtc.WithMediaEngine(media)),
		connections: make(map[ConnectionID]*ConnectionData),
	}, nil
}

func (p *WebRTCPool) SetICEServerURL(url string) {
	p.ICEServerURL = url
}

func (p *WebRTCPool) CreateConnection(connID ConnectionID) (*ConnectionData, error) {
	cd := p.GetConnectionData(connID)
	if cd != nil {
		return nil, fmt.Errorf("%s repeated.", connID)
	}

	var config webrtc.Configuration
	if p.ICEServerURL != "" {
		config = webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{p.ICEServerURL},
				},
			},
		}
	}
	peerConnection, err := p.API.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	cd, err = NewConnectionData(connID, peerConnection, p)
	if err != nil {
		return nil, err
	}

	p.mu.Lock()
	p.connections[connID] = cd
	p.mu.Unlock()

	peerConnection.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		cd.OnICEConnectionStateChange(state)
	})

	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		cd.OnConnectionStateChange(state)
	})

	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i != nil {
			//check(peerConnection.AddICECandidate(i.ToJSON()))
		}
	})

	peerConnection.OnDataChannel(func(channel *webrtc.DataChannel) {
		cd.OnDataChannel(channel)
	})
	return cd, nil
}

func (p *WebRTCPool) CloseConnection(connID ConnectionID) error {
	cd := p.GetConnectionData(connID)
	if cd == nil {
		return fmt.Errorf("%s not found", connID)
	}

	if err := cd.Connection.Close(); err != nil {
		return err
	}

	p.mu.Lock()
	delete(p.connections, connID)
	p.mu.Unlock()
	return nil
}

func (p *WebRTCPool) GetConnectionData(connID ConnectionID) *ConnectionData {
	p.mu.RLock()
	defer p.mu.RUnlock()
	cd, exists := p.connections[connID]
	if exists {
		return cd
	}
	return nil
}
