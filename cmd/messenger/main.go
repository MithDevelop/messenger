package main

import (
	"bufio"
	"fmt"
	"messenger/internal/chat"
	"messenger/internal/database"
	"messenger/internal/models"
	"messenger/internal/network"
	"messenger/internal/peer"
	"messenger/internal/state"
	"messenger/internal/utils"
	"os"
	"strings"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	mdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

func main() {

	priv, err := peer.LoadOrCreateIdentity()

	if err != nil {
		panic(err)
	}

	err = database.InitDatabase()

	if err != nil {
		panic(err)
	}

	node, err := libp2p.New(
		libp2p.Identity(priv),
	)

	if err != nil {
		panic(err)
	}

	node.SetStreamHandler(state.ProtocolID, network.HandleStream)
	state.LocalPeerID = node.ID().String()

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
		&network.DiscoveryNotifee{Node: node},
	)

	if err := service.Start(); err != nil {
		panic(err)
	}

	fmt.Println("\nmDNS discovery started!")
	fmt.Println("Waiting for peers...")
	fmt.Println("===================================")
	network.StartPresenceLoop()
	network.MonitorPeers()
	network.StartRetryLoop()

	stdReader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print("> ")

		text, err := stdReader.ReadString('\n')

		if err != nil {
			continue
		}

		text = strings.TrimSpace(text)

		if text == "" {
			continue
		}

		if chat.HandleCommand(text) {
			continue
		}

		state.Mu.Lock()
		peerCount := len(state.Peers)
		state.Mu.Unlock()

		if peerCount == 0 {
			fmt.Println("No peers connected.")
			continue
		}

		msg := models.Message{
			ID:        utils.GenerateMessageID(),
			Type:      "chat",
			From:      node.ID().String(),
			Username:  state.Username,
			Message:   text,
			Timestamp: time.Now().Unix(),
		}

		err = database.SaveMessage(msg)

		if err != nil {
			fmt.Println("DB save error:", err)
		}

		network.SendToAll(msg, true)
	}
}
