package core

import "fmt"

type Validator interface {
	ValidateBlock(*Block) error
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc: bc,
	}
}

func (v *BlockValidator) ValidateBlock(b *Block) error {
	if v.bc.HasBlock(b.Height) {
		return fmt.Errorf("failed to insert block with hash %s into chain as it already contains block %d", b.Hash(BlockHasher{}), b.Height)
	}

	hasher := BlockHasher{}
	blockHash := b.Hash(hasher)
	if b.Height != v.bc.Height()+1 {
		return fmt.Errorf("block (%s) is too high (%d) - next block should be %d", blockHash, b.Height, v.bc.Height()+1)
	}

	prevHeader, err := v.bc.GetHeader(v.bc.Height())
	if err != nil {
		return err
	}

	prevHash := hasher.Hash(prevHeader)
	if prevHash != b.PrevBlockHash {
		return fmt.Errorf("invalid previous block hash provided %s for block %s - expected %s", b.PrevBlockHash, blockHash, prevHash)
	}

	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}
