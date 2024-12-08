package webrtc

import (
	"fmt"
	"sync"

	"github.com/imssyang/gweb/internal/log"
	"github.com/pion/webrtc/v4"
)

type ConnectionID string

type ConnectionData struct {
	mu         sync.RWMutex
	ConnID     ConnectionID
	Connection *webrtc.PeerConnection
	Pool       *WebRTCPool
	channels   map[ChannelID]*ChannelData
	tracks     map[TrackID]*TrackData
}

func NewConnectionData(connID ConnectionID, connection *webrtc.PeerConnection, pool *WebRTCPool) (*ConnectionData, error) {
	if connection == nil || pool == nil {
		return nil, fmt.Errorf("%s invalid.", connID)
	}
	return &ConnectionData{
		ConnID:     connID,
		Connection: connection,
		Pool:       pool,
		channels:   make(map[ChannelID]*ChannelData),
		tracks:     make(map[TrackID]*TrackData),
	}, nil
}

func (d *ConnectionData) OnICEConnectionStateChange(state webrtc.ICEConnectionState) {
	// This will notify you when the peer has connected/disconnected
	log.Zap.Infof("ICEConnectionStateChange: %v\n", state)
}

func (d *ConnectionData) OnICEGatheringStateChange(state webrtc.ICEGatheringState) {
	log.Zap.Infof("ICEGatheringStateChange: %v\n", state)
}

func (d *ConnectionData) OnConnectionStateChange(state webrtc.PeerConnectionState) {
	// This will notify you when the peer has connected/disconnected
	log.Zap.Infof("Peer Connection State has changed: %s\n", state.String())
	if state == webrtc.PeerConnectionStateFailed {
		// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
		// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
		// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
		fmt.Println("Peer Connection has gone to failed exiting")
	}

	if state == webrtc.PeerConnectionStateClosed {
		// PeerConnection was explicitly closed. This usually happens from a DTLS CloseNotify
		fmt.Println("Peer Connection has gone to closed exiting")
	}
}

func (d *ConnectionData) OnICECandidate(candidate *webrtc.ICECandidate) {
	log.Zap.Infof("OnICECandidate: %v\n", candidate)
	if candidate != nil {
		//check(peerConnection.AddICECandidate(i.ToJSON()))
	}
}

func (d *ConnectionData) OnDataChannel(channel *webrtc.DataChannel) {
	log.Zap.Infof("New DataChannel %s %d\n", channel.Label(), channel.ID())

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

func (d *ConnectionData) WaitConnectionComplete() error {
	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(d.Connection)

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete
	return nil
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
