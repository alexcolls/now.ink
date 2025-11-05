# Architecture

**now.ink System Architecture**

This document describes the technical architecture of now.ink, a decentralized live-streaming platform with geo-tagged NFT minting on Solana.

---

## Table of Contents

- [High-Level Overview](#high-level-overview)
- [Frontend Architecture](#frontend-architecture)
- [Backend Architecture](#backend-architecture)
- [Blockchain Layer](#blockchain-layer)
- [Storage Layer](#storage-layer)
- [Streaming Architecture](#streaming-architecture)
- [Database Design](#database-design)
- [Authentication & Wallets](#authentication--wallets)
- [Deployment Architecture](#deployment-architecture)

---

## High-Level Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        CLIENT LAYER                              │
├─────────────────────────┬───────────────────────────────────────┤
│   Web (Nuxt 4)          │   Mobile (React Native)               │
│   - Map interface       │   - Live streaming                    │
│   - Time slider         │   - Camera capture                    │
│   - NFT discovery       │   - GPS tracking                      │
└─────────────────────────┴───────────────────────────────────────┘
                              ↕ HTTPS/WebSocket
┌─────────────────────────────────────────────────────────────────┐
│                      BACKEND LAYER (Go)                          │
├─────────────────────────────────────────────────────────────────┤
│   API Gateway (Fiber/Gin)                                       │
│   ├── REST endpoints                                            │
│   ├── WebSocket signaling (WebRTC)                              │
│   └── Authentication middleware                                 │
├─────────────────────────────────────────────────────────────────┤
│   Services                                                       │
│   ├── User Service                                              │
│   ├── Stream Service                                            │
│   ├── NFT Minting Service                                       │
│   ├── Geo Query Service                                         │
│   └── Follow/Social Graph Service                               │
└─────────────────────────────────────────────────────────────────┘
                              ↕
┌─────────────────────────────────────────────────────────────────┐
│                     PERSISTENCE LAYER                            │
├────────────────────────┬────────────────────────────────────────┤
│  Database              │  Cache                                 │
│  (PostgreSQL/MongoDB)  │  (Redis)                               │
│  - Users               │  - Sessions                            │
│  - Streams             │  - Active streams                      │
│  - NFT metadata        │  - Geo-spatial index                   │
└────────────────────────┴────────────────────────────────────────┘
                              ↕
┌─────────────────────────────────────────────────────────────────┐
│                   BLOCKCHAIN & STORAGE LAYER                     │
├────────────────────────┬────────────────────────────────────────┤
│  Solana Network        │  Arweave                               │
│  - Metaplex NFTs       │  - Video storage                       │
│  - On-chain metadata   │  - Permanent, decentralized            │
│  - Wallet integration  │  - Content-addressed                   │
└────────────────────────┴────────────────────────────────────────┘
```

---

## Frontend Architecture

### Web Application (Nuxt 4)

**Tech Stack:**
- Nuxt 4 (Vue 3 + TypeScript)
- Vite for build tooling
- TailwindCSS for styling
- Google Maps API for map interface
- Solana Web3.js + Wallet Adapter

**Key Components:**

```typescript
// app/
├── pages/
│   ├── index.vue              // Landing page
│   ├── map.vue                // Main map interface with time slider
│   ├── stream/[id].vue        // Watch a specific stream/NFT
│   ├── profile/[user].vue     // User profile & their NFTs
│   └── premium.vue            // Premium tier info
├── components/
│   ├── MapView.vue            // Google Maps integration
│   ├── TimeSlider.vue         // Time-range filter
│   ├── StreamPlayer.vue       // Video playback
│   ├── NFTCard.vue            // NFT display
│   └── WalletConnect.vue      // Phantom/Solflare integration
├── composables/
│   ├── useWallet.ts           // Wallet connection logic
│   ├── useGeolocation.ts      // GPS utilities
│   ├── useNFT.ts              // NFT fetching/display
│   └── useStream.ts           // Stream playback
└── stores/
    ├── auth.ts                // User auth state
    ├── map.ts                 // Map state & filters
    └── nfts.ts                // NFT cache
```

**Key Features:**
- Server-side rendering for SEO
- Progressive Web App (PWA) support
- Responsive design (mobile-first)
- Real-time updates via WebSocket

---

### Mobile App (React Native)

**Tech Stack:**
- React Native 0.73+ (TypeScript)
- Expo for rapid development
- React Native Camera for video capture
- Geolocation API
- Solana Mobile Wallet Adapter

**Key Components:**

```typescript
// src/
├── screens/
│   ├── HomeScreen.tsx         // Main camera/stream screen
│   ├── MapScreen.tsx          // Map view
│   ├── ProfileScreen.tsx      // User profile
│   └── SettingsScreen.tsx     // App settings
├── components/
│   ├── CameraView.tsx         // Live camera feed
│   ├── StreamControls.tsx     // Record/stop/save
│   ├── MintDialog.tsx         // NFT minting confirmation
│   └── MapMarkers.tsx         // Pin rendering
├── hooks/
│   ├── useCamera.ts           // Camera management
│   ├── useLocation.ts         // GPS tracking
│   ├── useWebRTC.ts           // Streaming logic
│   └── useSolana.ts           // Blockchain interactions
└── services/
    ├── api.ts                 // Backend API client
    ├── stream.ts              // WebRTC signaling
    └── nft.ts                 // Minting service
```

**Key Features:**
- Native camera access with real-time streaming
- Background location tracking (with permission)
- Push notifications for followers
- Offline mode (queue mints for later)

---

## Backend Architecture

### API Server (Go)

**Framework Choice:** Go Fiber (Express-like, fast) or Gin (lightweight, proven)

**Project Structure:**

```go
// backend/
├── cmd/
│   └── api/
│       └── main.go            // Entry point
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── auth.go        // Auth endpoints
│   │   │   ├── stream.go      // Streaming endpoints
│   │   │   ├── nft.go         // NFT endpoints
│   │   │   ├── geo.go         // Geo queries
│   │   │   └── social.go      // Follow/feed
│   │   ├── middleware/
│   │   │   ├── auth.go        // JWT validation
│   │   │   ├── cors.go        // CORS config
│   │   │   └── rate_limit.go  // Rate limiting
│   │   └── routes.go          // Route definitions
│   ├── services/
│   │   ├── user/              // User management
│   │   ├── stream/            // Stream processing
│   │   ├── nft/               // NFT minting logic
│   │   ├── geo/               // Spatial queries
│   │   └── social/            // Follow graph
│   ├── models/
│   │   ├── user.go
│   │   ├── stream.go
│   │   ├── nft.go
│   │   └── follow.go
│   ├── db/
│   │   ├── postgres.go        // DB connection
│   │   └── migrations/        // Schema migrations
│   └── blockchain/
│       ├── solana.go          // Solana RPC client
│       └── metaplex.go        // Metaplex integration
├── pkg/
│   ├── config/                // Config loading
│   ├── logger/                // Structured logging
│   └── utils/                 // Shared utilities
└── go.mod
```

**Key Services:**

1. **User Service**
   - Wallet-based authentication (sign message to prove ownership)
   - Profile management
   - Premium tier status

2. **Stream Service**
   - WebRTC signaling server
   - Live viewer counts
   - Stream metadata storage

3. **NFT Minting Service**
   - Upload video to Arweave
   - Create Metaplex NFT
   - Emit blockchain transaction
   - Take platform commission (% TBD)

4. **Geo Query Service**
   - Spatial indexing (PostGIS or MongoDB geospatial)
   - Proximity-based NFT discovery
   - Time-range filtering

5. **Social Graph Service**
   - Follow/unfollow logic
   - Chronological feed generation (no algorithm)
   - Public/private access control

---

## Blockchain Layer

### Solana Integration

**Components:**
- Solana RPC client (mainnet-beta, devnet for testing)
- Metaplex SDK for NFT minting
- Wallet adapters (Phantom, Solflare)

**NFT Metadata Structure:**

```json
{
  "name": "now.ink Moment #12345",
  "symbol": "NOWINK",
  "description": "Captured at [location] on [timestamp]",
  "image": "ar://[arweave_thumbnail_hash]",
  "animation_url": "ar://[arweave_video_hash]",
  "external_url": "https://now.ink/nft/12345",
  "attributes": [
    {
      "trait_type": "Latitude",
      "value": "40.7128"
    },
    {
      "trait_type": "Longitude",
      "value": "-74.0060"
    },
    {
      "trait_type": "Timestamp",
      "value": "2025-11-05T01:23:45Z"
    },
    {
      "trait_type": "Creator",
      "value": "user_wallet_address"
    },
    {
      "trait_type": "Duration",
      "value": "42"
    }
  ],
  "properties": {
    "files": [
      {
        "uri": "ar://[arweave_video_hash]",
        "type": "video/mp4"
      }
    ],
    "category": "video"
  }
}
```

**Minting Flow:**

```
User stops recording
    ↓
Video saved locally
    ↓
User confirms save
    ↓
Backend uploads video to Arweave → get TX hash
    ↓
Backend creates metadata JSON → upload to Arweave
    ↓
Backend calls Metaplex to mint NFT
    ↓
Platform takes commission (% TBD of mint fee)
    ↓
NFT minted, stored in DB with metadata
    ↓
User receives NFT in wallet
```

---

## Storage Layer

### Arweave Integration

**Why Arweave?**
- Permanent, decentralized storage
- One-time payment, data lives forever
- Content-addressed (immutable)
- Better for large files (video) than IPFS pinning services

**Upload Process:**

```typescript
// Pseudocode
async function uploadToArweave(videoBuffer: Buffer): Promise<string> {
  const arweave = Arweave.init({...});
  const wallet = await arweave.wallets.generate();
  
  const transaction = await arweave.createTransaction({
    data: videoBuffer
  }, wallet);
  
  transaction.addTag('Content-Type', 'video/mp4');
  transaction.addTag('App-Name', 'now.ink');
  
  await arweave.transactions.sign(transaction, wallet);
  await arweave.transactions.post(transaction);
  
  return transaction.id; // ar://[this_id]
}
```

**Cost Estimation:**
- ~$5 per GB (one-time, permanent)
- Average video: 30 seconds @ 10 Mbps = ~40 MB = $0.20/video
- User pays this in platform commission

---

## Streaming Architecture

### WebRTC Real-Time Streaming

**Flow:**

```
Mobile App (Broadcaster)
    ↓ [Offer SDP]
Backend Signaling Server (Go WebSocket)
    ↓ [Forward Offer]
Web/Mobile Viewers
    ↓ [Answer SDP]
Backend
    ↓ [Forward Answer]
Broadcaster
    ↓ [ICE Candidates exchanged]
P2P Connection Established (or via TURN server)
```

**Components:**
- **Signaling Server:** Go WebSocket server for SDP/ICE exchange
- **TURN Server:** (Optional) For NAT traversal (use coturn or Twilio)
- **Media Server:** (Optional) For scalability (Janus, mediasoup)

**For MVP:**
- Simple WebSocket signaling in Go
- Direct P2P connections (WebRTC peer-to-peer)
- For >10 viewers, consider SFU (Selective Forwarding Unit)

---

## Database Design

### Schema (PostgreSQL Example)

```sql
-- Users
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  wallet_address VARCHAR(44) UNIQUE NOT NULL,
  username VARCHAR(50) UNIQUE,
  bio TEXT,
  avatar_url TEXT,
  is_premium BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Streams (live or past)
CREATE TABLE streams (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  title TEXT,
  is_live BOOLEAN DEFAULT TRUE,
  started_at TIMESTAMP DEFAULT NOW(),
  ended_at TIMESTAMP,
  location GEOGRAPHY(POINT, 4326), -- PostGIS for geo queries
  nft_mint_address VARCHAR(44), -- Solana mint address if saved
  arweave_tx_id VARCHAR(43)
);

-- NFTs (saved streams)
CREATE TABLE nfts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  stream_id UUID REFERENCES streams(id),
  mint_address VARCHAR(44) UNIQUE NOT NULL,
  metadata_uri TEXT NOT NULL, -- Arweave URI
  creator_wallet VARCHAR(44) NOT NULL,
  latitude FLOAT,
  longitude FLOAT,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Follows (social graph)
CREATE TABLE follows (
  follower_id UUID REFERENCES users(id),
  following_id UUID REFERENCES users(id),
  created_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (follower_id, following_id)
);

-- Indexes
CREATE INDEX idx_streams_location ON streams USING GIST(location);
CREATE INDEX idx_nfts_created_at ON nfts(created_at);
CREATE INDEX idx_follows_follower ON follows(follower_id);
```

---

## Authentication & Wallets

### Wallet-Based Auth

**Flow:**

1. User connects wallet (Phantom/Solflare)
2. Frontend requests a nonce from backend
3. User signs nonce with private key
4. Backend verifies signature using public key (wallet address)
5. Backend issues JWT token
6. Subsequent requests use JWT in `Authorization` header

**Go Implementation (pseudocode):**

```go
// Generate nonce
func GenerateNonce() string {
  return uuid.New().String()
}

// Verify signature
func VerifyWalletSignature(walletAddress, message, signature string) bool {
  // Use Solana SDK to verify ed25519 signature
  publicKey := solana.MustPublicKeyFromBase58(walletAddress)
  return publicKey.Verify([]byte(message), signature)
}

// Issue JWT
func IssueJWT(walletAddress string) (string, error) {
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "wallet": walletAddress,
    "exp":    time.Now().Add(24 * time.Hour).Unix(),
  })
  return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
```

---

## Deployment Architecture

### Production Setup

```
┌────────────────────────────────────────┐
│  Cloudflare / CDN                      │
│  - Static assets (web)                 │
│  - DDoS protection                     │
└────────────────────────────────────────┘
                  ↓
┌────────────────────────────────────────┐
│  Load Balancer (NGINX / ALB)           │
└────────────────────────────────────────┘
                  ↓
┌────────────────────────────────────────┐
│  Go API Servers (Docker containers)    │
│  - Auto-scaling                        │
│  - Health checks                       │
└────────────────────────────────────────┘
                  ↓
┌─────────────────┬──────────────────────┐
│  PostgreSQL     │  Redis               │
│  (managed)      │  (caching)           │
└─────────────────┴──────────────────────┘
```

**Recommended Providers:**
- **Compute:** DigitalOcean, Fly.io, or AWS ECS
- **Database:** Managed PostgreSQL (AWS RDS, DigitalOcean, Supabase)
- **Redis:** Upstash or AWS ElastiCache
- **Monitoring:** Grafana, Prometheus, Sentry

---

## Security Considerations

1. **Rate Limiting:** Prevent spam minting (limit free users to X mints/day)
2. **CORS:** Whitelist only production domains
3. **Input Validation:** Sanitize all user inputs (video metadata, usernames)
4. **Wallet Security:** Never store private keys; only verify signatures
5. **Video Content Moderation:** (TBD) Community reports + manual review
6. **HTTPS Only:** Enforce TLS for all connections

---

## Performance Targets

- **API Latency:** < 200ms for geo queries
- **Streaming Latency:** < 2s for live viewers
- **NFT Mint Time:** < 30s from recording stop to on-chain
- **Map Load Time:** < 1s for 1000 pins
- **Database Queries:** < 50ms with proper indexing

---

## Next Steps

1. Set up Go project structure (see [SETUP.md](SETUP.md))
2. Initialize Nuxt 4 web app
3. Bootstrap React Native mobile app
4. Configure Solana devnet integration
5. Test WebRTC signaling locally
6. Deploy MVP to staging

---

**Questions?** See [CONTRIBUTING.md](CONTRIBUTING.md) or open an issue.
