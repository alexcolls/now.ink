# 🪙 Solana Production Wallet Setup

Complete guide for setting up Solana wallets for production deployment of now.ink.

---

## 📋 Overview

now.ink requires a Solana wallet to:
- Mint NFTs via Metaplex
- Pay transaction fees
- Receive platform commission (5%)
- Sign transactions on behalf of the platform

---

## 🔐 Step 1: Generate Platform Wallet

### Install Solana CLI
```bash
# Install Solana CLI
sh -c "$(curl -sSfL https://release.solana.com/stable/install)"

# Add to PATH
export PATH="/home/$USER/.local/share/solana/install/active_release/bin:$PATH"

# Verify installation
solana --version
```

### Generate New Wallet
```bash
# Create wallets directory
mkdir -p blockchain/wallets

# Generate new keypair for DEVNET testing
solana-keygen new --outfile blockchain/wallets/platform-devnet.json

# IMPORTANT: Save the seed phrase shown!
# Store it securely (password manager, hardware wallet backup)
```

### View Public Key
```bash
# Get wallet address
solana-keygen pubkey blockchain/wallets/platform-devnet.json

# Example output: 8xY3QJ9K7mN4pR6tS5vW8zA1bC2dE3fG4hH5iI6jJ7kK
```

---

## 💰 Step 2: Fund Wallet (Devnet)

### Configure Devnet
```bash
# Set Solana to devnet
solana config set --url https://api.devnet.solana.com

# Verify configuration
solana config get
```

### Request Airdrop
```bash
# Airdrop 2 SOL (devnet only)
solana airdrop 2 blockchain/wallets/platform-devnet.json

# Check balance
solana balance blockchain/wallets/platform-devnet.json

# Should show: 2 SOL
```

### If Airdrop Fails
Try these alternatives:
```bash
# Use different RPC
solana config set --url https://api.devnet.solana.com

# Or use web faucet
# Visit: https://faucet.solana.com
# Paste your wallet address and request airdrop
```

---

## 🧪 Step 3: Test Minting on Devnet

### Install Dependencies
```bash
cd blockchain
npm install
```

### Test Mint Script
```bash
# Export wallet path
export PLATFORM_WALLET_PATH=wallets/platform-devnet.json

# Test mint (replace with your values)
npx tsx scripts/mint-nft.ts \
  --metadata-uri "ar://test-metadata" \
  --name "Test NFT" \
  --creator-wallet "YOUR_TEST_WALLET_ADDRESS" \
  --network devnet \
  --output json

# Should output mint address if successful
```

### Verify on Solscan
```bash
# Visit Solscan devnet explorer
https://solscan.io/token/MINT_ADDRESS?cluster=devnet

# Check:
# - NFT exists
# - Metadata is correct
# - Creator is correct (95% share)
# - Platform has 5% share
```

---

## 🚀 Step 4: Mainnet Wallet (Production)

### Generate Mainnet Wallet
```bash
# Generate new keypair for MAINNET
solana-keygen new --outfile blockchain/wallets/platform-mainnet.json

# ⚠️ CRITICAL: Save seed phrase in multiple secure locations
# - Password manager (1Password, Bitwarden)
# - Hardware wallet backup
# - Encrypted USB drive in safe
```

### Fund Mainnet Wallet
```bash
# Switch to mainnet
solana config set --url https://api.mainnet-beta.solana.com

# Transfer SOL from your main wallet
solana transfer blockchain/wallets/platform-mainnet.json 0.5 \
  --from YOUR_MAIN_WALLET.json

# Verify balance
solana balance blockchain/wallets/platform-mainnet.json

# Should show: ~0.5 SOL
```

### How Much SOL Needed?

**Per NFT Mint:**
- Metaplex mint: ~0.001 SOL
- Metadata upload: ~0.0001 SOL
- **Total:** ~0.0015 SOL per mint

**Recommended Balance:**
- **100 mints:** 0.2 SOL
- **1,000 mints:** 1.5 SOL
- **Production start:** 0.5-1 SOL (covers ~300-600 mints)

---

## 🔑 Step 5: Secure Wallet Keys

### File Permissions
```bash
# Restrict access to wallet files
chmod 600 blockchain/wallets/*.json

# Verify permissions
ls -la blockchain/wallets/

# Should show: -rw------- (read/write for owner only)
```

### Environment Variables
```bash
# Add to .env.production
PLATFORM_WALLET_PRIVATE_KEY=$(cat blockchain/wallets/platform-mainnet.json | jq -c)

# Never commit wallet files to git!
echo "blockchain/wallets/*.json" >> .gitignore
```

### Backup Strategy
1. **Primary:** Password manager (encrypted)
2. **Secondary:** Hardware wallet backup
3. **Tertiary:** Encrypted USB in physical safe
4. **Never:** Plain text, cloud storage, screenshots

---

## 🔧 Step 6: Backend Configuration

### Update Environment
```bash
# Edit backend/.env
nano backend/.env
```

Add these variables:
```bash
# Solana Configuration
SOLANA_NETWORK=mainnet-beta
SOLANA_RPC_URL=https://api.mainnet-beta.solana.com
USE_REAL_MINTING=true

# Wallet (full JSON content)
PLATFORM_WALLET_PRIVATE_KEY='[1,2,3,...]'  # From wallet JSON
```

### Alternative: Use File Path
```bash
# Instead of inline JSON, use file path
PLATFORM_WALLET_PATH=./blockchain/wallets/platform-mainnet.json
```

---

## 🧪 Step 7: Production Testing

### Test on Mainnet (Small Amount First!)
```bash
# Switch to mainnet
solana config set --url https://api.mainnet-beta.solana.com

# Start backend with production config
cd backend
USE_REAL_MINTING=true go run cmd/api/main.go

# In another terminal, test mint
./test-e2e.sh

# ⚠️ This will mint a real NFT and cost ~0.0015 SOL
```

### Verify Mainnet Mint
```bash
# Check on Solscan (mainnet)
https://solscan.io/token/MINT_ADDRESS

# Verify:
# - NFT exists on mainnet
# - Creator royalties correct
# - Platform commission received
# - Video uploaded to Arweave
```

---

## 💳 Step 8: RPC Provider (Optional but Recommended)

Public Solana RPC has rate limits. For production, use a dedicated provider:

### Recommended Providers

**1. Helius**
- Free tier: 100 req/s
- Paid: Starts at $29/month
- Best for NFTs & Metaplex
```bash
SOLANA_RPC_URL=https://mainnet.helius-rpc.com/?api-key=YOUR_KEY
```

**2. QuickNode**
- Free tier: 1M requests/month
- Paid: Starts at $49/month
- Reliable & fast
```bash
SOLANA_RPC_URL=https://your-endpoint.solana-mainnet.quiknode.pro/YOUR_KEY/
```

**3. Alchemy**
- Free tier: 3M compute units/month
- Paid: Starts at $49/month
```bash
SOLANA_RPC_URL=https://solana-mainnet.g.alchemy.com/v2/YOUR_KEY
```

### Setup Custom RPC
```bash
# Get API key from provider
# Update .env.production
SOLANA_RPC_URL=https://your-custom-rpc-url

# Test connection
solana config set --url YOUR_RPC_URL
solana cluster-version
```

---

## 📊 Step 9: Monitoring

### Check Wallet Balance
```bash
# Script to check balance
cat > scripts/check-balance.sh << 'EOF'
#!/bin/bash
WALLET_PATH="blockchain/wallets/platform-mainnet.json"
BALANCE=$(solana balance $WALLET_PATH 2>&1)
echo "Platform Wallet Balance: $BALANCE"

# Alert if below threshold
THRESHOLD=0.1
CURRENT=$(echo $BALANCE | awk '{print $1}')
if (( $(echo "$CURRENT < $THRESHOLD" | bc -l) )); then
    echo "⚠️  WARNING: Balance below threshold!"
    echo "Please add more SOL to continue minting."
fi
EOF

chmod +x scripts/check-balance.sh
```

### Transaction Monitoring
```bash
# Get recent transactions
solana transaction-history \
  $(solana-keygen pubkey blockchain/wallets/platform-mainnet.json) \
  --limit 10

# Monitor in real-time
watch -n 10 'solana balance blockchain/wallets/platform-mainnet.json'
```

### Set Up Alerts
```bash
# Cron job to check balance daily
crontab -e

# Add line:
0 9 * * * /path/to/now.ink/scripts/check-balance.sh | mail -s "nowink Balance Check" admin@yourcompany.com
```

---

## 🚨 Troubleshooting

### Airdrop Failed (Devnet)
```bash
# Try different RPC
solana config set --url https://api.devnet.solana.com

# Or use web faucet
# https://faucet.solana.com
```

### Transaction Failed
```bash
# Check recent blockhash
solana block-height

# Increase timeout
solana config set --commitment confirmed

# Retry with higher priority fees
# (Modify mint script to add priorityFee)
```

### Wallet Not Found
```bash
# Verify file exists
ls -la blockchain/wallets/

# Check permissions
chmod 600 blockchain/wallets/platform-mainnet.json

# Verify JSON format
cat blockchain/wallets/platform-mainnet.json | jq
```

### Out of SOL
```bash
# Check balance
solana balance blockchain/wallets/platform-mainnet.json

# Transfer more SOL
solana transfer blockchain/wallets/platform-mainnet.json 1.0 \
  --from YOUR_FUNDING_WALLET.json
```

---

## 📋 Production Checklist

### Before Launch
- [ ] Devnet wallet generated and tested
- [ ] Mainnet wallet generated
- [ ] Seed phrase backed up (3 locations)
- [ ] Wallet funded with SOL (0.5+ SOL)
- [ ] Test mint successful on mainnet
- [ ] RPC provider configured (optional)
- [ ] Balance monitoring setup
- [ ] File permissions secured (chmod 600)
- [ ] Wallet NOT in git repository
- [ ] Environment variables configured

### Security
- [ ] Seed phrase encrypted
- [ ] Wallet files have restricted permissions
- [ ] Keys not in version control
- [ ] Production wallet different from dev
- [ ] Backup wallets stored securely
- [ ] Alert system for low balance

---

## 💰 Cost Estimates

### Transaction Fees
- **NFT Mint:** ~0.0015 SOL ($0.20-0.40)
- **Metadata Upload:** ~0.0001 SOL
- **Account Creation:** ~0.002 SOL (one-time)

### Monthly Costs (Minting)
- **100 mints/month:** ~0.15 SOL ($20-40)
- **1,000 mints/month:** ~1.5 SOL ($200-400)
- **10,000 mints/month:** ~15 SOL ($2,000-4,000)

### RPC Provider (Optional)
- **Free tier:** $0 (rate limited)
- **Basic plan:** $29-49/month
- **Pro plan:** $99-299/month

---

## 🔗 Useful Links

- **Solana Docs:** https://docs.solana.com
- **Metaplex Docs:** https://docs.metaplex.com
- **Solscan Explorer:** https://solscan.io
- **Devnet Faucet:** https://faucet.solana.com
- **Helius RPC:** https://helius.dev
- **QuickNode:** https://quicknode.com
- **Alchemy:** https://alchemy.com

---

## 📞 Support

Issues with Solana setup?
- Check troubleshooting section above
- Review Solana CLI docs
- Ask in GitHub Issues (tag: solana)

---

**Wallet Setup Complete!** ✅

Your now.ink platform is ready for mainnet NFT minting.

_Remember: Keep your seed phrase safe. If lost, funds cannot be recovered._
