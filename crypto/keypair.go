package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/andrei0427/go-blockchain/types"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

func NewPrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		key: key,
	}
}

func (k *PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		r: r,
		s: s,
	}, nil
}

func (pk *PrivateKey) NewPublicKey() *PublicKey {
	return &PublicKey{
		key: &pk.key.PublicKey,
	}
}

type PublicKey struct {
	key *ecdsa.PublicKey
}

func (k *PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.key, k.key.X, k.key.Y)
}

func (k *PublicKey) Address() types.Address {
	h := sha256.Sum256(k.ToSlice())

	return types.AddressFromBytes(h[len(h)-20:]) // last 20 bytes
}

type Signature struct {
	r, s *big.Int
}

func (sig *Signature) Verify(pubKey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubKey.key, data, sig.r, sig.s)
}
