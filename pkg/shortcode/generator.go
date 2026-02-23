// Package shortcode provides functions for generating short codes.
package shortcode

import (
	"crypto/rand"
	"math/big"
)

const (
	codeLength = 7
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func Generate() (string, error) {
	charsetLen := big.NewInt(int64(len(charset)))
	result := make([]byte, codeLength)

	for i := 0; i < codeLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result), nil
}
