package crypto

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
	"time"
)

// GetNonce generates a valid nonce given an int representing the difficulty, and a byte hash
func GetNonce(target uint32, hash []byte) []byte {
	h := sha256.New()

	rand.Seed(time.Now().UnixNano())

	nonce := make([]byte, 32)
	i := 0
	for {
		i++
		rand.Read(nonce)

		b := append(hash, nonce...)
		h.Write(b)
		value := h.Sum(nil)

		intHash := binary.BigEndian.Uint32(value)

		if intHash > target {
			println(i)
			return nonce
		}
		println(intHash)
	}
}
