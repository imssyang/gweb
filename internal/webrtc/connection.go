package webrtc

import (
	"fmt"
	"sync"

	"github.com/imssyang/gweb/internal/log"
	"github.com/pion/webrtc/v4"
)

type ConnectionID string

type ConnectionData struct {
	mu             sync.RWMutex
	ConnID         ConnectionID
	Connection     *webrtc.PeerConnection
	Pool           *PoolData
	GatheringState webrtc.ICEGatheringState
	candidates     []*webrtc.ICECandidate
	channels       map[ChannelID]*ChannelData
	streams        map[string]*StreamData

	onICECandidateHandler func(*webrtc.ICECandidate, webrtc.ICEGatheringState)
}

func NewConnectionData(connID ConnectionID, connection *webrtc.PeerConnection, pool *PoolData, handlers ...any) (*ConnectionData, error) {
	if connection == nil || pool == nil {
		return nil, fmt.Errorf("%s invalid.", connID)
	}
	connData := ConnectionData{
		ConnID:         connID,
		Connection:     connection,
		Pool:           pool,
		GatheringState: webrtc.ICEGatheringStateNew,
		candidates:     make([]*webrtc.ICECandidate, 0),
		channels:       make(map[ChannelID]*ChannelData),
		streams:        make(map[string]*StreamData),
	}

	for _, handler := range handlers {
		switch v := handler.(type) {
		case func(*webrtc.ICECandidate, webrtc.ICEGatheringState):
			connData.onICECandidateHandler = v
		}
	}
	return &connData, nil
}

func (d *ConnectionData) OnICEConnectionStateChange(state webrtc.ICEConnectionState) {
	log.Zap.Debugf("ConnID[%v] ICEConnectionStateChange: %v", d.ConnID, state)
	switch state {
	case webrtc.ICEConnectionStateNew:
	case webrtc.ICEConnectionStateChecking:
	case webrtc.ICEConnectionStateConnected:
		d.mu.RLock()
		defer d.mu.RUnlock()
		log.Zap.Infof("Connid: %s NumOfStream: %d", d.ConnID, len(d.streams))
		for _, stream := range d.streams {
			stream.OnDataTracksOpen()
		}
	case webrtc.ICEConnectionStateCompleted:
	case webrtc.ICEConnectionStateDisconnected:
	case webrtc.ICEConnectionStateFailed:
	case webrtc.ICEConnectionStateClosed:
	}
}

func (d *ConnectionData) OnICEGatheringStateChange(state webrtc.ICEGatheringState) {
	log.Zap.Debugf("ConnID[%v] ICEGatheringStateChange: %v", d.ConnID, state)
	d.GatheringState = state
}

func (d *ConnectionData) OnConnectionStateChange(state webrtc.PeerConnectionState) {
	log.Zap.Debugf("ConnID[%v] ConnectionStateChange: %v", d.ConnID, state)
	switch state {
	case webrtc.PeerConnectionStateNew:
	case webrtc.PeerConnectionStateConnecting:
	case webrtc.PeerConnectionStateConnected:
	case webrtc.PeerConnectionStateDisconnected:
		d.Pool.CloseConnection(d.ConnID)
	case webrtc.PeerConnectionStateFailed:
		d.Pool.CloseConnection(d.ConnID)
	case webrtc.PeerConnectionStateClosed:
		d.Pool.CloseConnection(d.ConnID)
	}
}

func (d *ConnectionData) OnICECandidate(candidate *webrtc.ICECandidate) {
	log.Zap.Debugf("ConnID[%v] OnICECandidate: %v", d.ConnID, candidate)
	handler := d.onICECandidateHandler
	if handler != nil {
		handler(candidate, d.GatheringState)
	} else {
		if candidate != nil {
			d.mu.Lock()
			d.candidates = append(d.candidates, candidate)
			d.mu.Unlock()
		}
	}
}

func (d *ConnectionData) OnDataChannel(channel *webrtc.DataChannel) {
	log.Zap.Debugf("ConnID[%v] DataChannel %s %d", d.ConnID, channel.Label(), channel.ID())

	chanID := ChannelID(channel.Label())
	cd, err := NewChannelData(chanID, channel, d)
	if err != nil {
		return
	}

	d.mu.Lock()
	d.channels[chanID] = cd
	d.mu.Unlock()

	channel.OnOpen(func() {
		cd.OnDataChannelOpen()
	})

	channel.OnClose(func() {
		d.mu.Lock()
		delete(d.channels, chanID)
		d.mu.Unlock()
	})

	channel.OnMessage(func(msg webrtc.DataChannelMessage) {
		cd.OnDataChannelMessage(msg)
	})
}

func (d *ConnectionData) SetStreamURLs(streamURLs []string) error {
	for _, url := range streamURLs {
		sd, err := NewStreamData(url, d)
		if err != nil {
			return err
		}

		d.mu.Lock()
		d.streams[url] = sd
		d.mu.Unlock()
	}
	return nil
}

func (d *ConnectionData) SetRemoteDescription(desc webrtc.SessionDescription) error {
	if err := d.Connection.SetRemoteDescription(desc); err != nil {
		return err
	}
	return nil
}

func (d *ConnectionData) SetLocalDescription(sdpType webrtc.SDPType, waitGatheringComplete bool) (webrtc.SessionDescription, error) {
	var desc webrtc.SessionDescription
	var err error
	switch sdpType {
	case webrtc.SDPTypeOffer:
		desc, err = d.Connection.CreateOffer(nil)
	case webrtc.SDPTypeAnswer:
		desc, err = d.Connection.CreateAnswer(nil)
	}
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	if err = d.Connection.SetLocalDescription(desc); err != nil {
		return webrtc.SessionDescription{}, err
	}

	if waitGatheringComplete {
		<-webrtc.GatheringCompletePromise(d.Connection)
		descWithCandidates := d.Connection.LocalDescription()
		return *descWithCandidates, nil
	} else {
		return desc, nil
	}
}

func (d *ConnectionData) AddRemoteCandidates(candidates ...webrtc.ICECandidateInit) error {
	for _, candidate := range candidates {
		err := d.Connection.AddICECandidate(candidate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *ConnectionData) FetchLocalCandidates() ([]*webrtc.ICECandidate, webrtc.ICEGatheringState) {
	d.mu.Lock()
	defer d.mu.Unlock()
	candidates := make([]*webrtc.ICECandidate, len(d.candidates))
	copy(candidates, d.candidates)
	d.candidates = d.candidates[:0]
	return candidates, d.GatheringState
}

func (d *ConnectionData) AddChannel(chanID ChannelID) (*ChannelData, error) {
	channel, err := d.Connection.CreateDataChannel(string(chanID), nil)
	if err != nil {
		return nil, err
	}

	d.OnDataChannel(channel)
	chanData := d.GetChannel(chanID)
	if chanData == nil {
		return nil, fmt.Errorf("%s invalid.", chanID)
	}
	return chanData, nil
}

func (d *ConnectionData) GetChannel(chanID ChannelID) *ChannelData {
	d.mu.RLock()
	defer d.mu.RUnlock()
	cd, exists := d.channels[chanID]
	if exists {
		return cd
	}
	return nil
}

func (d *ConnectionData) GetStream(streamID string) *StreamData {
	d.mu.RLock()
	defer d.mu.RUnlock()
	td, exists := d.streams[streamID]
	if exists {
		return td
	}
	return nil
}
