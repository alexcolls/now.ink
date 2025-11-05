# now.ink

**Your life, minted.**

![License](https://img.shields.io/badge/license-Dual%20License-blue)
![Status](https://img.shields.io/badge/status-pre--alpha-orange)
![Solana](https://img.shields.io/badge/blockchain-Solana-blueviolet)

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

### 🎥 Live Streaming
- Record video directly from your phone in real-time
- Stream live to your followers (or selected viewers)
- No uploads—everything captured in-app for authenticity

### 🪙 NFT Minting
- Save a clip? It auto-mints as a Solana NFT via Metaplex
- Embedded metadata: GPS coordinates, timestamp, user ID, video hash
- Every clip is a unique, verifiable moment in space and time
- Platform takes a small commission (% TBD) on mints—no ads, no privacy exploitation

### 🗺️ Time-Travel Map
- Google Maps-style interface with historical pins
- Time-range slider: filter by today, last week, month, year, or further back
- Zoom in to discover moments tied to specific locations
- Free local playback (if you're physically near the location)
- Premium tier unlocks global playback (pricing TBD)

### 🤝 Real Following, Zero Algorithm
- Follow actual humans—no bots, no ghosts
- Chronological feed: if you follow Alex, you see Alex's pins
- No shadowbans, no pay-to-win reach
- Public by default, auditable by design

### 🔐 Privacy & Control
- Choose public or private clips (friends-only)
- No AI filters or edits—raw footage only
- Users own their NFTs—trade, sell, or keep them forever

---

## 🛠️ Tech Stack

| Layer | Technology |
|-------|-----------|
| **Web Frontend** | Nuxt 4 (TypeScript) |
| **Mobile App** | React Native (TypeScript) |
| **Backend API** | Go (Fiber/Gin framework) |
| **Blockchain** | Solana + Metaplex NFT Standard |
| **Storage** | Arweave (decentralized, permanent) |
| **Streaming** | WebRTC (low-latency) |
| **Database** | TBD (MongoDB/PostgreSQL) |
| **Wallets** | Phantom, Solflare integration |

---

## 📦 Quick Start

### Prerequisites
- Node.js 20+ (for Nuxt 4 & React Native)
- Go 1.21+ (for backend)
- Solana CLI (for blockchain integration)
- Docker (optional, for local dev)

### Clone & Install
```bash
git clone https://github.com/yourusername/now.ink.git
cd now.ink

# Web frontend
cd web
npm install
cp .env.sample .env
npm run dev

# Mobile app
cd ../mobile
npm install
cp .env.sample .env
npm run ios # or npm run android

# Backend API
cd ../backend
go mod download
cp .env.sample .env
go run cmd/api/main.go
```

See **[docs/SETUP.md](docs/SETUP.md)** for detailed setup instructions.

---

## 📁 Repository Structure

```
now.ink/
├── web/                # Nuxt 4 web application
├── mobile/             # React Native mobile app
├── backend/            # Go API server
├── blockchain/         # Solana/Metaplex contracts & scripts
├── docs/               # Comprehensive documentation
│   ├── ARCHITECTURE.md
│   ├── API.md
│   ├── BLOCKCHAIN.md
│   ├── CONTRIBUTING.md
│   ├── SETUP.md
│   ├── USER_GUIDE.md
│   ├── ROADMAP.md
│   └── LICENSE.md
└── README.md
```

---

## 📜 License

**Dual License:**
- **Free for personal use** under [GPL-3.0](docs/LICENSE.md) (or similar)
- **Commercial license required** for profit-making entities

See **[docs/LICENSE.md](docs/LICENSE.md)** for full details.

---

## 🤝 Contributing

We're open-source by design. Fork it, fix it, ship it. See **[docs/CONTRIBUTING.md](docs/CONTRIBUTING.md)** for guidelines.

**Why open-source?**  
Because if we don't open-source the truth, who will? No black-box algorithms, no hidden feeds. Every line of code is auditable. Every decision is yours.

---

## 🚀 Roadmap

**Phase 1 (MVP):** Core recording, minting, basic map  
**Phase 2:** Social features (following, public feed)  
**Phase 3:** Premium tier & monetization  
**Phase 4:** Advanced discovery, marketplace

See **[docs/ROADMAP.md](docs/ROADMAP.md)** for detailed milestones.

---

## 🌟 Philosophy

We fight the split: the kids in bubbles, the echo-chambers, the thirst traps engineered by suits who want you addicted, not awake. This is ethical by design—Solana for speed, Arweave for forever, always decentralized, always yours.

**Anti-AI.** If it wasn't caught live, it doesn't exist here.  
**Anti-algorithm.** You follow who you want. You see what they post. End of story.  
**Anti-corporate.** No ads. No data mining. Just a small cut on mints to keep the lights on.

---

## 📞 Links

- **Website:** [now.ink](https://now.ink)
- **Docs:** [docs/](docs/)
- **Issues:** [GitHub Issues](https://github.com/yourusername/now.ink/issues)
- **Discord:** TBD

---

**Build it. Launch it. Watch people finally look up.**
