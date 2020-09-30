package chain

type Block struct {
	previousHash []byte
	rootHash     []byte
	timestamp    int64
	target       []byte
	nonce        []byte
	trie         interface{}
}
