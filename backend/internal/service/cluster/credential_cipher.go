package cluster

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

const cipherEnvKey = "CREDENTIAL_CIPHER_KEY"

// CredentialCipher abstracts credential encryption/decryption behavior.
type CredentialCipher interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type AESCredentialCipher struct {
	aead cipher.AEAD
}

func NewCredentialCipher() CredentialCipher {
	keySource := os.Getenv(cipherEnvKey)
	if keySource == "" {
		keySource = "kbmanage-dev-credential-cipher-key-change-me"
	}
	key := sha256.Sum256([]byte(keySource))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	return AESCredentialCipher{aead: aead}
}

func (c AESCredentialCipher) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", errors.New("credential plaintext must not be empty")
	}
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := c.aead.Seal(nil, nonce, []byte(plaintext), nil)
	payload := append(nonce, sealed...)
	return base64.StdEncoding.EncodeToString(payload), nil
}

func (c AESCredentialCipher) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", errors.New("credential ciphertext must not be empty")
	}
	payload, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	nonceSize := c.aead.NonceSize()
	if len(payload) <= nonceSize {
		return "", errors.New("invalid credential ciphertext payload")
	}
	nonce := payload[:nonceSize]
	sealed := payload[nonceSize:]
	plain, err := c.aead.Open(nil, nonce, sealed, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
