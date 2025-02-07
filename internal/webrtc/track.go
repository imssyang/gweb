package webrtc

import (
	"fmt"
	"sync"

	"github.com/imssyang/gweb/internal/log"
	media_ "github.com/imssyang/gweb/internal/media"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
)

type TrackParam struct {
	URI      string
	ID       string
	MimeType string
}

func NewTrackParam(uri string, id int32, mimeType string) TrackParam {
	return TrackParam{
		URI:      uri,
		ID:       fmt.Sprint(id),
		MimeType: mimeType,
	}
}

type TrackData struct {
	mu          sync.RWMutex
	Param       TrackParam
	Media       *media_.Stream
	LocalSample *webrtc.TrackLocalStaticSample
	RtpSender   *webrtc.RTPSender
	ConnData    *ConnectionData
}

func NewTrackData(param TrackParam, media *media_.Stream, localSample *webrtc.TrackLocalStaticSample, rtpSender *webrtc.RTPSender, connData *ConnectionData) (*TrackData, error) {
	if media == nil || localSample == nil || rtpSender == nil || connData == nil {
		return nil, fmt.Errorf("%s invalid.\n", param)
	}
	return &TrackData{
		Param:       param,
		Media:       media,
		LocalSample: localSample,
		RtpSender:   rtpSender,
		ConnData:    connData,
	}, nil
}

func (d *TrackData) OnRtcpRead() ([]rtcp.Packet, error) {
	pkts, _, err := d.RtpSender.ReadRTCP()
	if err != nil {
		log.Zap.Errorf("%v read rtcp fail: %v", d.Param, err)
		return nil, err
	}

	log.Zap.Infof("%v read rtcp: %+v", d.Param, pkts)
	return pkts, nil
}

func (d *TrackData) OnPacketWrite(packet *media_.Packet) error {
	err := d.LocalSample.WriteSample(media.Sample{
		Data:     packet.Data,
		Duration: packet.DurationT(),
	})
	if err != nil {
		return err
	}

	log.Zap.Infof("OnPacketWrite %v:%v:%v:%v:%d:%+v",
		d.Param.URI,
		d.Param.ID,
		d.Media.Codec.Type,
		d.Media.Codec.Name,
		len(packet.Data),
		packet.DurationT())
	return nil
}
