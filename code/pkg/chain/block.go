package chain

import (
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

func (block *Block) String(printTrie bool) string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "BEGIN HEADER\n")
	fmt.Fprintf(&b, "[previousHash]  '%x'\n", block.previousHash)
	fmt.Fprintf(&b, "[rootHash]      '%x'\n", block.rootHash)
	fmt.Fprintf(&b, "[timestamp]     '%d'\n", block.timestamp)
	fmt.Fprintf(&b, "[nonce]         '%x'\n", block.nonce)
	fmt.Fprintf(&b, "END HEADER\n")
	fmt.Fprintf(&b, "[mpt]\n")
	if printTrie {
		fmt.Fprintf(&b, block.trie.String())
	}
	return b.String()
}
