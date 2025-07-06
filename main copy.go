package main

import (
	"hello-go/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.POST("/wallet", handlers.CreateWallet)
	r.GET("/wallet/:address", handlers.GetBalance)
	r.POST("/transfer", handlers.Transfer)
	// 你可以继续添加区块链相关接口
	r.Run(":8080")

	// 给 wallet充值
	//bc.TopUpWallet(wallet1.Address)
	// 转账
	//err = bc.Transfer(wallet1.Address, wallet2.Address, 50)
	//余额查询
	//balanc, _ := bc.GetBalance(wallet1.Address)
	//log.Panicln(balanc)
}
