package webrtc

import (
	"fmt"
	"sync"

	"github.com/pion/webrtc/v4"
)

type TrackID struct {
	ID       string
	StreamID string
	MimeType string
}

func NewTrackID(id, streamID, mimeType string) TrackID {
	return TrackID{
		ID:       id,
		StreamID: streamID,
		MimeType: mimeType,
	}
}

type TrackData struct {
	mu         sync.RWMutex
	TrackID    TrackID
	TrackLocal *webrtc.TrackLocalStaticSample
	RtpSender  *webrtc.RTPSender
	ConnData   *ConnectionData
}

func NewTrackData(trackID TrackID, trackLocal *webrtc.TrackLocalStaticSample, rtpSender *webrtc.RTPSender, connData *ConnectionData) (*TrackData, error) {
	if trackLocal == nil || rtpSender == nil || connData == nil {
		return nil, fmt.Errorf("%s invalid.\n", trackID)
	}
	return &TrackData{
		TrackID:    trackID,
		TrackLocal: trackLocal,
		RtpSender:  rtpSender,
		ConnData:   connData,
	}, nil
}
