package chain

import (
	"time"
)

type Block struct {
	previousHash []byte
	rootHash     []byte
	timestamp    time.Time
	target       [32]byte
	nonce        []byte
	trie         interface{}
}
