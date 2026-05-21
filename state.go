package main

import "sync"

var peers = make(map[string]*PeerInfo)

var localPeerID string

var processedMessages = make(map[string]bool)

var pendingMessages = make(map[string]PendingMessage)

var mu sync.Mutex

var username = "anonymous"

const ProtocolID = "/messenger/1.0.0"
