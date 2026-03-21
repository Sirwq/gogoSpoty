package crypto

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandState() string {
	k := make([]byte, 32)
	rand.Read(k)
	return hex.EncodeToString(k)
}
