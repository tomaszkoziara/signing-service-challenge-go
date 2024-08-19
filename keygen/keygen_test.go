package keygen_test

import (
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/keygen"
	"github.com/stretchr/testify/assert"
)

func TestRSAHappyPath(t *testing.T) {
	public, private, err := keygen.RSA()
	assert.NoError(t, err)
	assert.NotEmpty(t, public)
	assert.NotEmpty(t, private)
}

func TestECCHappyPath(t *testing.T) {
	public, private, err := keygen.ECC()
	assert.NoError(t, err)
	assert.NotEmpty(t, public)
	assert.NotEmpty(t, private)
}