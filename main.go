package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	libp2p "github.com/libp2p/go-libp2p"

	"github.com/libp2p/go-libp2p/core/host"
	network "github.com/libp2p/go-libp2p/core/network"
	peer "github.com/libp2p/go-libp2p/core/peer"

	mdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

const ProtocolID = "/messenger/1.0.0"

var chatStream network.Stream

type DiscoveryNotifee struct {
	node host.Host
}

func (n *DiscoveryNotifee) HandlePeerFound(info peer.AddrInfo) {

	// Не подключаться к самому себе
	if info.ID == n.node.ID() {
		return
	}

	fmt.Println("\nFounSd peer:", info.ID)

	// Просто подключаемся
	err := n.node.Connect(context.Background(), info)

	if err != nil {
		fmt.Println("Connection failed:", err)
		return
	}

	fmt.Println("Connected to:", info.ID)

	// Создаём stream ТОЛЬКО если его ещё нет
	if chatStream == nil {

		stream, err := n.node.NewStream(
			context.Background(),
			info.ID,
			ProtocolID,
		)

		if err != nil {
			fmt.Println("Stream error:", err)
			return
		}

		chatStream = stream

		fmt.Println("Chat stream created!")

		// Запускаем чтение сообщений
		go readMessages(stream)
	}
}

func readMessages(stream network.Stream) {

	reader := bufio.NewReader(stream)

	for {

		msg, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("\nConnection closed")
			chatStream = nil
			return
		}

		fmt.Printf("\nFriend: %s", msg)
		fmt.Print("> ")
	}
}

func handleStream(stream network.Stream) {

	fmt.Println("\nIncoming chat connection!")

	// Если stream уже есть — игнорируем дубликат
	//if chatStream != nil {
	//	stream.Close()
	//	return
	//}

	chatStream = stream

	go readMessages(stream)
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

	fmt.Println("\nListening addresses:")

	for _, addr := range node.Addrs() {
		fmt.Printf("- %s/p2p/%s\n", addr, node.ID())
	}

	service := mdns.NewMdnsService(
		node,
		"messenger-mdns",
		&DiscoveryNotifee{
			node: node,
		},
	)

	err = service.Start()

	if err != nil {
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

		if chatStream == nil {
			fmt.Println("No peer connected.")
			continue
		}

		writer := bufio.NewWriter(chatStream)

		_, err = writer.WriteString(text)

		if err != nil {
			fmt.Println("Send error:", err)
			continue
		}

		writer.Flush()
	}
}
