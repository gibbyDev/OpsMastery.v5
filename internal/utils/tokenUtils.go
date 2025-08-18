package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomToken() string {
	bytes := make([]byte, 16) // 16 bytes = 128 bits
	_, err := rand.Read(bytes)
	if err != nil {
		panic("Failed to generate random token")
	}
	return hex.EncodeToString(bytes)
}
