package network

import (
	"sync"

	"github.com/andrei0427/go-blockchain/core"
	"github.com/andrei0427/go-blockchain/types"
)

type TxPool struct {
	all     *TxSortedMap
	pending *TxSortedMap

	// When all.Length() reaches maxLength, the oldest transaction is removed to make room
	maxLength int
}

func NewTxPool(maxLength int) *TxPool {
	return &TxPool{
		all:       NewTxSortedMap(),
		pending:   NewTxSortedMap(),
		maxLength: maxLength,
	}
}

func (p *TxPool) Add(tx *core.Transaction) {
	// remove oldest transaction if limit is reached
	if p.all.Count() >= p.maxLength {
		oldest := p.all.First()
		p.all.Remove(oldest.Hash(core.TxHasher{}))
	}

	if !p.Contains(tx.Hash(core.TxHasher{})) {
		p.all.Add(tx)
		p.pending.Add(tx)
	}
}

func (p *TxPool) Contains(hash types.Hash) bool {
	return p.all.Contains(hash)
}

func (p *TxPool) Pending() []*core.Transaction {
	return p.pending.txx.Data
}

func (p *TxPool) ClearPending() {
	p.pending.Clear()
}

func (p *TxPool) PendingCount() int {
	return p.pending.Count()
}

type TxMap map[types.Hash]*core.Transaction
type TxSortedMap struct {
	lock   sync.RWMutex
	lookup TxMap
	txx    *types.List[*core.Transaction]
}

func NewTxSortedMap() *TxSortedMap {
	return &TxSortedMap{
		lookup: make(TxMap),
		txx:    types.NewList[*core.Transaction](),
	}
}

func (m *TxSortedMap) First() *core.Transaction {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if m.Count() == 0 {
		return nil
	}

	first := m.txx.Get(0)
	return m.lookup[first.Hash(core.TxHasher{})]
}

func (m *TxSortedMap) Get(h types.Hash) *core.Transaction {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.lookup[h]
}

func (m *TxSortedMap) Add(tx *core.Transaction) {
	hash := tx.Hash(core.TxHasher{})

	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.lookup[hash]; !ok {
		m.lookup[hash] = tx
		m.txx.Insert(tx)
	}
}

func (m *TxSortedMap) Remove(h types.Hash) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.txx.Remove(m.lookup[h])
	delete(m.lookup, h)
}

func (m *TxSortedMap) Count() int {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return len(m.lookup)
}

func (m *TxSortedMap) Contains(h types.Hash) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	_, ok := m.lookup[h]
	return ok
}

func (m *TxSortedMap) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.lookup = make(TxMap)
	m.txx.Clear()
}
