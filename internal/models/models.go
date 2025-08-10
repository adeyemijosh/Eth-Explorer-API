package models

import (
	"time"
)

// Block represents an Ethereum block
type Block struct {
	Number       string    `json:"number"`
	Hash         string    `json:"hash"`
	ParentHash   string    `json:"parent_hash"`
	Timestamp    time.Time `json:"timestamp"`
	Miner        string    `json:"miner"`
	GasLimit     string    `json:"gas_limit"`
	GasUsed      string    `json:"gas_used"`
	Difficulty   string    `json:"difficulty"`
	Size         string    `json:"size"`
	Transactions []string  `json:"transactions"`
}

// Transaction represents an Ethereum transaction
type Transaction struct {
	Hash             string `json:"hash"`
	BlockNumber      string `json:"block_number"`
	BlockHash        string `json:"block_hash"`
	TransactionIndex string `json:"transaction_index"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gas_price"`
	GasUsed          string `json:"gas_used,omitempty"`
	Status           string `json:"status,omitempty"`
	Nonce            string `json:"nonce"`
	Input            string `json:"input"`
}

// Balance represents wallet balance information
type Balance struct {
	Address    string `json:"address"`
	Balance    string `json:"balance"`
	BalanceWei string `json:"balance_wei"`
}

// GasPrice represents current gas price
type GasPrice struct {
	GasPrice    string `json:"gas_price"`
	GasPriceWei string `json:"gas_price_wei"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
