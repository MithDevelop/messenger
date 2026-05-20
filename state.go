package main

import "sync"

var peers = make(map[string]*PeerInfo)

var mu sync.Mutex

var username = "anonymous"

const ProtocolID = "/messenger/1.0.0"
