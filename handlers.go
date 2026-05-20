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

	peerInfo, exists := peers[peerID]

	if exists {
		peerInfo.Username = msg.Username
		peerInfo.LastSeen = time.Now()
	}

	mu.Unlock()

	fmt.Printf(
		"\nPeer %s is known as %s\n",
		peerID,
		msg.Username,
	)

	fmt.Print("> ")
}

func handleMessage(peerID string, msg Message) {

	switch msg.Type {

	case "chat":
		handleChat(msg)

	case "private":
		handlePrivate(msg)

	case "user_info":
		handleUserInfo(peerID, msg)

	default:
		fmt.Println("Unknown message type:", msg.Type)
	}
}
