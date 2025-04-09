package network

import (
	"encoding/json"
	"log"
	"net/http"

	"namereg-chain/core"
)

type APIServer struct {
	Chain *core.Blockchain
}

func NewAPIServer(chain *core.Blockchain) *APIServer {
	return &APIServer{
		Chain: chain,
	}
}

func (api *APIServer) Start(addr string) {
	http.HandleFunc("/register", api.handleRegister)
	http.HandleFunc("/lookup", api.handleLookup)
	http.HandleFunc("/chain", api.handleChain)
	http.HandleFunc("/nonce", api.handleNonce)

	log.Printf("API server listening on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// POST /register
// Body: { "from": "user1", "name": "alice", "address": "0xabc", "nonce": 0 }
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

	// Validate tx without applying
	if err := api.Chain.State.ValidateTransaction(tx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Apply it via block addition only
	block := api.Chain.AddBlock([]core.Transaction{tx})

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

	resp := map[string]string{
		"name":    name,
		"address": address,
	}
	writeJSON(w, resp)
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
