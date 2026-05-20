package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	libp2p "github.com/libp2p/go-libp2p"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	mdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

const ProtocolID = "/messenger/1.0.0"

// peerID -> stream
var chatStreams = make(map[string]network.Stream)
var mu sync.Mutex

// global name
var username = "anonymous"

type DiscoveryNotifee struct {
	node host.Host
}

func (n *DiscoveryNotifee) HandlePeerFound(info peer.AddrInfo) {
	if info.ID == n.node.ID() {
		return
	}

	if n.node.ID().String() < info.ID.String() {
		return
	}
	fmt.Println("\nFound peer:", info.ID)

	err := n.node.Connect(context.Background(), info)
	if err != nil {
		fmt.Println("Connection failed:", err)
		return
	}

	mu.Lock()
	if _, exists := chatStreams[info.ID.String()]; exists {
		mu.Unlock()
		return
	}
	mu.Unlock()

	stream, err := n.node.NewStream(context.Background(), info.ID, ProtocolID)
	if err != nil {
		fmt.Println("Stream error:", err)
		return
	}

	mu.Lock()
	chatStreams[info.ID.String()] = stream
	mu.Unlock()

	fmt.Println("Chat stream created with:", info.ID)

	go readMessages(info.ID.String(), stream)
}

func readMessages(peerID string, stream network.Stream) {

	decoder := json.NewDecoder(stream)

	for {

		var msg Message

		err := decoder.Decode(&msg)

		if err != nil {
			fmt.Println("\nConnection closed with:", peerID)

			mu.Lock()
			delete(chatStreams, peerID)
			mu.Unlock()

			return
		}

		handleMessage(msg)
	}
}

func handleStream(stream network.Stream) {
	peerID := stream.Conn().RemotePeer().String()

	fmt.Println("\nIncoming chat connection from:", peerID)

	mu.Lock()
	if _, exists := chatStreams[peerID]; !exists {
		chatStreams[peerID] = stream
	}
	mu.Unlock()

	go readMessages(peerID, stream)
}

func sendToAll(msg Message) {

	mu.Lock()

	streams := make(map[string]network.Stream)

	for id, stream := range chatStreams {
		streams[id] = stream
	}

	mu.Unlock()

	for id, stream := range streams {

		encoder := json.NewEncoder(stream)

		err := encoder.Encode(msg)

		if err != nil {
			fmt.Println("Send error to", id, ":", err)
			continue
		}
	}
}
func handleChat(msg Message) {
	fmt.Printf(
		"\n%s: %s\n",
		msg.Username,
		msg.Message,
	)
}

func handleMessage(msg Message) {

	switch msg.Type {

	case "chat":
		handleChat(msg)

	case "system":
		//handleSystem(msg)

	case "ping":
		//handlePing(msg)

	default:
		fmt.Println("Unknown message type:", msg.Type)
	}
}

func main() {
	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	node.SetStreamHandler(ProtocolID, handleStream)

	fmt.Println("===================================")
	fmt.Println(" Messenger started!")
	fmt.Println("===================================")
	fmt.Println("Peer ID:", node.ID())

	for _, addr := range node.Addrs() {
		fmt.Printf("- %s/p2p/%s\n", addr, node.ID())
	}

	service := mdns.NewMdnsService(
		node,
		"messenger-mdns",
		&DiscoveryNotifee{node: node},
	)

	if err := service.Start(); err != nil {
		panic(err)
	}

	fmt.Println("\nmDNS discovery started!")
	fmt.Println("Waiting for peers...")
	fmt.Println("===================================")

	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		text, err := stdReader.ReadString('\n')
		if handleCommand(text) {
			continue
		}
		if err != nil {
			continue
		}

		if len(chatStreams) == 0 {
			fmt.Println("No peers connected.")
			continue
		}

		msg := Message{
			Type:      "chat",
			From:      node.ID().String(),
			Username:  username,
			Message:   text,
			Timestamp: time.Now().Unix(),
		}

		sendToAll(msg)
	}
}
