package media

import (
	"fmt"

	"github.com/imssyang/gweb/pkg/ffmpeg"
)

type Media struct {
	ID   uint32
	URIs map[string]*Format
}

func NewMedia() (*Media, error) {
	mediaID := ffmpeg.NewMedia()
	if mediaID == 0 {
		return nil, fmt.Errorf("ffmpeg.NewMedia failed")
	}

	return &Media{
		ID:   mediaID,
		URIs: make(map[string]*Format),
	}, nil
}

func (m *Media) AddDemuxer(uri string) (*Format, error) {
	ffFormat, err := ffmpeg.NewDemuxer(m.ID, uri)
	if err != nil {
		return nil, fmt.Errorf("ffmpeg.AddDemuxer(%v) failed", uri)
	}

	return newFormatByFFmpeg(ffFormat), nil
}

func (m *Media) AddMuxer(uri, muxFmt string) (*Format, error) {
	ffFormat, err := ffmpeg.NewMuxer(m.ID, uri, muxFmt)
	if err != nil {
		return nil, fmt.Errorf("ffmpeg.AddMuxer(%v) failed", uri)
	}

	return newFormatByFFmpeg(ffFormat), nil
}
