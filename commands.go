package main

import (
	"fmt"
	"strings"
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
		fmt.Println("/help")
		fmt.Println("/nick <name>")
		fmt.Println("/peers")

	case "/nick":

		if len(parts) < 2 {
			fmt.Println("Usage: /nick <name>")
			return true
		}

		username = parts[1]

		fmt.Println("Username changed to:", username)

	case "/peers":

		mu.Lock()

		fmt.Println("Connected peers:")

		for id := range chatStreams {
			fmt.Println("-", id)
		}

		mu.Unlock()

	default:
		fmt.Println("Unknown command")
	}

	return true
}
