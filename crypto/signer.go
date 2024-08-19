package crypto

import (
	"crypto"
	"crypto/rand"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

type SignerFactory func(privateKey []byte) (Signer, error)

type RSASigner struct {
	KeyPair *RSAKeyPair
}

func (s *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	return s.KeyPair.Private.Sign(rand.Reader, dataToBeSigned, crypto.SHA256)
}

func RSASignerFactory(privateKey []byte) (Signer, error) {
	marshaler := NewRSAMarshaler()
	keyPair, err := marshaler.Unmarshal(privateKey)
	if err != nil {
		return &RSASigner{}, err
	}

	return &RSASigner{
		KeyPair: keyPair,
	}, nil
}

type ECCSigner struct {
	KeyPair *ECCKeyPair
}

func (s *ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	return s.KeyPair.Private.Sign(rand.Reader, dataToBeSigned, crypto.SHA256)
}

func ECCSignerFactory(privateKey []byte) (Signer, error) {
	marshaler := NewECCMarshaler()
	keyPair, err := marshaler.Decode(privateKey)
	if err != nil {
		return &ECCSigner{}, err
	}

	return &ECCSigner{
		KeyPair: keyPair,
	}, nil
}