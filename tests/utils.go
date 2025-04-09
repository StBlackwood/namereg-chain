package tests

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"namereg-chain/core"
	"namereg-chain/network"
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

func start2Nodes() {
	// Setup node1
	chain1 := core.NewBlockchain()
	peer1 := network.NewPeerClient([]string{"http://localhost:8082"})
	api1 := network.NewAPIServerWithPeers(chain1, peer1)
	go api1.Start("localhost:8081")

	// Setup node2
	chain2 := core.NewBlockchain()
	peer2 := network.NewPeerClient([]string{"http://localhost:8081"})
	api2 := network.NewAPIServerWithPeers(chain2, peer2)
	go api2.Start("localhost:8082")
}
