# 📱 Mobile App Production Deployment

Complete guide for building and deploying the now.ink mobile app to iOS App Store and Google Play Store.

---

## 📋 Prerequisites

### Required Tools
```bash
# Node.js 20+
node --version

# Expo CLI
npm install -g eas-cli

# Login to Expo
eas login
```

### Required Accounts
- ✅ Expo account (free tier works)
- ✅ Apple Developer Program ($99/year for iOS)
- ✅ Google Play Console ($25 one-time for Android)
- ✅ GitHub account (for version control)

---

## 🔧 Step 1: Production Configuration

### 1.1 Update API URL
```bash
cd mobile

# Create production environment
cat > .env.production << 'EOF'
API_URL=https://api.yourdomain.com/api/v1
ENVIRONMENT=production
EOF

# Update app.json/app.config.js to use env variables
```

### 1.2 Configure app.json
```bash
nano app.json
```

Update key fields:
```json
{
  "expo": {
    "name": "now.ink",
    "slug": "nowink",
    "version": "0.2.0",
    "orientation": "portrait",
    "icon": "./assets/icon.png",
    "userInterfaceStyle": "dark",
    "splash": {
      "image": "./assets/splash.png",
      "resizeMode": "contain",
      "backgroundColor": "#000000"
    },
    "ios": {
      "supportsTablet": false,
      "bundleIdentifier": "com.nowink.app",
      "buildNumber": "1",
      "infoPlist": {
        "NSCameraUsageDescription": "now.ink needs camera access to record moments",
        "NSLocationWhenInUseUsageDescription": "now.ink needs your location to tag moments on the map",
        "NSPhotoLibraryUsageDescription": "now.ink needs access to save videos"
      }
    },
    "android": {
      "adaptiveIcon": {
        "foregroundImage": "./assets/adaptive-icon.png",
        "backgroundColor": "#000000"
      },
      "package": "com.nowink.app",
      "versionCode": 1,
      "permissions": [
        "android.permission.CAMERA",
        "android.permission.ACCESS_FINE_LOCATION",
        "android.permission.ACCESS_COARSE_LOCATION",
        "android.permission.WRITE_EXTERNAL_STORAGE"
      ]
    },
    "extra": {
      "eas": {
        "projectId": "YOUR_EAS_PROJECT_ID"
      }
    }
  }
}
```

### 1.3 Create EAS Configuration
```bash
eas build:configure
```

This creates `eas.json`:
```json
{
  "cli": {
    "version": ">= 5.9.0"
  },
  "build": {
    "development": {
      "developmentClient": true,
      "distribution": "internal",
      "env": {
        "API_URL": "http://localhost:8080/api/v1"
      }
    },
    "preview": {
      "distribution": "internal",
      "env": {
        "API_URL": "https://staging-api.yourdomain.com/api/v1"
      }
    },
    "production": {
      "env": {
        "API_URL": "https://api.yourdomain.com/api/v1"
      }
    }
  },
  "submit": {
    "production": {}
  }
}
```

---

## 🎨 Step 2: App Assets

### 2.1 App Icon
Create 1024x1024px icon:
```bash
# Place at mobile/assets/icon.png
# Requirements:
# - 1024x1024px PNG
# - No transparency
# - Square corners (iOS adds rounded corners)
```

### 2.2 Splash Screen
Create 2048x2048px splash:
```bash
# Place at mobile/assets/splash.png
# Requirements:
# - 2048x2048px PNG or JPG
# - Center 1024x1024 is "safe zone"
```

### 2.3 Adaptive Icon (Android)
Create 1024x1024px foreground:
```bash
# Place at mobile/assets/adaptive-icon.png
# Requirements:
# - 1024x1024px PNG
# - Transparency supported
# - Design fits in center 66%
```

---

## 🍎 Step 3: iOS Build (TestFlight)

### 3.1 Apple Developer Setup
1. Join Apple Developer Program
2. Create App ID: `com.nowink.app`
3. Create App Store Connect app
4. Generate certificates (EAS handles this)

### 3.2 Build for iOS
```bash
cd mobile

# First build (creates credentials)
eas build --platform ios --profile production

# Follow prompts:
# - Allow EAS to manage credentials: Yes
# - Generate new Apple Distribution Certificate: Yes
# - Generate new Push Notification key: Yes (if needed)
```

### 3.3 Submit to TestFlight
```bash
# After successful build
eas submit --platform ios --latest

# Or manually:
# 1. Download .ipa from Expo dashboard
# 2. Upload via Transporter app
# 3. Submit for review in App Store Connect
```

### 3.4 TestFlight Setup
1. Go to App Store Connect
2. Select your app → TestFlight
3. Add Internal Testers (up to 100)
4. Add External Testers (requires beta review)
5. Share TestFlight link with testers

---

## 🤖 Step 4: Android Build (Play Store)

### 4.1 Google Play Console Setup
1. Pay $25 one-time fee
2. Create new app
3. Fill in store listing
4. Set up content rating
5. Complete privacy policy (use docs/PRIVACY-POLICY.md)

### 4.2 Build for Android
```bash
cd mobile

# First build (creates keystore)
eas build --platform android --profile production

# EAS will generate and manage your keystore
# Store backup safely!
```

### 4.3 Submit to Play Store
```bash
# After successful build
eas submit --platform android --latest

# Or manually:
# 1. Download .aab from Expo dashboard
# 2. Upload to Play Console
# 3. Submit for review
```

### 4.4 Play Console Configuration
1. **App content:**
   - Privacy policy URL: https://yourdomain.com/privacy
   - Ads: No ads in app
   - Target audience: Ages 13+
   - Content rating: Apply for rating

2. **Store presence:**
   - App name, description
   - Screenshots (see next section)
   - Feature graphic
   - App category

3. **Release:**
   - Internal testing track (first)
   - Closed testing (beta)
   - Open testing (public beta)
   - Production (final release)

---

## 📸 Step 5: Store Listing Assets

### 5.1 Screenshots (Required)

**iOS (Required sizes):**
- iPhone 6.7" (1290 x 2796 pixels) - 3 screenshots min
- iPhone 6.5" (1284 x 2778 pixels) - 3 screenshots min
- iPad Pro 12.9" (2048 x 2732 pixels) - 2 screenshots min

**Android (Required sizes):**
- Phone (16:9) - 1080 x 1920 pixels - 2 screenshots min
- 7-inch tablet (optional)
- 10-inch tablet (optional)

**Pro tip:** Use screenshot tools:
```bash
# iOS: Xcode Simulator + CMD+S
# Android: Android Studio Emulator + Screenshot button
# Or: Expo Go on physical device
```

### 5.2 Store Descriptions

**Short Description (80 chars):**
```
Record real moments, mint as NFTs. No filters, no algorithms.
```

**Full Description (iOS: 4000 chars, Android: 4000 chars):**
```
now.ink - Your life, minted.

RECORD REAL MOMENTS
Press record and capture authentic moments from your phone. GPS-tagged, time-stamped, forever.

MINT AS NFTs
Save your clip and it auto-mints as a Solana NFT. You own it. Forever. On blockchain.

DISCOVER ON THE MAP
Explore geo-tagged moments on an interactive map. Tap a pin, watch a video, see history unfold.

KEY FEATURES:
✓ Wallet-based login (no passwords)
✓ GPS-tagged video recording
✓ Automatic NFT minting
✓ Permanent Arweave storage
✓ Interactive world map
✓ User profiles

PHILOSOPHY:
• Anti-algorithm: Chronological feed, no recommendations
• Anti-AI: Only real, live-captured content
• User-owned: Your NFTs on Solana blockchain
• No ads: Revenue from 5% platform commission

REQUIREMENTS:
• Solana wallet (Phantom, Solflare)
• Age 13+
• Location services

now.ink isn't social media—it's archaeology you lived through.

Build it. Ship it. Watch people finally look up.
```

### 5.3 Keywords (iOS)
```
nft, blockchain, solana, video, moments, map, crypto, web3, arweave, decentralized
```

### 5.4 Promotional Assets

**iOS App Preview Video (optional but recommended):**
- 15-30 seconds
- Show app in action
- No external branding

**Android Feature Graphic (required):**
- 1024 x 500 pixels
- PNG or JPG
- Showcases app

---

## 🧪 Step 6: Testing

### 6.1 Internal Testing
```bash
# Build internal version
eas build --platform all --profile preview

# Distribute to team
# iOS: TestFlight internal
# Android: Internal testing track
```

### 6.2 Beta Testing Checklist
- [ ] Wallet connection works
- [ ] Camera recording functional
- [ ] GPS tagging accurate
- [ ] Video upload successful
- [ ] NFT minting works (mock mode)
- [ ] Map displays NFTs correctly
- [ ] Profile shows user NFTs
- [ ] Performance acceptable
- [ ] No crashes

### 6.3 Production Testing
Before production release:
- [ ] Test with real backend URL
- [ ] Test real wallet (small amount)
- [ ] Mint test NFT on mainnet
- [ ] Verify on Solscan
- [ ] Check all permissions work
- [ ] Test on multiple devices
- [ ] Verify Terms/Privacy links work

---

## 🚀 Step 7: Production Release

### 7.1 Pre-Launch Checklist
- [ ] Backend deployed to production
- [ ] API URL updated in app
- [ ] SSL certificate configured
- [ ] Database backed up
- [ ] Monitoring setup
- [ ] Support email configured
- [ ] Privacy policy published
- [ ] Terms of service published

### 7.2 iOS Release
```bash
# Build production
eas build --platform ios --profile production

# Submit
eas submit --platform ios --latest

# In App Store Connect:
# 1. Add version info
# 2. Submit for review
# 3. Wait 24-48 hours
# 4. If approved, release manually or auto
```

### 7.3 Android Release
```bash
# Build production
eas build --platform android --profile production

# Submit
eas submit --platform android --latest

# In Play Console:
# 1. Upload to production track
# 2. Complete rollout form
# 3. Submit for review
# 4. Wait 24-72 hours
# 5. If approved, release (rollout or instant)
```

---

## 📊 Step 8: Post-Launch

### 8.1 Monitor Metrics
- Download count
- Active users
- Crash reports
- User reviews
- Support tickets

### 8.2 Update Process
```bash
# Increment version
# In app.json: "version": "0.2.1"
# iOS: "buildNumber": "2"
# Android: "versionCode": 2

# Build and submit
eas build --platform all --profile production
eas submit --platform all --latest
```

### 8.3 Respond to Reviews
- Reply to user feedback (iOS & Android)
- Fix critical bugs quickly
- Release patches as needed

---

## 🔧 Troubleshooting

### Build Fails
```bash
# Clear cache
eas build:cancel

# Check logs
eas build:list
eas build:view BUILD_ID

# Common issues:
# - Missing credentials → Run eas build again
# - Invalid app.json → Validate JSON syntax
# - Dependency errors → Run npm install
```

### Submission Rejected
**Common reasons:**
- Missing privacy policy link
- Incomplete permissions justification
- App crashes on review
- Content policy violations
- Missing age rating

### Certificate Issues
```bash
# View credentials
eas credentials

# Reset (if needed)
eas credentials --platform ios
# Select: Remove all credentials
# Then rebuild
```

---

## 💰 Costs Summary

### Development
- **Expo:** Free (EAS Build has free tier: 30 builds/month)
- **Paid Expo:** $29/month (unlimited builds + priority)

### Distribution
- **Apple Developer:** $99/year
- **Google Play:** $25 one-time
- **Total:** ~$124 first year, $99/year after

### Optional
- **App Store Optimization:** $0-500/month
- **Beta testing tools:** $0 (built-in)
- **Analytics:** $0 (basic), $50+/month (advanced)

---

## 🔗 Useful Links

- **EAS Build:** https://docs.expo.dev/build/introduction/
- **EAS Submit:** https://docs.expo.dev/submit/introduction/
- **App Store Connect:** https://appstoreconnect.apple.com
- **Google Play Console:** https://play.google.com/console
- **App Store Guidelines:** https://developer.apple.com/app-store/review/guidelines/
- **Play Store Policies:** https://play.google.com/about/developer-content-policy/

---

## 📞 Support

Issues with deployment?
- Check EAS build logs
- Review Expo documentation
- Ask in GitHub Issues (tag: mobile)
- Expo Discord: https://chat.expo.dev

---

**Mobile Deployment Complete!** ✅

Your now.ink mobile app is ready for the world.

_Build it. Ship it. Watch people finally look up._ ✨
