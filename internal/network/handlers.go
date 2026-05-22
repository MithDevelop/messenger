package network

import (
	"fmt"
	"messenger/internal/database"
	"messenger/internal/models"
	"messenger/internal/state"
	"messenger/internal/utils"
	"time"
)

func HandleChat(msg models.Message) {

	name := msg.Username

	if name == "" {
		name = msg.From
	}

	fmt.Printf(
		"\n%s: %s\n",
		name,
		msg.Message,
	)

	ack := models.Message{
		ID:        utils.GenerateMessageID(),
		Type:      "ack",
		From:      state.LocalPeerID,
		To:        msg.From,
		ReplyTo:   msg.ID,
		Timestamp: time.Now().Unix(),
	}

	SendToPeer(msg.From, ack, false)

	err := database.SaveMessage(msg)

	if err != nil {
		fmt.Println("DB save error:", err)
	}

	fmt.Print("> ")
}

func HandlePrivate(msg models.Message) {

	err := database.SaveMessage(msg)

	if err != nil {
		fmt.Println("DB save error:", err)
	}

	fmt.Printf(
		"\n[PRIVATE] %s: %s\n",
		msg.Username,
		msg.Message,
	)

	fmt.Print("> ")
}

func HandleUserInfo(peerID string, msg models.Message) {

	isNew := false

	state.Mu.Lock()

	contact, exists := state.Contacts[peerID]

	if !exists {

		state.Contacts[peerID] = &database.Contact{
			PeerID:   peerID,
			Username: msg.Username,
			AddedAt:  time.Now(),
			LastSeen: time.Now(),
			Trusted:  false,
		}

		isNew = true

	} else {

		contact.Username = msg.Username
		contact.LastSeen = time.Now()
	}

	state.Mu.Unlock()

	if isNew {
		fmt.Println("\nNew contact added:", msg.Username)
		fmt.Print("> ")
	}
}

func HandleMessage(peerID string, msg models.Message) {

	//idempotent processing
	if msg.ID != "" {

		state.Mu.Lock()

		if state.ProcessedMessages[msg.ID] {

			state.Mu.Unlock()
			return
		}

		state.ProcessedMessages[msg.ID] = true

		state.Mu.Unlock()
	}

	switch msg.Type {

	case "chat":
		HandleChat(msg)

	case "private":
		HandlePrivate(msg)

	case "user_info":
		HandleUserInfo(peerID, msg)

	case "ping":
		HandlePing(msg)

	case "pong":
		HandlePong(peerID, msg)

	case "ack":
		HandleAck(msg)

	default:
		fmt.Println("Unknown message type:", msg.Type)
	}
}

func HandlePing(msg models.Message) {

	reply := models.Message{
		ID:        utils.GenerateMessageID(),
		Type:      "pong",
		From:      "",
		To:        msg.From,
		Username:  state.Username,
		Message:   "",
		Timestamp: time.Now().UnixMilli(),
	}

	SendToPeer(msg.From, reply, false)
}

func HandlePong(peerID string, msg models.Message) {

	latency := time.Now().UnixMilli() - msg.Timestamp

	state.Mu.Lock()

	peer, exists := state.Peers[peerID]

	if exists {

		peer.LastSeen = time.Now()
		peer.Latency = latency
		peer.Online = true
	}

	state.Mu.Unlock()
}

func HandleAck(msg models.Message) {

	fmt.Println(
		"\nMessage delivered:",
		msg.ReplyTo,
	)
	state.Mu.Lock()

	_, exists := state.PendingMessages[msg.ReplyTo]

	if exists {
		delete(state.PendingMessages, msg.ReplyTo)
	}

	state.Mu.Unlock()

	if exists {
		fmt.Println("\nMessage delivered:", msg.ReplyTo)
		fmt.Print("> ")
	}
}
