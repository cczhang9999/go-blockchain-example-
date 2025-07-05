package database

import (
	"database/sql"
	"fmt"
	"hello-go/models"
	"time"
)

type BlockchainMySQL struct {
	db *sql.DB
}

func NewBlockchainMySQL(db *sql.DB) *BlockchainMySQL {
	return &BlockchainMySQL{db: db}
}

// 保存区块
func (b *BlockchainMySQL) SaveBlock(block *models.Block) error {
	query := `INSERT INTO blocks (index_num, hash, prev_hash, data, timestamp, nonce, difficulty) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := b.db.Exec(query, block.Index, block.Hash, block.PrevHash,
		block.Data, block.Timestamp, block.Nonce, block.Difficulty)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	block.ID = id
	return nil
}

// 根据索引获取区块
func (b *BlockchainMySQL) GetBlockByIndex(index int) (*models.Block, error) {
	query := `SELECT id, index_num, hash, prev_hash, data, timestamp, nonce, difficulty 
              FROM blocks WHERE index_num = ?`

	block := &models.Block{}
	err := b.db.QueryRow(query, index).Scan(
		&block.ID, &block.Index, &block.Hash, &block.PrevHash,
		&block.Data, &block.Timestamp, &block.Nonce, &block.Difficulty)

	if err != nil {
		return nil, err
	}

	return block, nil
}

// 获取最新区块
func (b *BlockchainMySQL) GetLatestBlock() (*models.Block, error) {
	query := `SELECT id, index_num, hash, prev_hash, data, timestamp, nonce, difficulty 
              FROM blocks ORDER BY index_num DESC LIMIT 1`

	block := &models.Block{}
	err := b.db.QueryRow(query).Scan(
		&block.ID, &block.Index, &block.Hash, &block.PrevHash,
		&block.Data, &block.Timestamp, &block.Nonce, &block.Difficulty)

	if err != nil {
		return nil, err
	}

	return block, nil
}

// 获取所有区块
func (b *BlockchainMySQL) GetAllBlocks() ([]*models.Block, error) {
	query := `SELECT id, index_num, hash, prev_hash, data, timestamp, nonce, difficulty 
              FROM blocks ORDER BY index_num`
	rows, err := b.db.Query(query)
	fmt.Printf("执行SQL: %s [参数: %v]\n", query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []*models.Block
	for rows.Next() {
		block := &models.Block{}
		err := rows.Scan(
			&block.ID, &block.Index, &block.Hash, &block.PrevHash,
			&block.Data, &block.Timestamp, &block.Nonce, &block.Difficulty)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

// 保存交易
func (b *BlockchainMySQL) SaveTransaction(tx *models.Transaction) error {
	query := `INSERT INTO transactions (block_id, from_addr, to_addr, amount, timestamp) 
              VALUES (?, ?, ?, ?, ?)`

	result, err := b.db.Exec(query, tx.BlockID, tx.FromAddr, tx.ToAddr, tx.Amount, tx.Timestamp)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	tx.ID = id
	return nil
}

// 获取区块的所有交易
func (b *BlockchainMySQL) GetTransactionsByBlockID(blockID int64) ([]*models.Transaction, error) {
	query := `SELECT id, block_id, from_addr, to_addr, amount, timestamp 
              FROM transactions WHERE block_id = ?`

	rows, err := b.db.Query(query, blockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		tx := &models.Transaction{}
		err := rows.Scan(&tx.ID, &tx.BlockID, &tx.FromAddr, &tx.ToAddr, &tx.Amount, &tx.Timestamp)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (b *BlockchainMySQL) SaveWallet(wallet *models.Wallet) error {
	query := `INSERT INTO wallets (address, balance,private_key,creat_time) VALUES (?,?,?,?)`
	result, err := b.db.Exec(query,
		wallet.Address, wallet.Balance, wallet.PrivateKey, time.Now())
	fmt.Printf("执行SQL: %s [参数: %v]\n", query, wallet.Address, wallet.Balance, wallet.PrivateKey)
	if err != nil {
		fmt.Println("插入失败:", err)
	}
	rows, _ := result.RowsAffected()
	fmt.Println("插入行数:", rows)
	return nil
}

// 给 wallet 充值
func (b *BlockchainMySQL) TopUpWallet(address string) error {
	query := `UPDATE wallets SET balance = 1000 WHERE address = ?`
	result, err := b.db.Exec(query, address)
	fmt.Printf("执行SQL: %s [参数: %v]\n", query, address)
	if err != nil {
		fmt.Println("update fi:", err)
	}
	rows, _ := result.RowsAffected()
	fmt.Println("upadte rows:", rows)
	return nil
}

// 转账
func (b *BlockchainMySQL) Transfer(from, to string, amount float64) error {
	// 检查余额
	var fromBalance float64
	err := b.db.QueryRow("SELECT balance FROM wallets WHERE address = ?", from).Scan(&fromBalance)
	if err != nil {
		return err
	}
	if fromBalance < amount {
		return fmt.Errorf("余额不足")
	}
	// 扣减和增加余额
	_, err = b.db.Exec("UPDATE wallets SET balance = balance - ? WHERE address = ?", amount, from)
	if err != nil {
		return err
	}
	_, err = b.db.Exec("UPDATE wallets SET balance = balance + ? WHERE address = ?", amount, to)
	if err != nil {
		return err
	}
	// 记录交易
	_, err = b.db.Exec("INSERT INTO transactions (from_addr, to_addr, amount) VALUES (?, ?, ?)", from, to, amount)
	return err
}

// 查询钱包余额
func (b *BlockchainMySQL) GetBalance(address string) (float64, error) {
	var balance float64
	err := b.db.QueryRow("SELECT balance FROM wallets WHERE address = ?", address).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
