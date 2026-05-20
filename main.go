package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	libp2p "github.com/libp2p/go-libp2p"

	mdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

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

		if handleCommand(text) {
			continue
		}

		if len(peers) == 0 {
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
