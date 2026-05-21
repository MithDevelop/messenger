package main

import (
	"fmt"
	"time"
)

func startPresenceLoop() {

	for {

		time.Sleep(5 * time.Second)

		sendPingToAll()
	}
}

func sendPingToAll() {

	mu.Lock()

	peerIDs := make([]string, 0)

	for id := range peers {
		peerIDs = append(peerIDs, id)
	}

	mu.Unlock()

	for _, peerID := range peerIDs {

		msg := Message{
			ID:        generateMessageID(),
			Type:      "ping",
			From:      localPeerID,
			To:        peerID,
			Username:  username,
			Message:   "",
			Timestamp: time.Now().UnixMilli(),
		}

		sendToPeer(peerID, msg, false)
	}
}

func monitorPeers() {

	for {

		time.Sleep(10 * time.Second)

		mu.Lock()

		for id, peer := range peers {

			if time.Since(peer.LastSeen) > 15*time.Second {

				fmt.Println("\nPeer timed out:", id)

				peer.Online = false

				delete(peers, id)
			}
		}

		mu.Unlock()
	}
}
