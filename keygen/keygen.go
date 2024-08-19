package keygen

import "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"

// RSA returns public and private keys.
func RSA() ([]byte, []byte, error) {
	generator := crypto.RSAGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		return nil, nil, err
	}

	marshaler := crypto.NewRSAMarshaler()
	return marshaler.Marshal(*keyPair)
}

// ECC returns public and private keys.
func ECC() ([]byte, []byte, error) {
	generator := crypto.ECCGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		return nil, nil, err
	}

	marshaler := crypto.NewECCMarshaler()
	return marshaler.Encode(*keyPair)
}