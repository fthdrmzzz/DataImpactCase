package utils

import (
	"crypto/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateToken generates a random token of the specified length
func GenerateToken(length int) string {
	token := make([]byte, length)
	randomBytes := make([]byte, length)

	for i := 0; i < length; i++ {
		if _, err := rand.Read(randomBytes); err != nil {
			panic(err)
		}
		token[i] = charset[int(randomBytes[i])%len(charset)]
	}

	return string(token)
}

// GenerateTokenWithDefaultLength generates a random token of default length 32
func GenerateTokenWithDefaultLength() string {
	return GenerateToken(32)
}
