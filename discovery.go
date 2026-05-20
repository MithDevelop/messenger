package main

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type DiscoveryNotifee struct {
	node host.Host
}

func (n *DiscoveryNotifee) HandlePeerFound(info peer.AddrInfo) {

	if info.ID == n.node.ID() {
		return
	}

	if n.node.ID().String() < info.ID.String() {
		return
	}

	fmt.Println("\nFound peer:", info.ID)

	err := n.node.Connect(context.Background(), info)

	if err != nil {
		fmt.Println("Connection failed:", err)
		return
	}

	mu.Lock()

	if _, exists := peers[info.ID.String()]; exists {
		mu.Unlock()
		return
	}

	mu.Unlock()

	stream, err := n.node.NewStream(
		context.Background(),
		info.ID,
		ProtocolID,
	)

	if err != nil {
		fmt.Println("Stream error:", err)
		return
	}

	mu.Lock()

	peers[info.ID.String()] = &PeerInfo{
		ID:        info.ID.String(),
		Stream:    stream,
		Connected: time.Now(),
		LastSeen:  time.Now(),
	}

	mu.Unlock()

	fmt.Println("Chat stream created with:", info.ID)

	go readMessages(info.ID.String(), stream)

	sendUserInfo(info.ID.String())
}
