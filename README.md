# now.ink

**Your life, minted.**

![License](https://img.shields.io/badge/license-Dual%20License-blue)
![Status](https://img.shields.io/badge/status-MVP%20Complete-green)
![Solana](https://img.shields.io/badge/blockchain-Solana-blueviolet)
![Version](https://img.shields.io/badge/version-0.2.0-blue)

---

## 🌍 The Vision

**No filters. No fakes. No algorithms. Just raw moments streaming live from your phone, stamped forever on Solana's blockchain.**

Press record, wander real streets, and drop a pin of what's happening. Friends watch in real time; when you stop, save it or let it vanish. If you save? Boom—an NFT drops right there, pinned to that corner, that second, that heartbeat.

Open the map tomorrow—or ten years from now—and slide back through time. Last week? Last summer? The bar you kissed her in 'twenty-two? All glowing dots, waiting. Stand close, play them free; pay once, see everything on Earth.

**No uploads. No AI. No bullshit.**

If it wasn't caught bleeding, sweating, laughing in front of this lens, it doesn't exist. We're the anti-feed. The quiet rebellion against scroll-rot.

**now.ink isn't social media—it's archaeology you lived through.**

---

## ✨ Core Features

### 🎥 Video Recording & Upload
- ✅ Record video directly from your phone with GPS tagging
- ✅ Real-time location capture
- ✅ Secure multipart file upload (max 100MB)
- ✅ No uploads—everything captured in-app for authenticity

### 🪙 NFT Minting (Metaplex)
- ✅ Auto-mint as Solana NFT via Metaplex SDK
- ✅ Embedded metadata: GPS coordinates, timestamp, creator wallet
- ✅ Permanent storage on Arweave
- ✅ Platform takes 5% commission, creator gets 95%
- ✅ Mock + production minting modes

### 🗺️ Interactive Map
- ✅ Google Maps with NFT markers
- ✅ Custom camera icon pins
- ✅ Tap markers to play videos
- ✅ Auto-center on NFT locations
- ✅ Stats bar with NFT count
- 🔜 Time-range slider filtering
- 🔜 Radius-based discovery

### 👤 User Profiles
- ✅ Wallet-based authentication (JWT)
- ✅ Personal NFT grid (2-column layout)
- ✅ Stats counter (Moments/Followers/Following)
- ✅ Filtered by connected wallet
- 🔜 Follow/unfollow functionality
- 🔜 Chronological feed

### 🔐 Privacy & Control
- ✅ Wallet-based authentication
- ✅ No AI filters or edits—raw footage only
- ✅ Users own their NFTs—on Solana blockchain
- ✅ Public by default (private coming soon)

---

## 🛠️ Tech Stack

| Layer | Technology |
|-------|-----------|
| **Mobile App** | React Native + Expo (TypeScript) ✅ |
| **Backend API** | Go + Fiber framework ✅ |
| **Blockchain** | Solana + Metaplex NFT Standard ✅ |
| **Storage** | Arweave (permanent, decentralized) ✅ |
| **Database** | PostgreSQL 16 + PostGIS ✅ |
| **Wallets** | Solana Mobile Wallet Adapter ✅ |
| **Maps** | React Native Maps (Google) ✅ |
| **Video** | Expo AV + Expo Camera ✅ |
| **Web Frontend** | Nuxt 4 (coming soon) 🔜 |

---

## 📦 Quick Start

### Prerequisites
- Node.js 20+ (for Nuxt 4 & React Native)
- Go 1.21+ (for backend)
- Solana CLI (for blockchain integration)
- Docker (optional, for local dev)

### Clone & Install
```bash
git clone https://github.com/alexcolls/now.ink.git
cd now.ink

# Backend API
cd backend
go mod download
cp .env.sample .env
go run cmd/api/main.go

# Mobile app (in separate terminal)
cd mobile
npm install
npm start
```

See **[MVP-COMPLETE.md](MVP-COMPLETE.md)** for testing guide and **[DEPLOYMENT.md](DEPLOYMENT.md)** for production deployment.

---

## 📁 Repository Structure

```
now.ink/
├── mobile/             # React Native + Expo app ✅
├── backend/            # Go + Fiber API server ✅
├── blockchain/         # Solana/Metaplex scripts ✅
├── nginx/              # Reverse proxy config ✅
├── scripts/            # Deployment automation ✅
├── docs/               # Documentation
├── DEPLOYMENT.md       # Production deployment guide ✅
├── MOBILE-COMPLETE.md  # Mobile feature completion ✅
├── MVP-COMPLETE.md     # MVP testing guide ✅
├── FINAL-STATUS.md     # Project status summary ✅
├── CHANGELOG.md        # Version history ✅
├── LICENSE             # Dual license ✅
├── CONTRIBUTING.md     # Contribution guidelines ✅
└── README.md           # This file
```

---

## 📜 License

**Dual License:**
- **Free for personal use** - Individual users can use, modify, and share freely
- **Commercial license required** for businesses generating revenue

See **[LICENSE](LICENSE)** for full details.

---

## 🤝 Contributing

We're open-source by design. Fork it, fix it, ship it. See **[CONTRIBUTING.md](CONTRIBUTING.md)** for guidelines.

**Why open-source?**  
Because if we don't open-source the truth, who will? No black-box algorithms, no hidden feeds. Every line of code is auditable. Every decision is yours.

**Pull requests welcome for:**
- Bug fixes
- Performance improvements
- Documentation updates
- New features (discuss in issues first)

---

## 🚀 Roadmap

### ✅ Phase 1: MVP (Complete)
- Video recording + GPS tagging
- NFT minting on Solana
- Interactive map with markers
- User profiles with NFT grid
- Arweave permanent storage
- Docker deployment ready

### 🔄 Phase 2: Launch Prep (In Progress)
- Production server deployment
- Mobile app store submissions
- SSL & domain configuration
- Beta testing with users

### 🔜 Phase 3: Social Features
- Follow/unfollow functionality
- Chronological feed
- User search
- Notifications system

### 🔜 Phase 4: Premium & Growth
- Premium tier (global playback)
- Web app (Nuxt 4)
- Advanced map filtering
- NFT marketplace integration

See **[DEPLOYMENT-READY.md](DEPLOYMENT-READY.md)** for detailed next steps.

---

## 🌟 Philosophy

We fight the split: the kids in bubbles, the echo-chambers, the thirst traps engineered by suits who want you addicted, not awake. This is ethical by design—Solana for speed, Arweave for forever, always decentralized, always yours.

**Anti-AI.** If it wasn't caught live, it doesn't exist here.  
**Anti-algorithm.** You follow who you want. You see what they post. End of story.  
**Anti-corporate.** No ads. No data mining. Just a small cut on mints to keep the lights on.

---

## 📞 Links

- **GitHub:** [github.com/alexcolls/now.ink](https://github.com/alexcolls/now.ink)
- **Documentation:** See root directory markdown files
- **Issues:** [GitHub Issues](https://github.com/alexcolls/now.ink/issues)
- **Version:** 0.2.0 (MVP Complete)
- **Status:** Production-ready, deployment pending

---

**Build it. Launch it. Watch people finally look up.**
