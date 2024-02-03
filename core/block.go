package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/andrei0427/go-blockchain/crypto"
	"github.com/andrei0427/go-blockchain/types"
)

type Header struct {
	Version       uint32
	DataHash      types.Hash // Transactions hash
	PrevBlockHash types.Hash
	Timestamp     int64
	Height        uint32
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)

	return buf.Bytes()
}

type Block struct {
	*Header
	Transactions []*Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature

	// Cached version of header-hash
	hash types.Hash
}

func NewBlock(h *Header, tx []*Transaction) (*Block, error) {
	return &Block{
		Header:       h,
		Transactions: tx,
	}, nil
}

func NewBlockFromPrevHeader(prev *Header, txs []*Transaction) (*Block, error) {
	dataHash, err := CalculateDataHash(txs)
	if err != nil {
		return nil, err
	}

	header := &Header{
		Version:       1,
		Height:        prev.Height + 1,
		DataHash:      dataHash,
		PrevBlockHash: BlockHasher{}.Hash(prev),
		Timestamp:     time.Now().UnixNano(),
	}

	return NewBlock(header, txs)

}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)
}

func (b *Block) Sign(pk crypto.PrivateKey) error {
	sig, err := pk.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}

	b.Validator = *pk.PublicKey()
	b.Signature = sig
	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block signature is missing")
	}

	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("signature verification failed")
	}

	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}

	dataHash, err := CalculateDataHash(b.Transactions)
	if err != nil {
		return err
	}
	if dataHash != b.DataHash {
		return fmt.Errorf("block (%s) has invalid data hash", b.Hash(BlockHasher{}))
	}

	return nil
}

func (b *Block) Decode(d Decoder[*Block]) error {
	return d.Decode(b)
}

func (b *Block) Encode(e Encoder[*Block]) error {
	return e.Encode(b)
}

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}

	return b.hash
}

func CalculateDataHash(txs []*Transaction) (types.Hash, error) {
	var (
		buf  = &bytes.Buffer{}
		hash types.Hash
	)

	for _, tx := range txs {
		if err := tx.Encode(NewGobTxEncoder(buf)); err != nil {
			return hash, err
		}
	}

	hash = sha256.Sum256(buf.Bytes())

	return hash, nil
}
