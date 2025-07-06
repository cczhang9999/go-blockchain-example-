package handlers

import (
	"hello-go/blockchain"
	"hello-go/config"
	"hello-go/database"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	bcInstance *blockchain.Blockchain
	once       sync.Once
)

// Response 统一响应结构
type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

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

// sendResponse 发送统一格式的响应
func sendResponse(c *gin.Context, success bool, message string, data interface{}, errMsg string) {
	response := Response{
		Success:   success,
		Message:   message,
		Data:      data,
		Error:     errMsg,
		Timestamp: time.Now(),
	}

	if success {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusInternalServerError, response)
	}
}

// CreateWallet 创建钱包
func CreateWallet(c *gin.Context) {
	bc := getBlockchainInstance()

	wallet, err := bc.CreateNewWallet()
	if err != nil {
		sendResponse(c, false, "", nil, "Failed to create wallet: "+err.Error())
		return
	}

	walletData := gin.H{
		"address":     wallet.Address,
		"private_key": wallet.PrivateKey,
		"balance":     wallet.Balance,
	}

	sendResponse(c, true, "Wallet created successfully", walletData, "")
}

// GetBalance 查询余额
func GetBalance(c *gin.Context) {
	bc := getBlockchainInstance()
	address := c.Param("address")

	if address == "" {
		sendResponse(c, false, "", nil, "Address parameter is required")
		return
	}

	balance, err := bc.GetBalance(address)
	if err != nil {
		sendResponse(c, false, "", nil, "Failed to get balance: "+err.Error())
		return
	}

	balanceData := gin.H{
		"address": address,
		"balance": balance,
	}

	sendResponse(c, true, "Balance retrieved successfully", balanceData, "")
}

// Transfer 转账
func Transfer(c *gin.Context) {
	var transferRequest struct {
		FromAddress string  `json:"from_address" binding:"required"`
		ToAddress   string  `json:"to_address" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&transferRequest); err != nil {
		sendResponse(c, false, "", nil, "Invalid request data: "+err.Error())
		return
	}

	bc := getBlockchainInstance()

	// 执行转账
	err := bc.Transfer(transferRequest.FromAddress, transferRequest.ToAddress, transferRequest.Amount)
	if err != nil {
		sendResponse(c, false, "", nil, "Transfer failed: "+err.Error())
		return
	}

	transferData := gin.H{
		"from_address": transferRequest.FromAddress,
		"to_address":   transferRequest.ToAddress,
		"amount":       transferRequest.Amount,
		"timestamp":    time.Now(),
	}

	sendResponse(c, true, "Transfer completed successfully", transferData, "")
}

// GetTransactionHistory 获取交易历史
func GetTransactionHistory(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		sendResponse(c, false, "", nil, "Address parameter is required")
		return
	}

	// 获取数据库连接
	dbConfig := config.GetDBConfig()
	db, err := database.NewMySQLDB(dbConfig)
	if err != nil {
		sendResponse(c, false, "", nil, "Failed to connect to database: "+err.Error())
		return
	}
	defer db.Close()

	// 查询该地址的所有交易记录
	query := `
		SELECT id, from_addr, to_addr, amount, timestamp 
		FROM transactions 
		WHERE from_addr = ? OR to_addr = ? 
		ORDER BY timestamp DESC
	`

	rows, err := db.Query(query, address, address)
	if err != nil {
		sendResponse(c, false, "", nil, "Failed to get transaction history: "+err.Error())
		return
	}
	defer rows.Close()

	var transactions []gin.H
	for rows.Next() {
		var id int64
		var fromAddr, toAddr string
		var amount float64
		var timestamp time.Time

		err := rows.Scan(&id, &fromAddr, &toAddr, &amount, &timestamp)
		if err != nil {
			continue
		}

		transactionType := "sent"
		if toAddr == address {
			transactionType = "received"
		}

		transactions = append(transactions, gin.H{
			"id":               id,
			"from_address":     fromAddr,
			"to_address":       toAddr,
			"amount":           amount,
			"timestamp":        timestamp,
			"transaction_type": transactionType,
		})
	}

	historyData := gin.H{
		"address":      address,
		"transactions": transactions,
		"total_count":  len(transactions),
	}

	sendResponse(c, true, "Transaction history retrieved successfully", historyData, "")
}

// GetTransactionsByBlock 获取指定区块的交易
func GetTransactionsByBlock(c *gin.Context) {
	blockIDStr := c.Param("block_id")
	blockID, err := strconv.ParseInt(blockIDStr, 10, 64)
	if err != nil {
		sendResponse(c, false, "", nil, "Invalid block ID")
		return
	}

	// 获取数据库连接
	dbConfig := config.GetDBConfig()
	db, err := database.NewMySQLDB(dbConfig)
	if err != nil {
		sendResponse(c, false, "", nil, "Failed to connect to database: "+err.Error())
		return
	}
	defer db.Close()

	blockchainDB := database.NewBlockchainMySQL(db)

	transactions, err := blockchainDB.GetTransactionsByBlockID(blockID)
	if err != nil {
		sendResponse(c, false, "", nil, "Failed to get transactions: "+err.Error())
		return
	}

	blockData := gin.H{
		"block_id":     blockID,
		"transactions": transactions,
		"count":        len(transactions),
	}

	sendResponse(c, true, "Block transactions retrieved successfully", blockData, "")
}

// GetAllTransactions 获取所有交易记录
func GetAllTransactions(c *gin.Context) {
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// 获取数据库连接
	dbConfig := config.GetDBConfig()
	db, err := database.NewMySQLDB(dbConfig)
	if err != nil {
		sendResponse(c, false, "", nil, "Failed to connect to database: "+err.Error())
		return
	}
	defer db.Close()

	// 查询总数量
	var total int
	err = db.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&total)
	if err != nil {
		sendResponse(c, false, "", nil, "Failed to get transaction count: "+err.Error())
		return
	}

	// 查询交易记录
	query := `
		SELECT id, from_addr, to_addr, amount, timestamp 
		FROM transactions 
		ORDER BY timestamp DESC 
		LIMIT ? OFFSET ?
	`

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		sendResponse(c, false, "", nil, "Failed to get transactions: "+err.Error())
		return
	}
	defer rows.Close()

	var transactions []gin.H
	for rows.Next() {
		var id int64
		var fromAddr, toAddr string
		var amount float64
		var timestamp time.Time

		err := rows.Scan(&id, &fromAddr, &toAddr, &amount, &timestamp)
		if err != nil {
			continue
		}

		transactions = append(transactions, gin.H{
			"id":           id,
			"from_address": fromAddr,
			"to_address":   toAddr,
			"amount":       amount,
			"timestamp":    timestamp,
		})
	}

	paginationData := gin.H{
		"transactions": transactions,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	}

	sendResponse(c, true, "All transactions retrieved successfully", paginationData, "")
}

// GetBlockchainInfo 获取区块链信息
func GetBlockchainInfo(c *gin.Context) {
	bc := getBlockchainInstance()

	// 验证区块链
	isValid, err := bc.ValidateChain()
	if err != nil {
		sendResponse(c, false, "", nil, "Failed to validate chain: "+err.Error())
		return
	}

	// 获取所有区块 - 通过数据库接口获取
	dbConfig := config.GetDBConfig()
	db, err := database.NewMySQLDB(dbConfig)
	if err != nil {
		sendResponse(c, false, "", nil, "Failed to connect to database: "+err.Error())
		return
	}
	defer db.Close()

	blockchainDB := database.NewBlockchainMySQL(db)
	blocks, err := blockchainDB.GetAllBlocks()
	if err != nil {
		sendResponse(c, false, "", nil, "Failed to get blocks: "+err.Error())
		return
	}

	blockchainData := gin.H{
		"is_valid":     isValid,
		"blocks":       blocks,
		"block_count":  len(blocks),
		"last_updated": time.Now(),
	}

	sendResponse(c, true, "Blockchain information retrieved successfully", blockchainData, "")
}

// 交易记录
