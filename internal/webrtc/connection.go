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
	Pool           *WebRTCPool
	GatheringState webrtc.ICEGatheringState
	candidates     []*webrtc.ICECandidate
	channels       map[ChannelID]*ChannelData
	tracks         map[TrackID]*TrackData

	onICECandidateHandler func(*webrtc.ICECandidate, webrtc.ICEGatheringState)
}

func NewConnectionData(connID ConnectionID, connection *webrtc.PeerConnection, pool *WebRTCPool, handlers ...any) (*ConnectionData, error) {
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
		tracks:         make(map[TrackID]*TrackData),
	}

	for _, handler := range handlers {
		// 使用类型断言判断类型并处理
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
	case webrtc.PeerConnectionStateFailed:
		// Wait until PeerConnection has had no network activity for 30 seconds or another failure.
		// It may be reconnected using an ICE Restart.
		// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
		// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
	case webrtc.PeerConnectionStateClosed:
		// PeerConnection was explicitly closed. This usually happens from a DTLS CloseNotify
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

func (d *ConnectionData) SetRemoteDescription(desc webrtc.SessionDescription) error {
	if err := d.Connection.SetRemoteDescription(desc); err != nil {
		return err
	}
	return nil
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
	chanData := d.GetChannelData(chanID)
	if chanData == nil {
		return nil, fmt.Errorf("%s invalid.", chanID)
	}
	return chanData, nil
}

func (d *ConnectionData) AddTrack(id, streamID, mimeType string) (*TrackData, error) {
	trackLocal, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: mimeType}, id, streamID)
	if err != nil {
		return nil, err
	}

	rtpSender, err := d.Connection.AddTrack(trackLocal)
	if err != nil {
		return nil, err
	}

	trackID := NewTrackID(id, streamID, mimeType)
	trackData, err := NewTrackData(trackID, trackLocal, rtpSender, d)
	if err != nil {
		return nil, err
	}

	// Read incoming RTCP packets
	// Before these packets are returned they are processed by interceptors. For things
	// like NACK this needs to be called.
	//go func() {
	//	rtcpBuf := make([]byte, 1500)
	//	for {
	//		if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
	//			return
	//		}
	//	}
	//}()
	return trackData, nil
}

func (d *ConnectionData) GetChannelData(chanID ChannelID) *ChannelData {
	d.mu.RLock()
	defer d.mu.RUnlock()
	cd, exists := d.channels[chanID]
	if exists {
		return cd
	}
	return nil
}

func (d *ConnectionData) GetTrackData(trackID TrackID) *TrackData {
	d.mu.RLock()
	defer d.mu.RUnlock()
	td, exists := d.tracks[trackID]
	if exists {
		return td
	}
	return nil
}
