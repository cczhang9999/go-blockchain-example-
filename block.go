package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// 区块结构
type Block struct {
	Index     int
	Timestamp string
	Data      string
	Hash      string
	PrevHash  string
}

// 计算区块哈希
func calculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%s%s", block.Index, block.Timestamp, block.Data, block.PrevHash)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// 创建创世区块
func GenesisBlock() Block {
	return Block{
		Index:     0,
		Timestamp: time.Now().String(),
		Data:      "Genesis Block",
		PrevHash:  "",
	}
}

// 创建新区块
func NewBlock(prevBlock Block, data string) Block {
	block := Block{
		Index:     prevBlock.Index + 1,
		Timestamp: time.Now().String(),
		Data:      data,
		PrevHash:  prevBlock.Hash,
	}
	block.Hash = calculateHash(block)
	return block
}

func main() {
	// 创建区块链
	blockchain := []Block{GenesisBlock()}
	blockchain[0].Hash = calculateHash(blockchain[0])

	// 添加新区块
	blockchain = append(blockchain, NewBlock(blockchain[0], "First Block Data"))
	blockchain = append(blockchain, NewBlock(blockchain[1], "Second Block Data"))

	// 打印区块链
	for _, block := range blockchain {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %s\n", block.Timestamp)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("PrevHash: %s\n", block.PrevHash)
		fmt.Println("---")
	}
}
