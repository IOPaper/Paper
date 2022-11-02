package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"io"
	"math/big"
)

type EcdsaKeypair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

func LoadSecp256k1(priKey []byte) (*EcdsaKeypair, error) {
	pri, err := crypto.ToECDSA(priKey)
	if err != nil {
		return nil, err
	}
	return &EcdsaKeypair{
		PrivateKey: pri,
		PublicKey:  &pri.PublicKey,
	}, nil
}

func NewSecp256k1() (*EcdsaKeypair, error) {
	priKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &EcdsaKeypair{
		PrivateKey: priKey,
		PublicKey:  &priKey.PublicKey,
	}, nil
}

func (k *EcdsaKeypair) Sign(msg io.Reader) ([]byte, error) {
	if k.PrivateKey == nil {
		return nil, errors.New("private key is nil")
	}
	hash256 := sha256.New()
	io.Copy(hash256, msg)
	r, s, err := ecdsa.Sign(rand.Reader, k.PrivateKey, hash256.Sum(nil))
	if err != nil {
		return nil, err
	}
	cobs := k.PrivateKey.Curve.Params().P.BitLen() / 8
	sign := make([]byte, cobs*2)
	rByte, sByte := r.Bytes(), s.Bytes()
	copy(sign[cobs-len(rByte):], rByte)
	copy(sign[cobs*2-len(sByte):], sByte)
	return sign, nil
}

func (k *EcdsaKeypair) Verify(msg io.Reader, sign []byte) bool {
	if k.PublicKey == nil {
		return false
	}
	hash256 := sha256.New()
	io.Copy(hash256, msg)
	cobs := k.PublicKey.Params().P.BitLen() / 8
	r, s := new(big.Int), new(big.Int)
	r.SetBytes(sign[:cobs])
	s.SetBytes(sign[cobs:])
	return ecdsa.Verify(k.PublicKey, hash256.Sum(nil), r, s)
}

func (k *EcdsaKeypair) ExportKeypair() ([]byte, []byte) {
	return crypto.FromECDSA(k.PrivateKey), crypto.FromECDSAPub(k.PublicKey)
}
