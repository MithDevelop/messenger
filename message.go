package main

type Message struct {
	Type      string `json:"type"`
	From      string `json:"from"`
	Username  string `json:"username"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
