package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	Data         []string
	Hash         string
	PreviousHash string
	Timestamp    time.Time
	Nonce        int
}

func calculateHash(block Block) string {
	data := ""
	for _, transaction := range block.Data {
		data += transaction
	}
	timestamp := string(rune(block.Timestamp.Unix()))
	headers := block.PreviousHash + data + timestamp + string(rune(block.Nonce))
	hash := sha256.New()
	hash.Write([]byte(headers))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

func createNewBlock(previousBlock Block, transactions []string) Block {
	var newBlock Block
	newBlock.Timestamp = time.Now()
	newBlock.PreviousHash = previousBlock.Hash
	newBlock.Data = transactions
	newBlock.Nonce = 0
	newBlock.Hash = calculateHash(newBlock) //new block hash

	return newBlock
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  string
}

func calculateHashMerkle(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

func newMerkleNode(left, right *MerkleNode, data string) *MerkleNode {
	node := &MerkleNode{}
	if left == nil && right == nil {
		node.Hash = calculateHashMerkle(data)
	} else {
		hash := sha256.New()
		hash.Write([]byte(left.Hash + right.Hash))
		hashed := hash.Sum(nil)
		node.Hash = hex.EncodeToString(hashed)
	}
	node.Left = left
	node.Right = right
	return node
}

func buildMerkleTree(data []string) *MerkleNode {
	var nodes []*MerkleNode
	for _, transaction := range data {
		nodes = append(nodes, newMerkleNode(nil, nil, transaction))
	}
	for len(nodes) > 1 {
		var newLevel []*MerkleNode
		for i := 0; i < len(nodes); i += 2 {
			left := nodes[i]
			var right *MerkleNode
			if i+1 < len(nodes) {
				right = nodes[i+1]
			}
			node := newMerkleNode(left, right, "")
			newLevel = append(newLevel, node)
		}
		nodes = newLevel
	}

	return nodes[0]
}

func printMerkleTree(node *MerkleNode, depth int) {
	if node == nil {
		return
	}
	for i := 0; i < depth; i++ {
		fmt.Print("  ")
	}
	fmt.Println(node.Hash)
	printMerkleTree(node.Left, depth+1)
	printMerkleTree(node.Right, depth+1)
}

func main() {
	//do we need to set our own previous block hash or is there another way?
	previousBlock := Block{
		Data:         []string{"Transaction 1", "Transaction 2"},
		Hash:         "000032adfa39c23...",
		PreviousHash: "00002a7cd92c0fa...",
		Timestamp:    time.Now(),
		Nonce:        12345,
	}

	transactions := []string{"Transaction 3", "Transaction 4"}

	newBlock := createNewBlock(previousBlock, transactions)
	println("New Block Hash:", newBlock.Hash)

	transactions1 := []string{"Transaction 1", "Transaction 2", "Transaction 3", "Transaction 4"}
	root := buildMerkleTree(transactions1)
	fmt.Println("Merkle Tree:")
	printMerkleTree(root, 0)
}
