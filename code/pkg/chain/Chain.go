package chain

import (
	"ESI/pkg/trie"
	"crypto/sha256"
	"time"
)

var target int = 1 << 31

type Chain struct {
	genesis *Block
}

func NewChain() *Chain {
	return &Chain{nil}
}

// Insert a node into the trie
func (chain *Chain) Insert(trie *trie.MerklePatriciaTrie) {
	if chain.genesis == nil {
		h := sha256.New()
		h.Write(trie.GetRoot().GetHash())
		timestamp := time.Now().UTC().Unix()
		chain.genesis = &Block{make([]byte, 0), h.Sum(nil), timestamp, make([]byte, target), []byte{}, trie}
	}
	// else {
	// 	insert(trie.root, value, 0)
	// }
}

// func insert(genesis *Block)
