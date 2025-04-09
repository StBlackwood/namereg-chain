package core

import "fmt"

type Blockchain struct {
	Blocks []*Block
	State  *State
}

func NewBlockchain() *Blockchain {
	genesis := NewBlock(0, "", []Transaction{})
	return &Blockchain{
		Blocks: []*Block{genesis},
		State:  NewState(),
	}
}

func (bc *Blockchain) AddBlock(txs []Transaction) *Block {
	latest := bc.Blocks[len(bc.Blocks)-1]
	block := NewBlock(latest.Height+1, latest.Hash, txs)

	// Apply all txs to state
	for _, tx := range txs {
		err := bc.State.ApplyTransaction(tx)
		if err != nil {
			fmt.Printf("error applying transaction %v\n", err.Error())
			continue
		}
	}

	bc.Blocks = append(bc.Blocks, block)
	return block
}

func (bc *Blockchain) LatestBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}
