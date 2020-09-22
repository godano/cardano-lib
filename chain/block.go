package chain

import "github.com/godano/cardano-lib/time"

type Block struct {
	_                 struct{} `cbor:",toarray"`
	Header            BlockHeader
	TransactionBodies []TransactionBody
}

type BlockHeader struct {
	Hash          string // Hash of the block to which this header belongs.
	Height        uint64 // Height indicates
	Slot          uint64 // Slot number in which
	StakePoolVKey string //
}

// Transactions returns an array of transactions that are contained in this block.
//
func (block *Block) Transactions() []Transaction {
	n := len(block.TransactionBodies)
	trans := make([]Transaction, n)
	for i := 0; i < n; i = + 1 {
		trans[i] = Transaction{
			Body: block.TransactionBodies[i],
		}
	}
	return trans
}

func (block *Block) SlotDate() *time.AbstractSlotDate {
	return nil
}