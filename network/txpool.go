package network

import (
	"sort"

	"github.com/andrei0427/go-blockchain/core"
	"github.com/andrei0427/go-blockchain/types"
)

type TxMapSorter struct {
	transactions []*core.Transaction
}

func NewTxMapSorter(txMap map[types.Hash]*core.Transaction) *TxMapSorter {
	txx := make([]*core.Transaction, len(txMap))

	i := 0
	for _, val := range txMap {
		txx[i] = val
		i++
	}

	s := &TxMapSorter{txx}

	sort.Sort(s)
	return s
}

func (s *TxMapSorter) Len() int {
	return len(s.transactions)
}
func (s *TxMapSorter) Swap(i, j int) {
	s.transactions[i], s.transactions[j] = s.transactions[j], s.transactions[i]
}
func (s *TxMapSorter) Less(i, j int) bool {
	return s.transactions[i].FirstSeenOn() < s.transactions[j].FirstSeenOn()
}

type TxPool struct {
	transactions map[types.Hash]*core.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

func (p *TxPool) Transactions() []*core.Transaction {
	sorter := NewTxMapSorter(p.transactions)
	return sorter.transactions
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
