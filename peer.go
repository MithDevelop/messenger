package main

import (
	"time"

	"github.com/libp2p/go-libp2p/core/network"
)

type PeerInfo struct {
	ID        string
	Username  string
	Latency   int64
	Online    bool
	Stream    network.Stream
	Connected time.Time
	LastSeen  time.Time
}
