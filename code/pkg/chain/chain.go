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
	blocks []Block
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
func (chain *Chain) InChain(str string, strInChain bool) (bool, [][]byte, [][]byte) {
	var blockHashes [][]byte
	var trieHashes [][]byte
	var inTrieResult bool

	val := []byte(str)
	var i int = 0
	for _, block := range chain.blocks {
		inTrieResult, trieHashes = block.trie.InTrie(block.trie.GetRoot(), val, 0, trieHashes)
		if inTrieResult == true {
			strInChain = true
			break
		}
		i++
	}

	for index := i; index < len(chain.blocks); index++ {
		blockHashes = append(blockHashes, chain.blocks[i].rootHash)
	}

	if len(blockHashes) == 0 {
		strInChain = false
	} else { //reverse trieHashes
		for j := 0; j < len(trieHashes)/2; j++ {
			k := len(trieHashes) - j - 1
			trieHashes[j], trieHashes[k] = trieHashes[k], trieHashes[j]
		}
	}

	return inTrieResult, trieHashes, blockHashes
}

// Validate validates a blockchain for hash correctness
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
