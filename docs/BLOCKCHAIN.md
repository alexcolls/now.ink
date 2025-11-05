# Blockchain Integration

**now.ink - Solana, Metaplex & Arweave Technical Documentation**

This document covers the blockchain and decentralized storage implementation for now.ink.

---

## Table of Contents

- [Overview](#overview)
- [Solana Integration](#solana-integration)
- [Metaplex NFT Standard](#metaplex-nft-standard)
- [Arweave Storage](#arweave-storage)
- [Minting Flow](#minting-flow)
- [Platform Commission](#platform-commission)
- [Wallet Integration](#wallet-integration)
- [Smart Contracts](#smart-contracts)
- [Testing on Devnet](#testing-on-devnet)

---

## Overview

now.ink uses a three-layer blockchain architecture:

1. **Solana** - Fast, low-cost blockchain for NFT minting
2. **Metaplex** - NFT standard and tooling (Token Metadata Program)
3. **Arweave** - Permanent, decentralized storage for video files and metadata

**Why this stack?**
- **Solana:** Sub-second finality, ~$0.00025/transaction
- **Metaplex:** Industry-standard NFT format, broad wallet/marketplace support
- **Arweave:** One-time payment for permanent storage (~$5/GB)

---

## Solana Integration

### Connection Setup

```typescript
import { Connection, clusterApiUrl } from '@solana/web3.js';

// Mainnet
const connection = new Connection(
  process.env.SOLANA_RPC_URL || clusterApiUrl('mainnet-beta'),
  'confirmed'
);

// Devnet (for testing)
const connectionDevnet = new Connection(
  clusterApiUrl('devnet'),
  'confirmed'
);
```

### Environment Variables

```bash
SOLANA_RPC_URL=https://api.mainnet-beta.solana.com
# Or use a paid RPC for better reliability:
# SOLANA_RPC_URL=https://solana-mainnet.g.alchemy.com/v2/YOUR_KEY

SOLANA_NETWORK=mainnet-beta  # or 'devnet' for testing
PLATFORM_WALLET_PRIVATE_KEY=base58_encoded_private_key
```

### Key Management

**IMPORTANT:** Never commit private keys to version control!

- **Platform Wallet:** Signs mint transactions and collects commissions
- **User Wallets:** Connected via Phantom/Solflare (users sign their own transactions)

```go
// Go backend - load platform wallet
import (
  "github.com/gagliardetto/solana-go"
)

func LoadPlatformWallet() solana.PrivateKey {
  privateKeyStr := os.Getenv("PLATFORM_WALLET_PRIVATE_KEY")
  privateKey, err := solana.PrivateKeyFromBase58(privateKeyStr)
  if err != nil {
    log.Fatal("Invalid platform wallet key")
  }
  return privateKey
}
```

---

## Metaplex NFT Standard

### Token Metadata Program

Metaplex uses the **Token Metadata Program** (address: `metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s`)

### NFT Structure

```typescript
interface MetaplexNFT {
  name: string;           // "now.ink Moment #12345"
  symbol: string;         // "NOWINK"
  uri: string;            // Arweave URI (ar://metadata_hash)
  sellerFeeBasisPoints: number;  // 0 (no resale royalties for MVP)
  creators: Creator[];    // Platform + user
  collection?: Collection;  // Optional collection grouping
  uses?: Uses;            // Optional usage restrictions
}

interface Creator {
  address: string;        // Wallet address
  verified: boolean;      // True for platform, false for user (until they verify)
  share: number;          // Percentage (0-100)
}
```

### Metadata JSON (stored on Arweave)

```json
{
  "name": "now.ink Moment #12345",
  "symbol": "NOWINK",
  "description": "Captured at Times Square, New York on 2025-11-05 at 13:45:30 UTC",
  "image": "ar://thumbnail_hash",
  "animation_url": "ar://video_hash",
  "external_url": "https://now.ink/nft/8xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU",
  "attributes": [
    { "trait_type": "Latitude", "value": "40.758896" },
    { "trait_type": "Longitude", "value": "-73.985130" },
    { "trait_type": "Timestamp", "value": "2025-11-05T13:45:30Z" },
    { "trait_type": "Creator", "value": "7xKXtg..." },
    { "trait_type": "Duration (seconds)", "value": "42" },
    { "trait_type": "Location", "value": "Times Square, New York, USA" },
    { "trait_type": "App", "value": "now.ink" }
  ],
  "properties": {
    "files": [
      {
        "uri": "ar://video_hash",
        "type": "video/mp4"
      },
      {
        "uri": "ar://thumbnail_hash",
        "type": "image/jpeg"
      }
    ],
    "category": "video"
  }
}
```

### Minting with Metaplex (TypeScript SDK)

```typescript
import { Metaplex, keypairIdentity, bundlrStorage } from '@metaplex-foundation/js';
import { Connection, Keypair } from '@solana/web3.js';

const connection = new Connection(clusterApiUrl('mainnet-beta'));
const platformWallet = Keypair.fromSecretKey(/* ... */);

const metaplex = Metaplex.make(connection)
  .use(keypairIdentity(platformWallet))
  .use(bundlrStorage({
    address: 'https://node1.bundlr.network',
    providerUrl: process.env.SOLANA_RPC_URL,
    timeout: 60000,
  }));

// Upload metadata to Arweave via Bundlr
const { uri } = await metaplex.nfts().uploadMetadata({
  name: "now.ink Moment #12345",
  symbol: "NOWINK",
  description: "...",
  image: "ar://thumbnail_hash",
  animation_url: "ar://video_hash",
  attributes: [ /* ... */ ],
  properties: { /* ... */ }
});

// Mint NFT
const { nft } = await metaplex.nfts().create({
  uri: uri,
  name: "now.ink Moment #12345",
  sellerFeeBasisPoints: 0,  // No royalties
  creators: [
    {
      address: platformWallet.publicKey,
      share: 5,  // Platform gets 5% (commission TBD)
      verified: true
    },
    {
      address: userWalletAddress,
      share: 95,  // User gets 95%
      verified: false
    }
  ]
});

console.log('Minted NFT:', nft.address.toString());
```

---

## Arweave Storage

### Why Arweave?

- **Permanent:** Pay once, data stored forever (200+ years guarantee)
- **Decentralized:** Replicated across 1000+ nodes
- **Content-addressed:** Immutable, verifiable (ar://hash)
- **Cost-effective:** ~$5 per GB, one-time

### Upload Process

```typescript
import Arweave from 'arweave';
import fs from 'fs';

const arweave = Arweave.init({
  host: 'arweave.net',
  port: 443,
  protocol: 'https'
});

async function uploadVideo(videoBuffer: Buffer): Promise<string> {
  const wallet = JSON.parse(fs.readFileSync('arweave-wallet.json', 'utf8'));
  
  const transaction = await arweave.createTransaction({
    data: videoBuffer
  }, wallet);
  
  // Add tags for searchability
  transaction.addTag('Content-Type', 'video/mp4');
  transaction.addTag('App-Name', 'now.ink');
  transaction.addTag('App-Version', '0.1.0');
  transaction.addTag('Type', 'video');
  
  await arweave.transactions.sign(transaction, wallet);
  const uploader = await arweave.transactions.getUploader(transaction);
  
  while (!uploader.isComplete) {
    await uploader.uploadChunk();
    console.log(`${uploader.pctComplete}% complete`);
  }
  
  return transaction.id;  // Returns TX hash for ar://
}

// Usage
const videoHash = await uploadVideo(videoBuffer);
const arweaveUrl = `ar://${videoHash}`;
```

### Cost Calculation

```typescript
async function estimateArweaveCost(fileSizeBytes: number): Promise<number> {
  const arweave = Arweave.init({ /* ... */ });
  const priceWinston = await arweave.transactions.getPrice(fileSizeBytes);
  const priceAR = arweave.ar.winstonToAr(priceWinston);
  
  // Convert AR to USD (fetch current price from API)
  const arPriceUSD = await fetchARPrice();  // e.g., $8.50
  return parseFloat(priceAR) * arPriceUSD;
}

// Example: 40 MB video
const cost = await estimateArweaveCost(40 * 1024 * 1024);
console.log(`Upload cost: $${cost.toFixed(4)}`);
// Output: Upload cost: $0.1700
```

### Bundlr Integration (Recommended)

For faster uploads and better UX, use Bundlr (batches uploads to Arweave):

```typescript
import { WebBundlr } from '@bundlr-network/client';

const bundlr = new WebBundlr(
  'https://node1.bundlr.network',
  'solana',
  platformWalletPrivateKey,
  { providerUrl: process.env.SOLANA_RPC_URL }
);

await bundlr.ready();

const tags = [
  { name: 'Content-Type', value: 'video/mp4' },
  { name: 'App-Name', value: 'now.ink' }
];

const transaction = bundlr.createTransaction(videoBuffer, { tags });
await transaction.sign();
await transaction.upload();

const arweaveId = transaction.id;
```

---

## Minting Flow

### Complete End-to-End Process

```
1. User stops recording
   ├── Video saved locally on device
   └── GPS coordinates + timestamp captured

2. User taps "Save as NFT"
   ├── App uploads video to backend
   └── Backend validates (size, format, user auth)

3. Backend uploads to Arweave
   ├── Video uploaded → get video TX hash
   ├── Generate thumbnail → upload → get thumbnail TX hash
   └── Create metadata JSON → upload → get metadata URI

4. Backend mints NFT via Metaplex
   ├── Platform wallet signs transaction
   ├── User wallet set as primary creator (95% share)
   ├── Platform wallet gets 5% commission (TBD)
   └── NFT minted on Solana

5. Backend stores NFT metadata in database
   ├── mint_address, arweave_tx_id, lat/lon, timestamp
   └── Link to user and stream records

6. Push notification sent to user
   └── "Your moment is now on-chain! View at now.ink/nft/{mint_address}"

7. NFT appears on map
   └── Users can discover via geo queries or time slider
```

### Transaction Cost Breakdown

| Component | Cost (approximate) |
|-----------|-------------------|
| Arweave video upload (40 MB) | $0.17 |
| Arweave thumbnail (200 KB) | $0.001 |
| Arweave metadata JSON (5 KB) | $0.0001 |
| Solana NFT mint | $0.001 |
| **Total** | **~$0.17** |

**Commission Model (TBD):**
- Option 1: Fixed fee per mint (e.g., $0.50)
- Option 2: Percentage-based (e.g., 10% of total cost → $0.017)
- Option 3: Subscription-only (free mints for premium users)

---

## Platform Commission

### Implementation Approaches

#### Option 1: Creator Share (On-Chain)

Set platform as co-creator with X% share:

```typescript
creators: [
  { address: platformWallet, share: 5, verified: true },
  { address: userWallet, share: 95, verified: false }
]
```

**Pros:** Enforced on-chain, automatic royalty splits  
**Cons:** Only applies to NFT sales, not minting

#### Option 2: User Pays Upfront (Recommended for MVP)

User pays platform before minting:

```typescript
// Frontend: User approves SOL transfer
const transaction = new Transaction().add(
  SystemProgram.transfer({
    fromPubkey: userWallet,
    toPubkey: platformWallet,
    lamports: COMMISSION_LAMPORTS  // e.g., 0.01 SOL
  })
);

await wallet.sendTransaction(transaction, connection);

// Then backend mints NFT
```

**Pros:** Simple, guaranteed revenue  
**Cons:** Requires two transactions (payment + mint)

#### Option 3: Charge for Arweave Upload

Platform pays Arweave costs, user reimburses via Solana:

```typescript
const arweaveCost = await estimateArweaveCost(videoSize);
const commissionSOL = arweaveCost * 1.2;  // 20% markup

// User sends SOL to platform
// Platform uploads to Arweave
```

---

## Wallet Integration

### Supported Wallets

- **Phantom** (recommended, most popular)
- **Solflare**
- **Backpack** (optional)

### Web (Nuxt 4)

```typescript
import { WalletAdapterNetwork } from '@solana/wallet-adapter-base';
import {
  PhantomWalletAdapter,
  SolflareWalletAdapter
} from '@solana/wallet-adapter-wallets';

const wallets = [
  new PhantomWalletAdapter(),
  new SolflareWalletAdapter({ network: WalletAdapterNetwork.Mainnet })
];

// User connects wallet
await wallet.connect();
const walletAddress = wallet.publicKey.toString();
```

### Mobile (React Native)

```typescript
import { transact } from '@solana-mobile/mobile-wallet-adapter-protocol';

const result = await transact(async (wallet) => {
  const authResult = await wallet.authorize({
    cluster: 'mainnet-beta',
    identity: { name: 'now.ink' }
  });
  
  return authResult.accounts[0].address;
});
```

---

## Smart Contracts

### Do We Need Custom Smart Contracts?

**For MVP: NO**

Metaplex handles NFT minting out-of-the-box. Custom smart contracts only needed for:

- Advanced royalty logic
- On-chain governance
- Complex commission distribution

### Future Considerations (Post-MVP)

If implementing custom logic:

```rust
// Rust (Anchor framework)
use anchor_lang::prelude::*;

#[program]
pub mod nowink {
    use super::*;

    pub fn mint_moment(
        ctx: Context<MintMoment>,
        latitude: f64,
        longitude: f64,
        timestamp: i64,
        arweave_uri: String
    ) -> Result<()> {
        // Custom minting logic
        // Validate geo-coordinates
        // Store on-chain metadata
        Ok(())
    }
}
```

---

## Testing on Devnet

### Setup

```bash
# Install Solana CLI
sh -c "$(curl -sSfL https://release.solana.com/v1.17.0/install)"

# Configure for devnet
solana config set --url https://api.devnet.solana.com

# Create test wallet
solana-keygen new --outfile ~/.config/solana/devnet.json

# Airdrop SOL (for testing)
solana airdrop 2
```

### Test Minting Script

```typescript
// test-mint.ts
import { Connection, clusterApiUrl, Keypair } from '@solana/web3.js';
import { Metaplex, keypairIdentity } from '@metaplex-foundation/js';

async function testMint() {
  const connection = new Connection(clusterApiUrl('devnet'));
  const wallet = Keypair.fromSecretKey(/* devnet wallet */);
  
  const metaplex = Metaplex.make(connection)
    .use(keypairIdentity(wallet));
  
  const { nft } = await metaplex.nfts().create({
    uri: 'ar://test_metadata_hash',
    name: 'Test now.ink Moment',
    sellerFeeBasisPoints: 0
  });
  
  console.log('Test NFT minted:', nft.address.toString());
  console.log('View on Solscan:', `https://solscan.io/token/${nft.address}?cluster=devnet`);
}

testMint();
```

### Devnet Faucets

- **Solana:** `solana airdrop 2` (CLI)
- **Arweave:** Use testnet at `https://arweave.net/tx/test` (no actual storage)

---

## Resources

- **Solana Docs:** https://docs.solana.com
- **Metaplex Docs:** https://docs.metaplex.com
- **Arweave Docs:** https://docs.arweave.org
- **Bundlr Docs:** https://docs.bundlr.network
- **Solana Wallet Adapter:** https://github.com/solana-labs/wallet-adapter

---

**Next:** See [SETUP.md](SETUP.md) for local development environment configuration.
