package media

import (
	"math/big"

	"github.com/imssyang/gweb/pkg/libmedia"
)

type Stream struct {
	MediaID   uint32
	URI       string
	Index     int32
	TimeBase  big.Rat
	StartTime int64
	Duration  int64
	FrameNum  int64
	Codec     *Codec
}

func NewStreamByFFmpeg(ffStream *libmedia.Stream) *Stream {
	return &Stream{
		MediaID:   ffStream.MediaID,
		URI:       ffStream.URI,
		Index:     ffStream.Index,
		TimeBase:  ffStream.TimeBase,
		StartTime: ffStream.StartTime,
		Duration:  ffStream.StartTime,
		FrameNum:  ffStream.FrameNum,
		Codec:     NewCodecByFFmpeg(ffStream.Codec),
	}
}
