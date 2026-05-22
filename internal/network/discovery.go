package network

import (
	"context"
	"fmt"
	"time"

	"messenger/internal/peer"
	"messenger/internal/state"

	"github.com/libp2p/go-libp2p/core/host"

	libpeer "github.com/libp2p/go-libp2p/core/peer"
)

type DiscoveryNotifee struct {
	Node host.Host
}

func (n *DiscoveryNotifee) HandlePeerFound(info libpeer.AddrInfo) {

	if info.ID == n.Node.ID() {
		return
	}

	if n.Node.ID().String() < info.ID.String() {
		return
	}

	fmt.Println("\nFound peer:", info.ID)

	err := n.Node.Connect(context.Background(), info)

	if err != nil {
		fmt.Println("Connection failed:", err)
		return
	}

	state.Mu.Lock()

	if _, exists := state.Peers[info.ID.String()]; exists {
		state.Mu.Unlock()
		return
	}

	state.Mu.Unlock()

	stream, err := n.Node.NewStream(
		context.Background(),
		info.ID,
		state.ProtocolID,
	)

	if err != nil {
		fmt.Println("Stream error:", err)
		return
	}

	state.Mu.Lock()

	state.Peers[info.ID.String()] = &peer.PeerInfo{
		ID:        info.ID.String(),
		Stream:    stream,
		Connected: time.Now(),
		LastSeen:  time.Now(),
		Latency:   0,
		Online:    true,
	}

	state.Mu.Unlock()

	fmt.Println("Chat stream created with:", info.ID)

	go ReadMessages(info.ID.String(), stream)

	SendUserInfo(info.ID.String())
}
