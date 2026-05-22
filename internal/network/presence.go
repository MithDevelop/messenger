package network

import (
	"fmt"
	"messenger/internal/models"
	"messenger/internal/state"
	"messenger/internal/utils"
	"time"
)

func StartPresenceLoop() {

	for {

		time.Sleep(5 * time.Second)

		SendPingToAll()
	}
}

func SendPingToAll() {

	state.Mu.Lock()

	peerIDs := make([]string, 0)

	for id := range state.Peers {
		peerIDs = append(peerIDs, id)
	}

	state.Mu.Unlock()

	for _, peerID := range peerIDs {

		msg := models.Message{
			ID:        utils.GenerateMessageID(),
			Type:      "ping",
			From:      state.LocalPeerID,
			To:        peerID,
			Username:  state.Username,
			Message:   "",
			Timestamp: time.Now().UnixMilli(),
		}

		SendToPeer(peerID, msg, false)
	}
}

func MonitorPeers() {

	for {

		time.Sleep(10 * time.Second)

		state.Mu.Lock()

		for id, peer := range state.Peers {

			if time.Since(peer.LastSeen) > 15*time.Second {

				fmt.Println("\nPeer timed out:", id)

				peer.Online = false

				delete(state.Peers, id)
			}
		}

		state.Mu.Unlock()
	}
}
