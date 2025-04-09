package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

type Block struct {
	Height       int
	Timestamp    int64
	PrevHash     string
	Transactions []Transaction
	Hash         string
}

func NewBlock(height int, prevHash string, txs []Transaction) *Block {
	block := &Block{
		Height:       height,
		Timestamp:    time.Now().Unix(),
		PrevHash:     prevHash,
		Transactions: txs,
	}
	block.Hash = block.CalculateHash()
	return block
}

func (b *Block) CalculateHash() string {
	data, _ := json.Marshal(struct {
		Height       int
		Timestamp    int64
		PrevHash     string
		Transactions []Transaction
	}{
		b.Height,
		b.Timestamp,
		b.PrevHash,
		b.Transactions,
	})
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
