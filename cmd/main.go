package main

import (
	"namereg-chain/core"
	"namereg-chain/network"
)

func main() {
	// Initialize the blockchain and genesis state
	chain := core.NewBlockchain()

	// Start the HTTP API server
	api := network.NewAPIServer(chain)
	api.Start("localhost:8080")
}
