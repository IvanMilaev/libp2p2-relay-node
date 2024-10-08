package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	libp2p "github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func main() {
	// Parse command-line arguments
	var ip string
	var peerID string
	flag.StringVar(&ip, "ip", "", "IP address of the relay node")
	flag.StringVar(&peerID, "peerid", "", "Peer ID of the relay node")
	flag.Parse()

	if ip == "" || peerID == "" {
		log.Fatalln("Usage: go run main.go -ip <relay_ip> -peerid <relay_peer_id>")
	}

	ctx := context.Background()

	// Create a new libp2p host
	host, err := libp2p.New()
	if err != nil {
		log.Fatalln("Failed to create libp2p host:", err)
	}

	// Construct the relay node's multiaddress using the provided IP and Peer ID
	relayAddrStr := fmt.Sprintf("/ip4/%s/tcp/4001/p2p/%s", ip, peerID)
	relayAddr, err := multiaddr.NewMultiaddr(relayAddrStr)
	if err != nil {
		log.Fatalln("Invalid multiaddress:", err)
	}

	relayInfo, err := peerstore.AddrInfoFromP2pAddr(relayAddr)
	if err != nil {
		log.Fatalln("Failed to get AddrInfo from multiaddress:", err)
	}

	// Connect to the relay node
	if err := host.Connect(ctx, *relayInfo); err != nil {
		log.Fatalln("Failed to connect to relay node:", err)
	}
	fmt.Println("Connected to relay node")
}
