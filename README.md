# Ethereum Blockchain Explorer API

A high-performance REST API built with Go and Gin framework for exploring Ethereum blockchain data. This API provides endpoints to fetch transaction details, block information, wallet balances, and current network statistics.

## 🚀 Features

- **Block Information**: Fetch detailed block data by number.
- **Transaction Details**: Get comprehensive transaction information by hash.
- **Wallet Balances**: Check ETH balance for any address.
- **Latest Block**: Get the most recent block data.
- **Gas Price**: Get the current network gas price.
- **Health Check**: Endpoint to check the status of the API.

## 📁 Project Structure

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

## 🛠 Setup Instructions

### Prerequisites

- Go 1.21 or higher
- An Ethereum node URL (e.g., from Infura, Alchemy, or a local node)

### 1. Clone the Repository

```bash
git clone <repository-url>
cd eth-explorer-api
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Configure Environment Variables

Create a `.env` file in the root directory with the following content:

```
# Server configuration
PORT=8080

# Ethereum Node URL
ETH_NODE_URL=https://mainnet.infura.io/v3/YOUR_PROJECT_ID
```

Replace `YOUR_PROJECT_ID` with your actual Ethereum node project ID.

### 4. Run the Application

```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`.

## 📚 API Endpoints

The base URL for all endpoints is `http://localhost:8080/api/v1`.

### Get Block Information

`GET /eth/block/:number`

- **`:number`**: The block number (e.g., `18500000`) or `"latest"`.

### Get Transaction Details

`GET /eth/transaction/:hash`

- **`:hash`**: The transaction hash.

### Get Wallet Balance

`GET /eth/balance/:address`

- **`:address`**: The Ethereum wallet address.

### Get Latest Block

`GET /eth/latest-block`

### Get Current Gas Price

`GET /eth/gas-price`

### Health Check

`GET /health`
