package webrtc

import (
	"fmt"
	"time"

	"github.com/imssyang/gweb/internal/conf"
	"github.com/imssyang/gweb/internal/log"
)

var pools map[WebRTCID]*WebRTCPool

func Init() {
	pools = make(map[WebRTCID]*WebRTCPool)
	confPools := conf.App.WebRTC.Pools
	confICEServerURLs := conf.App.WebRTC.ICEServers
	for i, confPool := range confPools {
		webrtcID := NewWebRTCID(
			NetworkType(confPool.Network),
			confPool.Address,
			true,
		)
		webrtcPool, err := NewWebRTCPool(webrtcID, confICEServerURLs)
		if err != nil {
			log.Zap.Fatalf("[%d] NewWebRTCPool(%+v) error: %s", i, webrtcID, err)
			return
		}

		pools[webrtcID] = webrtcPool
		log.Zap.Debugf("[%d/%d] WebRTCPool(%+v) running", i+1, len(confPools), webrtcID)
		break // TODO
	}
}

func Pool(network NetworkType, bindPort bool) (*WebRTCPool, error) {
	for webrtcID, pool := range pools {
		if webrtcID.Network == network {
			if bindPort && webrtcID.LocalPort == 0 {
				log.Zap.Debugf("Ignore WebRTCPool: %+v", webrtcID)
				continue
			}
			return pool, nil
		}
	}
	return nil, fmt.Errorf("No valid pool for %s,%v", network, bindPort)
}

func GetConnection(connID ConnectionID, timeout time.Duration) (*ConnectionData, error) {
	timeoutTimer := time.NewTimer(timeout)
	ticker := time.NewTicker(50 * time.Millisecond)

	for {
		select {
		case <-timeoutTimer.C:
			return nil, fmt.Errorf("Timeout to get connect for ConnID[%v]", connID)
		case <-ticker.C:
			for _, pool := range pools {
				connData := pool.GetConnectionData(connID)
				if connData != nil {
					return connData, nil
				}
			}
		}
	}
}
