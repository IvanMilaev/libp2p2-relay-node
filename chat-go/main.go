package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	network "github.com/libp2p/go-libp2p/core/network"
	peer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func main() {
	// Parse command-line arguments
	var (
		help        bool
		port        int
		destination string
		relayIP     string
		relayPeerID string
	)

	flag.BoolVar(&help, "h", false, "Display Help")
	flag.IntVar(&port, "p", 0, "Local TCP port to listen on")
	flag.StringVar(&destination, "d", "", "Destination multiaddress string")
	flag.StringVar(&relayIP, "relay-ip", "", "IP address of the relay node")
	flag.StringVar(&relayPeerID, "relay-peerid", "", "Peer ID of the relay node")
	flag.Parse()

	if help || port == 0 || relayIP == "" || relayPeerID == "" {
		fmt.Printf("Usage: %s -p <PORT> [OPTIONS]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	ctx := context.Background()

	// Generate a random peer ID
	privKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		log.Fatalln(err)
	}

	// Construct the relay node's multiaddress
	relayAddrStr := fmt.Sprintf("/ip4/%s/tcp/4001/p2p/%s", relayIP, relayPeerID)
	relayAddr, err := multiaddr.NewMultiaddr(relayAddrStr)
	if err != nil {
		log.Fatalf("Invalid relay multiaddress: %v", err)
	}

	relayInfo, err := peer.AddrInfoFromP2pAddr(relayAddr)
	if err != nil {
		log.Fatalf("Failed to get relay AddrInfo: %v", err)
	}

	if destination == "" {
		// Acting as listener

		// Create a new libp2p host with AutoRelay and Static Relays
		host, err := libp2p.New(
			libp2p.Identity(privKey),
			libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)),
			libp2p.EnableAutoRelay(
				autorelay.WithStaticRelays([]peer.AddrInfo{*relayInfo}),
			),
		)
		if err != nil {
			log.Fatalln(err)
		}

		// Wait for the host to obtain a relay address
		fmt.Println("Waiting for relay address...")
		var relayAddrFound bool
		for i := 0; i < 10; i++ {
			addrs := host.Addrs()
			for _, addr := range addrs {
				if strings.Contains(addr.String(), "p2p-circuit") {
					relayAddrFound = true
					break
				}
			}
			if relayAddrFound {
				break
			}
			// Wait a moment before checking again
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
		}
		if !relayAddrFound {
			log.Println("Failed to obtain relay address. Exiting.")
			return
		}

		// Display the listener's relay address
		fullAddrs := host.Addrs()
		fmt.Printf("Listening for connections. To connect, run:\n")
		for _, addr := range fullAddrs {
			if strings.Contains(addr.String(), "p2p-circuit") {
				fmt.Printf("%s -p <PORT> -d '%s/p2p/%s' --relay-ip %s --relay-peerid %s\n",
					os.Args[0],
					addr.String(),
					host.ID().String(),
					relayIP,
					relayPeerID,
				)
				break
			}
		}

		host.SetStreamHandler("/chat/1.0.0", handleStream)

		select {} // Hang forever
	} else {
		// Acting as dialer

		// Create a new libp2p host
		host, err := libp2p.New(
			libp2p.Identity(privKey),
			libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)),
		)
		if err != nil {
			log.Fatalln(err)
		}

		// Connect to the relay node
		if err := host.Connect(ctx, *relayInfo); err != nil {
			log.Fatalf("Failed to connect to relay node: %v", err)
		}
		fmt.Println("Connected to relay node:", relayAddrStr)

		destAddr, err := multiaddr.NewMultiaddr(destination)
		if err != nil {
			log.Fatalf("Invalid destination multiaddress: %v", err)
		}

		destInfo, err := peer.AddrInfoFromP2pAddr(destAddr)
		if err != nil {
			log.Fatalf("Failed to get destination AddrInfo: %v", err)
		}

		// Connect to the destination
		if err := host.Connect(ctx, *destInfo); err != nil {
			log.Fatalf("Failed to connect to destination: %v", err)
		}
		fmt.Println("Connected to destination:", destination)
		fmt.Printf("Chat started. Type messages and press Enter to send.\n")
		fmt.Printf("Your Peer ID: %s\n", host.ID().String())

		// Open a stream
		stream, err := host.NewStream(ctx, destInfo.ID, "/chat/1.0.0")
		if err != nil {
			log.Fatalf("Failed to open stream: %v", err)
		}

		// Start chat
		go readData(stream)
		writeData(stream)
	}
}

func handleStream(s network.Stream) {
	fmt.Println("New chat stream opened")
	go readData(s)
	writeData(s)
}

func readData(s network.Stream) {
	reader := bufio.NewReader(s)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stream:", err)
			s.Close()
			return
		}
		fmt.Printf("\x1b[32m%s\x1b[0m", str)
	}
}

func writeData(s network.Stream) {
	reader := bufio.NewReader(os.Stdin)
	for {
		userInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			s.Close()
			return
		}
		_, err = s.Write([]byte(userInput))
		if err != nil {
			fmt.Println("Error writing to stream:", err)
			s.Close()
			return
		}
	}
}
