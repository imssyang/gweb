package ffmpeg

// #include "libmedia/ffmpeg.h"
import "C"
import (
	"fmt"
	"unsafe"
)

type FormatFuncs interface {
	ReadPacket() (*Packet, error)
}

type Format struct {
	MediaID   uint32
	URI       string
	Name      string
	StartTime int64
	Duration  int64
	BitRate   int64
	Streams   map[int32]*Stream
}

func NewFormat(mediaID uint32, uri string) (*Format, error) {
	cMediaID := C.uint32_t(mediaID)
	cURI := C.CString(uri)
	defer C.free(unsafe.Pointer(cURI))

	cCtx := C.GetFormatContext(cMediaID, cURI)
	if cCtx == nil {
		return nil, fmt.Errorf("C.GetFormatContext(%v, %v) fail", mediaID, uri)
	}

	var name string
	if cCtx.iformat != nil {
		name = C.GoString(cCtx.iformat.name)
	} else if cCtx.oformat != nil {
		name = C.GoString(cCtx.oformat.name)
	}

	streams := make(map[int32]*Stream)
	for i := 0; i < int(cCtx.nb_streams); i++ {
		stream, err := NewStream(mediaID, uri, int32(i))
		if err != nil {
			return nil, fmt.Errorf("NewStream(%v, %v, %v) fail", mediaID, uri, i)
		}
		streams[stream.Index] = stream
	}

	return &Format{
		MediaID:   mediaID,
		URI:       uri,
		Name:      name,
		StartTime: int64(cCtx.start_time),
		Duration:  int64(cCtx.duration),
		BitRate:   int64(cCtx.bit_rate),
		Streams:   streams,
	}, nil
}

func (f *Format) ReadPacket() (*Packet, error) {
	cMediaID := C.uint32_t(f.MediaID)
	cURI := C.CString(f.URI)
	defer C.free(unsafe.Pointer(cURI))

	cPacket := C.ReadPacket(cMediaID, cURI)
	if cPacket == nil {
		return nil, fmt.Errorf("C.ReadPacket(%v, %v) fail", f.MediaID, f.URI)
	}

	defer C.FreePacket(cPacket)
	return NewPacket(f.MediaID, f.URI, cPacket)
}
