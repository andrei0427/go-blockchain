package core

import (
	"bytes"
	"testing"
	"time"

	"github.com/andrei0427/go-blockchain/crypto"
	"github.com/andrei0427/go-blockchain/types"
	"github.com/stretchr/testify/assert"
)

func TestSignBlock(t *testing.T) {
	pk := crypto.NewPrivateKey()
	b := randomBlock(t, 0, types.Hash{})

	assert.Nil(t, b.Sign(pk))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	pk := crypto.NewPrivateKey()
	b := randomBlock(t, 0, types.Hash{})

	assert.Nil(t, b.Sign(pk))
	assert.Nil(t, b.Verify())

	otherPk := crypto.NewPrivateKey()
	b.Validator = *otherPk.PublicKey()

	assert.NotNil(t, b.Verify())
}

func TestEncodeDecodeBlock(t *testing.T) {
	b := randomBlock(t, 1, types.Hash{})
	buf := &bytes.Buffer{}

	assert.Nil(t, b.Encode(NewGobBlockEncoder(buf)))

	decoded := new(Block)
	assert.Nil(t, decoded.Decode(NewGobBlockDecoder(buf)))

	assert.Equal(t, b, decoded)
}

func randomBlock(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	pk := crypto.NewPrivateKey()
	tx := randomSignedTx(t)
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}

	b, err := NewBlock(header, []*Transaction{tx})
	assert.Nil(t, err)

	dataHash, err := CalculateDataHash(b.Transactions)
	assert.Nil(t, err)
	b.Header.DataHash = dataHash

	assert.Nil(t, b.Sign(pk))

	return b
}
