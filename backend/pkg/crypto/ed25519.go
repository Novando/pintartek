package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
)

// GenerateKeyPairEd25519 generates a new Ed25519 public/private key pair.
func GenerateKeyPairEd25519() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return pubKey, privKey, nil
}

// ValidateKeyPairEd25519 checks if the provided private key matches the public key.
func ValidateKeyPairEd25519(pubKey ed25519.PublicKey, privKey ed25519.PrivateKey) error {
	if !pubKey.Equal(privKey.Public().(ed25519.PublicKey)) {
		return errors.New("private key does not match the public key")
	}
	return nil
}
