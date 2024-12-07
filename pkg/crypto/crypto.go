package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"hash"

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
