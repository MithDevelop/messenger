package models

type Message struct {
	ID        string `json:"id"`
	ReplyTo   string `json:"reply_to"`
	Type      string `json:"type"`
	From      string `json:"from"`
	To        string `json:"to"`
	Username  string `json:"username"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
