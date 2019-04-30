package handler

import (
	"github.com/mmpei/gossip/model"
	"sync"
	"time"
)

var Manager = NewPeerManager()

type PeerManager struct {
	sync.Mutex
	peerMap   map[int]*model.PeerInfo
}

func NewPeerManager() *PeerManager {
	return &PeerManager{
		peerMap: make(map[int]*model.PeerInfo),
	}
}

func (pm *PeerManager) Set(peerInfo *model.PeerInfo) error {
	pm.Lock()
	defer pm.Unlock()

	if peer, ok := pm.peerMap[peerInfo.PeerId]; ok {
		if peer.Version < peerInfo.Version {
			pm.peerMap[peerInfo.PeerId] = peerInfo
		}
	} else {
		pm.peerMap[peerInfo.PeerId] = peerInfo
	}
	return nil
}

func (pm *PeerManager) Get() (list []model.PeerInfo) {
	pm.Lock()
	defer pm.Unlock()

	for id, peer := range pm.peerMap {
		if id == model.Self.PeerId {
			continue
		}
		list = append(list, *peer)
	}
	return
}

func (pm *PeerManager) SetDown(peerId int) error {
	pm.Lock()
	defer pm.Unlock()

	if peer, ok := pm.peerMap[peerId]; ok {
		peer.Alive = false
		peer.SyncTime = time.Now()
	} else {
		pm.peerMap[peerId] = &model.PeerInfo{
			PeerId: peerId,
			Version: 0,
			Alive: false,
			SyncTime: time.Now(),
			Seed: false,
		}
	}
	return nil
}
