package main

import (
	"fmt"
	"hello-go/blockchain"
	"hello-go/config"
	"hello-go/database"
	"log"
)

func main() {
	// 初始化数据库连接
	dbConfig := config.GetDBConfig()
	db, err := database.NewMySQLDB(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// 初始化区块链数据访问层
	blockchainDB := database.NewBlockchainMySQL(db)

	// 创建区块链实例
	bc := blockchain.NewBlockchain(blockchainDB)

	// 检查是否有创世区块，如果没有则创建
	latestBlock, err := blockchainDB.GetLatestBlock()
	if err != nil {
		// 创建创世区块
		genesis, err := bc.CreateGenesisBlock()
		if err != nil {
			log.Fatal("Failed to create genesis block:", err)
		}
		fmt.Printf("Created genesis block: %s\n", genesis.Hash)
	} else {
		fmt.Printf("Latest block: Index=%d, Hash=%s\n", latestBlock.Index, latestBlock.Hash)
	}

	// 创建新区块
	newBlock, err := bc.CreateNewBlock("Transaction data 1", 4)
	if err != nil {
		log.Fatal("Failed to create new block:", err)
	}
	fmt.Printf("Created new block: Index=%d, Hash=%s\n", newBlock.Index, newBlock.Hash)

	// 验证区块链
	isValid, err := bc.ValidateChain()
	if err != nil {
		log.Fatal("Failed to validate chain:", err)
	}
	fmt.Printf("Blockchain is valid: %t\n", isValid)

	// 显示所有区块
	blocks, err := blockchainDB.GetAllBlocks()
	if err != nil {
		log.Fatal("Failed to get all blocks:", err)
	}

	fmt.Println("\nAll blocks:")
	for _, block := range blocks {
		fmt.Printf("Index: %d, Hash: %s, Data: %s\n",
			block.Index, block.Hash, block.Data)
	}
}
