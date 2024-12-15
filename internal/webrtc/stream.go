package webrtc

import (
	"fmt"
	"sync"

	"github.com/imssyang/gweb/internal/log"
	"github.com/imssyang/gweb/internal/media"
	"github.com/pion/webrtc/v4"
)

type StreamData struct {
	mu       sync.RWMutex
	URL      string
	Media    *media.Media
	ConnData *ConnectionData
	tracks   map[TrackID]*TrackData
}

func NewStreamData(uRL string, connData *ConnectionData) (*StreamData, error) {
	m, err := media.NewMedia(uRL)
	if err != nil {
		return nil, err
	}

	d := &StreamData{
		URL:      uRL,
		Media:    m,
		ConnData: connData,
		tracks:   make(map[TrackID]*TrackData),
	}

	for i, parser := range m.Parsers {
		trackID := TrackID{
			ID:       fmt.Sprint(i),
			StreamID: uRL,
			MimeType: parser.MimeType(),
		}

		_, err := d.AddTrack(trackID, parser)
		if err != nil {
			return nil, err
		}
	}

	return d, nil
}

func (d *StreamData) OnDataTracksOpen() {
	d.mu.RLock()
	defer d.mu.RUnlock()
	log.Zap.Infof("URL: %s NumOfTrack: %d", d.URL, len(d.tracks))
	for _, track := range d.tracks {
		go func(t *TrackData) {
			t.OnPacketRead()
		}(track)
		go func(t *TrackData) {
			t.OnPacketWrite()
		}(track)
	}
}

func (d *StreamData) AddTrack(trackID TrackID, trackParser media.TrackParser) (*TrackData, error) {
	trackLocal, err := webrtc.NewTrackLocalStaticSample(
		webrtc.RTPCodecCapability{MimeType: trackID.MimeType},
		trackID.ID,
		trackID.StreamID)
	if err != nil {
		return nil, err
	}

	rtpSender, err := d.ConnData.Connection.AddTrack(trackLocal)
	if err != nil {
		return nil, err
	}

	trackData, err := NewTrackData(trackID, trackParser, trackLocal, rtpSender, d.ConnData)
	if err != nil {
		return nil, err
	}

	d.mu.Lock()
	d.tracks[trackID] = trackData
	d.mu.Unlock()

	return trackData, nil
}

func (d *StreamData) GetTrack(trackID TrackID) *TrackData {
	d.mu.RLock()
	defer d.mu.RUnlock()
	td, exists := d.tracks[trackID]
	if exists {
		return td
	}
	return nil
}
