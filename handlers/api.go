package handlers

import (
	"hello-go/blockchain"
	"hello-go/config"
	"hello-go/database"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	bcInstance *blockchain.Blockchain
	once       sync.Once
)

// getBlockchainInstance 获取区块链实例（单例模式）
func getBlockchainInstance() *blockchain.Blockchain {
	once.Do(func() {
		// 初始化数据库连接
		dbConfig := config.GetDBConfig()
		db, err := database.NewMySQLDB(dbConfig)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}

		// 初始化区块链数据访问层
		blockchainDB := database.NewBlockchainMySQL(db)

		// 创建区块链实例
		bcInstance = blockchain.NewBlockchain(blockchainDB)

		// 检查是否有创世区块，如果没有则创建
		latestBlock, err := blockchainDB.GetLatestBlock()
		if err != nil {
			// 创建创世区块
			genesis, err := bcInstance.CreateGenesisBlock()
			if err != nil {
				log.Fatal("Failed to create genesis block:", err)
			}
			log.Printf("Created genesis block: %s", genesis.Hash)
		} else {
			log.Printf("Latest block: Index=%d, Hash=%s", latestBlock.Index, latestBlock.Hash)
		}
	})
	return bcInstance
}

// CreateWallet 创建钱包
func CreateWallet(c *gin.Context) {
	bc := getBlockchainInstance()

	wallet, err := bc.CreateNewWallet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create wallet: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"address":     wallet.Address,
			"private_key": wallet.PrivateKey,
		},
	})
}

// GetBalance 查询余额
func GetBalance(c *gin.Context) {
	bc := getBlockchainInstance()
	address := c.Param("address")

	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Address parameter is required",
		})
		return
	}

	balance, err := bc.GetBalance(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get balance: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"address": address,
			"balance": balance,
		},
	})
}

// Transfer 转账
func Transfer(c *gin.Context) {
	var transferRequest struct {
		FromAddress string  `json:"from_address" binding:"required"`
		ToAddress   string  `json:"to_address" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&transferRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data: " + err.Error(),
		})
		return
	}

	bc := getBlockchainInstance()

	// 执行转账
	err := bc.Transfer(transferRequest.FromAddress, transferRequest.ToAddress, transferRequest.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Transfer failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transfer completed successfully",
		"data": gin.H{
			"from_address": transferRequest.FromAddress,
			"to_address":   transferRequest.ToAddress,
			"amount":       transferRequest.Amount,
		},
	})
}

// GetBlockchainInfo 获取区块链信息
func GetBlockchainInfo(c *gin.Context) {
	bc := getBlockchainInstance()

	// 验证区块链
	isValid, err := bc.ValidateChain()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to validate chain: " + err.Error(),
		})
		return
	}

	// 获取所有区块 - 通过数据库接口获取
	dbConfig := config.GetDBConfig()
	db, err := database.NewMySQLDB(dbConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to database: " + err.Error(),
		})
		return
	}
	defer db.Close()

	blockchainDB := database.NewBlockchainMySQL(db)
	blocks, err := blockchainDB.GetAllBlocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get blocks: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"is_valid": isValid,
			"blocks":   blocks,
			"count":    len(blocks),
		},
	})
}
