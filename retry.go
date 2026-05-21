package main

import (
	"time"
)

func checkPendingMessages() {

	mu.Lock()

	pending := make([]PendingMessage, 0)

	for _, p := range pendingMessages {
		pending = append(pending, p)
	}

	mu.Unlock()

	for _, p := range pending {

		if time.Since(p.SentAt) < 5*time.Second {
			continue
		}

		if p.RetryCount >= 3 {

			mu.Lock()
			delete(pendingMessages, p.Message.ID)
			mu.Unlock()

			continue
		}

		p.RetryCount++
		p.SentAt = time.Now()

		mu.Lock()
		pendingMessages[p.Message.ID] = p
		mu.Unlock()

		sendToPeer(p.PeerID, p.Message, false)
	}
}

func startRetryLoop() {

	for {

		time.Sleep(3 * time.Second)

		checkPendingMessages()
	}
}
