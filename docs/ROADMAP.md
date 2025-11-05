# Roadmap

**now.ink Development Roadmap**

---

## Version Strategy

We follow semantic versioning with a twist:
- **v0.x.x** = Pre-1.0 MVP and iterations
- **v1.0.0** = First production-ready release (when explicitly decided)
- Minor/patch bumps happen organically as features ship

**Current Version:** v0.1.0 (planning phase)

---

## Phase 1: MVP (v0.1.0 - v0.5.0)

**Timeline:** Q1-Q2 2025  
**Goal:** Proof of concept with core features

### v0.1.0 - Foundation
- [ ] Project setup (Nuxt 4, React Native, Go backend)
- [ ] Database schema (PostgreSQL with PostGIS)
- [ ] Solana devnet integration
- [ ] Basic auth (wallet-based)

### v0.2.0 - Recording & Streaming
- [ ] Mobile app: camera access + GPS tracking
- [ ] WebRTC signaling server (Go WebSocket)
- [ ] Live streaming (broadcaster → viewers)
- [ ] Video save locally

### v0.3.0 - NFT Minting
- [ ] Arweave integration (video upload)
- [ ] Metaplex NFT minting
- [ ] Backend minting service
- [ ] Platform commission logic (TBD %)

### v0.4.0 - Map Interface
- [ ] Web app: Google Maps integration
- [ ] Display NFT pins on map
- [ ] Time slider (filter by date range)
- [ ] Proximity-based playback

### v0.5.0 - MVP Polish
- [ ] User profiles
- [ ] NFT detail pages
- [ ] Search & filters
- [ ] Basic error handling & UX

**Deliverable:** Functional app where users can record, mint, and discover geo-tagged video NFTs.

---

## Phase 2: Social Features (v0.6.0 - v0.8.0)

**Timeline:** Q3 2025  
**Goal:** Build the anti-algorithm social layer

### v0.6.0 - Following
- [ ] Follow/unfollow users
- [ ] Chronological feed (no algorithm!)
- [ ] Push notifications for new posts

### v0.7.0 - Interactions
- [ ] Comments on NFTs (optional, TBD)
- [ ] Share NFTs (external links)
- [ ] Privacy settings (public/private/friends-only)

### v0.8.0 - Discovery
- [ ] Explore nearby moments
- [ ] Trending locations (simple count, not algorithmic)
- [ ] User search

**Deliverable:** Social platform where users follow real people, see real moments, chronologically.

---

## Phase 3: Premium & Monetization (v0.9.0 - v0.12.0)

**Timeline:** Q4 2025  
**Goal:** Sustainable revenue without ads

### v0.9.0 - Premium Tier
- [ ] Subscription system (pricing TBD)
- [ ] Global playback (bypass proximity)
- [ ] Unlimited mints
- [ ] Solana payment integration

### v0.10.0 - NFT Marketplace
- [ ] Buy/sell NFTs (peer-to-peer)
- [ ] Platform takes small commission (2-5%)
- [ ] Wallet-to-wallet transfers

### v0.11.0 - Creator Tools
- [ ] Analytics dashboard (views, followers)
- [ ] Batch minting
- [ ] Custom collections

### v0.12.0 - Monetization Optimization
- [ ] A/B test commission rates
- [ ] Referral program
- [ ] Premium upsell flows

**Deliverable:** Self-sustaining platform with clear revenue model.

---

## Phase 4: Advanced Features (v0.13.0+)

**Timeline:** 2026+  
**Goal:** Enhancements & community-driven features

### Potential Features
- [ ] AR overlays (view historical moments in AR)
- [ ] Cross-chain support (Ethereum, Polygon)
- [ ] DAO governance (community votes on features)
- [ ] Desktop apps (Electron)
- [ ] Advanced filtering (weather, time of day)
- [ ] Collaborative moments (multi-user streams)
- [ ] NFT royalties (optional creator royalties)
- [ ] Integration with other dApps

**Deliverable:** Feature-rich platform shaped by community feedback.

---

## Future Considerations

### Technical Debt
- Migrate to custom smart contract (if needed)
- Optimize database queries (indexing, caching)
- Implement CDN for faster video delivery
- Horizontal scaling for API servers

### Community Requests
- Will be tracked via GitHub Issues & Discord polls
- Quarterly roadmap review based on user feedback

### Open Questions
- **Commission %:** TBD based on user testing
- **Premium Pricing:** TBD (likely $5-10/month)
- **Content Moderation:** Community-driven vs. manual review?

---

## Version History

| Version | Date | Highlights |
|---------|------|------------|
| v0.1.0  | TBD  | Initial project setup |
| v0.2.0  | TBD  | Live streaming works |
| v0.3.0  | TBD  | First NFT minted on devnet |
| v0.4.0  | TBD  | Map interface live |
| v0.5.0  | TBD  | MVP complete |

---

## Contributing to the Roadmap

Have ideas? Open a [GitHub Discussion](https://github.com/now.ink/discussions) or join our Discord (TBD).

**Remember:** We stay in v0.x.x until we collectively decide it's time for v1.0.0. No rush. Quality over hype.

---

**Build it. Ship it. Watch people finally look up.**
