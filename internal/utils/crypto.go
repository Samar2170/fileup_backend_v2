package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fileupbackendv2/config"
)

func GenerateKey(n int) (string, error) {
	key := make([]byte, n)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	apiKey := base64.RawURLEncoding.EncodeToString(key)
	return apiKey, nil
}
func HashKey(apiKey string) string {
	combined := append([]byte(apiKey), []byte(config.SecretKey)...)
	hash := sha256.New()
	hash.Write(combined)
	hashedBytes := hash.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}
