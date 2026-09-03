package media

import (
	"fmt"

	"github.com/imssyang/gweb/pkg/libmedia"
)

type Media struct {
	ID   uint32
	URIs map[string]*Format
}

func NewMedia() (*Media, error) {
	mediaID := libmedia.NewMedia()
	if mediaID == 0 {
		return nil, fmt.Errorf("ffmpeg.NewMedia failed")
	}

	return &Media{
		ID:   mediaID,
		URIs: make(map[string]*Format),
	}, nil
}

func (m *Media) AddDemuxer(uri string) (*Format, error) {
	ffFormat, err := libmedia.NewDemuxer(m.ID, uri)
	if err != nil {
		return nil, fmt.Errorf("libmedia.NewDemuxer(%v) failed", uri)
	}

	return newFormatByFFmpeg(ffFormat), nil
}

func (m *Media) AddMuxer(uri, muxFmt string) (*Format, error) {
	ffFormat, err := libmedia.NewMuxer(m.ID, uri, muxFmt)
	if err != nil {
		return nil, fmt.Errorf("libmedia.NewMuxer(%v) failed", uri)
	}

	return newFormatByFFmpeg(ffFormat), nil
}
