package util

import (
	"crypto/rand"
	"testing"

	"github.com/andrei0427/go-blockchain/core"
	"github.com/andrei0427/go-blockchain/crypto"
	"github.com/andrei0427/go-blockchain/types"
	"github.com/stretchr/testify/assert"
)

func RandomBytes(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

func RandomHash() types.Hash {
	return types.HashFromBytes(RandomBytes(32))
}

func NewRandomTransaction(size int) *core.Transaction {
	return core.NewTransaction(RandomBytes(size))
}

func NewRandomSignedTx(t *testing.T, pk crypto.PrivateKey, size int) *core.Transaction {
	tx := NewRandomTransaction(size)
	assert.Nil(t, tx.Sign(pk))
	return tx
}
