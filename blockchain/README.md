# now.ink Blockchain Integration

**Solana + Metaplex for NFT Minting**

---

## Approach

For MVP (v0.1.0 - v0.5.0), we use **Metaplex Token Metadata Program** directly without custom smart contracts.

**Why?**
- Fast to market
- Battle-tested NFT standard
- Broad marketplace support
- Lower complexity

**Future:** Custom Solana programs (v0.6.0+) if needed for:
- Advanced royalty logic
- On-chain governance
- Complex commission distribution

---

## Directory Structure

```
blockchain/
├── scripts/          # Testing & deployment scripts
│   ├── mint-test.ts  # Test NFT minting on devnet
│   ├── airdrop.ts    # SOL airdrop for testing
│   └── setup.ts      # Initialize wallets
├── tests/            # Integration tests
├── wallets/          # Test wallets (gitignored)
└── README.md
```

---

## Setup

### 1. Install Dependencies

```bash
cd /home/quantium/labs/now.ink/blockchain
npm init -y
npm install @solana/web3.js @metaplex-foundation/js @metaplex-foundation/mpl-token-metadata arweave
npm install -D typescript @types/node tsx
```

### 2. Configure Solana CLI

```bash
# Set to devnet
solana config set --url https://api.devnet.solana.com

# Generate platform wallet
solana-keygen new --outfile wallets/platform-wallet.json

# Airdrop SOL
solana airdrop 2 -k wallets/platform-wallet.json
```

### 3. Test Minting

```bash
cd /home/quantium/labs/now.ink/blockchain
npm run mint-test
```

---

## NFT Metadata Standard

See [docs/BLOCKCHAIN.md](../docs/BLOCKCHAIN.md) for full details.

**Key fields:**
- `name`: "now.ink Moment #12345"
- `symbol`: "NOWINK"
- `animation_url`: `ar://video_hash`
- `attributes`: GPS coords, timestamp, creator

---

## Platform Commission

**Option 1: Creator Share (Implemented)**
```typescript
creators: [
  { address: platformWallet, share: 5, verified: true },
  { address: userWallet, share: 95, verified: false }
]
```

**Option 2: Upfront Payment** (Future)
- User sends SOL to platform before minting
- More predictable revenue

---

## Testing

All tests use **devnet**. Never use mainnet wallets in code!

```bash
# Test minting
npm run mint-test

# Test Arweave upload
npm run upload-test

# Full integration
npm test
```

---

## Next Steps

1. ✅ Setup complete (see scripts/)
2. 🔄 Test minting on devnet
3. 🔌 Integrate with backend API
4. 🚀 Deploy to mainnet (v0.5.0+)

---

**For implementation details, see [docs/BLOCKCHAIN.md](../docs/BLOCKCHAIN.md)**
