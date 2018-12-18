package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
)

// GetPasswordHash return password hash for the login,
func GetPasswordHash(login string) ([]byte, error) {

	h := hmac.New(sha256.New, []byte("secret"))
	h.Write([]byte(login))

	return h.Sum(nil), nil
}
