# 📣 Marketing Materials

Complete marketing assets and launch materials for now.ink.

---

## 📂 Contents

### 1. Landing Page
- **File**: `landing-page.html`
- **Status**: ✅ Complete
- **Description**: Single-page HTML landing page with embedded CSS
- **Features**:
  - Responsive design
  - App Store & Play Store download buttons
  - Feature showcase
  - Philosophy section
  - "How It Works" guide
  - Email beta signup CTA
  - Social media links
  - Footer with legal links

**Deployment**: Host on static hosting (Netlify, Vercel, GitHub Pages, or custom domain)

### 2. Marketing Assets Guide
- **File**: `MARKETING-ASSETS.md`
- **Status**: ✅ Complete (606 lines)
- **Includes**:
  - App Store screenshot requirements & templates
  - Google Play screenshot requirements & templates
  - Feature graphic specifications
  - App icon design guidelines
  - Social media graphics (OG image, Twitter cards, Instagram)
  - Press kit structure
  - Store listing copy (iOS & Android)
  - Launch email template
  - Video content scripts
  - Analytics & tracking setup
  - Pre-launch checklist
  - Launch day plan
  - Post-launch content calendar

### 3. Press Release
- **File**: `PRESS-RELEASE.md`
- **Status**: ✅ Complete (272 lines)
- **Includes**:
  - Official press release
  - Fact sheet
  - Beta tester quotes (templates)
  - Social media handles
  - Media contact information
  - Interview request information
  - Distribution list (tech, crypto, developer media)
  - Boilerplate company description

---

## 🚀 Quick Start

### Deploying Landing Page

#### Option 1: Netlify (Recommended - Easiest)
```bash
# Install Netlify CLI
npm install -g netlify-cli

# Deploy from marketing folder
cd marketing
netlify deploy

# For production
netlify deploy --prod
```

#### Option 2: Vercel
```bash
# Install Vercel CLI
npm install -g vercel

# Deploy
cd marketing
vercel
```

#### Option 3: GitHub Pages
```bash
# Copy landing page to docs/ or gh-pages branch
cp marketing/landing-page.html docs/index.html
git add docs/index.html
git commit -m "Add landing page"
git push origin main

# Enable GitHub Pages in repository settings
```

#### Option 4: Custom Server
```bash
# Use nginx or Apache to serve landing-page.html
# Place in /var/www/html/ or equivalent
cp landing-page.html /var/www/html/index.html
```

### Customizing Landing Page

1. **Update App Store Links** (Lines 359, 365)
```html
<a href="YOUR_IOS_APP_STORE_LINK" class="btn btn-primary">
<a href="YOUR_ANDROID_PLAY_STORE_LINK" class="btn btn-secondary">
```

2. **Update Beta Email** (Line 477)
```html
<a href="mailto:YOUR_EMAIL@now.ink?subject=Beta%20Access%20Request">
```

3. **Update Social Links** (Lines 495-509)
```html
<a href="YOUR_TWITTER_LINK">
<a href="YOUR_DISCORD_LINK">
<a href="YOUR_GITHUB_LINK">
```

4. **Add Analytics** (Before `</head>`)
```html
<!-- Google Analytics -->
<script async src="https://www.googletagmanager.com/gtag/js?id=GA_MEASUREMENT_ID"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());
  gtag('config', 'GA_MEASUREMENT_ID');
</script>
```

---

## ✅ Pre-Launch Checklist

### Graphics (To Create)
- [ ] App icon (1024x1024 iOS, 512x512 Android)
- [ ] iOS screenshots (5 per device size)
- [ ] Android screenshots (8 total)
- [ ] Feature graphic for Play Store (1024x500)
- [ ] OG image for social sharing (1200x630)
- [ ] Twitter card (1200x675)
- [ ] Instagram posts (1080x1080)
- [ ] Instagram stories (1080x1920)

### Copy (Complete ✅)
- [x] Landing page HTML
- [x] App Store listing text
- [x] Play Store listing text
- [x] Press release
- [x] Launch email template
- [x] Social media post templates

### Video (To Create)
- [ ] App Store preview video (30s)
- [ ] Social media teaser (15s)
- [ ] Demo video for YouTube
- [ ] Tutorial videos

### Setup
- [ ] Domain configured (now.ink)
- [ ] Landing page deployed
- [ ] Analytics installed (Google Analytics, Plausible, etc.)
- [ ] Social media accounts created
  - [ ] Twitter/X: @nowink_app
  - [ ] Instagram: @nowink_app
  - [ ] Discord server
  - [ ] TikTok: @nowink_app
- [ ] Email setup (hello@, press@, beta@now.ink)
- [ ] App Store developer accounts
  - [ ] Apple Developer ($99/year)
  - [ ] Google Play Console ($25 one-time)

---

## 📧 Email Addresses Needed

Setup these email forwards/aliases:

- **hello@now.ink** - General inquiries
- **press@now.ink** - Media/press inquiries
- **beta@now.ink** - Beta signup requests
- **support@now.ink** - User support
- **security@now.ink** - Security reports

---

## 🎨 Design Assets Needed

### Priority 1 (Required for Launch)
1. **App Icon** - Professional design
2. **5 iOS Screenshots** - With overlay text
3. **8 Android Screenshots** - With overlay text
4. **OG Image** - For social sharing

### Priority 2 (Nice to Have)
5. Feature graphic for Play Store
6. Twitter card image
7. Instagram post templates
8. Press kit logo package

### Tools for Creating Assets

**Free Options**:
- **Figma** - Design tool (free tier)
- **Canva** - Templates (free tier)
- **Pixelmator** - macOS image editor ($39 one-time)
- **GIMP** - Free Photoshop alternative

**Paid Options**:
- **Adobe Photoshop** - Industry standard
- **Sketch** - macOS design tool
- **Affinity Designer** - One-time purchase alternative

**Screenshot Tools**:
- **Screenshot Maker** (https://www.screenshot.guru/) - Free mockups
- **Previewed** (https://previewed.app/) - Device mockups
- **Rotato** - 3D device mockups (macOS)

---

## 📱 Social Media Strategy

### Twitter/X Launch
```
🎉 We're live!

now.ink is now available on iOS & Android.

📍 Record GPS-tagged moments
🪙 Mint as NFTs on Solana
🗺️ Explore the world through real videos

No algorithms. No AI. Just authentic moments you own forever.

Download now: [link]

#NowInk #Web3 #Solana
```

### Product Hunt Launch
- **Title**: now.ink - Own your video moments as NFTs
- **Tagline**: GPS-tagged videos minted as blockchain NFTs. No algorithms. No AI.
- **Description**: Use text from PRESS-RELEASE.md
- **Tags**: blockchain, social-media, web3, nft, mobile-app

### Hacker News "Show HN"
```
Show HN: now.ink – GPS-tagged video NFTs with no algorithms

We built a decentralized social platform where users record 
GPS-tagged videos and mint them as NFTs on Solana. 

Key principles:
- No algorithms (chronological feed only)
- No AI-generated content
- Users own their content via blockchain
- Open source

Tech stack: Go + Fiber backend, React Native mobile, 
Solana + Arweave for blockchain storage.

Would love feedback from the HN community!

GitHub: https://github.com/alexcolls/now.ink
```

### Reddit Posts
- r/solana - Focus on blockchain tech
- r/cryptocurrency - Focus on NFT ownership
- r/reactnative - Focus on tech stack
- r/opensource - Focus on open source aspect

---

## 📊 Metrics to Track

### Week 1 Goals
- 100 landing page visits
- 50 app downloads
- 10 NFTs minted
- 5 social media mentions

### Month 1 Goals
- 1,000 landing page visits
- 200 app downloads
- 100 NFTs minted
- 50 active users
- 10 press mentions

### Quarter 1 Goals
- 10,000 landing page visits
- 1,000 app downloads
- 500 active users
- 2,500 NFTs minted
- Product Hunt feature

---

## 🎯 Launch Distribution Plan

### Day 1 (Launch Day)
- **9:00 AM**: Apps go live on stores
- **9:15 AM**: Update landing page with download links
- **9:30 AM**: Twitter announcement + thread
- **10:00 AM**: Product Hunt launch
- **10:30 AM**: Email list announcement
- **11:00 AM**: Reddit posts (r/solana, r/cryptocurrency)
- **2:00 PM**: Hacker News "Show HN"
- **Throughout day**: Respond to comments, monitor metrics

### Days 2-7
- Daily: Share user-generated content
- Daily: Engage with comments/mentions
- Day 2: Feature highlight (GPS tagging)
- Day 3: "How it works" tutorial
- Day 4: User spotlight
- Day 5: Behind the scenes
- Day 6: Feature highlight (NFTs)
- Day 7: Week 1 recap + metrics

### Weeks 2-4
- Weekly feature highlights
- User testimonials
- Development updates
- Community challenges
- Partnership announcements

---

## 💡 Content Ideas

### Tutorial Videos
1. "How to mint your first video NFT on now.ink"
2. "Exploring the world map"
3. "Setting up your Solana wallet"
4. "Privacy settings explained"

### Blog Posts
1. "Why we built now.ink"
2. "The problem with algorithmic feeds"
3. "True ownership through blockchain"
4. "How NFT minting works"
5. "Our open-source philosophy"

### Social Media Posts
1. Feature spotlights (one per week)
2. User stories
3. Development updates
4. Behind-the-scenes
5. Tips & tricks
6. Community highlights

---

## 🔗 Important Links

- **Website**: https://now.ink
- **GitHub**: https://github.com/alexcolls/now.ink
- **App Store**: [To be added]
- **Play Store**: [To be added]
- **Twitter**: https://twitter.com/nowink_app
- **Discord**: https://discord.gg/nowink
- **Email**: hello@now.ink

---

## 📄 License & Usage

All marketing materials are part of the now.ink project and follow the same dual-license model:
- **Personal use**: Free
- **Commercial use**: Requires commercial license

Press and media may use materials for coverage of now.ink without prior approval.

---

**Created**: November 5, 2025  
**Version**: 1.0.0  
**Status**: Ready for Launch
