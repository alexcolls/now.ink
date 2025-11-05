# 🚀 now.ink Deployment Guide

Complete guide for deploying now.ink to production.

---

## 📋 Prerequisites

### Required
- **VPS/Cloud Server** (min 2GB RAM, 2 vCPUs, 20GB storage)
- **Domain name** with DNS access
- **Solana wallet** with SOL for minting fees
- **Arweave wallet** with AR for storage (or use Bundlr)
- **Docker** & **Docker Compose** installed on server

### Recommended Providers
- **Server**: DigitalOcean, AWS, Hetzner, Linode
- **RPC**: Helius, QuickNode, or Alchemy (Solana)
- **Domain**: Namecheap, Cloudflare

---

## 🔐 Step 1: Generate Secrets

### 1.1 JWT Secret
```bash
openssl rand -base64 32
```
Save this for `.env.production`

### 1.2 Database Password
```bash
openssl rand -base64 24
```
Save this for `.env.production`

### 1.3 Solana Wallet
```bash
# On local machine
cd blockchain
solana-keygen new --outfile wallets/platform-mainnet.json

# View public key
solana-keygen pubkey wallets/platform-mainnet.json

# Fund wallet (needs ~0.5 SOL for operations)
# Transfer SOL from your main wallet
```

### 1.4 Arweave Wallet
```bash
# Generate using Arweave wallet CLI or use existing wallet
# Download keyfile from https://arweave.app
```

---

## 🖥️ Step 2: Server Setup

### 2.1 Initial Server Configuration
```bash
# SSH into server
ssh root@your-server-ip

# Update system
apt update && apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Install Docker Compose
apt install docker-compose-plugin -y

# Create deployment user
adduser nowink
usermod -aG docker nowink
su - nowink
```

### 2.2 Clone Repository
```bash
git clone https://github.com/alexcolls/now.ink.git
cd now.ink
```

### 2.3 Configure Environment
```bash
# Copy production env template
cp .env.production.sample .env.production

# Edit with your values
nano .env.production
```

Fill in:
- `DB_PASSWORD` - Your database password
- `JWT_SECRET` - Your JWT secret
- `PLATFORM_WALLET_PRIVATE_KEY` - Your Solana wallet private key
- `CORS_ALLOWED_ORIGINS` - Your domain URLs

### 2.4 Upload Wallet Files
```bash
# On local machine, copy wallets to server
scp blockchain/wallets/platform-mainnet.json nowink@your-server:/home/nowink/now.ink/
scp arweave-keyfile.json nowink@your-server:/home/nowink/now.ink/
```

---

## 🐳 Step 3: Deploy with Docker

### 3.1 Build and Start Services
```bash
# Load environment
export $(cat .env.production | xargs)

# Build and start
docker-compose up -d --build

# Check status
docker-compose ps

# View logs
docker-compose logs -f api
```

### 3.2 Verify Deployment
```bash
# Check health
curl http://localhost/health

# Test API
curl http://localhost/api/v1/nfts
```

---

## 🌐 Step 4: Domain & SSL

### 4.1 DNS Configuration
Point your domain to server IP:
```
Type: A
Name: api
Value: your-server-ip
TTL: 3600
```

### 4.2 SSL with Let's Encrypt
```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx -y

# Get certificate
sudo certbot --nginx -d api.your-domain.com

# Auto-renewal is configured automatically
```

### 4.3 Update Nginx for HTTPS
Edit `nginx/nginx.conf` and uncomment the HTTPS server block, then:
```bash
docker-compose restart nginx
```

---

## 📱 Step 5: Mobile App Configuration

### 5.1 Update API URL
```bash
cd mobile

# Edit or create .env
echo "API_URL=https://api.your-domain.com/api/v1" > .env
```

### 5.2 Build for Production

#### iOS (TestFlight)
```bash
eas build --platform ios
eas submit --platform ios
```

#### Android (Play Store)
```bash
eas build --platform android
eas submit --platform android
```

---

## 🧪 Step 6: Testing

### 6.1 Backend Tests
```bash
# On server
./test-e2e.sh
```

### 6.2 Mobile Tests
- Install from TestFlight/Play Store
- Test wallet connection
- Record and mint a test NFT
- Verify on map and profile

### 6.3 Verify on Blockchain
```bash
# Check Solana transaction
solscan.io/tx/YOUR_TX_ID

# Check Arweave upload
arweave.net/YOUR_TX_ID
```

---

## 🔧 Step 7: Maintenance

### 7.1 Monitoring
```bash
# View logs
docker-compose logs -f

# Resource usage
docker stats

# Disk space
df -h
```

### 7.2 Backups
```bash
# Backup database
docker exec nowink-postgres pg_dump -U nowink_user nowink > backup-$(date +%Y%m%d).sql

# Backup to S3/Backblaze (recommended)
```

### 7.3 Updates
```bash
# Pull latest code
git pull origin main

# Rebuild and restart
docker-compose up -d --build

# Verify
docker-compose ps
```

---

## 🚨 Troubleshooting

### API not starting
```bash
# Check logs
docker-compose logs api

# Common issues:
# - Database connection failed → Check DB_PASSWORD
# - Port already in use → Check for conflicts
# - Out of memory → Upgrade server
```

### Database issues
```bash
# Restart database
docker-compose restart postgres

# Check migrations
docker exec nowink-postgres psql -U nowink_user -d nowink -c "\\dt"
```

### Minting failures
```bash
# Check wallet balance
solana balance -k platform-mainnet.json

# Check RPC health
curl -X POST YOUR_RPC_URL -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","id":1,"method":"getHealth"}'
```

---

## 📊 Production Checklist

### Before Launch
- [ ] Server provisioned and secured
- [ ] Domain configured with SSL
- [ ] Database backed up regularly
- [ ] Monitoring setup (Grafana/Prometheus)
- [ ] Error tracking (Sentry)
- [ ] Wallet funded with SOL
- [ ] Arweave wallet funded with AR
- [ ] Mobile apps submitted to stores

### Security
- [ ] SSH key authentication only
- [ ] Firewall configured (UFW)
- [ ] Rate limiting enabled
- [ ] CORS properly configured
- [ ] Secrets not in git
- [ ] Database password strong
- [ ] JWT secret secure

### Performance
- [ ] CDN setup (Cloudflare)
- [ ] Redis caching enabled
- [ ] Database indexes verified
- [ ] API response times < 500ms
- [ ] Video upload works reliably

---

## 💰 Cost Estimates

### Monthly Costs (approximate)
- **Server** (DigitalOcean): $12-24/month
- **Domain**: $1-2/month
- **Solana RPC** (if custom): $0-50/month
- **Minting fees**: ~0.001 SOL per mint
- **Arweave storage**: ~$0.01 per MB (one-time)

### Scaling Costs
- 100 mints/day: ~$10/month in SOL
- 1,000 mints/day: ~$100/month in SOL
- Storage: $10 per 1GB permanent (Arweave)

---

## 🔗 Useful Commands

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f [service]

# Restart service
docker-compose restart [service]

# Execute command in container
docker-compose exec api /bin/sh

# Check database
docker-compose exec postgres psql -U nowink_user -d nowink

# Backup database
docker-compose exec postgres pg_dump -U nowink_user nowink > backup.sql

# Restore database
docker-compose exec -T postgres psql -U nowink_user nowink < backup.sql
```

---

## 📞 Support

- **GitHub Issues**: https://github.com/alexcolls/now.ink/issues
- **Docs**: `/home/quantium/labs/now.ink/docs`
- **Status**: Check `/health` endpoint

---

**Deployment Complete! 🎉**

Your now.ink platform is now live and ready to onboard users.

_Build it. Ship it. Watch people finally look up._ ✨
