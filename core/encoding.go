package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"io"
	"math/big"

	"github.com/andrei0427/go-blockchain/crypto"
)

type Encoder[T any] interface {
	Encode(T) error
}

type Decoder[T any] interface {
	Decode(T) error
}

type GobTxEncoder struct {
	w io.Writer
}

// Need this to fix private properties within ecdsa.PublicKey package for elliptic.P256() curve
type EncodablePublicKey struct {
	X, Y *big.Int
}
type EncodableTransaction struct {
	TransactionWithoutPublicKey
	PublicKey EncodablePublicKey
}

func NewGobTxEncoder(w io.Writer) *GobTxEncoder {
	return &GobTxEncoder{
		w: w,
	}
}

func (e *GobTxEncoder) Encode(tx *Transaction) error {
	encTx := EncodableTransaction{
		TransactionWithoutPublicKey: tx.TransactionWithoutPublicKey,
		PublicKey: EncodablePublicKey{
			X: tx.From.Key.X,
			Y: tx.From.Key.Y,
		},
	}

	return gob.NewEncoder(e.w).Encode(encTx)
}

type GobTxDecoder struct {
	r     io.Reader
	Curve elliptic.Curve
}

func NewGobTxDecoder(r io.Reader) *GobTxDecoder {
	return &GobTxDecoder{
		r:     r,
		Curve: elliptic.P256(),
	}
}

func (d *GobTxDecoder) Decode(tx *Transaction) error {
	decTx := new(EncodableTransaction)
	err := gob.NewDecoder(d.r).Decode(decTx)
	if err != nil {
		return err
	}

	tx.TransactionWithoutPublicKey = decTx.TransactionWithoutPublicKey
	tx.From = crypto.PublicKey{
		Key: ecdsa.PublicKey{
			Curve: d.Curve,
			X:     decTx.PublicKey.X,
			Y:     decTx.PublicKey.Y,
		},
	}

	return nil
}
