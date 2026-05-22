package network

import (
	"encoding/json"
	"fmt"
	"messenger/internal/database"
	"messenger/internal/models"
	"messenger/internal/peer"
	"messenger/internal/state"
	"messenger/internal/utils"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
)

func ReadMessages(peerID string, stream network.Stream) {

	decoder := json.NewDecoder(stream)

	for {

		var msg models.Message

		err := decoder.Decode(&msg)

		if err != nil {

			fmt.Println("\nConnection closed with:", peerID)

			state.Mu.Lock()
			delete(state.Peers, peerID)
			state.Mu.Unlock()

			return
		}

		state.Mu.Lock()

		peerInfo, exists := state.Peers[peerID]

		if exists {
			peerInfo.Username = msg.Username
			peerInfo.LastSeen = time.Now()
		}

		state.Mu.Unlock()

		HandleMessage(peerID, msg)
	}
}

func SendUserInfo(peerID string) {

	msg := models.Message{
		ID:        utils.GenerateMessageID(),
		Type:      "user_info",
		From:      state.LocalPeerID,
		To:        peerID,
		Username:  state.Username,
		Timestamp: time.Now().Unix(),
	}

	SendToPeer(peerID, msg, true)
}

func HandleStream(stream network.Stream) {

	peerID := stream.Conn().RemotePeer().String()

	fmt.Println("\nIncoming chat connection from:", peerID)

	state.Mu.Lock()

	if _, exists := state.Peers[peerID]; !exists {

		state.Peers[peerID] = &peer.PeerInfo{
			ID:        peerID,
			Stream:    stream,
			Connected: time.Now(),
			LastSeen:  time.Now(),
			Latency:   0,
			Online:    true,
		}
	}

	state.Mu.Unlock()

	go ReadMessages(peerID, stream)

	SendUserInfo(peerID)
}

func SendToAll(msg models.Message, track bool) {

	state.Mu.Lock()

	streams := make(map[string]network.Stream)

	for id, peer := range state.Peers {
		streams[id] = peer.Stream
	}

	state.Mu.Unlock()

	for id, stream := range streams {

		encoder := json.NewEncoder(stream)

		err := encoder.Encode(msg)

		if err != nil {
			fmt.Println("Send error to", id, ":", err)
			continue
		}
	}
}

func SendToPeer(peerID string, msg models.Message, track bool) {

	state.Mu.Lock()

	peer, exists := state.Peers[peerID]

	state.Mu.Unlock()

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

		state.Mu.Lock()

		state.PendingMessages[msg.ID] = database.PendingMessage{
			Message:    msg,
			PeerID:     peerID,
			SentAt:     time.Now(),
			RetryCount: 0,
		}

		state.Mu.Unlock()
	}
}
