package state

import (
	"messenger/internal/database"
	"messenger/internal/peer"
	"sync"
)

var Peers = make(map[string]*peer.PeerInfo)

var LocalPeerID string

var ProcessedMessages = make(map[string]bool)

var PendingMessages = make(map[string]database.PendingMessage)

var Contacts = make(map[string]*database.Contact)

var Mu sync.Mutex

var Username = "anonymous"

const ProtocolID = "/messenger/1.0.0"
