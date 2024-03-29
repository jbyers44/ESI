package chain

import (
	"ESI/pkg/crypto"
	"ESI/pkg/trie"
	"bytes"
	"fmt"
)

// Block is the data structure representing a block in the blockchain
type Block struct {
	previousHash []byte
	rootHash     []byte
	timestamp    int64
	target       uint32
	nonce        []byte
	trie         *trie.MerklePatriciaTrie
}

// Validate validates an individual block fo hash correctness
func (block *Block) Validate() bool {
	if crypto.CheckNonce(block.target, block.rootHash, block.nonce) {
		return block.trie.Validate()
	}
	return false
}

// GetRootHash gets the root hash for the block's MerklePatriciaTrie
func (block *Block) GetRootHash() []byte {
	return block.rootHash
}

func (block *Block) String(printTrie bool) string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "BEGIN HEADER\n")
	fmt.Fprintf(&b, "%x\n", block.previousHash)
	fmt.Fprintf(&b, "%x\n", block.rootHash)
	fmt.Fprintf(&b, "%d\n", block.timestamp)
	fmt.Fprintf(&b, "%d\n", block.target)
	fmt.Fprintf(&b, "%x\n", block.nonce)
	fmt.Fprintf(&b, "END HEADER\n")
	if printTrie {
		fmt.Fprintf(&b, block.trie.String())
	}
	return b.String()
}
