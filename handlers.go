package main

import (
	"fmt"
	"time"
)

func handleChat(msg Message) {

	name := msg.Username

	if name == "" {
		name = msg.From
	}

	fmt.Printf(
		"\n%s: %s\n",
		name,
		msg.Message,
	)

	ack := Message{
		ID:        generateMessageID(),
		Type:      "ack",
		From:      localPeerID,
		To:        msg.From,
		ReplyTo:   msg.ID,
		Timestamp: time.Now().Unix(),
	}

	sendToPeer(msg.From, ack, false)

	err := saveMessage(msg)

	if err != nil {
		fmt.Println("DB save error:", err)
	}

	fmt.Print("> ")
}

func handlePrivate(msg Message) {

	err := saveMessage(msg)

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

func handleUserInfo(peerID string, msg Message) {

	isNew := false

	mu.Lock()

	contact, exists := contacts[peerID]

	if !exists {

		contacts[peerID] = &Contact{
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

	mu.Unlock()

	if isNew {
		fmt.Println("\nNew contact added:", msg.Username)
		fmt.Print("> ")
	}
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
		ID:        generateMessageID(),
		Type:      "pong",
		From:      "",
		To:        msg.From,
		Username:  username,
		Message:   "",
		Timestamp: time.Now().UnixMilli(),
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

	_, exists := pendingMessages[msg.ReplyTo]

	if exists {
		delete(pendingMessages, msg.ReplyTo)
	}

	mu.Unlock()

	if exists {
		fmt.Println("\nMessage delivered:", msg.ReplyTo)
		fmt.Print("> ")
	}
}
