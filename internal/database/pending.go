package database

import (
	"messenger/internal/models"
	"time"
)

type PendingMessage struct {
	Message    models.Message
	PeerID     string
	SentAt     time.Time
	RetryCount int
}
