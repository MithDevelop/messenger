package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
)

func readMessages(peerID string, stream network.Stream) {

	decoder := json.NewDecoder(stream)

	for {

		var msg Message

		err := decoder.Decode(&msg)

		if err != nil {

			fmt.Println("\nConnection closed with:", peerID)

			mu.Lock()
			delete(peers, peerID)
			mu.Unlock()

			return
		}

		mu.Lock()

		peerInfo, exists := peers[peerID]

		if exists {
			peerInfo.Username = msg.Username
			peerInfo.LastSeen = time.Now()
		}

		mu.Unlock()

		handleMessage(peerID, msg)
	}
}

func sendUserInfo(peerID string) {

	msg := Message{
		Type:      "user_info",
		From:      "",
		To:        peerID,
		Username:  username,
		Message:   "",
		Timestamp: time.Now().Unix(),
	}

	sendToPeer(peerID, msg)
}

func handleStream(stream network.Stream) {

	peerID := stream.Conn().RemotePeer().String()

	fmt.Println("\nIncoming chat connection from:", peerID)

	mu.Lock()

	if _, exists := peers[peerID]; !exists {

		peers[peerID] = &PeerInfo{
			ID:        peerID,
			Stream:    stream,
			Connected: time.Now(),
			LastSeen:  time.Now(),
		}
	}

	mu.Unlock()

	go readMessages(peerID, stream)

	sendUserInfo(peerID)
}

func sendToAll(msg Message) {

	mu.Lock()

	streams := make(map[string]network.Stream)

	for id, peer := range peers {
		streams[id] = peer.Stream
	}

	mu.Unlock()

	for id, stream := range streams {

		encoder := json.NewEncoder(stream)

		err := encoder.Encode(msg)

		if err != nil {
			fmt.Println("Send error to", id, ":", err)
			continue
		}
	}
}

func sendToPeer(peerID string, msg Message) {

	mu.Lock()

	peer, exists := peers[peerID]

	mu.Unlock()

	if !exists {
		fmt.Println("Peer not found")
		return
	}

	encoder := json.NewEncoder(peer.Stream)

	err := encoder.Encode(msg)

	if err != nil {
		fmt.Println("Send error:", err)
	}
}
