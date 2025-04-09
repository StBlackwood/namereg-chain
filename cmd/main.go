package main

import (
	"log"
	"namereg-chain/config"
	"namereg-chain/core"
	"namereg-chain/network"
)

func main() {
	// Load config file
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Print info
	log.Printf("Starting node %s on port %s\n", cfg.NodeID, cfg.Port)
	log.Printf("Known peers: %v\n", cfg.Peers)

	// Create blockchain and peer client
	chain := core.NewBlockchain()
	peerClient := network.NewPeerClient(cfg.Peers)

	// Start the API server (with peer awareness)
	api := network.NewAPIServerWithPeers(chain, peerClient)
	api.Start(":" + cfg.Port)
}
