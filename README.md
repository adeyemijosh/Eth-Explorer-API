# Ethereum Blockchain Explorer API

A high-performance REST API built with Go and Gin framework for exploring Ethereum blockchain data. This API provides endpoints to fetch transaction details, block information, wallet balances, and current network statistics.

## Features

- **Block Information**: Fetch detailed block data by number.
- **Transaction Details**: Get comprehensive transaction information by hash.
- **Wallet Balances**: Check ETH balance for any address.
- **Latest Block**: Get the most recent block data.
- **Gas Price**: Get the current network gas price
- **Health Check**: Endpoint to check the status of the API.

##  Project Structure

```
eth-explorer-api/
├── cmd/
│   └── main.go          # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go    # Configuration management
│   ├── handlers/
│   │   └── eth.go       # HTTP request handlers
│   ├── services/
│   │   └── eth_service.go # Ethereum blockchain service
│   └── models/
│       └── models.go    # Data models and structures
├── .env                 # Environment variables
├── go.mod               # Go module dependencies
└── README.md           # This file
```

## Setup Instructions

###
