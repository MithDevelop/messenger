package database

import "time"

type Contact struct {
	PeerID   string    `json:"peer_id"`
	Username string    `json:"username"`
	AddedAt  time.Time `json:"added_at"`
	LastSeen time.Time `json:"last_seen"`
	Trusted  bool      `json:"trusted"`
}
