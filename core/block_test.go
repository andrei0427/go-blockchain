package core

import (
	"testing"
	"time"

	"github.com/andrei0427/go-blockchain/crypto"
	"github.com/andrei0427/go-blockchain/types"
	"github.com/stretchr/testify/assert"
)

func randomBlock(height uint32) *Block {
	header := &Header{
		Version:       1,
		PrevBlockHash: types.RandomHash(),
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}

	tx := Transaction{
		Data: []byte("foo"),
	}

	return NewBlock(header, []Transaction{tx})
}

func randomBlockSigned(t *testing.T, height uint32) *Block {
	pk := crypto.NewPrivateKey()
	b := randomBlock(height)
	assert.Nil(t, b.Sign(pk))

	return b
}

func TestSignBlock(t *testing.T) {
	pk := crypto.NewPrivateKey()
	b := randomBlock(0)

	assert.Nil(t, b.Sign(pk))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	pk := crypto.NewPrivateKey()
	b := randomBlock(0)

	assert.Nil(t, b.Sign(pk))
	assert.Nil(t, b.Verify())

	otherPk := crypto.NewPrivateKey()
	b.Validator = *otherPk.NewPublicKey()

	assert.NotNil(t, b.Verify())
}
