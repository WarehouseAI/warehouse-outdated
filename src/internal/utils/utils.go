package utils

import (
	"encoding/base64"
	"math/rand"
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

func GenerateCode(length int) string {
	charset := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	batch := make([]byte, length)

	for i := range batch {
		batch[i] = charset[rand.Intn(len(charset))]
	}

	return string(batch)
}
