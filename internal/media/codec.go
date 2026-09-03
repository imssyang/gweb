package media

import (
	"math/big"

	"github.com/imssyang/gweb/pkg/libmedia"
	"github.com/pion/webrtc/v4"
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

func NewOnlyVideoCodecByFFmpeg(ffCodec *libmedia.OnlyVideoCodec) *OnlyVideoCodec {
	if ffCodec == nil {
		return nil
	}

	return &OnlyVideoCodec{
		Width:                       ffCodec.Width,
		Height:                      ffCodec.Height,
		SampleAspectRatio:           ffCodec.SampleAspectRatio,
		FrameRate:                   ffCodec.FrameRate,
		PixelFormat:                 ffCodec.PixelFormat,
		FieldOrder:                  ffCodec.FieldOrder,
		ColorRange:                  ffCodec.ColorRange,
		ColorPrimaries:              ffCodec.ColorPrimaries,
		ColorTransferCharacteristic: ffCodec.ColorTransferCharacteristic,
		ColorSpace:                  ffCodec.ColorSpace,
		ChromaLocation:              ffCodec.ChromaLocation,
		FrameDelay:                  ffCodec.FrameDelay,
	}
}

type OnlyAudioCodec struct {
	SampleFormat  string
	SampleRate    int32
	ChannelNum    int32
	ChannelLayout string
	FrameSize     int32
}

func NewOnlyAudioCodecByFFmpeg(ffCodec *libmedia.OnlyAudioCodec) *OnlyAudioCodec {
	if ffCodec == nil {
		return nil
	}

	return &OnlyAudioCodec{
		SampleFormat:  ffCodec.SampleFormat,
		SampleRate:    ffCodec.SampleRate,
		ChannelNum:    ffCodec.ChannelNum,
		ChannelLayout: ffCodec.ChannelLayout,
		FrameSize:     ffCodec.FrameSize,
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

func NewCodecByFFmpeg(ffCodec *libmedia.Codec) *Codec {
	sideDatas := make([]*PacketSideData, 0)
	for _, sd := range ffCodec.SideDatas {
		sideDatas = append(sideDatas, NewPacketSideDataByFFmpeg(sd))
	}

	return &Codec{
		MediaID:     ffCodec.MediaID,
		URI:         ffCodec.URI,
		StreamIndex: ffCodec.StreamIndex,
		Type:        ffCodec.Type,
		Name:        ffCodec.Name,
		Tag:         ffCodec.Tag,
		Profile:     ffCodec.Profile,
		Level:       ffCodec.Level,
		BitRate:     ffCodec.BitRate,
		ExtraData:   ffCodec.ExtraData,
		SideDatas:   sideDatas,
		Video:       NewOnlyVideoCodecByFFmpeg(ffCodec.Video),
		Audio:       NewOnlyAudioCodecByFFmpeg(ffCodec.Audio),
	}
}

func (c *Codec) MimeType() string {
	switch c.Name {
	case "h264":
		return webrtc.MimeTypeH264
	case "hevc":
		return webrtc.MimeTypeH265
	case "opus":
		return webrtc.MimeTypeOpus
	case "vp8":
		return webrtc.MimeTypeVP8
	case "vp9":
		return webrtc.MimeTypeVP9
	case "av1":
		return webrtc.MimeTypeAV1
	case "g722":
		return webrtc.MimeTypeG722
	case "pcm_mulaw":
		return webrtc.MimeTypePCMU
	case "pcm_alaw":
		return webrtc.MimeTypePCMA
	default:
		return ""
	}
}
