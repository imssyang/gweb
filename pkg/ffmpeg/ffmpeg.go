package ffmpeg

// #include "libmedia/ffmpeg.h"
import "C"
import (
	"fmt"
	"unsafe"
)

func NewMedia() uint32 {
	return uint32(C.NewMedia())
}

func NewDemuxer(mediaID uint32, uri string) (*Format, error) {
	cMediaID := C.uint32_t(mediaID)
	cURI := C.CString(uri)
	defer C.free(unsafe.Pointer(cURI))

	ok := bool(C.AddDemuxer(cMediaID, cURI))
	if !ok {
		return nil, fmt.Errorf("C.AddDemuxer(%v, %v) fail", mediaID, uri)
	}

	return NewFormat(mediaID, uri)
}

func NewMuxer(mediaID uint32, uri, muxFmt string) (*Format, error) {
	cMediaID := C.uint32_t(mediaID)
	cURI := C.CString(uri)
	defer C.free(unsafe.Pointer(cURI))
	cMuxFmt := C.CString(muxFmt)
	defer C.free(unsafe.Pointer(cMuxFmt))

	ok := bool(C.AddMuxer(cMediaID, cURI, cMuxFmt))
	if !ok {
		return nil, fmt.Errorf("C.AddMuxer(%v, %v, %v) fail", mediaID, uri, muxFmt)
	}

	return NewFormat(mediaID, uri)
}
