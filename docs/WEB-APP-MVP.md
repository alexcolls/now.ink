# 🌐 Web App MVP Guide

Complete guide for building the now.ink web application with Nuxt 4.

---

## 📋 Overview

The now.ink web app provides:
- Interactive map view of all GPS-tagged NFTs
- Chronological feed from followed users
- User profiles and social features
- Solana wallet connection
- NFT viewing and details
- Responsive design (desktop + mobile web)

**Tech Stack:**
- **Framework**: Nuxt 4 (Vue 3 + Vite)
- **Styling**: TailwindCSS 3
- **Maps**: MapLibre GL + Maptiler
- **Wallet**: Solana Wallet Adapter
- **State**: Pinia
- **HTTP**: ofetch (Nuxt native)

---

## 🚀 Project Setup

### Initialize Nuxt 4 Project

```bash
cd /home/quantium/labs/now.ink
npx nuxi@latest init web
cd web
```

**Select options:**
- Package manager: npm
- Initialize git: No (already in monorepo)

### Install Dependencies

```bash
npm install
npm install -D tailwindcss postcss autoprefixer
npm install @solana/wallet-adapter-base @solana/wallet-adapter-wallets @solana/wallet-adapter-vue @solana/web3.js
npm install maplibre-gl @maptiler/sdk
npm install pinia @pinia/nuxt
npm install @vueuse/core
npm install dayjs
```

### Initialize TailwindCSS

```bash
npx tailwindcss init -p
```

---

## 📁 Project Structure

```
web/
├── app.vue                 # Root component
├── nuxt.config.ts         # Nuxt configuration
├── tailwind.config.js     # Tailwind configuration
├── tsconfig.json          # TypeScript config
├── package.json
├── .env.sample
├── public/
│   ├── favicon.ico
│   └── og-image.png
├── assets/
│   └── css/
│       └── main.css       # Global styles
├── components/
│   ├── Map/
│   │   ├── InteractiveMap.vue
│   │   ├── MarkerCluster.vue
│   │   └── VideoMarker.vue
│   ├── Feed/
│   │   ├── FeedList.vue
│   │   ├── FeedItem.vue
│   │   └── VideoPlayer.vue
│   ├── Profile/
│   │   ├── UserCard.vue
│   │   ├── NFTGrid.vue
│   │   └── StatsBar.vue
│   ├── Wallet/
│   │   ├── WalletButton.vue
│   │   └── WalletModal.vue
│   ├── Common/
│   │   ├── Button.vue
│   │   ├── Modal.vue
│   │   └── Loading.vue
│   └── Layout/
│       ├── Header.vue
│       ├── Sidebar.vue
│       └── Footer.vue
├── composables/
│   ├── useApi.ts          # API client
│   ├── useAuth.ts         # Authentication
│   ├── useWallet.ts       # Wallet connection
│   ├── useNFT.ts          # NFT data
│   └── useFeed.ts         # Feed data
├── stores/
│   ├── auth.ts            # Auth state
│   ├── nfts.ts            # NFT state
│   └── user.ts            # User state
├── pages/
│   ├── index.vue          # Home page (map)
│   ├── feed.vue           # Feed page
│   ├── profile/
│   │   └── [id].vue       # User profile
│   ├── nft/
│   │   └── [mint].vue     # NFT details
│   ├── search.vue         # User search
│   └── about.vue          # About page
├── middleware/
│   └── auth.ts            # Auth middleware
├── types/
│   ├── nft.ts
│   ├── user.ts
│   └── api.ts
└── utils/
    ├── constants.ts
    └── helpers.ts
```

---

## ⚙️ Configuration Files

### `nuxt.config.ts`

```typescript
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  
  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
    '@vueuse/nuxt',
  ],

  css: ['~/assets/css/main.css'],

  app: {
    head: {
      title: 'now.ink - Real Moments, Real Ownership',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        {
          name: 'description',
          content: 'Record GPS-tagged video moments and mint them as NFTs. No algorithms, no AI, just authentic moments.',
        },
        { property: 'og:title', content: 'now.ink - Real Moments, Real Ownership' },
        { property: 'og:type', content: 'website' },
        { property: 'og:url', content: 'https://now.ink' },
        { property: 'og:image', content: 'https://now.ink/og-image.png' },
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
      ],
    },
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080',
      solanaNetwork: process.env.NUXT_PUBLIC_SOLANA_NETWORK || 'devnet',
      maptilerApiKey: process.env.NUXT_PUBLIC_MAPTILER_API_KEY || '',
    },
  },

  vite: {
    optimizeDeps: {
      include: ['@solana/web3.js', '@solana/wallet-adapter-base'],
    },
  },
});
```

### `tailwind.config.js`

```javascript
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './components/**/*.{js,vue,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './plugins/**/*.{js,ts}',
    './app.vue',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#667eea',
          50: '#f5f7ff',
          100: '#ebf0ff',
          200: '#d6dfff',
          300: '#b9c9ff',
          400: '#94abff',
          500: '#667eea',
          600: '#5568d3',
          700: '#4553b8',
          800: '#374199',
          900: '#2d347a',
        },
        secondary: {
          DEFAULT: '#764ba2',
          50: '#faf6ff',
          100: '#f5edff',
          200: '#e9d9ff',
          300: '#d7b8ff',
          400: '#bd8eff',
          500: '#764ba2',
          600: '#653f8a',
          700: '#543372',
          800: '#44285c',
          900: '#36204a',
        },
      },
    },
  },
  plugins: [],
};
```

### `.env.sample`

```env
# API Configuration
NUXT_PUBLIC_API_BASE=http://localhost:8080

# Solana Configuration
NUXT_PUBLIC_SOLANA_NETWORK=devnet
# Options: devnet, testnet, mainnet-beta

# Maptiler API Key
NUXT_PUBLIC_MAPTILER_API_KEY=your_maptiler_api_key_here

# Optional: Analytics
NUXT_PUBLIC_PLAUSIBLE_DOMAIN=now.ink
```

---

## 🧩 Core Components

### Interactive Map Component

**`components/Map/InteractiveMap.vue`:**

```vue
<template>
  <div class="relative w-full h-full">
    <div ref="mapContainer" class="w-full h-full rounded-lg" />
    
    <div v-if="loading" class="absolute inset-0 flex items-center justify-center bg-gray-900/50">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-white" />
    </div>
  </div>
</template>

<script setup lang="ts">
import maplibregl from 'maplibre-gl';
import 'maplibre-gl/dist/maplibre-gl.css';

const props = defineProps<{
  nfts: any[];
  center?: [number, number];
  zoom?: number;
}>();

const emit = defineEmits<{
  markerClick: [nft: any];
}>();

const mapContainer = ref<HTMLDivElement>();
const map = ref<maplibregl.Map>();
const loading = ref(true);
const config = useRuntimeConfig();

onMounted(() => {
  initMap();
});

onBeforeUnmount(() => {
  map.value?.remove();
});

watch(() => props.nfts, (newNfts) => {
  if (map.value && newNfts) {
    updateMarkers(newNfts);
  }
});

function initMap() {
  if (!mapContainer.value) return;

  map.value = new maplibregl.Map({
    container: mapContainer.value,
    style: `https://api.maptiler.com/maps/streets-v2/style.json?key=${config.public.maptilerApiKey}`,
    center: props.center || [-74.006, 40.7128], // NYC default
    zoom: props.zoom || 10,
  });

  map.value.on('load', () => {
    loading.value = false;
    if (props.nfts) {
      updateMarkers(props.nfts);
    }
  });

  // Add navigation controls
  map.value.addControl(new maplibregl.NavigationControl(), 'top-right');
}

function updateMarkers(nfts: any[]) {
  if (!map.value) return;

  // Remove existing markers
  const existingMarkers = document.querySelectorAll('.custom-marker');
  existingMarkers.forEach(m => m.remove());

  // Add new markers
  nfts.forEach((nft) => {
    const el = document.createElement('div');
    el.className = 'custom-marker';
    el.innerHTML = '📍';
    el.style.fontSize = '24px';
    el.style.cursor = 'pointer';

    el.addEventListener('click', () => {
      emit('markerClick', nft);
    });

    new maplibregl.Marker({ element: el })
      .setLngLat([nft.longitude, nft.latitude])
      .addTo(map.value!);
  });
}
</script>
```

### Feed List Component

**`components/Feed/FeedList.vue`:**

```vue
<template>
  <div class="space-y-4">
    <div v-if="loading" class="text-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto" />
    </div>

    <div v-else-if="feed.length === 0" class="text-center py-12">
      <p class="text-gray-500">No moments in your feed yet.</p>
      <p class="text-gray-400 text-sm mt-2">Follow users to see their moments here!</p>
    </div>

    <div v-else>
      <FeedItem
        v-for="item in feed"
        :key="item.mint_address"
        :nft="item"
        @click="handleClick(item)"
      />
    </div>

    <button
      v-if="hasMore"
      @click="loadMore"
      class="w-full py-3 bg-gray-100 hover:bg-gray-200 rounded-lg transition"
    >
      Load More
    </button>
  </div>
</template>

<script setup lang="ts">
const { feed, loading, hasMore, loadMore } = useFeed();

const router = useRouter();

function handleClick(nft: any) {
  router.push(`/nft/${nft.mint_address}`);
}
</script>
```

### Wallet Connection Component

**`components/Wallet/WalletButton.vue`:**

```vue
<template>
  <div>
    <button
      v-if="!connected"
      @click="showModal = true"
      class="px-6 py-2 bg-primary hover:bg-primary-600 text-white rounded-lg transition"
    >
      Connect Wallet
    </button>

    <div v-else class="flex items-center gap-3">
      <span class="text-sm text-gray-600">
        {{ shortenAddress(publicKey) }}
      </span>
      <button
        @click="disconnect"
        class="px-4 py-2 bg-gray-100 hover:bg-gray-200 rounded-lg text-sm transition"
      >
        Disconnect
      </button>
    </div>

    <WalletModal v-model="showModal" @connect="handleConnect" />
  </div>
</template>

<script setup lang="ts">
const { connected, publicKey, disconnect } = useWallet();
const showModal = ref(false);

function shortenAddress(address: string | null) {
  if (!address) return '';
  return `${address.slice(0, 4)}...${address.slice(-4)}`;
}

function handleConnect() {
  showModal.value = false;
}
</script>
```

---

## 🔧 Composables

### API Client (`composables/useApi.ts`)

```typescript
export const useApi = () => {
  const config = useRuntimeConfig();
  const baseURL = config.public.apiBase;
  const { token } = useAuth();

  const api = $fetch.create({
    baseURL,
    headers: computed(() => ({
      'Content-Type': 'application/json',
      ...(token.value ? { Authorization: `Bearer ${token.value}` } : {}),
    })),
  });

  return {
    // NFT endpoints
    async getNFTs(limit = 50, offset = 0) {
      return api('/api/v1/nfts', {
        params: { limit, offset },
      });
    },

    async getNFT(mintAddress: string) {
      return api(`/api/v1/nfts/${mintAddress}`);
    },

    // Feed endpoint
    async getFeed(limit = 20, offset = 0) {
      return api('/api/v1/social/feed', {
        params: { limit, offset },
      });
    },

    // User endpoints
    async getUserProfile(userId: string) {
      return api(`/api/v1/users/${userId}`);
    },

    async searchUsers(query: string, limit = 20) {
      return api('/api/v1/users/search', {
        params: { q: query, limit },
      });
    },

    // Social endpoints
    async followUser(userId: string) {
      return api(`/api/v1/social/follow/${userId}`, {
        method: 'POST',
      });
    },

    async unfollowUser(userId: string) {
      return api(`/api/v1/social/follow/${userId}`, {
        method: 'DELETE',
      });
    },
  };
};
```

### Wallet Composable (`composables/useWallet.ts`)

```typescript
import { ref, computed } from 'vue';

export const useWallet = () => {
  const connected = ref(false);
  const publicKey = ref<string | null>(null);
  const wallet = ref<any>(null);

  async function connect() {
    try {
      // Detect Phantom wallet
      const { solana } = window as any;
      
      if (!solana?.isPhantom) {
        alert('Please install Phantom wallet!');
        window.open('https://phantom.app/', '_blank');
        return;
      }

      const response = await solana.connect();
      publicKey.value = response.publicKey.toString();
      wallet.value = solana;
      connected.value = true;

      // Store in localStorage
      localStorage.setItem('walletConnected', 'true');
      localStorage.setItem('walletPublicKey', publicKey.value);
    } catch (error) {
      console.error('Wallet connection error:', error);
    }
  }

  async function disconnect() {
    if (wallet.value) {
      await wallet.value.disconnect();
    }
    
    connected.value = false;
    publicKey.value = null;
    wallet.value = null;
    
    localStorage.removeItem('walletConnected');
    localStorage.removeItem('walletPublicKey');
  }

  // Auto-reconnect on page load
  onMounted(() => {
    const wasConnected = localStorage.getItem('walletConnected');
    if (wasConnected === 'true') {
      connect();
    }
  });

  return {
    connected,
    publicKey,
    wallet,
    connect,
    disconnect,
  };
};
```

---

## 📄 Pages

### Home Page (`pages/index.vue`)

```vue
<template>
  <div class="h-screen flex flex-col">
    <Header />
    
    <div class="flex-1 relative">
      <InteractiveMap
        :nfts="nfts"
        @marker-click="handleMarkerClick"
      />
    </div>

    <!-- NFT Detail Modal -->
    <Modal v-model="showDetail" v-if="selectedNFT">
      <NFTDetail :nft="selectedNFT" />
    </Modal>
  </div>
</template>

<script setup lang="ts">
const { nfts, loading } = useNFT();
const selectedNFT = ref(null);
const showDetail = ref(false);

function handleMarkerClick(nft: any) {
  selectedNFT.value = nft;
  showDetail.value = true;
}

useSeoMeta({
  title: 'now.ink - Interactive Map',
  description: 'Explore GPS-tagged video moments from around the world',
});
</script>
```

### Feed Page (`pages/feed.vue`)

```vue
<template>
  <div class="min-h-screen bg-gray-50">
    <Header />
    
    <div class="max-w-3xl mx-auto py-8 px-4">
      <h1 class="text-3xl font-bold mb-8">Your Feed</h1>
      <FeedList />
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
});

useSeoMeta({
  title: 'Feed - now.ink',
  description: 'Your chronological feed of moments from people you follow',
});
</script>
```

---

## 🚀 Running the Project

### Development

```bash
cd web
npm run dev
```

Open http://localhost:3000

### Build for Production

```bash
npm run build
npm run preview
```

### Deploy

**Vercel (Recommended):**
```bash
npm install -g vercel
vercel
```

**Netlify:**
```bash
npm run generate
# Upload .output/public to Netlify
```

---

## ✅ Implementation Checklist

### Setup
- [ ] Initialize Nuxt 4 project
- [ ] Install dependencies
- [ ] Configure Tailwind CSS
- [ ] Setup environment variables
- [ ] Create project structure

### Components
- [ ] Interactive map with markers
- [ ] Feed list with infinite scroll
- [ ] NFT detail view
- [ ] Wallet connection modal
- [ ] User profile cards
- [ ] Header/footer layout

### Features
- [ ] Map view with NFT markers
- [ ] Chronological feed
- [ ] User profiles
- [ ] Search functionality
- [ ] Follow/unfollow
- [ ] Wallet authentication
- [ ] Responsive design

### Testing
- [ ] Test on desktop browsers
- [ ] Test on mobile browsers
- [ ] Test wallet connection
- [ ] Test API integration
- [ ] Performance testing

---

**Created**: November 5, 2025  
**Version**: 1.0.0  
**Status**: Ready for Implementation
