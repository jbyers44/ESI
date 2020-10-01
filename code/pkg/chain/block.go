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
	trie         interface{}
}

func (block *Block) String(printTrie bool) string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "BEGIN HEADER\n")
	fmt.Fprintf(&b, "[previousHash]      '%b'\n", block.previousHash)
	fmt.Fprintf(&b, "[rootHash]          '%b'\n", block.rootHash)
	fmt.Fprintf(&b, "[timestamp]         '%d'\n", block.timestamp)
	fmt.Fprintf(&b, "[nonce]             '%b'\n", block.nonce)
	fmt.Fprintf(&b, "END HEADER\n")
	fmt.Fprintf(&b, "[mpt]\n")
	if printTrie {
		trie := block.trie.(trie.MerklePatriciaTrie)
		fmt.Fprintf(&b, trie.String())
	}
	return b.String()
}
