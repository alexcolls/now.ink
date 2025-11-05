# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.2] - 2025-11-05

### Added - Marketing & Launch Assets
- Complete landing page (HTML + embedded CSS, responsive)
- Marketing assets guide (606 lines) with:
  - App Store screenshot templates & specifications
  - Google Play screenshot templates & specifications
  - Feature graphic design guidelines
  - App icon requirements
  - Social media graphics (OG image, Twitter cards, Instagram)
  - Press kit structure and brand guidelines
  - Store listing copy (iOS & Android)
  - Launch email templates
  - Video content scripts
  - Analytics & tracking setup
  - Pre-launch checklist
  - Launch day plan
  - Post-launch content calendar
- Official press release (272 lines) including:
  - Media-ready announcement
  - Fact sheet
  - Beta tester quotes templates
  - Interview request information
  - Distribution list (tech, crypto, developer media)
  - Boilerplate company description
- Marketing README with deployment instructions
- PROJECT-STATUS.md - Comprehensive project overview

### Changed
- Project now 100% ready for public launch
- All marketing materials prepared and documented
- Version bump to 0.3.2

## [0.3.1] - 2025-11-05

### Added - Developer Tools & CI/CD
- GitHub Actions CI/CD pipeline with 8 jobs:
  - Backend Go tests with PostgreSQL service
  - golangci-lint for backend code quality
  - TypeScript checks for mobile + blockchain
  - Prettier formatting validation
  - Docker build testing (main branch only)
  - Trivy security scanning with SARIF upload
  - Markdown link checking
  - Deploy notification step
- Pre-commit hooks configuration:
  - Go formatting and linting
  - TypeScript checks
  - Markdown formatting
  - Conventional commit messages
- Enhanced .gitignore with wallet protection

### Added - Mobile Deployment
- Complete mobile deployment guide (535 lines)
- EAS build configuration documented
- iOS TestFlight setup instructions
- Android Play Store setup instructions
- Store listing templates
- Screenshot requirements

### Added - Business Essentials
- Privacy Policy (269 lines) - GDPR compliant
- Terms of Service (354 lines) - App store ready
- Age restrictions documented (13+)
- Data handling policies
- Blockchain permanence disclosures

### Added - Documentation Updates
- CHANGELOG.md created with full version history
- LICENSE with dual licensing model (free personal / paid commercial)
- CONTRIBUTING.md with 332-line contribution guidelines
- README.md updated with badges and roadmap

## [0.3.0] - 2025-11-05

### Added - Solana Production Setup
- Complete Solana production guide (458 lines):
  - Wallet generation instructions
  - Devnet & mainnet configuration
  - RPC provider options and comparison
  - Cost estimates for operations
  - Security best practices
  - Monitoring scripts and balance checking
  - Emergency procedures

### Added - Production Database
- Automated PostgreSQL + PostGIS setup script
- Database backup and restore procedures
- Performance optimization settings

## [0.2.0] - 2025-11-05

### Added - Deployment Infrastructure
- Docker containerization with multi-stage builds
- Docker Compose orchestration (API, PostgreSQL, Redis, Nginx)
- Nginx reverse proxy with rate limiting and security headers
- Production environment configuration templates
- Automated PostgreSQL + PostGIS setup script
- Comprehensive 377-line deployment guide (DEPLOYMENT.md)
- Deployment readiness document (DEPLOYMENT-READY.md)
- Production-ready infrastructure for VPS/cloud deployment

### Added - Mobile UI Completion
- VideoPlayer component for Arweave video playback
- Interactive MapScreen with Google Maps integration
- NFT markers on map with custom camera icons
- Tap-to-play video functionality from markers
- ProfileScreen with user NFT grid (2-column layout)
- User stats display (Moments/Followers/Following)
- Filter NFTs by connected wallet
- Loading and empty states for all screens
- End-to-end testing script (test-e2e.sh)

### Added - Documentation
- MOBILE-COMPLETE.md - Complete mobile feature guide
- SESSION-SUMMARY.md - Development session recap
- DEPLOYMENT-READY.md - Production readiness checklist
- Updated README with current project status
- CHANGELOG.md (this file)
- LICENSE with dual licensing model
- CONTRIBUTING.md guidelines

### Changed
- Mobile app now 100% feature complete for MVP
- All placeholder screens replaced with functional UI
- Production deployment fully automated and documented
- Project status: MVP Complete → Production Ready

### Technical Details
- Added react-native-maps for map functionality
- Implemented Expo AV for video playback
- Created reusable VideoPlayer component
- Database persistence for all NFT data
- Mock/real minting mode support

## [0.1.0] - 2025-11-04

### Added - Core MVP Features
- Go backend API with Fiber framework
- PostgreSQL 16 database with PostGIS extension
- JWT authentication with wallet-based nonces
- React Native + Expo mobile app
- Solana wallet integration (Phantom, Solflare)
- Video recording with GPS tagging
- Arweave integration for permanent storage
- Metaplex NFT minting (mock + real modes)
- Database persistence for streams and NFTs
- Multi-user support with wallet authentication

### Added - Backend Services
- Stream service (start, end, save)
- NFT service (mint, list, get)
- User service (authentication, profile)
- Arweave storage client
- Solana blockchain client
- TypeScript minting script for Metaplex

### Added - Mobile Screens
- HomeScreen with wallet connection
- CameraScreen with recording and minting
- MapScreen (placeholder)
- ProfileScreen (placeholder)
- API integration layer

### Added - Database Schema
- users table with wallet addresses
- streams table with GPS locations (PostGIS)
- nfts table with mint addresses
- follows table for social graph
- sessions table for auth nonces

### Added - Blockchain Integration
- Solana devnet configuration
- Metaplex SDK integration
- Platform wallet setup
- 5%/95% commission split
- NFT metadata standards

### Technical Details
- 5,500+ lines of production code
- 45+ files created
- 15 API endpoints
- 5 database tables
- Mock mode for development
- Production mode for real minting

## [0.0.1] - 2025-11-03

### Added - Initial Setup
- Project structure created
- Repository initialized
- Core documentation started
- Technology stack decisions made
- Development environment setup

---

## Legend

- **Added** - New features
- **Changed** - Changes in existing functionality
- **Deprecated** - Soon-to-be removed features
- **Removed** - Removed features
- **Fixed** - Bug fixes
- **Security** - Vulnerability fixes

## Version Numbering

- **0.x.x** - Pre-1.0 development (current)
- **1.0.0** - First public release (upcoming)
- **1.x.x** - Feature additions (backward compatible)
- **x.0.0** - Breaking changes

---

**Current Version:** 0.2.0  
**Status:** Production-ready, deployment pending  
**Next Release:** 1.0.0 (public launch)
