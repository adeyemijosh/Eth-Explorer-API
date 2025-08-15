package handlers

import (
	"net/http"

	"eth-explorer-api/internal/models"
	"eth-explorer-api/internal/services"

	"github.com/gin-gonic/gin"
)

type EthHandler struct {
	ethService *services.EthService
}

func NewEthHandler(ethService *services.EthService) *EthHandler {
	return &EthHandler{
		ethService: ethService,
	}
}

func (h *EthHandler) GetBlock(c *gin.Context) {
	blockNumber := c.Param("number")

	block, err := h.ethService.GetBlock(blockNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to fetch block",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, block)
}

func (h *EthHandler) GetTransaction(c *gin.Context) {
	txHash := c.Param("hash")

	transaction, err := h.ethService.GetTransaction(txHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to fetch transaction",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *EthHandler) GetBalance(c *gin.Context) {
	address := c.Param("address")

	balance, err := h.ethService.GetBalance(address)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to fetch balance",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, balance)
}

func (h *EthHandler) GetLatestBlock(c *gin.Context) {
	block, err := h.ethService.GetLatestBlock()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to fetch latest block",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, block)
}

// GetGasPrice handles GET /api/v1/eth/gas-price
func (h *EthHandler) GetGasPrice(c *gin.Context) {
	gasPrice, err := h.ethService.GetGasPrice()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to fetch gas price",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gasPrice)
}

func (h *EthHandler) GetTransactionHistory(c *gin.Context) {
	address := c.Param("address")

	history, err := h.ethService.GetTransactionHistory(address)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to fetch transaction history",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h *EthHandler) GetTokenBalance(c *gin.Context) {
	userAddress := c.Param("address")
	tokenAddress := c.Param("tokenAddress")

	balance, err := h.ethService.GetTokenBalance(userAddress, tokenAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to fetch token balance",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, balance)
}

func (h *EthHandler) GetTokenTransfers(c *gin.Context) {
	address := c.Param("address")

	transfers, err := h.ethService.GetTokenTransfers(address)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to fetch token transfers",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transfers)
}

func (h *EthHandler) GetContractABI(c *gin.Context) {
	address := c.Param("address")

	abi, err := h.ethService.GetContractABI(address)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to fetch contract ABI",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, abi)
}

func (h *EthHandler) GetContractSource(c *gin.Context) {
	address := c.Param("address")

	source, err := h.ethService.GetContractSource(address)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to fetch contract source",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, source)
}

func (h *EthHandler) GetEventLogs(c *gin.Context) {
	address := c.Param("address")
	topics := c.QueryArray("topics")

	logs, err := h.ethService.GetEventLogs(address, topics)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to fetch event logs",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, logs)
}
