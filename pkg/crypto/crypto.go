package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"hash"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

func CreateRandomSalt(saltLen int) (salt []byte, err error) {
	salt, err = generateRandomBytes(saltLen)

	return
}

func DeriveKeySha256(masterPasswordHash []byte, salt []byte, iter int, keyLen int) (derivedKey []byte) {
	derivedKey = deriveKey(masterPasswordHash, salt, iter, keyLen, sha256.New)

	return
}

func generateRandomBytes(len int) (randBytes []byte, err error) {
	randBytes = make([]byte, len)
	_, err = rand.Read(randBytes)

	return
}

func deriveKey(masterPasswordHash []byte, salt []byte, iter int, keyLen int, hashFunc func() hash.Hash) (derivedKey []byte) {
	derivedKey = pbkdf2.Key(masterPasswordHash, salt, iter, keyLen, hashFunc)

	return
}

func EncryptAES(data []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("Key length must be 32 bytes for AES-256")
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	paddedData := PKCS5Padding(data, aes.BlockSize)

	ciphertext := make([]byte, len(paddedData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedData)

	finalCiphertext := append(iv, ciphertext...)

	return finalCiphertext, nil
}

func DecryptAES(ciphertext []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("Key length must be 32 bytes for AES-256")
	}

	iv := ciphertext[:aes.BlockSize]
	data := ciphertext[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, data)

	plaintext = PKCS5Unpadding(plaintext)

	return plaintext, nil
}

func PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(data, padtext...)
}

func PKCS5Unpadding(data []byte) []byte {
	padding := data[len(data)-1]

	return data[:len(data)-int(padding)]
}
