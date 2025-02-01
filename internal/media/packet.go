package media

import (
	"math/big"
	"time"

	"github.com/imssyang/gweb/pkg/ffmpeg"
)

type PacketSideData struct {
	Data []byte
	Type string
}

func NewPacketSideDataByFFmpeg(ffSD *ffmpeg.PacketSideData) *PacketSideData {
	return &PacketSideData{
		Data: ffSD.Data,
		Type: ffSD.Type,
	}
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

func NewPacketByFFmpeg(ffPacket *ffmpeg.Packet) *Packet {
	sideDatas := make([]*PacketSideData, 0)
	for _, sd := range ffPacket.SideDatas {
		sideDatas = append(sideDatas, NewPacketSideDataByFFmpeg(sd))
	}

	return &Packet{
		MediaID:     ffPacket.MediaID,
		URI:         ffPacket.URI,
		StreamIndex: ffPacket.StreamIndex,
		Pts:         ffPacket.Pts,
		Dts:         ffPacket.Dts,
		Data:        ffPacket.Data,
		Flags:       ffPacket.Flags,
		SideDatas:   sideDatas,
		Duration:    ffPacket.Duration,
		Pos:         ffPacket.Pos,
		TimeBase:    ffPacket.TimeBase,
	}
}

func (p *Packet) toTimestamp(timestamp int64) time.Duration {
	ptsSec := new(big.Rat).Mul(big.NewRat(timestamp, 1), &p.TimeBase)
	ptsMicrosec := ptsSec.Mul(ptsSec, big.NewRat(1_000_000, 1))
	pts, _ := ptsMicrosec.Float64()
	return time.Duration(int64(pts)) * time.Microsecond
}

func (p *Packet) PtsT() time.Duration {
	return p.toTimestamp(p.Pts)
}

func (p *Packet) DtsT() time.Duration {
	return p.toTimestamp(p.Dts)
}

func (p *Packet) DurationT() time.Duration {
	return p.toTimestamp(p.Duration)
}
