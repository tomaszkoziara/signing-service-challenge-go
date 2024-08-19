package main

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/signature"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/keygen"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/store/inmemory"
)

const (
	ListenAddress = ":8080"
	// TODO: add further configuration parameters here ...
)

func main() {
	store := inmemory.New()
	service := signature.New(store, map[string]signature.KeyGenerator{
		"RSA": keygen.RSA,
		"ECC": keygen.ECC,
	}, map[string]crypto.SignerFactory{
		"RSA": crypto.RSASignerFactory,
		"ECC": crypto.ECCSignerFactory,
	})
	server := api.NewServer(ListenAddress, service)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
