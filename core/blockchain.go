package core

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type Blockchain struct {
	store     Storage
	validator Validator

	lock    sync.RWMutex
	headers []*Header
}

func NewBlockchain(genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		headers: []*Header{},
		store:   NewMemoryStore(),
	}

	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)

	return bc, err
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	bc.addBlockWithoutValidation(b)

	return nil
}

func (bc *Blockchain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high, current height is %d", height, bc.Height())
	}

	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return bc.headers[height], nil
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.lock.Unlock()

	logrus.WithFields(logrus.Fields{
		"height": b.Height,
		"hash":   b.Hash(BlockHasher{}),
	}).Info("adding new block")

	return bc.store.Put(b)
}
