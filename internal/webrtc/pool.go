package webrtc

import (
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/imssyang/gweb/internal/log"
	"github.com/pion/ice/v4"
	"github.com/pion/webrtc/v4"
)

const (
	NetTypeUDP = "udp"
	NetTypeTCP = "tcp"
)

type NetType string

func (n NetType) String() string {
	return string(n)
}

type WebRTCID struct {
	NetType      NetType
	LocalIP      string
	LocalPort    int
	EnableDetach bool
}

func NewWebRTCID(netType NetType, address string, enableDetach bool) WebRTCID {
	var (
		localIp   string
		localPort int
	)

	if len(address) != 0 {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			log.Zap.Fatalf("SplitHostPort(%s) error: %s", address, err)
			panic(err)
		}

		localIp = host
		localPort, err = strconv.Atoi(port)
		if err != nil {
			log.Zap.Fatalf("Atoi(%d) error: %s", port, err)
			panic(err)
		}
	}

	return WebRTCID{
		NetType:      netType,
		LocalIP:      localIp,
		LocalPort:    localPort,
		EnableDetach: enableDetach,
	}
}

type WebRTCPool struct {
	mu            sync.RWMutex
	ID            WebRTCID
	API           *webrtc.API
	ICEServerURLs []string
	connections   map[ConnectionID]*ConnectionData
}

func NewWebRTCPool(webrtcID WebRTCID, iceServerURLs []string) (*WebRTCPool, error) {
	settingEngine := webrtc.SettingEngine{}
	if webrtcID.EnableDetach {
		settingEngine.DetachDataChannels()
	}

	localIP := webrtcID.LocalIP
	localPort := webrtcID.LocalPort
	switch webrtcID.NetType {
	case NetTypeUDP:
		if localPort != 0 {
			if len(localIP) == 0 {
				mux, err := ice.NewMultiUDPMuxFromPort(localPort)
				if err != nil {
					return nil, err
				}
				settingEngine.SetICEUDPMux(mux)
			} else {
				localAddress := fmt.Sprintf("%s:%d", localIP, localPort)
				conn, err := net.ListenPacket("udp", localAddress)
				if err != nil {
					return nil, err
				}
				webrtc.NewICEUDPMux(nil, conn)
			}
		}
	case NetTypeTCP:
		if localPort == 0 {
			return nil, fmt.Errorf("Invalid parameter: %+v", webrtcID)
		}

		settingEngine.SetNetworkTypes([]webrtc.NetworkType{
			webrtc.NetworkTypeTCP4,
			webrtc.NetworkTypeTCP6,
		})
		tcpListener, err := net.ListenTCP("tcp", &net.TCPAddr{
			IP:   net.ParseIP(localIP),
			Port: localPort,
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
		ID: webrtcID,
		API: webrtc.NewAPI(
			webrtc.WithSettingEngine(settingEngine),
			webrtc.WithMediaEngine(media)),
		ICEServerURLs: iceServerURLs,
		connections:   make(map[ConnectionID]*ConnectionData),
	}, nil
}

func (p *WebRTCPool) CreateConnection(connID ConnectionID, iceServerURLs []string, handlers ...any) (*ConnectionData, error) {
	cd := p.GetConnectionData(connID)
	if cd != nil {
		return nil, fmt.Errorf("%s repeated.", connID)
	}

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: p.ChooseICEServerURLs(iceServerURLs, p.ICEServerURLs),
			},
		},
	}
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	cd, err = NewConnectionData(connID, peerConnection, p, handlers...)
	if err != nil {
		return nil, err
	}

	p.mu.Lock()
	p.connections[connID] = cd
	p.mu.Unlock()

	peerConnection.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		cd.OnICEConnectionStateChange(state)
	})

	peerConnection.OnICEGatheringStateChange(func(state webrtc.ICEGatheringState) {
		cd.OnICEGatheringStateChange(state)
	})

	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		cd.OnConnectionStateChange(state)
	})

	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		cd.OnICECandidate(candidate)
	})

	peerConnection.OnDataChannel(func(channel *webrtc.DataChannel) {
		cd.OnDataChannel(channel)
	})
	return cd, nil
}

func (p *WebRTCPool) ChooseICEServerURLs(connURLs []string, confURLs []string) []string {
	if len(connURLs) == 0 {
		return []string{} // confURLs
	}

	urlSet := make(map[string]struct{}, len(confURLs))
	for _, url := range confURLs {
		urlSet[url] = struct{}{}
	}

	var mixURLs []string
	for _, url := range connURLs {
		if _, exists := urlSet[url]; exists {
			mixURLs = append(mixURLs, url)
		}
	}
	if len(mixURLs) != 0 {
		return []string{} //  mixURLs
	}

	return []string{} //  confURLs
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
