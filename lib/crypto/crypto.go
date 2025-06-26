package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

type Encryptor interface {
	EncryptMessage(message string) (string, error)
}

type Decryptor interface {
	DecryptMessage(message string) (string, error)
}

type EncryptorImpl struct {
	cb cipher.Block
}

func NewEncryptor() (*EncryptorImpl, error) {
	key, ok := os.LookupEnv("SECRET_KEY")
	if !ok {
		return nil, fmt.Errorf("SECRET_KEY environment variable not set")
	}
	if len(key) != 16 { // AES-128 requires a 16-byte key
		return nil, fmt.Errorf("key must be 16 bytes long")
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("could not create new cipher: %v", err)
	}
	return &EncryptorImpl{cb: block}, nil
}

func (e *EncryptorImpl) EncryptMessage(message string) (string, error) {
	byteMsg := []byte(message)
	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("could not encrypt: %v", err)
	}

	stream := cipher.NewCFBEncrypter(e.cb, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

type DecryptorImpl struct {
	key []byte
}

func NewDecryptor() (Decryptor, error) {
	key, ok := os.LookupEnv("SECRET_KEY")
	if !ok {
		return nil, fmt.Errorf("SECRET_KEY environment variable not set")
	}
	if len(key) != 16 { // AES-128 requires a 16-byte key
		return nil, fmt.Errorf("key must be 16 bytes long")
	}
	return &DecryptorImpl{key: []byte(key)}, nil
}

func (d *DecryptorImpl) DecryptMessage(message string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("could not base64 decode: %v", err)
	}

	block, err := aes.NewCipher(d.key)
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext block size")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
