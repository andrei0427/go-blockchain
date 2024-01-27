package core

import (
	"bytes"
	"testing"
	"time"

	"github.com/andrei0427/go-blockchain/types"
	"github.com/stretchr/testify/assert"
)

func TestHeader_Encode_Decode(t *testing.T) {
	h := &Header{
		Version:   1,
		PrevBlock: types.RandomHash(),
		Timestamp: time.Now().UnixNano(),
		Height:    10,
		Nonce:     42069,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, h.EncodeBinary(buf))

	hDecode := &Header{}
	assert.Nil(t, hDecode.DecodeBinary(buf))
	assert.Equal(t, h, hDecode)
}

func TestBlock_Encode_Decode(t *testing.T) {
	block := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			Timestamp: time.Now().UnixNano(),
			Height:    10,
			Nonce:     42069,
		},
		Transactions: nil,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, block.EncodeBinary(buf))

	hDecode := &Block{}
	assert.Nil(t, hDecode.DecodeBinary(buf))
	assert.Equal(t, block, hDecode)

}

func TestBlockHash(t *testing.T) {
	b := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			Timestamp: time.Now().UnixNano(),
			Height:    10,
			Nonce:     42069,
		},
		Transactions: nil,
	}

	h := b.Hash()
	assert.False(t, h.IsZero())
}
