package core

import (
	"fmt"

	"github.com/andrei0427/go-blockchain/crypto"
)

type Transaction struct {
	Data []byte

	From      crypto.PublicKey
	Signature *crypto.Signature
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
