package main

import (
    "fmt"
    "log"

    libp2p "github.com/libp2p/go-libp2p"
    relayv2 "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

func main() {
    // Create a new libp2p host
    host, err := libp2p.New(
        libp2p.ListenAddrStrings(
            "/ip4/0.0.0.0/tcp/4001",
            "/ip6/::/tcp/4001",
        ),
    )
    if err != nil {
        log.Fatalf("Failed to create libp2p host: %v", err)
    }

    // Start a relay service on the host
    _, err = relayv2.New(host, relayv2.WithLimit(nil))
    if err != nil {
        log.Fatalf("Failed to enable relay service: %v", err)
    }

    // Print the host's multiaddresses
    fmt.Println("Relay node is running. Peer ID:", host.ID())
    for _, addr := range host.Addrs() {
        fmt.Printf("Listening on: %s/p2p/%s\n", addr, host.ID())
    }

    select {} // Keep the process running
}

