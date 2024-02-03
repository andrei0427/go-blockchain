package core

import (
	"fmt"

	"github.com/andrei0427/go-blockchain/crypto"
	"github.com/andrei0427/go-blockchain/types"
)

type TransactionWithoutPublicKey struct {
	Data      []byte
	Signature *crypto.Signature

	// cached
	hash types.Hash
}

type Transaction struct {
	TransactionWithoutPublicKey
	From crypto.PublicKey
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		TransactionWithoutPublicKey: TransactionWithoutPublicKey{
			Data: data,
		},
	}
}

func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() {
		tx.hash = hasher.Hash(tx)
	}

	return tx.hash
}

func (tx *Transaction) Sign(pk crypto.PrivateKey) error {
	sig, err := pk.Sign(tx.Data)
	if err != nil {
		return err
	}

	tx.From = *pk.PublicKey()
	tx.Signature = sig

	return nil
}

func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}

func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx)
}

func (tx *Transaction) Encode(enc Encoder[*Transaction]) error {
	return enc.Encode(tx)
}
