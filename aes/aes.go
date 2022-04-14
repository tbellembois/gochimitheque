// https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes
package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func GenerateAESKey() (key string, err error) {
	k := make([]byte, 16)
	_, err = rand.Read(k)

	key = hex.EncodeToString(k)

	return
}

func Encrypt(stringToEncrypt string, keyString string) (encryptedString string, err error) {
	plaintext := []byte(stringToEncrypt)

	// Since the key is in string, we need to convert decode it to bytes
	var key []byte
	key, err = hex.DecodeString(keyString)

	if err != nil {
		return
	}

	// Create a new Cipher Block from the key
	var block cipher.Block
	block, err = aes.NewCipher(key)

	if err != nil {
		return
	}

	// Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	// https://golang.org/pkg/crypto/cipher/#NewGCM
	var aesGCM cipher.AEAD
	aesGCM, err = cipher.NewGCM(block)

	if err != nil {
		return
	}

	// Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	// Encrypt the data using aesGCM.Seal
	// Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	encryptedString = fmt.Sprintf("%x", ciphertext)

	return
}

func Decrypt(encryptedString string, keyString string) (decryptedString string) {
	var (
		key, enc []byte
		err      error
	)

	if key, err = hex.DecodeString(keyString); err != nil {
		return
	}

	if enc, err = hex.DecodeString(encryptedString); err != nil {
		return
	}

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	// Get the nonce size
	nonceSize := aesGCM.NonceSize()

	// Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return
	}

	decryptedString = string(plaintext)

	return
}
