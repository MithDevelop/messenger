package network

import (
	"messenger/internal/database"
	"messenger/internal/state"
	"time"
)

func checkPendingMessages() {

	state.Mu.Lock()

	pending := make([]database.PendingMessage, 0)

	for _, p := range state.PendingMessages {
		pending = append(pending, p)
	}

	state.Mu.Unlock()

	for _, p := range pending {

		if time.Since(p.SentAt) < 5*time.Second {
			continue
		}

		if p.RetryCount >= 3 {

			state.Mu.Lock()
			delete(state.PendingMessages, p.Message.ID)
			state.Mu.Unlock()

			continue
		}

		p.RetryCount++
		p.SentAt = time.Now()

		state.Mu.Lock()
		state.PendingMessages[p.Message.ID] = p
		state.Mu.Unlock()

		SendToPeer(p.PeerID, p.Message, false)
	}
}

func StartRetryLoop() {

	for {

		time.Sleep(3 * time.Second)

		checkPendingMessages()
	}
}
