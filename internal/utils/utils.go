package utils

import (
	"crypto/rand"
	"io"
)

func GenerateVerificationCode() (string, error) {
	table := []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	buffer := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, buffer); err != nil {
		return "", err
	}

	for i := 0; i < len(buffer); i++ {
		buffer[i] = table[int(buffer[i])%len(table)]
	}
	return string(buffer), nil
}
