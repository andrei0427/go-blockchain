package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyPair_Sign_Verify_Valid(t *testing.T) {
	pk := NewPrivateKey()
	pubKey := pk.NewPublicKey()
	msg := []byte("hello world")

	sig, err := pk.Sign(msg)
	assert.Nil(t, err)

	assert.True(t, sig.Verify(*pubKey, msg))
}

func TestKeyPair_Sign_Verify_Invalid(t *testing.T) {
	pk := NewPrivateKey()
	pubKey := pk.NewPublicKey()
	msg := []byte("hello world")

	sig, err := pk.Sign(msg)
	assert.Nil(t, err)

	otherPk := NewPrivateKey()
	otherPubKey := otherPk.NewPublicKey()

	assert.False(t, sig.Verify(*otherPubKey, msg))
	assert.False(t, sig.Verify(*pubKey, []byte("world hello")))
}
