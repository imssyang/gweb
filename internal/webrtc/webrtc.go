package webrtc

import (
	"fmt"

	"github.com/imssyang/gweb/internal/conf"
	"github.com/imssyang/gweb/internal/log"
)

var pools map[WebRTCID]*WebRTCPool

func Pool(netType NetType, bindPort bool) (*WebRTCPool, error) {
	for webrtcID, pool := range pools {
		if webrtcID.NetType == netType {
			if bindPort && webrtcID.LocalPort == 0 {
				log.Zap.Debugf("Ignore WebRTCPool: %+v", webrtcID)
				continue
			}
			return pool, nil
		}
	}
	return nil, fmt.Errorf("No valid pool for %s,%v", netType, bindPort)
}

func Init() {
	pools = make(map[WebRTCID]*WebRTCPool)
	iceServerURLs := conf.App.WebRTC.ICEServers
	for i, confPool := range conf.App.WebRTC.Pools {
		webrtcID := NewWebRTCID(
			NetType(confPool.NetType),
			confPool.Address,
			true,
		)
		webrtcPool, err := NewWebRTCPool(webrtcID, iceServerURLs)
		if err != nil {
			log.Zap.Fatalf("[%d] NewWebRTCPool(%+v) error: %s", i, webrtcID, err)
			return
		}

		pools[webrtcID] = webrtcPool
		log.Zap.Debugf("[%d/%d] WebRTCPool(%+v) running", i, len(pools), webrtcID)
		break // TODO
	}
}
