package main

import "time"

type PendingMessage struct {
	Message    Message
	PeerID     string
	SentAt     time.Time
	RetryCount int
}
