package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

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

type DiscoveryNotifee struct {
	node host.Host
}

func (n *DiscoveryNotifee) HandlePeerFound(info peer.AddrInfo) {
	if info.ID == n.node.ID() {
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
	reader := bufio.NewReader(stream)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("\nConnection closed with:", peerID)

			mu.Lock()
			delete(chatStreams, peerID)
			mu.Unlock()

			return
		}

		fmt.Printf("\nFriend (%s): %s", peerID, msg)
		fmt.Print("> ")
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

func sendToAll(msg string) {
	mu.Lock()
	defer mu.Unlock()

	for id, stream := range chatStreams {
		writer := bufio.NewWriter(stream)
		_, err := writer.WriteString(msg)
		if err != nil {
			fmt.Println("Send error to", id, ":", err)
			continue
		}
		writer.Flush()
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
		if err != nil {
			continue
		}

		if len(chatStreams) == 0 {
			fmt.Println("No peers connected.")
			continue
		}

		sendToAll(text)
	}
}
