package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hello-go/models"
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

type Blockchain struct {
	db Database
}

type Database interface {
	SaveBlock(block *models.Block) error
	GetBlockByIndex(index int) (*models.Block, error)
	GetLatestBlock() (*models.Block, error)
	GetAllBlocks() ([]*models.Block, error)
	SaveTransaction(tx *models.Transaction) error
	GetTransactionsByBlockID(blockID int64) ([]*models.Transaction, error)
	SaveWallet(*models.Wallet) error
	TopUpWallet(address string) error
	Transfer(addressFrom string, addressTo string, balance float64) error
	GetBalance(address string) (float64, error)
}

func NewBlockchain(db Database) *Blockchain {
	return &Blockchain{db: db}
}

// 计算区块哈希
func calculateHash(block *models.Block) string {
	record := fmt.Sprintf("%d%s%s%s%d%d",
		block.Index, block.Timestamp, block.Data,
		block.PrevHash, block.Nonce, block.Difficulty)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// 创建创世区块
func (bc *Blockchain) CreateGenesisBlock() (*models.Block, error) {
	genesis := &models.Block{
		Index:      0,
		Hash:       "",
		PrevHash:   "",
		Data:       "Genesis Block",
		Timestamp:  time.Now(),
		Nonce:      0,
		Difficulty: 4,
	}
	genesis.Hash = calculateHash(genesis)

	if err := bc.db.SaveBlock(genesis); err != nil {
		return nil, err
	}

	return genesis, nil
}

// 创建新区块
func (bc *Blockchain) CreateNewBlock(data string, difficulty int) (*models.Block, error) {
	prevBlock, err := bc.db.GetLatestBlock()
	if err != nil {
		return nil, err
	}

	block := &models.Block{
		Index:      prevBlock.Index + 1,
		PrevHash:   prevBlock.Hash,
		Data:       data,
		Timestamp:  time.Now(),
		Nonce:      0,
		Difficulty: difficulty,
	}

	// 挖矿
	block = bc.mineBlock(block)

	if err := bc.db.SaveBlock(block); err != nil {
		return nil, err
	}

	return block, nil
}

// 挖矿
func (bc *Blockchain) mineBlock(block *models.Block) *models.Block {
	target := strings.Repeat("0", block.Difficulty)

	for {
		block.Hash = calculateHash(block)
		if strings.HasPrefix(block.Hash, target) {
			break
		}
		block.Nonce++
	}

	return block
}

// 验证区块链
func (bc *Blockchain) ValidateChain() (bool, error) {
	blocks, err := bc.db.GetAllBlocks()
	if err != nil {
		return false, err
	}

	for i := 1; i < len(blocks); i++ {
		currentBlock := blocks[i]
		prevBlock := blocks[i-1]

		// 验证哈希
		if currentBlock.PrevHash != prevBlock.Hash {
			return false, nil
		}

		// 验证当前区块哈希
		if currentBlock.Hash != calculateHash(currentBlock) {
			return false, nil
		}
	}

	return true, nil
}

func (bc *Blockchain) CreateNewWallet() (*models.Wallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	wallet := &models.Wallet{
		Address:    address,
		PrivateKey: fmt.Sprintf("%x", privateKeyBytes),
		Balance:    0,
	}

	if err := bc.db.SaveWallet(wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (bc *Blockchain) TopUpWallet(address string) error {
	err := bc.db.TopUpWallet(address)
	if err != nil {
		fmt.Printf("TopUpWallet error: %v\n", err)
		return err
	}
	return nil
}

func (bc *Blockchain) Transfer(addressFrom string, addressTo string, balance float64) error {
	err := bc.db.Transfer(addressFrom, addressTo, balance)
	if err != nil {
		log.Println("转账失败:", err)
	} else {
		log.Println("转账成功")
	}
	return nil
}

func (bc *Blockchain) GetBalance(address string) (float64, error) {
	balance, err := bc.db.GetBalance(address)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
