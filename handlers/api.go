package handlers

import (
	"hello-go/blockchain"
	"hello-go/config"
	"hello-go/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建钱包
func CreateWallet(c *gin.Context) {
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
	wallet, err := bc.CreateNewWallet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"address": wallet.Address, "private_key": wallet.PrivateKey})
}

// 查询余额
func GetBalance(c *gin.Context) {
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
	address := c.Param("address")
	balanc, _ := bc.GetBalance(address)
	c.JSON(http.StatusOK, gin.H{"address": address, "balance": balanc})
}

// 转账（简化版，不含签名校验）
func Transfer(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"msg": "转账成功"})
}
