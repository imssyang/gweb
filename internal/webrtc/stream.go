package webrtc

import (
	"net/url"
	"path/filepath"
	"sync"

	"github.com/pion/webrtc/v4"
)

type StreamData struct {
	mu       sync.RWMutex
	URL      string
	ConnData *ConnectionData
	tracks   map[TrackID]*TrackData
}

func NewStreamData(url string, connData *ConnectionData) (*StreamData, error) {
	d := &StreamData{
		URL:      url,
		ConnData: connData,
		tracks:   make(map[TrackID]*TrackData),
	}

	return d, nil
}

func (d *StreamData) ConnectURL() error {
	parsedURL, err := url.Parse(d.URL)
	if err != nil {
		return err
	}

	if len(parsedURL.Scheme) == 0 || parsedURL.Scheme == "file" {
		path := parsedURL.Path
		if len(parsedURL.Host) > 0 {
			path = filepath.Join(string(filepath.Separator), parsedURL.Host, path)
		}

	}

	return nil
}

func (d *StreamData) AddTrack(id, streamID, mimeType string) (*TrackData, error) {
	trackLocal, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: mimeType}, id, streamID)
	if err != nil {
		return nil, err
	}

	rtpSender, err := d.ConnData.Connection.AddTrack(trackLocal)
	if err != nil {
		return nil, err
	}

	trackID := NewTrackID(id, streamID, mimeType)
	trackData, err := NewTrackData(trackID, trackLocal, rtpSender, d.ConnData)
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
