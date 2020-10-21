package crypto

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
	"time"
)

// GetNonce generates a valid nonce given an int representing the difficulty, and a byte hash
func GetNonce(target uint32, hash []byte) []byte {
	rand.Seed(time.Now().UnixNano())

	nonce := make([]byte, 32)
	for {
		h := sha256.New()
		rand.Read(nonce)

		b := append(hash, nonce...)
		h.Write(b)
		value := h.Sum(nil)

		intHash := binary.BigEndian.Uint32(value)

		if intHash > target {
			return nonce
		}
	}
}

// CheckNonce checks to see if a given nonce and hash meet the difficulty target
func CheckNonce(target uint32, hash []byte, nonce []byte) bool {
	h := sha256.New()

	b := append(hash, nonce...)
	h.Write(b)
	value := h.Sum(nil)

	intHash := binary.BigEndian.Uint32(value)

	if intHash > target {
		return true
	}
	return false
}
