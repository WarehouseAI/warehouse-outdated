package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateKey(length int) (string, error) {
	randomBytes := make([]byte, length)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	key := base64.URLEncoding.EncodeToString(randomBytes)
	key = key[:length]

	return key, nil
}
