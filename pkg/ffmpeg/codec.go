package ffmpeg

// #include "libmedia/ffmpeg.h"
import "C"
import (
	"fmt"
	"math/big"
	"unsafe"
)

type OnlyVideoCodec struct {
	Width                       int32
	Height                      int32
	SampleAspectRatio           big.Rat
	FrameRate                   big.Rat
	PixelFormat                 string
	FieldOrder                  string
	ColorRange                  string
	ColorPrimaries              string
	ColorTransferCharacteristic string
	ColorSpace                  string
	ChromaLocation              string
	FrameDelay                  int32
}

func NewOnlyVideoCodec(cParam *C.AVCodecParameters) *OnlyVideoCodec {
	cPixelFormat := C.av_get_pix_fmt_name(C.enum_AVPixelFormat(cParam.format))
	cFieldOrder := C.GetFieldOrderStr(cParam.field_order)
	cColorRange := C.av_color_range_name(cParam.color_range)
	cColorPrimaries := C.av_color_primaries_name(cParam.color_primaries)
	cColorTransferCharacteristic := C.av_color_transfer_name(cParam.color_trc)
	cColorSpace := C.av_color_space_name(cParam.color_space)
	cChromaLocation := C.av_chroma_location_name(cParam.chroma_location)

	return &OnlyVideoCodec{
		Width:  int32(cParam.width),
		Height: int32(cParam.height),
		SampleAspectRatio: *big.NewRat(
			int64(cParam.sample_aspect_ratio.num),
			int64(cParam.sample_aspect_ratio.den)),
		FrameRate: *big.NewRat(
			int64(cParam.framerate.num),
			int64(cParam.framerate.den)),
		PixelFormat:                 C.GoString(cPixelFormat),
		FieldOrder:                  C.GoString(cFieldOrder),
		ColorRange:                  C.GoString(cColorRange),
		ColorPrimaries:              C.GoString(cColorPrimaries),
		ColorTransferCharacteristic: C.GoString(cColorTransferCharacteristic),
		ColorSpace:                  C.GoString(cColorSpace),
		ChromaLocation:              C.GoString(cChromaLocation),
		FrameDelay:                  int32(cParam.video_delay),
	}
}

type OnlyAudioCodec struct {
	SampleFormat  string
	SampleRate    int32
	ChannelNum    int32
	ChannelLayout string
	FrameSize     int32
}

func NewOnlyAudioCodec(cParam *C.AVCodecParameters) *OnlyAudioCodec {
	cSampleFormat := C.av_get_sample_fmt_name(C.enum_AVSampleFormat(cParam.format))

	var channelLayout string
	cChannelLayout := C.AllocChannelLayoutStr(&cParam.ch_layout)
	if cChannelLayout != nil {
		channelLayout = C.GoString(cChannelLayout)
		C.free(unsafe.Pointer(cChannelLayout))
	}

	return &OnlyAudioCodec{
		SampleFormat:  C.GoString(cSampleFormat),
		SampleRate:    int32(cParam.sample_rate),
		ChannelNum:    int32(cParam.ch_layout.nb_channels),
		ChannelLayout: channelLayout,
		FrameSize:     int32(cParam.frame_size),
	}
}

type Codec struct {
	MediaID     uint32
	URI         string
	StreamIndex int32
	Type        string
	Name        string
	Tag         uint32
	Profile     int32
	Level       int32
	BitRate     int64
	ExtraData   []byte
	SideDatas   []*PacketSideData
	Video       *OnlyVideoCodec
	Audio       *OnlyAudioCodec
}

func NewCodec(mediaID uint32, uri string, streamIndex int32) (*Codec, error) {
	cMediaID := C.uint32_t(mediaID)
	cURI := C.CString(uri)
	cStreamIndex := C.int32_t(streamIndex)
	defer C.free(unsafe.Pointer(cURI))

	cCodecpar := C.GetCodecParameters(cMediaID, cURI, cStreamIndex)
	if cCodecpar == nil {
		return nil, fmt.Errorf("C.GetCodecParameters(%v, %v, %v) fail", mediaID, uri, streamIndex)
	}

	extraData := make([]byte, 0)
	if cCodecpar.extradata != nil && int32(cCodecpar.extradata_size) > 0 {
		extraData = C.GoBytes(
			unsafe.Pointer(cCodecpar.extradata),
			C.int(cCodecpar.extradata_size))
	}

	sideDatas := make([]*PacketSideData, 0)
	cSideDatas := unsafe.Slice(
		(*C.AVPacketSideData)(cCodecpar.coded_side_data),
		C.int(cCodecpar.nb_coded_side_data))
	for _, cSideData := range cSideDatas {
		sideData, err := NewPacketSideData(&cSideData)
		if err != nil {
			return nil, fmt.Errorf("NewPacketSideData fail")
		}
		sideDatas = append(sideDatas, sideData)
	}

	return &Codec{
		MediaID:     mediaID,
		URI:         uri,
		StreamIndex: streamIndex,
		Type:        C.GoString(C.av_get_media_type_string(cCodecpar.codec_type)),
		Name:        C.GoString(C.avcodec_get_name(cCodecpar.codec_id)),
		Tag:         uint32(cCodecpar.codec_tag),
		Profile:     int32(cCodecpar.profile),
		Level:       int32(cCodecpar.level),
		BitRate:     int64(cCodecpar.bit_rate),
		ExtraData:   extraData,
		SideDatas:   sideDatas,
		Video:       NewOnlyVideoCodec(cCodecpar),
		Audio:       NewOnlyAudioCodec(cCodecpar),
	}, nil
}
