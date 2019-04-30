package model

import (
	"time"
)

var Self *PeerInfo

type PeerInfo struct {
	PeerId   int
	Version  int
	SyncTime time.Time
	Alive    bool
	Seed     bool
}

func NewPeer(port int) *PeerInfo {
	return &PeerInfo{
		PeerId: port,
		Version: 0,
		Alive: true,
		SyncTime: time.Now(),
	}
}
