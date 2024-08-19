package signature

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/store"
	"github.com/go-playground/validator"
)

type SignatureDeviceService interface {
	CreateSignatureDevice(ctx context.Context, newSignDev NewSignatureDevice) error
	GetSignatureDevice(ctx context.Context, id string) (SignatureDevice, error)
	SignData(ctx context.Context, id string, dataToSign string) (Signature, error)
}

type KeyGenerator func()([]byte, []byte, error)

type NewSignatureDevice struct {
	ID string `validate:"required"`
	Tenant string `validate:"required"`
	SignatureAlg string `validate:"required,oneof=ECC RSA"`
	Label string
}

type SignatureDevice struct {
	ID string // are these unique or should we have our own IDs?
	SignatureAlg string
	Label string
	SignatureCounter int
}

type Signature struct {
	Signature string
	SignedData string
}

type Service struct {
	store store.Store
	keyGenerators map[string]KeyGenerator
	signers map[string]crypto.SignerFactory // TODO: should be on the same structure with key generators in order to prevent misalignment
}

func (s *Service) CreateSignatureDevice(ctx context.Context, newSignDev NewSignatureDevice) error {
	validate := validator.New()
	if err := validate.Struct(newSignDev); err != nil {
		return err
	}

	var publicKey []byte
	var privateKey []byte
	var err error
	if keyGenerator, found := s.keyGenerators[newSignDev.SignatureAlg]; found {
		publicKey, privateKey, err = keyGenerator()
		if err != nil {
			return fmt.Errorf("error while generating a new key pair: %w", err)
		}
	} else {
		return fmt.Errorf("missing key pair generator for algorithm '%v'", newSignDev.SignatureAlg)
	}

	err = s.store.CreateSignatureDevice(ctx, store.SignatureDevice{
		ID: newSignDev.ID,
		Tenant: newSignDev.Tenant,
		SignatureAlg: newSignDev.SignatureAlg,
		Label: newSignDev.Label,
		PublicKey: publicKey,
		PrivateKey: privateKey,
		SignatureCounter: 0,
		LastSignature: base64.StdEncoding.EncodeToString([]byte(newSignDev.ID)),
	})
	if err != nil {
		return fmt.Errorf("error creating signature device: %w", err) // should be as domain error rather than store error
	}

	return nil
}

func (s *Service) GetSignatureDevice(ctx context.Context, id string) (SignatureDevice, error) {
	signDevice, err := s.store.GetSignatureDevice(ctx, id)
	if err != nil {
		return SignatureDevice{}, err // should be as domain error rather than store error
	}

	return SignatureDevice{
		ID: signDevice.ID,
		SignatureAlg: signDevice.SignatureAlg,
		Label: signDevice.Label,
		SignatureCounter: 0,
	}, nil
}

func (s *Service) SignData(ctx context.Context, id string, dataToSign string) (Signature, error) {
	signDevice, err := s.store.GetSignatureDevice(ctx, id)
	if err != nil {
		return Signature{}, fmt.Errorf("error getting signature: %w", err)
	}

	signerFactory, found := s.signers[signDevice.SignatureAlg]
	if !found {
		return Signature{}, fmt.Errorf("missing signer factory for algorithm '%v'", signDevice.SignatureAlg)
	}

	signer, err := signerFactory(signDevice.PrivateKey)
	if err != nil {
		return Signature{}, fmt.Errorf("error creating signer for algorithm '%v': %w", signDevice.SignatureAlg, err)
	}

	dataToBeSigned := fmt.Sprintf("%v_%v_%v", signDevice.SignatureCounter, dataToSign, base64.StdEncoding.EncodeToString([]byte(signDevice.LastSignature)))
	dataToBeSignedHashed := sha256.Sum256([]byte(dataToBeSigned))
	signature, err := signer.Sign(dataToBeSignedHashed[:])
	if err != nil {
		return Signature{}, fmt.Errorf("error signing data: %w", err)
	}

	signatureBase64 := base64.StdEncoding.EncodeToString(signature)
	s.store.UpdateSignatureDevice(ctx, id, store.UpdateSignatureDevice{
		SignatureCounter: signDevice.SignatureCounter + 1,
		LastSignature: signatureBase64,
		Version: signDevice.Version,
	})

	return Signature{
		Signature: signatureBase64,
		SignedData: dataToBeSigned,
	}, nil
}

func New(store store.Store, keyGenerators map[string]KeyGenerator, signers map[string]crypto.SignerFactory) *Service {
	return &Service{
		store: store,
		keyGenerators: keyGenerators,
		signers: signers,
	}
}