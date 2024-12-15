package webrtc

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/imssyang/gweb/internal/log"
	media_ "github.com/imssyang/gweb/internal/media"
	"github.com/imssyang/gweb/internal/media/ivf"
	"github.com/imssyang/gweb/internal/media/ogg"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
	"github.com/pion/webrtc/v4/pkg/media/oggreader"
)

type TrackID struct {
	ID       string
	StreamID string
	MimeType string
}

type TrackData struct {
	mu          sync.RWMutex
	TrackID     TrackID
	TrackParser media_.TrackParser
	TrackLocal  *webrtc.TrackLocalStaticSample
	RtpSender   *webrtc.RTPSender
	ConnData    *ConnectionData
}

func NewTrackData(trackID TrackID, trackParser media_.TrackParser, trackLocal *webrtc.TrackLocalStaticSample, rtpSender *webrtc.RTPSender, connData *ConnectionData) (*TrackData, error) {
	if trackLocal == nil || rtpSender == nil || connData == nil {
		return nil, fmt.Errorf("%s invalid.\n", trackID)
	}
	return &TrackData{
		TrackID:     trackID,
		TrackParser: trackParser,
		TrackLocal:  trackLocal,
		RtpSender:   rtpSender,
		ConnData:    connData,
	}, nil
}

func (d *TrackData) OnPacketRead() {
	log.Zap.Infof("OnPacketRead: %v", d.TrackID)
	rtcpBuf := make([]byte, 1500)
	for {
		if n, _, rtcpErr := d.RtpSender.Read(rtcpBuf); rtcpErr != nil {
			log.Zap.Errorf("RtpSender read failed: %v", rtcpErr)
		} else {
			log.Zap.Infof("RtpSender read RTCP: %d", n)
		}
	}
}

func (d *TrackData) OnPacketWrite() {
	log.Zap.Infof("OnPacketWrite: %v", d.TrackID)
	switch parser := d.TrackParser.(type) {
	case ivf.IVFParser:
		timeSecond := float32(parser.Header.TimebaseNumerator) / float32(parser.Header.TimebaseDenominator)
		ticker := time.NewTicker(time.Millisecond * time.Duration(timeSecond*1000))
		defer ticker.Stop()
		log.Zap.Infof("OnPacketWrite IVFHeader: %+v", parser.Header)
		for ; true; <-ticker.C {
			frame, _, ivfErr := parser.NextFrame()
			if errors.Is(ivfErr, io.EOF) {
				fmt.Printf("All video frames parsed and sent")
				os.Exit(0)
			}

			if ivfErr != nil {
				panic(ivfErr)
			}

			if ivfErr = d.TrackLocal.WriteSample(media.Sample{Data: frame, Duration: time.Second}); ivfErr != nil {
				panic(ivfErr)
			} else {
				log.Zap.Infof("OnPacketWrite ivfFrame: %d", len(frame))
			}
		}
	case ogg.OggParser:
		const oggPageDuration = time.Millisecond * 20
		ticker := time.NewTicker(oggPageDuration)
		defer ticker.Stop()
		log.Zap.Infof("OnPacketWrite OGGHeader: %+v", parser.Header)

		// Keep track of last granule, the difference is the amount of samples in the buffer
		var lastGranule uint64
		for ; true; <-ticker.C {
			pageData, pageHeader_, oggErr := parser.NextFrame()
			if errors.Is(oggErr, io.EOF) {
				fmt.Printf("All audio pages parsed and sent")
				os.Exit(0)
			}

			if oggErr != nil {
				panic(oggErr)
			}

			// The amount of samples is the difference between the last and current timestamp
			pageHeader := pageHeader_.(*oggreader.OggPageHeader)
			sampleCount := float64(pageHeader.GranulePosition - lastGranule)
			lastGranule = pageHeader.GranulePosition
			sampleDuration := time.Duration((sampleCount/48000)*1000) * time.Millisecond

			if oggErr = d.TrackLocal.WriteSample(media.Sample{Data: pageData, Duration: sampleDuration}); oggErr != nil {
				panic(oggErr)
			} else {
				log.Zap.Infof("OnPacketWrite oggFrame: %d", len(pageData))
			}
		}
	default:
		log.Zap.Errorf("UnsupportFormat: %v", parser)
		os.Exit(0)
	}
}
