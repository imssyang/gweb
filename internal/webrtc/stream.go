package webrtc

import (
	"fmt"
	"sync"
	"time"

	"github.com/imssyang/gweb/internal/log"
	"github.com/imssyang/gweb/internal/media"
	"github.com/pion/webrtc/v4"
)

type StreamData struct {
	mu       sync.RWMutex
	URI      string
	Media    *media.Format
	ConnData *ConnectionData
	tracks   map[string]*TrackData
}

func NewStreamData(uri string, connData *ConnectionData) (*StreamData, error) {
	demuxer, err := connData.Media.AddDemuxer(uri)
	if err != nil {
		return nil, err
	}

	d := &StreamData{
		URI:      uri,
		Media:    demuxer,
		ConnData: connData,
		tracks:   make(map[string]*TrackData),
	}

	for streamIndex, stream := range d.Media.Streams {
		trackID := NewTrackParam(uri, streamIndex, stream.Codec.MimeType())
		_, err := d.AddTrack(trackID, stream)
		if err != nil {
			return nil, err
		}
	}

	return d, nil
}

func (d *StreamData) OnDataTracksOpen() {
	d.mu.RLock()
	defer d.mu.RUnlock()
	log.Zap.Infof("URI: %s NumOfTrack: %d", d.URI, len(d.tracks))

	go func() {
		var startDts *time.Duration
		startPoint := time.Now()
		for {
			packet, err := d.Media.ReadPacket()
			if err != nil {
				log.Zap.Warnf("URI: %s read packet fail", d.URI)
				break
			}

			for {
				dts := packet.DtsT()
				if startDts == nil {
					startDts = &dts
					break
				}

				duration := packet.DurationT()
				elapsedDts := dts - *startDts
				elapsedT := time.Since(startPoint)
				if elapsedT+duration > elapsedDts {
					break
				}

				time.Sleep(duration / 5)
			}

			trackData := d.GetTrack(fmt.Sprint(packet.StreamIndex))
			err = trackData.OnPacketWrite(packet)
			if err != nil {
				log.Zap.Warnf("URI: %s write packet fail", d.URI)
				break
			}
		}
	}()

	for _, track := range d.tracks {
		go func(t *TrackData) {
			for {
				_, err := t.OnRtcpRead()
				if err != nil {
					break
				}
			}
		}(track)
	}
}

func (d *StreamData) AddTrack(param TrackParam, media *media.Stream) (*TrackData, error) {
	localSample, err := webrtc.NewTrackLocalStaticSample(
		webrtc.RTPCodecCapability{MimeType: param.MimeType},
		param.ID,
		param.URI)
	if err != nil {
		return nil, err
	}

	rtpSender, err := d.ConnData.Connection.AddTrack(localSample)
	if err != nil {
		return nil, err
	}

	trackData, err := NewTrackData(param, media, localSample, rtpSender, d.ConnData)
	if err != nil {
		return nil, err
	}

	d.mu.Lock()
	d.tracks[param.ID] = trackData
	d.mu.Unlock()

	return trackData, nil
}

func (d *StreamData) GetTrack(trackID string) *TrackData {
	d.mu.RLock()
	defer d.mu.RUnlock()
	td, exists := d.tracks[trackID]
	if exists {
		return td
	}
	return nil
}
