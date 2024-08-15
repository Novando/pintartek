package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// EncryptAES encrypts a plain text string using AES with a given key.
func EncryptAES(plainText, key string) (string, error) {
	// Create a new AES cipher using the secret key.
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Generate a new GCM (Galois/Counter Mode) cipher mode.
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce of size required by GCM.
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the plain text and prepend the nonce.
	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)

	// Encode the encrypted text as a base64 string.
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptAES decrypts an encrypted string using AES with a given key.
func DecryptAES(encryptedText, key string) (string, error) {
	// Decode the base64-encoded encrypted text.
	cipherText, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher using the secret key.
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Generate a new GCM (Galois/Counter Mode) cipher mode.
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract the nonce from the encrypted text.
	nonceSize := aesGCM.NonceSize()
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	// Decrypt the text.
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
