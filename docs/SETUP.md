# Development Setup

**Complete guide to set up now.ink locally**

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Initial Setup](#initial-setup)
- [Web Frontend (Nuxt 4)](#web-frontend-nuxt-4)
- [Mobile App (React Native)](#mobile-app-react-native)
- [Backend API (Go)](#backend-api-go)
- [Database Setup](#database-setup)
- [Blockchain Setup](#blockchain-setup)
- [Running the Stack](#running-the-stack)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Software

```bash
# Node.js 20+ (for Nuxt & React Native)
node --version  # v20.x.x or higher

# Go 1.21+ (for backend)
go version  # go1.21+ or higher

# PostgreSQL 15+ (database)
psql --version  # 15.x or higher

# Git
git --version

# Solana CLI (for blockchain)
solana --version  # 1.17.0 or higher
```

### Optional Tools

- **Docker** - For containerized development
- **Redis** - For caching (recommended)
- **Postman/Insomnia** - API testing

---

## Initial Setup

### 1. Clone Repository

```bash
git clone https://github.com/alexcolls/now.ink.git
cd now.ink
```

### 2. Install Global Tools

```bash
# Install Solana CLI
sh -c "$(curl -sSfL https://release.solana.com/v1.17.0/install)"

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="$HOME/.local/share/solana/install/active_release/bin:$PATH"

# Verify installation
solana --version
```

---

## Web Frontend (Nuxt 4)

### Setup

```bash
cd /home/quantium/labs/now.ink
mkdir -p web && cd web

# Initialize Nuxt 4 project
npx nuxi init . --package-manager npm

# Install dependencies
npm install

# Install additional packages
npm install @solana/web3.js @solana/wallet-adapter-base \
  @solana/wallet-adapter-wallets @solana/wallet-adapter-vue-ui \
  @googlemaps/js-api-loader
```

### Configuration

Create `web/.env.sample`:

```bash
# API Configuration
NUXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NUXT_PUBLIC_WS_URL=ws://localhost:8080/ws

# Solana Configuration
NUXT_PUBLIC_SOLANA_NETWORK=devnet
NUXT_PUBLIC_SOLANA_RPC_URL=https://api.devnet.solana.com

# Google Maps
NUXT_PUBLIC_GOOGLE_MAPS_API_KEY=your_google_maps_api_key_here

# App Configuration
NUXT_PUBLIC_APP_URL=http://localhost:3000
```

Copy to actual `.env`:

```bash
cp .env.sample .env
# Edit .env with your actual values
```

### Run Development Server

```bash
cd /home/quantium/labs/now.ink/web
npm run dev

# App runs on http://localhost:3000
```

---

## Mobile App (React Native)

### Setup

```bash
cd /home/quantium/labs/now.ink
mkdir -p mobile && cd mobile

# Initialize React Native with Expo
npx create-expo-app . --template blank-typescript

# Install dependencies
npm install

# Install additional packages
npm install @solana/web3.js @solana-mobile/mobile-wallet-adapter-protocol \
  react-native-camera expo-location axios
```

### Configuration

Create `mobile/.env.sample`:

```bash
# API Configuration
API_URL=http://localhost:8080/api/v1
WS_URL=ws://localhost:8080/ws

# Solana Configuration
SOLANA_NETWORK=devnet
SOLANA_RPC_URL=https://api.devnet.solana.com

# App Configuration
APP_NAME=now.ink
```

Copy to actual `.env`:

```bash
cp .env.sample .env
```

### Run on Device/Emulator

```bash
cd /home/quantium/labs/now.ink/mobile

# iOS (requires macOS)
npm run ios

# Android
npm run android

# Expo Go (scan QR code)
npm start
```

---

## Backend API (Go)

### Setup

```bash
cd /home/quantium/labs/now.ink
mkdir -p backend && cd backend

# Initialize Go module
go mod init github.com/alexcolls/now.ink/backend

# Create project structure
mkdir -p cmd/api internal/{api,services,models,db,blockchain} pkg/{config,logger,utils}

# Install dependencies
go get github.com/gofiber/fiber/v2
go get github.com/gofiber/websocket/v2
go get github.com/lib/pq
go get github.com/joho/godotenv
go get github.com/golang-jwt/jwt/v5
go get github.com/gagliardetto/solana-go
```

### Configuration

Create `backend/.env.sample`:

```bash
# Server Configuration
PORT=8080
ENV=development

# Database Configuration (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=nowink_dev
DB_SSLMODE=disable

# Redis (optional, for caching)
REDIS_URL=redis://localhost:6379

# JWT Secret (generate with: openssl rand -base64 32)
JWT_SECRET=your_jwt_secret_here_change_in_production

# Solana Configuration
SOLANA_NETWORK=devnet
SOLANA_RPC_URL=https://api.devnet.solana.com
PLATFORM_WALLET_PRIVATE_KEY=base58_encoded_key_here

# Arweave Configuration
ARWEAVE_WALLET_PATH=/home/quantium/labs/now.ink/backend/arweave-wallet.json

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:19006

# Rate Limiting
RATE_LIMIT_REQUESTS=10
RATE_LIMIT_WINDOW=1s
```

Copy to actual `.env`:

```bash
cp .env.sample .env
# Edit with your actual credentials
```

### Create Main Entry Point

Create `backend/cmd/api/main.go`:

```go
package main

import (
    "log"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/joho/godotenv"
)

func main() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    // Create Fiber app
    app := fiber.New(fiber.Config{
        AppName: "now.ink API v0.1.0",
    })

    // Middleware
    app.Use(logger.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins: os.Getenv("CORS_ALLOWED_ORIGINS"),
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    }))

    // Health check
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "healthy",
            "version": "0.1.0",
        })
    })

    // API routes
    api := app.Group("/api/v1")
    api.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "now.ink API",
            "version": "0.1.0",
        })
    })

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    if err := app.Listen(":" + port); err != nil {
        log.Fatal(err)
    }
}
```

### Run Backend

```bash
cd /home/quantium/labs/now.ink/backend

# Run directly
go run cmd/api/main.go

# Or build and run
go build -o bin/api cmd/api/main.go
./bin/api

# API runs on http://localhost:8080
```

---

## Database Setup

### PostgreSQL with PostGIS

```bash
# Install PostgreSQL (Ubuntu/Debian)
sudo apt update
sudo apt install postgresql postgresql-contrib postgis

# Start PostgreSQL
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Create database and user
sudo -u postgres psql

# In psql:
CREATE DATABASE nowink_dev;
CREATE USER nowink WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE nowink_dev TO nowink;
\c nowink_dev
CREATE EXTENSION postgis;
\q
```

### Run Migrations

Create `backend/internal/db/migrations/001_initial_schema.sql`:

```sql
-- Users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_address VARCHAR(44) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE,
    bio TEXT,
    avatar_url TEXT,
    is_premium BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Streams
CREATE TABLE streams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title TEXT,
    is_live BOOLEAN DEFAULT TRUE,
    is_public BOOLEAN DEFAULT TRUE,
    started_at TIMESTAMP DEFAULT NOW(),
    ended_at TIMESTAMP,
    location GEOGRAPHY(POINT, 4326),
    nft_mint_address VARCHAR(44),
    arweave_tx_id VARCHAR(43),
    viewer_count INT DEFAULT 0
);

-- NFTs
CREATE TABLE nfts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID REFERENCES streams(id) ON DELETE CASCADE,
    mint_address VARCHAR(44) UNIQUE NOT NULL,
    metadata_uri TEXT NOT NULL,
    creator_wallet VARCHAR(44) NOT NULL,
    latitude FLOAT,
    longitude FLOAT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Follows
CREATE TABLE follows (
    follower_id UUID REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id)
);

-- Indexes
CREATE INDEX idx_users_wallet ON users(wallet_address);
CREATE INDEX idx_streams_user ON streams(user_id);
CREATE INDEX idx_streams_location ON streams USING GIST(location);
CREATE INDEX idx_nfts_mint ON nfts(mint_address);
CREATE INDEX idx_nfts_creator ON nfts(creator_wallet);
CREATE INDEX idx_nfts_created ON nfts(created_at);
CREATE INDEX idx_follows_follower ON follows(follower_id);
CREATE INDEX idx_follows_following ON follows(following_id);
```

Apply migrations:

```bash
cd /home/quantium/labs/now.ink/backend
psql -U nowink -d nowink_dev -f internal/db/migrations/001_initial_schema.sql
```

---

## Blockchain Setup

### Solana Devnet

```bash
# Configure Solana CLI for devnet
solana config set --url https://api.devnet.solana.com

# Generate platform wallet (for testing)
solana-keygen new --outfile /home/quantium/labs/now.ink/backend/platform-wallet.json

# Get wallet address
solana address -k /home/quantium/labs/now.ink/backend/platform-wallet.json

# Airdrop SOL for testing
solana airdrop 2 -k /home/quantium/labs/now.ink/backend/platform-wallet.json

# Check balance
solana balance -k /home/quantium/labs/now.ink/backend/platform-wallet.json
```

### Arweave Wallet

```bash
# Generate Arweave wallet (for testing, use testnet)
# In production, fund with real AR tokens

cd /home/quantium/labs/now.ink/backend

# Create a test wallet file (for dev, this is just a placeholder)
echo '{}' > arweave-wallet.json

# NOTE: For production, generate proper wallet at https://arweave.app
```

---

## Running the Stack

### Development Mode (All Services)

Create `docker-compose.yml` in project root:

```yaml
version: '3.8'

services:
  postgres:
    image: postgis/postgis:15-3.3
    environment:
      POSTGRES_DB: nowink_dev
      POSTGRES_USER: nowink
      POSTGRES_PASSWORD: your_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

### Start Everything

```bash
# Terminal 1: Database & Redis
docker-compose up

# Terminal 2: Backend
cd /home/quantium/labs/now.ink/backend
go run cmd/api/main.go

# Terminal 3: Web Frontend
cd /home/quantium/labs/now.ink/web
npm run dev

# Terminal 4: Mobile App
cd /home/quantium/labs/now.ink/mobile
npm start
```

### Access Points

- **Web:** http://localhost:3000
- **API:** http://localhost:8080
- **Mobile:** Expo Go app (scan QR)

---

## Testing

### Backend Tests

```bash
cd /home/quantium/labs/now.ink/backend

# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/services/user
```

### Frontend Tests

```bash
cd /home/quantium/labs/now.ink/web
npm test

cd /home/quantium/labs/now.ink/mobile
npm test
```

---

## Troubleshooting

### Common Issues

**1. Port Already in Use**

```bash
# Find process using port 8080
lsof -ti:8080

# Kill process
kill -9 $(lsof -ti:8080)
```

**2. PostgreSQL Connection Refused**

```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Restart if needed
sudo systemctl restart postgresql

# Check logs
sudo journalctl -u postgresql -n 50
```

**3. Solana RPC Rate Limits**

- Use paid RPC provider (Alchemy, QuickNode)
- Or run local validator:

```bash
solana-test-validator
```

**4. Go Module Issues**

```bash
cd /home/quantium/labs/now.ink/backend
go mod tidy
go mod download
```

**5. Node Modules Issues**

```bash
# Web
cd /home/quantium/labs/now.ink/web
rm -rf node_modules package-lock.json
npm install

# Mobile
cd /home/quantium/labs/now.ink/mobile
rm -rf node_modules package-lock.json
npm install
```

### Debug Mode

**Backend:**

```bash
cd /home/quantium/labs/now.ink/backend
ENV=development go run cmd/api/main.go
```

**Frontend:**

```bash
cd /home/quantium/labs/now.ink/web
npm run dev -- --debug
```

---

## Next Steps

1. ✅ Complete setup checklist above
2. 📖 Read [ARCHITECTURE.md](ARCHITECTURE.md) for system design
3. 🔌 Review [API.md](API.md) for endpoint documentation
4. ⛓️ Study [BLOCKCHAIN.md](BLOCKCHAIN.md) for Solana integration
5. 🤝 Check [CONTRIBUTING.md](CONTRIBUTING.md) before coding

---

**Setup complete! Time to build the anti-feed.** 🚀
