package core

import (
	"bytes"
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

type EncodableBlock struct {
	BlockWithoutValidator
	EncodedTransactions [][]byte
	Validator           EncodablePublicKey
}

type GobBlockEncoder struct {
	w io.Writer
}

func NewGobBlockEncoder(w io.Writer) *GobBlockEncoder {
	return &GobBlockEncoder{
		w: w,
	}
}

func (enc *GobBlockEncoder) Encode(b *Block) error {
	encBlock := EncodableBlock{
		BlockWithoutValidator: b.BlockWithoutValidator,
		EncodedTransactions:   make([][]byte, len(b.Transactions)),
		Validator: EncodablePublicKey{
			X: b.Validator.Key.X,
			Y: b.Validator.Key.Y,
		},
	}

	for i := 0; i < len(encBlock.Transactions); i++ {
		txBuf := &bytes.Buffer{}
		txEncoder := NewGobTxEncoder(txBuf)
		txEncoder.Encode(encBlock.Transactions[i])
		encBlock.EncodedTransactions[i] = txBuf.Bytes()
	}

	// remove all transactions before encoding
	encBlock.Transactions = make([]*Transaction, 0)

	return gob.NewEncoder(enc.w).Encode(encBlock)
}

type GobBlockDecoder struct {
	r     io.Reader
	Curve elliptic.Curve
}

func NewGobBlockDecoder(r io.Reader) *GobBlockDecoder {
	return &GobBlockDecoder{
		r:     r,
		Curve: elliptic.P256(),
	}
}

func (dec *GobBlockDecoder) Decode(b *Block) error {
	decBlock := new(EncodableBlock)
	err := gob.NewDecoder(dec.r).Decode(decBlock)
	if err != nil {
		return err
	}

	b.BlockWithoutValidator = decBlock.BlockWithoutValidator
	b.Validator = crypto.PublicKey{
		Key: ecdsa.PublicKey{
			Curve: dec.Curve,
			X:     decBlock.Validator.X,
			Y:     decBlock.Validator.Y,
		},
	}

	b.Transactions = make([]*Transaction, len(decBlock.EncodedTransactions))
	for i := 0; i < len(decBlock.EncodedTransactions); i++ {
		txDecoder := NewGobTxDecoder(bytes.NewReader(decBlock.EncodedTransactions[i]))
		tx := new(Transaction)
		txDecoder.Decode(tx)

		b.Transactions[i] = tx
	}

	return nil
}
