package core

import (
	"testing"
	"time"

	"github.com/andrei0427/go-blockchain/crypto"
	"github.com/andrei0427/go-blockchain/types"
	"github.com/stretchr/testify/assert"
)

func TestSignBlock(t *testing.T) {
	pk := crypto.NewPrivateKey()
	b := randomBlock(0, types.Hash{})

	assert.Nil(t, b.Sign(pk))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	pk := crypto.NewPrivateKey()
	b := randomBlock(0, types.Hash{})

	assert.Nil(t, b.Sign(pk))
	assert.Nil(t, b.Verify())

	otherPk := crypto.NewPrivateKey()
	b.Validator = *otherPk.PublicKey()

	assert.NotNil(t, b.Verify())
}

func randomBlock(height uint32, prevBlockHash types.Hash) *Block {
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}

	return NewBlock(header, []Transaction{})
}

func randomBlockSigned(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	pk := crypto.NewPrivateKey()
	b := randomBlock(height, prevBlockHash)
	tx := randomSignedTx(t)
	b.AddTransaction(tx)

	assert.Nil(t, b.Sign(pk))

	return b
}
