
package main

import (
    "context"
    "fmt"
    "log"

    libp2p "github.com/libp2p/go-libp2p"
    relayv2 "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
    ma "github.com/multiformats/go-multiaddr"
)

func main() {
    ctx := context.Background()

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
    _, err = relayv2.New(ctx, host, relayv2.WithLimit(nil))
    if err != nil {
        log.Fatalf("Failed to enable relay service: %v", err)
    }

    // Print the host's multiaddresses
    fmt.Println("Relay node is running. Peer ID:", host.ID().Pretty())
    for _, addr := range host.Addrs() {
        fmt.Printf("Listening on: %s/p2p/%s\n", addr, host.ID().Pretty())
    }

    select {} // Keep the process running
}

