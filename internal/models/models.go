package models

import (
	"time"
)

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

type Transaction struct {
	Hash             string `json:"hash"`
	BlockNumber      string `json:"blockNumber"`
	BlockHash        string `json:"blockHash"`
	TransactionIndex string `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	GasUsed          string `json:"gasUsed,omitempty"`
	Status           string `json:"status,omitempty"`
	Nonce            string `json:"nonce"`
	Input            string `json:"input"`
}

type Balance struct {
	Address    string `json:"address"`
	Balance    string `json:"balance"`
	BalanceWei string `json:"balance_wei"`
}

type GasPrice struct {
	GasPrice    string `json:"gas_price"`
	GasPriceWei string `json:"gas_price_wei"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type TransactionHistory struct {
	Address      string        `json:"address"`
	Transactions []Transaction `json:"transactions"`
}

type TokenBalance struct {
	Address      string `json:"address"`
	TokenAddress string `json:"token_address"`
	Balance      string `json:"balance"`
}

type TokenTransfer struct {
	TokenAddress string `json:"token_address"`
	From         string `json:"from"`
	To           string `json:"to"`
	Value        string `json:"value"`
	BlockHash    string `json:"block_hash"`
	TxHash       string `json:"tx_hash"`
}

type ContractABI struct {
	Address string `json:"address"`
	ABI     string `json:"abi"`
}

type ContractSource struct {
	Address    string `json:"address"`
	SourceCode string `json:"source_code"`
}

type EventLog struct {
	Address     string   `json:"address"`
	Topics      []string `json:"topics"`
	Data        string   `json:"data"`
	BlockNumber uint64   `json:"block_number"`
	TxHash      string   `json:"tx_hash"`
	TxIndex     uint     `json:"tx_index"`
	BlockHash   string   `json:"block_hash"`
	Index       uint     `json:"index"`
	Removed     bool     `json:"removed"`
}
