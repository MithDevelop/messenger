package main

import (
	"fmt"
	"strings"
	"time"
)

func handleCommand(input string) bool {

	if len(input) == 0 {
		return true
	}

	if input[0] != '/' {
		return false
	}

	parts := strings.Fields(input)

	if len(parts) == 0 {
		return true
	}

	switch parts[0] {

	case "/help":
		fmt.Println("Commands:")
		fmt.Println("/nick <name> - введи своё имя")
		fmt.Println("/peers - участники")

	case "/nick":

		if len(parts) < 2 {
			fmt.Println("Usage: /nick <name>")
			return true
		}

		username = parts[1]

		fmt.Println("Username changed to:", username)

	case "/msg":
		if len(parts) < 3 {
			fmt.Println("Usage: /msg <peerID> <message>")
			return true
		}

		peerID := parts[1]

		text := strings.Join(parts[2:], " ")

		msg := Message{
			ID:        generateMessageID(),
			Type:      "private",
			From:      "",
			To:        peerID,
			Username:  username,
			Message:   text,
			Timestamp: time.Now().Unix(),
		}

		sendToPeer(peerID, msg, true)

		fmt.Println("Private message sent")
	case "/peers":

		mu.Lock()

		fmt.Println("Connected peers:")

		for _, peer := range peers {

			name := peer.Username

			if name == "" {
				name = peer.ID
			}

			status := "offline"

			if peer.Online {
				status = "online"
			}

			fmt.Printf(
				"- %s | %s | %dms\n",
				name,
				status,
				peer.Latency,
			)
		}

		mu.Unlock()

	default:
		fmt.Println("Unknown command")
	}

	return true
}
