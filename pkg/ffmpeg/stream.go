package ffmpeg

// #include "ffmpeg.h"
import "C"
import (
	"fmt"
	"math/big"
	"unsafe"
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

func NewStream(mediaID uint32, uri string, index int32) (*Stream, error) {
	cMediaID := C.uint32_t(mediaID)
	cURI := C.CString(uri)
	cIndex := C.int32_t(index)
	defer C.free(unsafe.Pointer(cURI))

	cStream := C.GetStream(cMediaID, cURI, cIndex)
	if cStream == nil {
		return nil, fmt.Errorf("C.GetStream(%v, %v, %v) fail", mediaID, uri, index)
	}

	codec, err := NewCodec(mediaID, uri, index)
	if err != nil {
		return nil, fmt.Errorf("NewCodec(%v, %v, %v) fail", mediaID, uri, index)
	}

	return &Stream{
		MediaID: mediaID,
		URI:     uri,
		Index:   index,
		TimeBase: *big.NewRat(
			int64(cStream.time_base.num),
			int64(cStream.time_base.den)),
		StartTime: int64(cStream.start_time),
		Duration:  int64(cStream.duration),
		FrameNum:  int64(cStream.nb_frames),
		Codec:     codec,
	}, nil
}
