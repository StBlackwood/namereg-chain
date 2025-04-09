// namereg-chain/main.go
package main

import (
	"fmt"
	"net/http"
	"time"

	"namereg/core"
	"namereg/network"
	"namereg/rpc"
	"namereg/state"
)

func main() {
	// Initialize state
	s := state.NewState()

	// Initialize blockchain
	bc := core.NewBlockchain(s)

	// Start networking
	n := network.NewNode(bc)
	go n.Start()

	// Start RPC server
	r := rpc.NewServer(bc)
	fmt.Println("RPC server running on :8080")
	http.ListenAndServe(":8080", r.Router())
}

// The rest of the implementation is in separate folders:
// core/ - blockchain, blocks, transaction
// state/ - global state and logic
// rpc/ - HTTP handlers
// network/ - peer discovery and gossip
