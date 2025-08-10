package services

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"eth-explorer-api/internal/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

type EthService struct {
	client *ethclient.Client
}

func NewEthService(nodeURL string) (*EthService, error) {
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	return &EthService{
		client: client,
	}, nil
}

// GetBlock fetches block information by block number
func (s *EthService) GetBlock(blockNumber string) (*models.Block, error) {
	ctx := context.Background()

	var blockNum *big.Int
	var err error

	if blockNumber == "latest" {
		blockNum = nil
	} else {
		blockNum, err = s.parseBlockNumber(blockNumber)
		if err != nil {
			return nil, fmt.Errorf("invalid block number: %w", err)
		}
	}

	block, err := s.client.BlockByNumber(ctx, blockNum)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block: %w", err)
	}

	return s.blockToModel(block), nil
}

// GetTransaction fetches transaction information by hash
func (s *EthService) GetTransaction(txHash string) (*models.Transaction, error) {
	ctx := context.Background()

	hash := common.HexToHash(txHash)
	tx, isPending, err := s.client.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}

	if isPending {
		return s.transactionToModel(tx, "", "", "", ""), nil
	}

	// Get transaction receipt for additional info
	receipt, err := s.client.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction receipt: %w", err)
	}

	status := "1"
	if receipt.Status == 0 {
		status = "0"
	}

	return s.transactionToModel(
		tx,
		receipt.BlockNumber.String(),
		receipt.BlockHash.Hex(),
		strconv.FormatUint(uint64(receipt.TransactionIndex), 10),
		status,
	), nil
}

// GetBalance fetches the balance of an Ethereum address
func (s *EthService) GetBalance(address string) (*models.Balance, error) {
	ctx := context.Background()

	addr := common.HexToAddress(address)
	balance, err := s.client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch balance: %w", err)
	}

	balanceEth := s.weiToEther(balance)

	return &models.Balance{
		Address:    address,
		Balance:    balanceEth,
		BalanceWei: balance.String(),
	}, nil
}

// GetLatestBlock fetches the latest block
func (s *EthService) GetLatestBlock() (*models.Block, error) {
	return s.GetBlock("latest")
}

// GetGasPrice fetches the current gas price
func (s *EthService) GetGasPrice() (*models.GasPrice, error) {
	ctx := context.Background()

	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gas price: %w", err)
	}

	gasPriceGwei := s.weiToGwei(gasPrice)

	return &models.GasPrice{
		GasPrice:    gasPriceGwei,
		GasPriceWei: gasPrice.String(),
	}, nil
}

// Helper functions
func (s *EthService) parseBlockNumber(blockNumber string) (*big.Int, error) {
	num, err := strconv.ParseInt(blockNumber, 10, 64)
	if err != nil {
		// Try parsing as hex
		if blockNumber[:2] == "0x" {
			num, err = strconv.ParseInt(blockNumber[2:], 16, 64)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return big.NewInt(num), nil
}

func (s *EthService) blockToModel(block *types.Block) *models.Block {
	transactions := make([]string, len(block.Transactions()))
	for i, tx := range block.Transactions() {
		transactions[i] = tx.Hash().Hex()
	}

	return &models.Block{
		Number:       block.Number().String(),
		Hash:         block.Hash().Hex(),
		ParentHash:   block.ParentHash().Hex(),
		Timestamp:    time.Unix(int64(block.Time()), 0),
		Miner:        block.Coinbase().Hex(),
		GasLimit:     strconv.FormatUint(block.GasLimit(), 10),
		GasUsed:      strconv.FormatUint(block.GasUsed(), 10),
		Difficulty:   block.Difficulty().String(),
		Size:         strconv.FormatUint(block.Size(), 10),
		Transactions: transactions,
	}
}

func (s *EthService) transactionToModel(tx *types.Transaction, blockNumber, blockHash, txIndex, status string) *models.Transaction {
	var to string
	if tx.To() != nil {
		to = tx.To().Hex()
	}

	// Get the sender address
	from := ""
	signer := types.LatestSignerForChainID(tx.ChainId())
	if tx.To() != nil {
		msgSender, err := types.Sender(signer, tx)
		if err == nil {
			from = msgSender.Hex()
		}
	}

	model := &models.Transaction{
		Hash:             tx.Hash().Hex(),
		BlockNumber:      blockNumber,
		BlockHash:        blockHash,
		TransactionIndex: txIndex,
		From:             from,
		To:               to,
		Value:            s.weiToEther(tx.Value()),
		Gas:              strconv.FormatUint(tx.Gas(), 10),
		GasPrice:         s.weiToGwei(tx.GasPrice()),
		Nonce:            strconv.FormatUint(tx.Nonce(), 10),
		Input:            fmt.Sprintf("0x%x", tx.Data()),
	}

	if status != "" {
		model.Status = status
	}

	return model
}

func (s *EthService) weiToEther(wei *big.Int) string {
	ether := new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetInt(big.NewInt(params.Ether)))
	return ether.Text('f', 18)
}

func (s *EthService) weiToGwei(wei *big.Int) string {
	gwei := new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetInt(big.NewInt(params.GWei)))
	return gwei.Text('f', 9)
}
