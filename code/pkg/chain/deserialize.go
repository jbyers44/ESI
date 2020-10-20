package chain

import (
	"ESI/pkg/trie"
	"bytes"
	"log"
	"strconv"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// DeserializeChain deserializes a valid, properly formatted chain into a Chain object
func DeserializeChain(lines [][]byte) *Chain {
	chain := NewChain()
	for i, j := 0, 0; j < len(lines); j++ {
		if bytes.Compare(lines[j], []byte("BEGIN BLOCK")) == 0 {
			i = j
		}
		if bytes.Compare(lines[j], []byte("END BLOCK")) == 0 {
			// Can skip the first two lines of the block as they will always be "BEGIN BLOCK" and "BEGIN HEADER"
			chain.InsertBlock(loadBlock(lines[i+2 : j]))
		}
	}
	return chain
}

func loadBlock(lines [][]byte) Block {
	// First 5 lines will always be the header components
	previousHash := lines[0]
	rootHash := lines[1]
	timestamp, err := strconv.ParseInt(string(lines[2]), 10, 64)
	check(err)
	target64, err := strconv.ParseUint(string(lines[3]), 10, 32)
	check(err)
	target := uint32(target64)
	nonce := lines[4]
	trie := loadTrie(lines[5:])

	block := Block{previousHash, rootHash, timestamp, target, nonce, trie}
	return block
}

func loadTrie(lines [][]byte) *trie.MerklePatriciaTrie {
	return trie.NewMerklePatriciaTrie()
}
