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
		Data:      []byte("foo"),
		PublicKey: *pk.NewPublicKey(),
	}

	assert.Nil(t, tx.Sign(pk))
	assert.Nil(t, tx.Verify())

	otherPk := crypto.NewPrivateKey()
	tx.PublicKey = *otherPk.NewPublicKey()

	assert.NotNil(t, tx.Verify())
}
