package tests

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"namereg-chain/core"
	"testing"
)

func generateKeyPair() (*ecdsa.PrivateKey, []byte, string) {
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pubKey := elliptic.MarshalCompressed(priv.PublicKey.Curve, priv.PublicKey.X, priv.PublicKey.Y)
	pubKeyHash := sha256.Sum256(pubKey)
	address := hex.EncodeToString(pubKeyHash[:])
	return priv, pubKey, address
}

func signTransaction(priv *ecdsa.PrivateKey, tx *core.Transaction) []byte {
	hash := tx.Hash()
	r, s, _ := ecdsa.Sign(rand.Reader, priv, hash)
	sig := append(r.Bytes(), s.Bytes()...)
	return sig
}

func TestState_ValidNameRegistration(t *testing.T) {
	state := core.NewState()
	priv, pubKey, address := generateKeyPair()

	tx := core.Transaction{
		Name:    "alice",
		Address: address,
		Nonce:   0,
		PubKey:  pubKey,
	}
	tx.Signature = signTransaction(priv, &tx)

	if err := state.ApplyTransaction(tx); err != nil {
		t.Fatalf("Expected successful registration, got error: %v", err)
	}

	addr, ok := state.GetAddressByName("alice")
	if !ok || addr != address {
		t.Fatalf("Expected address %s for name 'alice', got %s", address, addr)
	}
}

func TestState_DuplicateName(t *testing.T) {
	state := core.NewState()
	priv, pubKey, address := generateKeyPair()

	tx1 := core.Transaction{Name: "bob", Address: address, Nonce: 0, PubKey: pubKey}
	tx1.Signature = signTransaction(priv, &tx1)

	tx2 := core.Transaction{Name: "bob", Address: address, Nonce: 1, PubKey: pubKey}
	tx2.Signature = signTransaction(priv, &tx2)

	_ = state.ApplyTransaction(tx1)
	err := state.ApplyTransaction(tx2)

	if err == nil || err.Error() != "name already registered" {
		t.Fatalf("Expected duplicate name error, got %v", err)
	}
}

func TestState_ReusedOrInvalidNonce(t *testing.T) {
	state := core.NewState()
	priv, pubKey, address := generateKeyPair()

	tx1 := core.Transaction{Name: "charlie", Address: address, Nonce: 0, PubKey: pubKey}
	tx1.Signature = signTransaction(priv, &tx1)
	_ = state.ApplyTransaction(tx1)

	// Reuse nonce
	tx2 := core.Transaction{Name: "delta", Address: address, Nonce: 0, PubKey: pubKey}
	tx2.Signature = signTransaction(priv, &tx2)
	err := state.ApplyTransaction(tx2)

	if err == nil || err.Error() != "invalid nonce" {
		t.Fatalf("Expected nonce reuse error, got %v", err)
	}

	// Invalid future nonce
	tx3 := core.Transaction{Name: "echo", Address: address, Nonce: 5, PubKey: pubKey}
	tx3.Signature = signTransaction(priv, &tx3)
	err = state.ApplyTransaction(tx3)

	if err == nil || err.Error() != "invalid nonce" {
		t.Fatalf("Expected invalid nonce error, got %v", err)
	}
}

func TestState_UnauthenticatedAddress(t *testing.T) {
	state := core.NewState()

	// Generate keypair A
	priv1, pubKey1, address1 := generateKeyPair()
	// Generate fake address B (unrelated)
	priv2, pubKey2, _ := generateKeyPair()

	tx := core.Transaction{
		Name:    "alice",
		Address: address1,
		Nonce:   0,
		PubKey:  pubKey1,
	}
	tx.Signature = signTransaction(priv1, &tx)

	if err := state.ApplyTransaction(tx); err != nil {
		t.Fatalf("Expected successful registration, got error: %v", err)
	}

	// Try to sign a tx with B, to access address from A
	tx2 := core.Transaction{
		Name:    "alice2",
		Address: address1,
		Nonce:   1,
		PubKey:  pubKey2,
	}
	tx2.Signature = signTransaction(priv2, &tx2)

	err := state.ApplyTransaction(tx2)
	if err == nil || err.Error() != "address doesn't match public key" {
		t.Fatalf("Expected public key mismatch error, got %v", err)
	}
}
