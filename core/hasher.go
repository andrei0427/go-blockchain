package core

import (
	"crypto/sha256"

	"github.com/andrei0427/go-blockchain/types"
)

type Hasher[t any] interface {
	Hash(t) types.Hash
}

type BlockHasher struct{}

func (BlockHasher) Hash(b *Block) types.Hash {
	h := sha256.Sum256(b.HeaderData())
	return types.Hash(h)
}
