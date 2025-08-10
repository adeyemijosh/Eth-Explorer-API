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

// GetBlock handles GET /api/v1/eth/block/:number
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

// GetTransaction handles GET /api/v1/eth/transaction/:hash
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

// GetBalance handles GET /api/v1/eth/balance/:address
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

// GetLatestBlock handles GET /api/v1/eth/latest-block
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
