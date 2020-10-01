package chain

import (
	"ESI/pkg/crypto"
	"ESI/pkg/helpers"
	"ESI/pkg/trie"
	"crypto/sha256"
	"time"
)

var target = uint32(1 << 31)

type Chain struct {
	blocks []Block
}

func NewChain() *Chain {
	return &Chain{nil}
}

// Insert a node into the trie
func (chain *Chain) Insert(previousHash []byte, trie *trie.MerklePatriciaTrie) []byte {
	if len(chain.blocks) == 0 {
		h := sha256.New()
		h.Write(trie.GetRoot().GetHash())
		hashBytes := h.Sum(nil)
		timestamp := time.Now().UTC().Unix()
		chain.blocks = append(chain.blocks, Block{previousHash, hashBytes, timestamp, target, []byte{}, trie})
		return hashBytes
	} else {
		return chain.createBlock(previousHash, trie)
	}
}

func (chain *Chain) createBlock(previousHash []byte, trie *trie.MerklePatriciaTrie) []byte {
	h := sha256.New()
	h.Write(trie.GetRoot().GetHash())
	hashBytes := h.Sum(nil)
	timestamp := time.Now().UTC().Unix()

	nonce := crypto.GetNonce(target, hashBytes)
	newBlock := Block{
		previousHash,
		hashBytes,
		timestamp,
		target,
		nonce,
		trie,
	}
	chain.blocks = append(chain.blocks, newBlock)
	return hashBytes
}

// InsertBatch takes the not split content of multiple strings and inserts all the blocks, keeping track of the previous hash
// and constructing the tree for each block
func (chain *Chain) InsertBatch(rawContent [][]byte) {

	previousHash := make([]byte, 0)
	for _, values := range rawContent {
		mpt := trie.NewMerklePatriciaTrie()
		mpt.InsertBatch(helpers.SplitBytes(values))
		mpt.GenerateHashes()

		previousHash = chain.Insert(previousHash, mpt)
	}
}

func (chain *Chain) GetChain() []Block {
	return chain.blocks
}
