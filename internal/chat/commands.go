package chat

import (
	"fmt"
	"messenger/internal/models"
	"messenger/internal/network"
	"messenger/internal/state"
	"messenger/internal/utils"
	"strings"
	"time"
)

func HandleCommand(input string) bool {

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
		fmt.Println("/msg <peerID> - личное сообщение")
		fmt.Println("/contacts - контакты")
		fmt.Println("")

	case "/nick":

		if len(parts) < 2 {
			fmt.Println("Usage: /nick <name>")
			return true
		}

		state.Username = parts[1]

		fmt.Println("Username changed to:", state.Username)

	case "/msg":
		if len(parts) < 3 {
			fmt.Println("Usage: /msg <peerID> <message>")
			return true
		}

		peerID := parts[1]

		text := strings.Join(parts[2:], " ")

		msg := models.Message{
			ID:        utils.GenerateMessageID(),
			Type:      "private",
			From:      "",
			To:        peerID,
			Username:  state.Username,
			Message:   text,
			Timestamp: time.Now().Unix(),
		}

		network.SendToPeer(peerID, msg, true)

		fmt.Println("Private message sent")

	case "/peers":

		state.Mu.Lock()

		fmt.Println("Connected peers:")

		for _, peer := range state.Peers {

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
		state.Mu.Unlock()

	case "/contacts":

		state.Mu.Lock()

		fmt.Println("Contacts:")

		for _, contact := range state.Contacts {

			fmt.Printf(
				"- %s (%s)\n",
				contact.Username,
				contact.PeerID,
			)
		}

		state.Mu.Unlock()

	default:
		fmt.Println("Unknown command")
	}

	return true
}
