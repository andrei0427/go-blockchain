package network

import (
	"github.com/andrei0427/go-blockchain/core"
	"github.com/andrei0427/go-blockchain/types"
)

type TxPool struct {
	transactions map[types.Hash]*core.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

// Caller is responsible to check if hash is assigned a tx already
func (p *TxPool) Add(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	p.transactions[hash] = tx
	return nil
}

func (p *TxPool) Has(h types.Hash) bool {
	_, ok := p.transactions[h]
	return ok
}

func (p *TxPool) Len() int {
	return len(p.transactions)
}

func (p *TxPool) Flush() {
	p.transactions = make(map[types.Hash]*core.Transaction)
}
