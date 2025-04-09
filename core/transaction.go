package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
)

type Transaction struct {
	Name      string `json:"name"`
	Address   string `json:"address"` // The hex of public key hash (owner)
	Nonce     uint64 `json:"nonce"`
	Signature []byte `json:"signature"` // Signed hash of tx content
	PubKey    []byte `json:"pubKey"`    // Full public key for signature verification
}

func (tx *Transaction) Hash() []byte {
	data := []byte(tx.Name + tx.Address + string(tx.Nonce))
	hash := sha256.Sum256(data)
	return hash[:]
}

func (tx *Transaction) VerifySignature() error {
	if tx.PubKey == nil || tx.Signature == nil {
		return errors.New("missing public key or signature")
	}

	// Decode public key
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), tx.PubKey)
	if x == nil || y == nil {
		return errors.New("invalid public key")
	}
	pubKey := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	// Split signature into r and s
	if len(tx.Signature) != 64 {
		return errors.New("invalid signature length")
	}
	r := new(big.Int).SetBytes(tx.Signature[:32])
	s := new(big.Int).SetBytes(tx.Signature[32:])

	// Verify signature
	hash := tx.Hash()
	if !ecdsa.Verify(pubKey, hash, r, s) {
		return errors.New("invalid signature")
	}

	// Optional: validate that tx.Address == hex(publicKeyHash)
	pubKeyHash := sha256.Sum256(tx.PubKey)
	calculatedAddress := hex.EncodeToString(pubKeyHash[:])
	if calculatedAddress != tx.Address {
		return errors.New("address doesn't match public key")
	}

	return nil
}
