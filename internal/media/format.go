package media

import (
	"fmt"

	"github.com/imssyang/gweb/pkg/libmedia"
)

type Format struct {
	MediaID   uint32
	URI       string
	Name      string
	StartTime int64
	Duration  int64
	BitRate   int64
	Streams   map[int32]*Stream

	ffFuncs libmedia.FormatFuncs
}

func newFormatByFFmpeg(ffFormat *libmedia.Format) *Format {
	streams := make(map[int32]*Stream)
	for _, stream := range ffFormat.Streams {
		streams[stream.Index] = NewStreamByFFmpeg(stream)
	}

	return &Format{
		MediaID:   ffFormat.MediaID,
		URI:       ffFormat.URI,
		Name:      ffFormat.Name,
		StartTime: ffFormat.StartTime,
		Duration:  ffFormat.Duration,
		BitRate:   ffFormat.BitRate,
		Streams:   streams,
		ffFuncs:   ffFormat,
	}
}

func (f *Format) ReadPacket() (*Packet, error) {
	ffPacket, err := f.ffFuncs.ReadPacket()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg.ReadPacket(%v, %v) fail", f.MediaID, f.URI)
	}

	return NewPacketByFFmpeg(ffPacket), nil
}
