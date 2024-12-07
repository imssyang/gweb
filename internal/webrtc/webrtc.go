package webrtc

import (
	"net"
	"strconv"

	"github.com/imssyang/gweb/internal/conf"
	"github.com/imssyang/gweb/internal/log"
)

var Pool *WebRTCPool

func Init() {
	for i, webrtcConf := range conf.App.WebRTC {
		address := webrtcConf.Address
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			log.Zap.Fatalf("WebRTC[%d](%s) run error: %s", i, address, err)
			return
		}

		portInt, err := strconv.Atoi(port)
		netType := NetType(webrtcConf.NetType)
		pool, err := NewWebRTCPool(netType, host, portInt, false)
		if err != nil {
			log.Zap.Fatalf("NewWebRTCPool[%d](%s) run error: %s", i, address, err)
			return
		}

		Pool = pool
		log.Zap.Debugf("WebRTCPool: %+v", Pool)
	}
}
