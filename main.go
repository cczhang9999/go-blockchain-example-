package main

import (
	"hello-go/handlers"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 根据环境设置Gin模式
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin路由器
	r := gin.Default()

	// 添加中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// API路由组
	api := r.Group("/api/v1")
	{
		// 钱包相关接口
		api.POST("/wallet", handlers.CreateWallet)
		api.GET("/wallet/:address", handlers.GetBalance)
		api.POST("/transfer", handlers.Transfer)

		// 交易记录相关接口
		api.GET("/transactions", handlers.GetAllTransactions)
		api.GET("/transactions/history/:address", handlers.GetTransactionHistory)
		api.GET("/transactions/block/:block_id", handlers.GetTransactionsByBlock)

		// 区块链信息接口
		api.GET("/blockchain", handlers.GetBlockchainInfo)

		// 健康检查接口
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Blockchain API is running",
			})
		})
	}

	// 根路径
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Blockchain API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"create_wallet":           "POST /api/v1/wallet",
				"get_balance":             "GET /api/v1/wallet/:address",
				"transfer":                "POST /api/v1/transfer",
				"get_all_transactions":    "GET /api/v1/transactions",
				"get_transaction_history": "GET /api/v1/transactions/history/:address",
				"get_block_transactions":  "GET /api/v1/transactions/block/:block_id",
				"blockchain_info":         "GET /api/v1/blockchain",
				"health_check":            "GET /api/v1/health",
			},
		})
	})

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting blockchain server on port %s", port)
	log.Printf("API documentation available at http://localhost:%s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
