package core

import (
	"testing"

	"github.com/andrei0427/go-blockchain/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSignTransaction(t *testing.T) {
	pk := crypto.NewPrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(pk))
	assert.NotNil(t, tx.Signature)
}

func TestVerifyTransaction(t *testing.T) {
	pk := crypto.NewPrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
		From: *pk.PublicKey(),
	}

	assert.Nil(t, tx.Sign(pk))
	assert.Nil(t, tx.Verify())

	otherPk := crypto.NewPrivateKey()
	tx.From = *otherPk.PublicKey()

	assert.NotNil(t, tx.Verify())
}

func randomSignedTx(t *testing.T) *Transaction {
	pk := crypto.NewPrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(pk))

	return tx
}
