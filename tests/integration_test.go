package tests

import (
	"bytes"
	"encoding/json"
	"namereg-chain/core"
	"net/http"
	"testing"
	"time"
)

func TestInMemoryNodeInteraction(t *testing.T) {
	// initial setup starts 2 server nodes at port 8081 and 8082
	start2Nodes()

	time.Sleep(1 * time.Second) // give servers time to start

	// Create a keypair and signed transaction
	priv, pubKey, address := generateKeyPair()

	tx := core.Transaction{
		Name:    "alice",
		Address: address,
		Nonce:   0,
		PubKey:  pubKey,
	}
	tx.Signature = signTransaction(priv, &tx)

	// Send the transaction to node1
	body, _ := json.Marshal(tx)
	resp, err := http.Post("http://localhost:8081/register", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to send transaction: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Transaction rejected: %s", resp.Status)
	}

	time.Sleep(1 * time.Second) // allow block to sync to node2

	// Verify name exists on node2
	resp, err = http.Get("http://localhost:8082/lookup?name=alice")
	if err != nil {
		t.Fatalf("Failed to fetch from node2: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected name on node2 but got %s", resp.Status)
	}

	t.Log("Transaction registered on node1 and reflected on node2")
}
