package store

import (
	"context"
	"errors"
)

var (
	ErrDeviceNotFound = errors.New("device not found")
)

type SignatureDevice struct {
	ID string
	Tenant string
	SignatureAlg string
	Label string
	PublicKey []byte
	PrivateKey []byte
	SignatureCounter int
	LastSignature string
	Version string // this field should belong to the stored data, but I'm using this I/O struct also as stored data for simplicity
}

type UpdateSignatureDevice struct {
	SignatureCounter int
	LastSignature string
	Version string
}

type Store interface {
	CreateSignatureDevice(ctx context.Context, sigDevice SignatureDevice) error
	UpdateSignatureDevice(ctx context.Context, id string, updateSignDevice UpdateSignatureDevice) error
	GetSignatureDevice(ctx context.Context, id string) (SignatureDevice, error)
}