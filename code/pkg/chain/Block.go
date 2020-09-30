package chain

import (
	"time"
)

// Block is the data structure representing a block in the blockchain
type Block struct {
	previousHash []byte
	rootHash     []byte
	timestamp    time.Time
	target       []byte
	nonce        []byte
	trie         interface{}
}
