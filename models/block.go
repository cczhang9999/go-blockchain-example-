package models

import (
	"time"
)

type Block struct {
	ID         int64     `json:"id"`
	Index      int       `json:"index"`
	Hash       string    `json:"hash"`
	PrevHash   string    `json:"prev_hash"`
	Data       string    `json:"data"`
	Timestamp  time.Time `json:"timestamp"`
	Nonce      int       `json:"nonce"`
	Difficulty int       `json:"difficulty"`
}

type Transaction struct {
	ID        int64     `json:"id"`
	BlockID   int64     `json:"block_id"`
	FromAddr  string    `json:"from_addr"`
	ToAddr    string    `json:"to_addr"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}
