package signature_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/signature"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/store"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/store/storestub"
	"github.com/stretchr/testify/assert"
)

func TestCreateSignatureDeviceValidationErrors(t *testing.T) {
	tests := []struct{
		name string
		newSignatureDevice signature.NewSignatureDevice
		errMsg string
	}{
		{
			name: "valid signature device",
			newSignatureDevice: signature.NewSignatureDevice{
				ID: "some-id",
				Tenant: "some-tenant",
				SignatureAlg: "RSA",
				Label: "some-label",
			},
			errMsg: "",
		},
		{
			name: "label is optional",
			newSignatureDevice: signature.NewSignatureDevice{
				ID: "some-id",
				Tenant: "some-tenant",
				SignatureAlg: "RSA",
				Label: "",
			},
			errMsg: "",
		},
		{
			name: "signature device with missing ID",
			newSignatureDevice: signature.NewSignatureDevice{
				ID: "",
				Tenant: "some-tenant",
				SignatureAlg: "RSA",
				Label: "some-label",
			},
			errMsg: "Field validation for 'ID' failed on the 'required' tag",
		},
		{
			name: "signature device with missing tenant",
			newSignatureDevice: signature.NewSignatureDevice{
				ID: "some-id",
				Tenant: "",
				SignatureAlg: "RSA",
				Label: "some-label",
			},
			errMsg: "Field validation for 'Tenant' failed on the 'required' tag",
		},
		{
			name: "signature device with unsupported signing algorithm",
			newSignatureDevice: signature.NewSignatureDevice{
				ID: "some-id",
				Tenant: "some-tenant",
				SignatureAlg: "ABC123",
				Label: "some-label",
			},
			errMsg: "Field validation for 'SignatureAlg' failed on the 'oneof' tag",
		},
	}

	storeStub := storestub.New()
	storeStub.CreateSignatureDeviceFn = func(ctx context.Context, sigDevice store.SignatureDevice) error {
		return nil
	}
	service := signature.New(storeStub, map[string]signature.KeyGenerator{
		"RSA": func() ([]byte, []byte, error) {return []byte{}, []byte{}, nil},
		"ABC123": func() ([]byte, []byte, error) {return []byte{}, []byte{}, nil},
	}, map[string]crypto.SignerFactory{})
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := service.CreateSignatureDevice(context.Background(), tc.newSignatureDevice)
			if tc.errMsg != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestHappyPath(t *testing.T) {
	ctx := context.Background()

	newSignatureDevice := signature.NewSignatureDevice{
		ID: "some-id",
		Tenant: "some-tenant",
		SignatureAlg: "RSA",
		Label: "some-label",
	}
	publicKey := []byte{1,2,3}
	privateKey := []byte{4,5,6}

	storeStub := storestub.New()
	storeStub.CreateSignatureDeviceFn = func(ctx context.Context, sigDevice store.SignatureDevice) error {
		assert.Equal(t, newSignatureDevice.ID, sigDevice.ID)
		assert.Equal(t, newSignatureDevice.Tenant, sigDevice.Tenant)
		assert.Equal(t, newSignatureDevice.SignatureAlg, sigDevice.SignatureAlg)
		assert.Equal(t, newSignatureDevice.Label, sigDevice.Label)
		assert.Equal(t, publicKey, sigDevice.PublicKey)
		assert.Equal(t, privateKey, sigDevice.PrivateKey)

		return nil
	}
	service := signature.New(storeStub, map[string]signature.KeyGenerator{
		"RSA": func() ([]byte, []byte, error) {return publicKey, privateKey, nil},
	}, map[string]crypto.SignerFactory{})

	err := service.CreateSignatureDevice(ctx, newSignatureDevice)
	assert.NoError(t, err)
}

func TestReturnErrorOnSignatureError(t *testing.T) {
	ctx := context.Background()

	newSignatureDevice := signature.NewSignatureDevice{
		ID: "some-id",
		Tenant: "some-tenant",
		SignatureAlg: "RSA",
		Label: "some-label",
	}

	storeStub := storestub.New()
	storeStub.CreateSignatureDeviceFn = func(ctx context.Context, sigDevice store.SignatureDevice) error {
		return nil
	}
	service := signature.New(storeStub, map[string]signature.KeyGenerator{
		"RSA": func() ([]byte, []byte, error) {return nil, nil, errors.New("some error")},
	}, map[string]crypto.SignerFactory{})

	err := service.CreateSignatureDevice(ctx, newSignatureDevice)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "error while generating a new key pair")
}

func TestReturnErrorOnMissingKeyPairGeneratorFunction(t *testing.T) {
	ctx := context.Background()

	newSignatureDevice := signature.NewSignatureDevice{
		ID: "some-id",
		Tenant: "some-tenant",
		SignatureAlg: "RSA",
		Label: "some-label",
	}

	storeStub := storestub.New()
	storeStub.CreateSignatureDeviceFn = func(ctx context.Context, sigDevice store.SignatureDevice) error {
		return nil
	}
	service := signature.New(storeStub, map[string]signature.KeyGenerator{}, map[string]crypto.SignerFactory{})

	err := service.CreateSignatureDevice(ctx, newSignatureDevice)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing key pair generator for algorithm 'RSA'")
}

// TODO: add tests for GetSignatureDevice