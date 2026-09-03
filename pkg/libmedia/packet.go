package libmedia

// #include "libmedia/media.h"
import "C"
import (
	"fmt"
	"math/big"
	"unsafe"
)

type PacketSideData struct {
	Data []byte
	Type string
}

func NewPacketSideData(cSD *C.AVPacketSideData) (*PacketSideData, error) {
	data := make([]byte, 0)
	if cSD.data != nil && uint32(cSD.size) > 0 {
		data = C.GoBytes(
			unsafe.Pointer(cSD.data),
			C.int(cSD.size))
	}

	var SDType string
	cSDType := C.GetPacketSideDataTypeStr(cSD)
	if cSDType != nil {
		SDType = C.GoString(cSDType)
	}

	return &PacketSideData{
		Data: data,
		Type: SDType,
	}, nil
}

type Packet struct {
	MediaID     uint32
	URI         string
	Pts         int64
	Dts         int64
	Data        []byte
	StreamIndex int32
	Flags       int32
	SideDatas   []*PacketSideData
	Duration    int64
	Pos         int64
	TimeBase    big.Rat
}

func NewPacket(mediaID uint32, uri string, cPacket *C.AVPacket) (*Packet, error) {
	data := make([]byte, 0)
	if cPacket.data != nil && int32(cPacket.size) > 0 {
		data = C.GoBytes(
			unsafe.Pointer(cPacket.data),
			C.int(cPacket.size))
	}

	sideDatas := make([]*PacketSideData, 0)
	for i := 0; i < int(cPacket.side_data_elems); i++ {
		sideData, err := NewPacketSideData(cPacket.side_data)
		if err != nil {
			return nil, fmt.Errorf("NewPacketSideData fail")
		}
		sideDatas = append(sideDatas, sideData)
	}

	return &Packet{
		MediaID:     mediaID,
		URI:         uri,
		Pts:         int64(cPacket.pts),
		Dts:         int64(cPacket.dts),
		Data:        data,
		StreamIndex: int32(cPacket.stream_index),
		Flags:       int32(cPacket.flags),
		SideDatas:   sideDatas,
		Duration:    int64(cPacket.duration),
		Pos:         int64(cPacket.pos),
		TimeBase: *big.NewRat(
			int64(cPacket.time_base.num),
			int64(cPacket.time_base.den)),
	}, nil
}
