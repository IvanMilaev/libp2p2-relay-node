package main

import (
    "fmt"
    "log"
    "time"

    libp2p "github.com/libp2p/go-libp2p"
    crypto "github.com/libp2p/go-libp2p/core/crypto"
    relayv2 "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

func main() {
    // Generate a new key pair (optional but ensures consistent Peer ID)
    privKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
    if err != nil {
        log.Fatalf("Failed to generate key pair: %v", err)
    }

    // Create a new libp2p host with the generated key
    host, err := libp2p.New(
        libp2p.Identity(privKey),
        libp2p.ListenAddrStrings(
            "/ip4/0.0.0.0/tcp/4001",
            "/ip6/::/tcp/4001",
        ),
    )
    if err != nil {
        log.Fatalf("Failed to create libp2p host: %v", err)
    }

    // Start a relay service on the host with custom resource limits
    _, err = relayv2.New(host,
        relayv2.WithResources(relayv2.Resources{
            Limit: &relayv2.Limit{
                Reservations:   128,        // Max number of reservations
                Circuits:       1024,       // Max number of relay circuits
                Duration:       time.Hour,  // Reservation duration
                Data:           1 << 30,    // Data limit per reservation (1 GB)
                AllowAllClients: true,      // Allow all clients to reserve
            },
        }),
    )
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

