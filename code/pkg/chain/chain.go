package chain

import (
	"ESI/pkg/crypto"
	"ESI/pkg/helpers"
	"ESI/pkg/trie"
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

var target = uint32(1 << 31)

// Chain represents the chain of blocks
type Chain struct {
	Blocks []Block
}

// NewChain returns a new chain with the block array initialized to nil
func NewChain() *Chain {
	return &Chain{nil}
}

// Insert a node into the trie
func (chain *Chain) Insert(previousHash []byte, trie *trie.MerklePatriciaTrie) []byte {
	h := sha256.New()
	h.Write(trie.GetRoot().GetHash())
	rootHash := h.Sum(nil)

	timestamp := time.Now().UTC().Unix()

	nonce := crypto.GetNonce(target, rootHash)

	newBlock := Block{previousHash, rootHash, timestamp, target, nonce, trie}
	chain.Blocks = append(chain.Blocks, newBlock)
	return rootHash
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

//Membership contains a boolean representing whether or not a string is in a chain and a proof of membership if true
type Membership struct {
	member      bool
	proofTrie   [][]byte
	proofBlocks [][]byte
}

//inchain returns a Membership indicating whether or not str is in chain
func (chain *Chain) inchain(str string, strInChain bool) (bool, [][]byte, [][]byte) {
	result := Membership{
		member:      strInChain,
		proofTrie:   nil,
		proofBlocks: nil,
	}

	currentBlock := chain.Blocks[len(chain.Blocks)-1]
	currentTrie := currentBlock.Trie

	var blockHashes [][]byte
	var trieHashes [][]byte
	var inTrieResult bool

	val := []byte(str)
	for _, block := range chain.Blocks {
		inTrieResult, trieHashes = block.Trie.InTrie(block.Trie.GetRoot(), val, 0, trieHashes)
		if inTrieResult == true {
			break
		}
	}

	return inTrieResult, trieHashes, blockHashes
func (chain *Chain) Validate() bool {
	previousHash := []byte{}
	for _, block := range chain.blocks {

		blockResult := block.Validate()
		println(blockResult)
		if blockResult == false {
			return false
		} else if bytes.Compare(previousHash, block.previousHash) != 0 {
			println(previousHash)
			println(block.previousHash)
			return false
		}
		previousHash = block.GetRootHash()
	}
	return true
}

// String returns a string representation of the blockchain
func (chain *Chain) String(printTrie bool) string {
	var b bytes.Buffer
	for i := len(chain.Blocks) - 1; i >= 0; i-- {
		fmt.Fprintf(&b, "BEGIN BLOCK\n")
		fmt.Fprintf(&b, chain.Blocks[i].String(printTrie))
		fmt.Fprintf(&b, "END BLOCK\n")
		if i > 0 {
			fmt.Fprintf(&b, "\n")
		}
	}
	return b.String()
}
