package webrtc

import (
	"fmt"
	"sync"
	"time"

	"github.com/imssyang/gweb/internal/conf"
	"github.com/imssyang/gweb/internal/log"
)

type Instance struct {
	mu    sync.RWMutex
	pools map[PoolID]*PoolData
}

func NewInstanceByConf() *Instance {
	inst := &Instance{
		pools: make(map[PoolID]*PoolData),
	}

	pools := conf.App.WebRTC.Pools
	iceServerURLs := conf.App.WebRTC.ICEServers
	for i, pool := range pools {
		poolID := NewPoolID(
			NetworkType(pool.Network),
			pool.Address,
			true,
		)
		poolData, err := NewPoolData(poolID, iceServerURLs)
		if err != nil {
			log.Zap.Fatalf("[%d] NewPoolData(%+v) error: %s", i, poolID, err)
			panic(err)
		}

		inst.pools[poolID] = poolData
		log.Zap.Debugf("[%d/%d] PoolData(%+v) running", i+1, len(pools), poolID)
		break // TODO
	}
	return inst
}

func (i *Instance) GetPool(network NetworkType, bindPort bool) (*PoolData, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	for poolID, pool := range i.pools {
		if poolID.Network == network {
			if bindPort && poolID.LocalPort == 0 {
				log.Zap.Debugf("Ignore PoolData: %+v", poolID)
				continue
			}
			return pool, nil
		}
	}
	return nil, fmt.Errorf("No valid pool for %s,%v", network, bindPort)
}

func (i *Instance) GetConnection(connID ConnectionID) *ConnectionData {
	i.mu.RLock()
	defer i.mu.RUnlock()
	for _, pool := range i.pools {
		connData := pool.GetConnection(connID)
		if connData != nil {
			return connData
		}
	}
	return nil
}

func (i *Instance) GetConnectionWithTimeout(connID ConnectionID, timeout time.Duration) (*ConnectionData, error) {
	timeoutTimer := time.NewTimer(timeout)
	ticker := time.NewTicker(50 * time.Millisecond)

	for {
		select {
		case <-timeoutTimer.C:
			return nil, fmt.Errorf("Timeout to get connect for ConnID[%v]", connID)
		case <-ticker.C:
			connData := i.GetConnection(connID)
			if connData != nil {
				return connData, nil
			}
		}
	}
}

var instByConf *Instance

func Init() {
	instByConf = NewInstanceByConf()
}

func Pool(network NetworkType, bindPort bool) (*PoolData, error) {
	return instByConf.GetPool(network, bindPort)
}

func Connection(connID ConnectionID) *ConnectionData {
	return instByConf.GetConnection(connID)
}

func ConnectionWithTimeout(connID ConnectionID, timeout time.Duration) (*ConnectionData, error) {
	return instByConf.GetConnectionWithTimeout(connID, timeout)
}
