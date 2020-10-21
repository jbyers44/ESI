package chain

import (
	"ESI/pkg/crypto"
	"ESI/pkg/helpers"
	"ESI/pkg/trie"
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

var target = uint32(1 << 31)

// Chain represents the chain of blocks
type Chain struct {
	blocks []Block
}

// NewChain returns a new chain with the block array initialized to nil
func NewChain() *Chain {
	return &Chain{nil}
}

// Insert a node into the trie
func (chain *Chain) Insert(previousHash []byte, trie *trie.MerklePatriciaTrie) []byte {
	rootHash := trie.GetRoot().GetHash()

	timestamp := time.Now().UTC().Unix()

	nonce := crypto.GetNonce(target, rootHash)

	newBlock := Block{previousHash, rootHash, timestamp, target, nonce, trie}
	chain.blocks = append(chain.blocks, newBlock)
	return rootHash
}

// InsertBlock inserts a complete block to the chain (helper method for deserialization)
func (chain *Chain) InsertBlock(block Block) {
	chain.blocks = append(chain.blocks, block)
}

// InsertBatch takes the not split content of multiple strings and inserts all the blocks, keeping track of the previous hash
// and constructing the tree for each block
func (chain *Chain) InsertBatch(contents map[string][]byte) {
	previousHash := []byte{}
	for _, content := range contents {
		mpt := trie.NewMerklePatriciaTrie()
		mpt.InsertBatch(helpers.SplitBytes(content))
		mpt.GenerateHashes()

		previousHash = chain.Insert(previousHash, mpt)
	}
}

//InChain returns a boolean indicating whether or not str is in chain, and proofs of membership of its trie and block
func (chain *Chain) InChain(str string) (bool, [][]byte, [][]byte) {
	var blockHashes [][]byte
	var trieHashes [][]byte
	var inTrieResult bool

	val := []byte(str)
	var i int = 0
	for _, block := range chain.blocks {
		inTrieResult, trieHashes = block.trie.InTrie(val)
		if inTrieResult {
			break
		}
		i++
	}

	for _, block := range chain.blocks[i:] {
		blockHashes = append(blockHashes, block.rootHash)
	}

	if inTrieResult {
		return inTrieResult, trieHashes, blockHashes
	}
	return inTrieResult, nil, nil
}

// Validate validates a blockchain for hash correctness
func (chain *Chain) Validate() bool {
	previousHash := []byte{}
	for _, block := range chain.blocks {

		blockResult := block.Validate()
		if blockResult == false {
			return false
		} else if bytes.Compare(previousHash, block.previousHash) != 0 {
			return false
		}
		previousHash = block.GetRootHash()
	}
	return true
}

// Corrupt chooses a random block in the chain and corrupts it using block.Corrupt
func (chain *Chain) Corrupt() {
	rand.Seed(time.Now().UnixNano())
	block := &chain.blocks[rand.Intn(len(chain.blocks))]
	block.trie.Corrupt()
}

// String returns a string representation of the blockchain
func (chain *Chain) String(printTrie bool) string {
	var b bytes.Buffer
	for i := len(chain.blocks) - 1; i >= 0; i-- {
		fmt.Fprintf(&b, "BEGIN BLOCK\n")
		fmt.Fprintf(&b, chain.blocks[i].String(printTrie))
		fmt.Fprintf(&b, "END BLOCK\n")
		if i > 0 {
			fmt.Fprintf(&b, "\n")
		}
	}
	return b.String()
}
