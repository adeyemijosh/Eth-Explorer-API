package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"eth-explorer-api/internal/models"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

type EthService struct {
	client          *ethclient.Client
	etherscanAPIKey string
}

func NewEthService(nodeURL, etherscanAPIKey string) (*EthService, error) {
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	return &EthService{
		client:          client,
		etherscanAPIKey: etherscanAPIKey,
	}, nil
}

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

func (s *EthService) GetTransaction(txHash string) (*models.Transaction, error) {
	ctx := context.Background()

	hash := common.HexToHash(txHash)
	tx, isPending, err := s.client.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}

	if isPending {
		return s.transactionToModel(tx, "", "", "", "", ""), nil
	}

	receipt, err := s.client.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction receipt: %w", err)
	}

	status := "1"
	if receipt.Status == 0 {
		status = "0"
	}

	chainID, err := s.client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get network ID: %w", err)
	}
	signer := types.LatestSignerForChainID(chainID)
	from, err := types.Sender(signer, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender: %w", err)
	}

	return s.transactionToModel(
		tx,
		receipt.BlockNumber.String(),
		receipt.BlockHash.Hex(),
		strconv.FormatUint(uint64(receipt.TransactionIndex), 10),
		status,
		from.Hex(),
	), nil
}

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

func (s *EthService) GetLatestBlock() (*models.Block, error) {
	return s.GetBlock("latest")
}

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

func (s *EthService) parseBlockNumber(blockNumber string) (*big.Int, error) {
	num, err := strconv.ParseInt(blockNumber, 10, 64)
	if err != nil {

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

func (s *EthService) transactionToModel(tx *types.Transaction, blockNumber, blockHash, txIndex, status, from string) *models.Transaction {
	var to string
	if tx.To() != nil {
		to = tx.To().Hex()
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

// GetTransactionHistory retrieves the transaction history for a given address.
func (s *EthService) GetTransactionHistory(address string) (*models.TransactionHistory, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s", address, s.etherscanAPIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction history from Etherscan: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Status  string               `json:"status"`
		Message string               `json:"message"`
		Result  []models.Transaction `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("etherscan API error: %s", result.Message)
	}

	return &models.TransactionHistory{
		Address:      address,
		Transactions: result.Result,
	}, nil
}

// GetTokenBalance retrieves the balance of a specific ERC-20 token for a given wallet address.
func (s *EthService) GetTokenBalance(userAddress, tokenAddress string) (*models.TokenBalance, error) {
	ctx := context.Background()

	// The address of the user's wallet
	walletAddress := common.HexToAddress(userAddress)

	// The address of the ERC-20 token contract
	contractAddress := common.HexToAddress(tokenAddress)

	// The function signature for `balanceOf(address)` is `0x70a08231`
	methodID := []byte{0x70, 0xa0, 0x82, 0x31}

	// Pad the wallet address to 32 bytes
	paddedAddress := common.LeftPadBytes(walletAddress.Bytes(), 32)

	// The call data is the method ID followed by the padded address
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)

	// Make the call to the contract
	result, err := s.client.CallContract(ctx, ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	// The result is the balance in wei, so we convert it to a big.Int
	balance := new(big.Int)
	balance.SetBytes(result)

	return &models.TokenBalance{
		Address:      userAddress,
		TokenAddress: tokenAddress,
		Balance:      balance.String(),
	}, nil
}

// GetTokenTransfers retrieves the ERC-20 token transfer history for a given address.
// GetContractABI retrieves the ABI for a given smart contract address.
func (s *EthService) GetContractABI(address string) (*models.ContractABI, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=%s", address, s.etherscanAPIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ABI from Etherscan: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("etherscan API error: %s", result.Message)
	}

	return &models.ContractABI{
		Address: address,
		ABI:     result.Result,
	}, nil
}

// GetContractSource retrieves the source code for a given smart contract address.
func (s *EthService) GetContractSource(address string) (*models.ContractSource, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=contract&action=getsourcecode&address=%s&apikey=%s", address, s.etherscanAPIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch source code from Etherscan: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  []struct {
			SourceCode string `json:"SourceCode"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if result.Status != "1" || len(result.Result) == 0 {
		return nil, fmt.Errorf("etherscan API error: %s", result.Message)
	}

	return &models.ContractSource{
		Address:    address,
		SourceCode: result.Result[0].SourceCode,
	}, nil
}

func (s *EthService) GetEventLogs(address string, topics []string) ([]models.EventLog, error) {
	ctx := context.Background()

	contractAddress := common.HexToAddress(address)

	var topicHashes [][]common.Hash
	if len(topics) > 0 {
		topicHashes = make([][]common.Hash, len(topics))
		for i, t := range topics {
			topicHashes[i] = []common.Hash{common.HexToHash(t)}
		}
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    topicHashes,
	}

	logs, err := s.client.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to filter logs: %w", err)
	}

	var eventLogs []models.EventLog
	for _, vLog := range logs {
		var logTopics []string
		for _, t := range vLog.Topics {
			logTopics = append(logTopics, t.Hex())
		}

		eventLogs = append(eventLogs, models.EventLog{
			Address:     vLog.Address.Hex(),
			Topics:      logTopics,
			Data:        fmt.Sprintf("0x%x", vLog.Data),
			BlockNumber: vLog.BlockNumber,
			TxHash:      vLog.TxHash.Hex(),
			TxIndex:     vLog.TxIndex,
			BlockHash:   vLog.BlockHash.Hex(),
			Index:       vLog.Index,
			Removed:     vLog.Removed,
		})
	}

	return eventLogs, nil
}

func (s *EthService) GetTokenTransfers(address string) ([]models.TokenTransfer, error) {
	ctx := context.Background()

	// The address to filter by
	addr := common.HexToAddress(address)

	// The signature of the "Transfer" event
	// Transfer(address,address,uint256)
	transferEventSignature := []byte("Transfer(address,address,uint256)")
	transferEventTopic := crypto.Keccak256Hash(transferEventSignature)

	// We can filter for transfers from or to the address
	paddedAddress := common.LeftPadBytes(addr.Bytes(), 32)
	query := ethereum.FilterQuery{
		Topics: [][]common.Hash{
			{transferEventTopic},
			nil,                                 // from
			{common.BytesToHash(paddedAddress)}, // to
		},
	}

	logs, err := s.client.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to filter logs: %w", err)
	}

	var transfers []models.TokenTransfer
	for _, vLog := range logs {
		transfer := models.TokenTransfer{
			TokenAddress: vLog.Address.Hex(),
			From:         common.HexToAddress(vLog.Topics[1].Hex()).Hex(),
			To:           common.HexToAddress(vLog.Topics[2].Hex()).Hex(),
			Value:        new(big.Int).SetBytes(vLog.Data).String(),
			BlockHash:    vLog.BlockHash.Hex(),
			TxHash:       vLog.TxHash.Hex(),
		}
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}
