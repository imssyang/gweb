package webrtc

import (
	"fmt"
	"sync"

	"github.com/pion/webrtc/v4"
)

type ChannelID string

type ChannelData struct {
	mu       sync.RWMutex
	ChanID   ChannelID
	Channel  *webrtc.DataChannel
	ConnData *ConnectionData
}

func NewChannelData(chanID ChannelID, channel *webrtc.DataChannel, connData *ConnectionData) (*ChannelData, error) {
	if channel == nil || connData == nil {
		return nil, fmt.Errorf("%s invalid.", chanID)
	}
	return &ChannelData{
		ChanID:   chanID,
		Channel:  channel,
		ConnData: connData,
	}, nil
}

func (d *ChannelData) OnDataChannelOpen() {
	fmt.Printf("Data channel '%s'-'%d' open.\n", d.Channel.Label(), d.Channel.ID())

	//enableDetach := d.ConnData.Pool.ID.EnableDetach
	//if enableDetach {
	//	fmt.Printf("Data channel Detach '%s' open.\n", d.ChanID)
	//	rwIO, err := d.Channel.Detach()
	//	if err != nil {
	//		return
	//	}
	//	go func(d io.Reader) {
	//		fmt.Printf("Message from DataChannel\n")
	//	}(rwIO)
	//	go func(d io.Writer) {
	//		fmt.Printf("Sending Message\n")
	//	}(rwIO)
	//} else {
	//	// Send the current time via a DataChannel to the remote peer every 3 seconds
	//	for range time.Tick(time.Second * 3) {
	//		if err := d.Channel.SendText(time.Now().String()); err != nil {
	//			if errors.Is(err, io.ErrClosedPipe) {
	//				return
	//			}
	//		}
	//	}
	//}
}

func (d *ChannelData) OnDataChannelMessage(msg webrtc.DataChannelMessage) {
	fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Channel.Label(), string(msg.Data))
}
