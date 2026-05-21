package main

import (
	"fmt"
	"time"
)

func handleChat(msg Message) {

	fmt.Printf(
		"\n%s: %s\n",
		msg.Username,
		msg.Message,
	)

	fmt.Print("> ")
	ack := Message{
		ID:        generateMessageID(),
		Type:      "ack",
		From:      localPeerID,
		To:        msg.From,
		ReplyTo:   msg.ID,
		Timestamp: time.Now().Unix(),
	}

	sendToPeer(msg.From, ack, false)
}

func handlePrivate(msg Message) {

	fmt.Printf(
		"\n[PRIVATE] %s: %s\n",
		msg.Username,
		msg.Message,
	)

	fmt.Print("> ")
}

func handleUserInfo(peerID string, msg Message) {

	mu.Lock()
	defer mu.Unlock()

	contact, exists := contacts[peerID]

	if !exists {

		contacts[peerID] = &Contact{
			PeerID:   peerID,
			Username: msg.Username,
			AddedAt:  time.Now(),
			LastSeen: time.Now(),
			Trusted:  false,
		}

		fmt.Println("\nNew contact added:", msg.Username)

		return
	}

	contact.Username = msg.Username
	contact.LastSeen = time.Now()
}

func handleMessage(peerID string, msg Message) {

	//idempotent processing
	if msg.ID != "" {

		mu.Lock()

		if processedMessages[msg.ID] {

			mu.Unlock()
			return
		}

		processedMessages[msg.ID] = true

		mu.Unlock()
	}

	switch msg.Type {

	case "chat":
		handleChat(msg)

	case "private":
		handlePrivate(msg)

	case "user_info":
		handleUserInfo(peerID, msg)

	case "ping":
		handlePing(msg)

	case "pong":
		handlePong(peerID, msg)

	case "ack":
		handleAck(msg)

	default:
		fmt.Println("Unknown message type:", msg.Type)
	}
}

func handlePing(msg Message) {

	reply := Message{
		Type:      "pong",
		From:      "",
		To:        msg.From,
		Username:  username,
		Message:   "",
		Timestamp: msg.Timestamp,
	}

	sendToPeer(msg.From, reply, false)
}

func handlePong(peerID string, msg Message) {

	latency := time.Now().UnixMilli() - msg.Timestamp

	mu.Lock()

	peer, exists := peers[peerID]

	if exists {

		peer.LastSeen = time.Now()
		peer.Latency = latency
		peer.Online = true
	}

	mu.Unlock()
}

func handleAck(msg Message) {

	fmt.Println(
		"\nMessage delivered:",
		msg.ReplyTo,
	)
	mu.Lock()
	delete(pendingMessages, msg.ReplyTo)
	mu.Unlock()
}
