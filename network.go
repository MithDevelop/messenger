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
		ID:        generateMessageID(),
		Type:      "user_info",
		From:      localPeerID,
		To:        peerID,
		Username:  username,
		Timestamp: time.Now().Unix(),
	}

	sendToPeer(peerID, msg, true)
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
			Latency:   0,
			Online:    true,
		}
	}

	mu.Unlock()

	go readMessages(peerID, stream)

	sendUserInfo(peerID)
}

func sendToAll(msg Message, track bool) {

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

func sendToPeer(peerID string, msg Message, track bool) {

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
		return
	}

	if track && msg.Type != "ack" && msg.Type != "pong" {

		mu.Lock()

		pendingMessages[msg.ID] = PendingMessage{
			Message:    msg,
			PeerID:     peerID,
			SentAt:     time.Now(),
			RetryCount: 0,
		}

		mu.Unlock()
	}
}
