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
	NetworkUDP = "udp"
	NetworkTCP = "tcp"
)

type NetworkType string

func (n NetworkType) String() string {
	return string(n)
}

type PoolID struct {
	Network      NetworkType
	LocalIP      string
	LocalPort    int
	EnableDetach bool
}

func NewPoolID(network NetworkType, address string, enableDetach bool) PoolID {
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

	return PoolID{
		Network:      network,
		LocalIP:      localIp,
		LocalPort:    localPort,
		EnableDetach: enableDetach,
	}
}

type PoolData struct {
	mu            sync.RWMutex
	ID            PoolID
	API           *webrtc.API
	ICEServerURLs []string
	connections   map[string]*ConnectionData
}

func NewPoolData(poolID PoolID, iceServerURLs []string) (*PoolData, error) {
	settingEngine := webrtc.SettingEngine{}
	if poolID.EnableDetach {
		settingEngine.DetachDataChannels()
	}

	localIP := poolID.LocalIP
	localPort := poolID.LocalPort
	switch poolID.Network {
	case NetworkUDP:
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
	case NetworkTCP:
		if localPort == 0 {
			return nil, fmt.Errorf("Invalid parameter: %+v", poolID)
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

	return &PoolData{
		ID: poolID,
		API: webrtc.NewAPI(
			webrtc.WithSettingEngine(settingEngine),
			webrtc.WithMediaEngine(media)),
		ICEServerURLs: iceServerURLs,
		connections:   make(map[string]*ConnectionData),
	}, nil
}

func (p *PoolData) CreateConnection(connID string, iceServerURLs []string, handlers ...any) (*ConnectionData, error) {
	cd := p.GetConnection(connID)
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

func (p *PoolData) ChooseICEServerURLs(connURLs []string, confURLs []string) []string {
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

func (p *PoolData) CloseConnection(connID string) error {
	cd := p.GetConnection(connID)
	if cd == nil {
		return fmt.Errorf("%s not found", connID)
	}

	p.mu.Lock()
	delete(p.connections, connID)
	p.mu.Unlock()

	defer func() {
		if err := cd.Connection.Close(); err != nil {
			log.Zap.Errorf("Failed to close connection(%v): %v", connID, err)
		}
	}()

	return nil
}

func (p *PoolData) GetConnection(connID string) *ConnectionData {
	p.mu.RLock()
	defer p.mu.RUnlock()
	cd, exists := p.connections[connID]
	if exists {
		return cd
	}
	return nil
}
