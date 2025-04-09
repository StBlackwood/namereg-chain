package network

import (
	"encoding/json"
	"log"
	"net/http"

	"namereg-chain/core"
)

type APIServer struct {
	Chain       *core.Blockchain
	PeerClient  *PeerClient
	EnablePeers bool
}

func NewAPIServer(chain *core.Blockchain) *APIServer {
	return &APIServer{Chain: chain}
}

func NewAPIServerWithPeers(chain *core.Blockchain, peers *PeerClient) *APIServer {
	return &APIServer{
		Chain:       chain,
		PeerClient:  peers,
		EnablePeers: true,
	}
}

func (api *APIServer) Start(addr string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", api.handleRegister)
	mux.HandleFunc("/lookup", api.handleLookup)
	mux.HandleFunc("/chain", api.handleChain)
	mux.HandleFunc("/nonce", api.handleNonce)
	mux.HandleFunc("/receive-block", api.handleReceiveBlock)

	log.Printf("API server listening on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux)) // use custom mux
}

// POST /register
func (api *APIServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var tx core.Transaction
	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	// Validate tx
	err = api.Chain.State.ValidateTransaction(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add to chain
	block := api.Chain.AddBlock([]core.Transaction{tx})

	// Broadcast to peers
	if api.EnablePeers {
		api.PeerClient.BroadcastBlock(block)
	}

	resp := map[string]interface{}{
		"message": "Transaction accepted and added to block",
		"block":   block,
	}
	writeJSON(w, resp)
}

// GET /lookup?name=alice
func (api *APIServer) handleLookup(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	address, ok := api.Chain.State.GetAddressByName(name)
	if !ok {
		http.Error(w, "name not registered", http.StatusNotFound)
		return
	}
	writeJSON(w, map[string]string{
		"name":    name,
		"address": address,
	})
}

// GET /chain
func (api *APIServer) handleChain(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, api.Chain.Blocks)
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (api *APIServer) handleNonce(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "address required", http.StatusBadRequest)
		return
	}
	nonce := api.Chain.State.GetNonce(address)
	writeJSON(w, map[string]uint64{"nonce": nonce})
}

// POST /receive-block
func (api *APIServer) handleReceiveBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var incoming core.Block
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		http.Error(w, "invalid block data", http.StatusBadRequest)
		return
	}

	// Check if already have this block height
	if incoming.Height <= api.Chain.LatestBlock().Height {
		http.Error(w, "block already known or stale", http.StatusConflict)
		return
	}

	// Check previous hash
	expectedPrev := api.Chain.LatestBlock().Hash
	if incoming.PrevHash != expectedPrev {
		http.Error(w, "invalid previous hash", http.StatusBadRequest)
		return
	}

	// Validate and apply all transactions to a temp copy of state
	tempState := api.Chain.State.Copy()
	for _, tx := range incoming.Transactions {
		if err := tempState.ApplyTransaction(tx); err != nil {
			http.Error(w, "invalid transaction in block: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// All good, apply to real chain
	api.Chain.State = tempState
	api.Chain.Blocks = append(api.Chain.Blocks, &incoming)

	log.Printf("Accepted block %d from peer\n", incoming.Height)
	w.WriteHeader(http.StatusOK)
}
