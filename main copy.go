package main

import (
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
	// 创建钱包
	wallet1, _ := bc.CreateNewWallet()
	wallet2, _ := bc.CreateNewWallet()
	// 给 wallet充值
	bc.TopUpWallet(wallet1.Address)
	// 转账
	err = bc.Transfer(wallet1.Address, wallet2.Address, 50)
	//余额查询
	balanc, _ := bc.GetBalance(wallet1.Address)
	log.Panicln(balanc)
}
