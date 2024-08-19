package storestub

import (
	"context"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/store"
)

type Store struct {
	CreateSignatureDeviceFn CreateSignatureDeviceFn
	GetSignatureDeviceFn GetSignatureDeviceFn
	UpdateSignatureDeviceFn UpdateSignatureDeviceFn
}

type CreateSignatureDeviceFn func (ctx context.Context, sigDevice store.SignatureDevice) error
type GetSignatureDeviceFn func(ctx context.Context, id string) (store.SignatureDevice, error)
type UpdateSignatureDeviceFn func(ctx context.Context, id string, updateSignDevice store.UpdateSignatureDevice) error

var defaultCreateSignatureDeviceFn = func(ctx context.Context, sigDevice store.SignatureDevice) error {
	panic("not implemented")
}

var defaultGetSignatureDeviceFn = func(ctx context.Context, id string) (store.SignatureDevice, error) {
	panic("not implemented")
}

var defaultUpdateSignatureDeviceFn = func(ctx context.Context, id string, updateSignDevice store.UpdateSignatureDevice) error {
	panic("not implemented")
}

func (s *Store) CreateSignatureDevice(ctx context.Context, sigDevice store.SignatureDevice) error {
	return s.CreateSignatureDeviceFn(ctx, sigDevice)
}

func (s *Store) GetSignatureDevice(ctx context.Context, id string) (store.SignatureDevice, error) {
	return s.GetSignatureDeviceFn(ctx, id)
}

func (s *Store) UpdateSignatureDevice(ctx context.Context, id string, updateSignDevice store.UpdateSignatureDevice) error {
	return s.UpdateSignatureDeviceFn(ctx, id, updateSignDevice)
}

func New() *Store {
	return &Store{
		CreateSignatureDeviceFn: defaultCreateSignatureDeviceFn,
		GetSignatureDeviceFn: defaultGetSignatureDeviceFn,
		UpdateSignatureDeviceFn: defaultUpdateSignatureDeviceFn,
	}
}