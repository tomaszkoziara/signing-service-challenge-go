package inmemory

import (
	"context"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/store"
	"github.com/google/uuid"
)



type InMemoryStore struct {
	DB map[string]store.SignatureDevice
}

func (ims *InMemoryStore) CreateSignatureDevice(ctx context.Context, sigDevice store.SignatureDevice) error {
	sigDevice.Version = uuid.NewString() // this should be handled in the DB store on the stored data, it's just there to show an optimistic lock use
	ims.DB[sigDevice.ID] = sigDevice
	return nil
}

func (ims *InMemoryStore) GetSignatureDevice(ctx context.Context, id string) (store.SignatureDevice, error) {
	if signDevice, found := ims.DB[id]; found {
		return signDevice, nil
	}

	return store.SignatureDevice{}, store.ErrDeviceNotFound
}

func (ims *InMemoryStore) UpdateSignatureDevice(ctx context.Context, id string, updateSignDevice store.UpdateSignatureDevice) error {
	signDevice, found := ims.DB[id]
	if !found {
		return store.ErrDeviceNotFound
	}

	// that function would be executed with an optimistic lock and atomically on the DB
	// here I'm just pretending that it's happening atomically
	if signDevice.Version == updateSignDevice.Version {
		signDevice.SignatureCounter = updateSignDevice.SignatureCounter
		signDevice.LastSignature = updateSignDevice.LastSignature
		signDevice.Version = uuid.NewString()
		ims.DB[id] = signDevice
	}

	return nil
}

func New() *InMemoryStore {
	return &InMemoryStore{
		DB: map[string]store.SignatureDevice{},
	}
}