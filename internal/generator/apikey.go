package generator

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
)

func NewAPIKey() string {
	randomBytes := make([]byte, 32)
	_, _ = rand.Read(randomBytes)

	key := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(key))

	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(hash[:])
}
