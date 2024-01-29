package core

import (
	"testing"

	"github.com/andrei0427/go-blockchain/types"
	"github.com/stretchr/testify/assert"
)

func TestBlockchain(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
}

func TestHasBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.True(t, bc.HasBlock(0))
	assert.False(t, bc.HasBlock(1))
}

func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	lenBlocks := (1000)
	for i := 0; i < lenBlocks; i++ {
		block := randomBlockSigned(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint32(lenBlocks))
	assert.Equal(t, len(bc.headers), lenBlocks+1)

	assert.NotNil(t, bc.AddBlock(randomBlockSigned(t, 90, types.Hash{})))
}

func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	lenBlocks := 1000

	for i := 0; i < lenBlocks; i++ {
		block := randomBlockSigned(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))

		header, err := bc.GetHeader(block.Height)
		assert.Nil(t, err)
		if err != nil {
			break
		}

		assert.Equal(t, header, block.Header)
	}
}

func TestAddBlockTooHigh(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	assert.Nil(t, bc.AddBlock(randomBlockSigned(t, 1, getPrevBlockHash(t, bc, uint32(1)))))
	assert.NotNil(t, bc.AddBlock(randomBlockSigned(t, 10, types.Hash{})))
}

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(randomBlock(0, types.Hash{}))
	assert.Nil(t, err)

	return bc
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)

	return BlockHasher{}.Hash(prevHeader)
}
