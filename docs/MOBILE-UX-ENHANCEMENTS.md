# 📱 Mobile UX Enhancements Guide

Complete guide for implementing enhanced UX features in the now.ink mobile app.

---

## 📋 Overview

This guide covers implementing:
1. **Video Thumbnails** - Generate and cache video thumbnails
2. **Pull-to-Refresh** - Refresh feeds and maps
3. **Map Search & Filters** - Search locations, filter by date/user
4. **Local Caching** - Cache NFT data and images
5. **Sharing** - Share NFTs to other apps
6. **Performance Optimizations** - FlatList optimization, image caching

---

## 🎬 Video Thumbnails

### Installation

```bash
cd mobile
npx expo install expo-video-thumbnails expo-file-system
```

### Thumbnail Generator Utility

**Create `mobile/src/utils/thumbnailGenerator.ts`:**

```typescript
import * as VideoThumbnails from 'expo-video-thumbnails';
import * as FileSystem from 'expo-file-system';
import { Platform } from 'react-native';

interface ThumbnailOptions {
  videoUri: string;
  time?: number; // Time in seconds for thumbnail
  quality?: number; // 0-1, default 0.7
}

export class ThumbnailGenerator {
  private static cache = new Map<string, string>();
  private static cacheDir = `${FileSystem.cacheDirectory}thumbnails/`;

  static async initialize() {
    // Create thumbnail cache directory
    const dirInfo = await FileSystem.getInfoAsync(this.cacheDir);
    if (!dirInfo.exists) {
      await FileSystem.makeDirectoryAsync(this.cacheDir, { intermediates: true });
    }
  }

  static async generate(options: ThumbnailOptions): Promise<string | null> {
    const { videoUri, time = 1000, quality = 0.7 } = options;

    // Check cache first
    const cacheKey = `${videoUri}_${time}`;
    if (this.cache.has(cacheKey)) {
      return this.cache.get(cacheKey)!;
    }

    try {
      const { uri } = await VideoThumbnails.getThumbnailAsync(videoUri, {
        time,
        quality,
      });

      // Move to cache directory
      const filename = `thumb_${Date.now()}.jpg`;
      const cachedUri = `${this.cacheDir}${filename}`;
      await FileSystem.moveAsync({
        from: uri,
        to: cachedUri,
      });

      this.cache.set(cacheKey, cachedUri);
      return cachedUri;
    } catch (error) {
      console.error('Failed to generate thumbnail:', error);
      return null;
    }
  }

  static async clearCache() {
    try {
      await FileSystem.deleteAsync(this.cacheDir, { idempotent: true });
      await FileSystem.makeDirectoryAsync(this.cacheDir, { intermediates: true });
      this.cache.clear();
    } catch (error) {
      console.error('Failed to clear thumbnail cache:', error);
    }
  }

  static getCached(videoUri: string, time: number = 1000): string | null {
    const cacheKey = `${videoUri}_${time}`;
    return this.cache.get(cacheKey) || null;
  }
}
```

### Video Thumbnail Component

**Create `mobile/src/components/VideoThumbnail.tsx`:**

```typescript
import React, { useEffect, useState } from 'react';
import { View, Image, ActivityIndicator, StyleSheet } from 'react-native';
import { ThumbnailGenerator } from '../utils/thumbnailGenerator';

interface VideoThumbnailProps {
  videoUri: string;
  time?: number;
  style?: any;
  placeholderColor?: string;
}

export const VideoThumbnail: React.FC<VideoThumbnailProps> = ({
  videoUri,
  time = 1000,
  style,
  placeholderColor = '#667eea',
}) => {
  const [thumbnailUri, setThumbnailUri] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadThumbnail();
  }, [videoUri, time]);

  const loadThumbnail = async () => {
    setLoading(true);
    
    // Check cache first
    const cached = ThumbnailGenerator.getCached(videoUri, time);
    if (cached) {
      setThumbnailUri(cached);
      setLoading(false);
      return;
    }

    // Generate new thumbnail
    const uri = await ThumbnailGenerator.generate({ videoUri, time });
    setThumbnailUri(uri);
    setLoading(false);
  };

  if (loading) {
    return (
      <View style={[styles.container, style, { backgroundColor: placeholderColor }]}>
        <ActivityIndicator color="#fff" />
      </View>
    );
  }

  if (!thumbnailUri) {
    return (
      <View style={[styles.container, style, { backgroundColor: placeholderColor }]} />
    );
  }

  return (
    <Image
      source={{ uri: thumbnailUri }}
      style={[styles.image, style]}
      resizeMode="cover"
    />
  );
};

const styles = StyleSheet.create({
  container: {
    width: 150,
    height: 200,
    justifyContent: 'center',
    alignItems: 'center',
    borderRadius: 8,
  },
  image: {
    width: 150,
    height: 200,
    borderRadius: 8,
  },
});
```

---

## 🔄 Pull-to-Refresh

### Feed with Pull-to-Refresh

**Update `mobile/src/screens/FeedScreen.tsx`:**

```typescript
import React, { useState, useEffect, useCallback } from 'react';
import {
  View,
  FlatList,
  RefreshControl,
  StyleSheet,
  Text,
} from 'react-native';
import { api } from '../utils/api';
import { VideoThumbnail } from '../components/VideoThumbnail';

export const FeedScreen = () => {
  const [feed, setFeed] = useState([]);
  const [refreshing, setRefreshing] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadFeed();
  }, []);

  const loadFeed = async () => {
    try {
      const response = await api.getFeed(20, 0);
      setFeed(response.feed);
    } catch (error) {
      console.error('Failed to load feed:', error);
    } finally {
      setLoading(false);
    }
  };

  const onRefresh = useCallback(async () => {
    setRefreshing(true);
    await loadFeed();
    setRefreshing(false);
  }, []);

  const renderItem = ({ item }: any) => (
    <View style={styles.feedItem}>
      <VideoThumbnail
        videoUri={item.video_url}
        style={styles.thumbnail}
      />
      <View style={styles.itemInfo}>
        <Text style={styles.title}>{item.name}</Text>
        <Text style={styles.creator}>@{item.creator_username}</Text>
        <Text style={styles.timestamp}>
          {new Date(item.created_at).toLocaleDateString()}
        </Text>
      </View>
    </View>
  );

  return (
    <FlatList
      data={feed}
      renderItem={renderItem}
      keyExtractor={(item) => item.mint_address}
      refreshControl={
        <RefreshControl
          refreshing={refreshing}
          onRefresh={onRefresh}
          tintColor="#667eea"
          colors={['#667eea']}
        />
      }
      contentContainerStyle={styles.container}
    />
  );
};

const styles = StyleSheet.create({
  container: {
    padding: 16,
  },
  feedItem: {
    flexDirection: 'row',
    marginBottom: 16,
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 12,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  thumbnail: {
    width: 100,
    height: 100,
    marginRight: 12,
  },
  itemInfo: {
    flex: 1,
    justifyContent: 'center',
  },
  title: {
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 4,
  },
  creator: {
    fontSize: 14,
    color: '#667eea',
    marginBottom: 4,
  },
  timestamp: {
    fontSize: 12,
    color: '#999',
  },
});
```

---

## 🔍 Map Search & Filters

### Map Search Component

**Create `mobile/src/components/MapSearch.tsx`:**

```typescript
import React, { useState } from 'react';
import {
  View,
  TextInput,
  StyleSheet,
  TouchableOpacity,
  Text,
  Modal,
  ScrollView,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';

interface MapSearchProps {
  onSearch: (query: string) => void;
  onFilterChange: (filters: MapFilters) => void;
}

interface MapFilters {
  dateRange?: { start: Date; end: Date };
  users?: string[];
  radius?: number;
}

export const MapSearch: React.FC<MapSearchProps> = ({
  onSearch,
  onFilterChange,
}) => {
  const [query, setQuery] = useState('');
  const [showFilters, setShowFilters] = useState(false);
  const [filters, setFilters] = useState<MapFilters>({});

  const handleSearch = () => {
    onSearch(query);
  };

  const handleFilterApply = () => {
    onFilterChange(filters);
    setShowFilters(false);
  };

  return (
    <>
      <View style={styles.container}>
        <View style={styles.searchBar}>
          <Ionicons name="search" size={20} color="#999" />
          <TextInput
            style={styles.input}
            placeholder="Search location..."
            value={query}
            onChangeText={setQuery}
            onSubmitEditing={handleSearch}
            returnKeyType="search"
          />
          {query.length > 0 && (
            <TouchableOpacity onPress={() => setQuery('')}>
              <Ionicons name="close-circle" size={20} color="#999" />
            </TouchableOpacity>
          )}
        </View>
        
        <TouchableOpacity
          style={styles.filterButton}
          onPress={() => setShowFilters(true)}
        >
          <Ionicons name="filter" size={20} color="#667eea" />
        </TouchableOpacity>
      </View>

      {/* Filter Modal */}
      <Modal
        visible={showFilters}
        animationType="slide"
        transparent
        onRequestClose={() => setShowFilters(false)}
      >
        <View style={styles.modalOverlay}>
          <View style={styles.modalContent}>
            <View style={styles.modalHeader}>
              <Text style={styles.modalTitle}>Filters</Text>
              <TouchableOpacity onPress={() => setShowFilters(false)}>
                <Ionicons name="close" size={24} color="#333" />
              </TouchableOpacity>
            </View>

            <ScrollView style={styles.filterList}>
              {/* Date Range Filter */}
              <View style={styles.filterSection}>
                <Text style={styles.filterLabel}>Date Range</Text>
                <TouchableOpacity style={styles.filterOption}>
                  <Text>Today</Text>
                </TouchableOpacity>
                <TouchableOpacity style={styles.filterOption}>
                  <Text>This Week</Text>
                </TouchableOpacity>
                <TouchableOpacity style={styles.filterOption}>
                  <Text>This Month</Text>
                </TouchableOpacity>
                <TouchableOpacity style={styles.filterOption}>
                  <Text>Custom Range</Text>
                </TouchableOpacity>
              </View>

              {/* Radius Filter */}
              <View style={styles.filterSection}>
                <Text style={styles.filterLabel}>Search Radius</Text>
                <TouchableOpacity style={styles.filterOption}>
                  <Text>1 km</Text>
                </TouchableOpacity>
                <TouchableOpacity style={styles.filterOption}>
                  <Text>5 km</Text>
                </TouchableOpacity>
                <TouchableOpacity style={styles.filterOption}>
                  <Text>10 km</Text>
                </TouchableOpacity>
                <TouchableOpacity style={styles.filterOption}>
                  <Text>50 km</Text>
                </TouchableOpacity>
              </View>
            </ScrollView>

            <TouchableOpacity
              style={styles.applyButton}
              onPress={handleFilterApply}
            >
              <Text style={styles.applyButtonText}>Apply Filters</Text>
            </TouchableOpacity>
          </View>
        </View>
      </Modal>
    </>
  );
};

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    padding: 12,
    backgroundColor: '#fff',
    borderBottomWidth: 1,
    borderBottomColor: '#eee',
  },
  searchBar: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#f5f5f5',
    borderRadius: 8,
    paddingHorizontal: 12,
    height: 40,
  },
  input: {
    flex: 1,
    marginLeft: 8,
    fontSize: 16,
  },
  filterButton: {
    marginLeft: 12,
    justifyContent: 'center',
    alignItems: 'center',
    width: 40,
    height: 40,
    borderRadius: 8,
    backgroundColor: '#f5f5f5',
  },
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0,0,0,0.5)',
    justifyContent: 'flex-end',
  },
  modalContent: {
    backgroundColor: '#fff',
    borderTopLeftRadius: 20,
    borderTopRightRadius: 20,
    maxHeight: '80%',
  },
  modalHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 20,
    borderBottomWidth: 1,
    borderBottomColor: '#eee',
  },
  modalTitle: {
    fontSize: 20,
    fontWeight: '600',
  },
  filterList: {
    padding: 20,
  },
  filterSection: {
    marginBottom: 24,
  },
  filterLabel: {
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 12,
  },
  filterOption: {
    padding: 12,
    backgroundColor: '#f5f5f5',
    borderRadius: 8,
    marginBottom: 8,
  },
  applyButton: {
    margin: 20,
    backgroundColor: '#667eea',
    padding: 16,
    borderRadius: 12,
    alignItems: 'center',
  },
  applyButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
});
```

---

## 💾 Local Caching

### Cache Manager

**Create `mobile/src/utils/cacheManager.ts`:**

```typescript
import AsyncStorage from '@react-native-async-storage/async-storage';
import * as FileSystem from 'expo-file-system';

interface CacheOptions {
  ttl?: number; // Time to live in seconds, default 1 hour
}

export class CacheManager {
  private static readonly CACHE_PREFIX = '@nowink_cache:';
  private static readonly IMAGE_CACHE_DIR = `${FileSystem.cacheDirectory}images/`;

  static async initialize() {
    const dirInfo = await FileSystem.getInfoAsync(this.IMAGE_CACHE_DIR);
    if (!dirInfo.exists) {
      await FileSystem.makeDirectoryAsync(this.IMAGE_CACHE_DIR, {
        intermediates: true,
      });
    }
  }

  // Generic data caching
  static async set(key: string, data: any, options: CacheOptions = {}) {
    const { ttl = 3600 } = options;
    const cacheItem = {
      data,
      timestamp: Date.now(),
      ttl: ttl * 1000,
    };

    try {
      await AsyncStorage.setItem(
        `${this.CACHE_PREFIX}${key}`,
        JSON.stringify(cacheItem)
      );
    } catch (error) {
      console.error('Cache set error:', error);
    }
  }

  static async get<T>(key: string): Promise<T | null> {
    try {
      const item = await AsyncStorage.getItem(`${this.CACHE_PREFIX}${key}`);
      if (!item) return null;

      const cacheItem = JSON.parse(item);
      const age = Date.now() - cacheItem.timestamp;

      if (age > cacheItem.ttl) {
        // Cache expired
        await this.delete(key);
        return null;
      }

      return cacheItem.data;
    } catch (error) {
      console.error('Cache get error:', error);
      return null;
    }
  }

  static async delete(key: string) {
    try {
      await AsyncStorage.removeItem(`${this.CACHE_PREFIX}${key}`);
    } catch (error) {
      console.error('Cache delete error:', error);
    }
  }

  static async clear() {
    try {
      const keys = await AsyncStorage.getAllKeys();
      const cacheKeys = keys.filter((key) => key.startsWith(this.CACHE_PREFIX));
      await AsyncStorage.multiRemove(cacheKeys);
    } catch (error) {
      console.error('Cache clear error:', error);
    }
  }

  // Image caching
  static async cacheImage(url: string): Promise<string> {
    const filename = url.split('/').pop() || `image_${Date.now()}`;
    const localUri = `${this.IMAGE_CACHE_DIR}${filename}`;

    // Check if already cached
    const fileInfo = await FileSystem.getInfoAsync(localUri);
    if (fileInfo.exists) {
      return localUri;
    }

    // Download and cache
    try {
      const downloadResult = await FileSystem.downloadAsync(url, localUri);
      return downloadResult.uri;
    } catch (error) {
      console.error('Image cache error:', error);
      return url; // Fallback to original URL
    }
  }

  static async clearImageCache() {
    try {
      await FileSystem.deleteAsync(this.IMAGE_CACHE_DIR, { idempotent: true });
      await FileSystem.makeDirectoryAsync(this.IMAGE_CACHE_DIR, {
        intermediates: true,
      });
    } catch (error) {
      console.error('Image cache clear error:', error);
    }
  }

  // NFT data caching
  static async cacheNFT(mintAddress: string, data: any) {
    await this.set(`nft:${mintAddress}`, data, { ttl: 3600 }); // 1 hour
  }

  static async getCachedNFT(mintAddress: string) {
    return this.get(`nft:${mintAddress}`);
  }

  // Feed caching
  static async cacheFeed(feed: any[]) {
    await this.set('feed', feed, { ttl: 300 }); // 5 minutes
  }

  static async getCachedFeed() {
    return this.get<any[]>('feed');
  }
}
```

---

## 📤 Sharing

### Share Component

**Create `mobile/src/components/ShareButton.tsx`:**

```typescript
import React from 'react';
import { TouchableOpacity, Share, Platform, Alert, StyleSheet } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import * as Clipboard from 'expo-clipboard';

interface ShareButtonProps {
  nft: {
    mint_address: string;
    name: string;
    creator_username?: string;
  };
  style?: any;
}

export const ShareButton: React.FC<ShareButtonProps> = ({ nft, style }) => {
  const handleShare = async () => {
    const shareUrl = `https://now.ink/nft/${nft.mint_address}`;
    const message = `Check out "${nft.name}" by @${nft.creator_username || 'unknown'} on now.ink!`;

    try {
      const result = await Share.share({
        message: Platform.OS === 'ios' ? message : `${message}\n${shareUrl}`,
        url: shareUrl, // iOS only
        title: nft.name,
      });

      if (result.action === Share.sharedAction) {
        if (result.activityType) {
          console.log('Shared via:', result.activityType);
        } else {
          console.log('Shared successfully');
        }
      } else if (result.action === Share.dismissedAction) {
        console.log('Share dismissed');
      }
    } catch (error) {
      console.error('Share error:', error);
      Alert.alert('Error', 'Failed to share. Link copied to clipboard.');
      await Clipboard.setStringAsync(shareUrl);
    }
  };

  return (
    <TouchableOpacity style={[styles.button, style]} onPress={handleShare}>
      <Ionicons name="share-outline" size={24} color="#667eea" />
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  button: {
    padding: 8,
  },
});
```

---

## ⚡ Performance Optimizations

### Optimized FlatList

**Create `mobile/src/components/OptimizedFlatList.tsx`:**

```typescript
import React, { memo } from 'react';
import { FlatList, FlatListProps } from 'react-native';

interface OptimizedFlatListProps<T> extends FlatListProps<T> {
  data: T[];
}

function OptimizedFlatListComponent<T>({
  data,
  renderItem,
  keyExtractor,
  ...props
}: OptimizedFlatListProps<T>) {
  return (
    <FlatList
      data={data}
      renderItem={renderItem}
      keyExtractor={keyExtractor}
      // Performance optimizations
      removeClippedSubviews={true}
      maxToRenderPerBatch={10}
      updateCellsBatchingPeriod={50}
      initialNumToRender={10}
      windowSize={10}
      getItemLayout={
        props.getItemLayout ||
        ((data, index) => ({
          length: 200, // Approximate item height
          offset: 200 * index,
          index,
        }))
      }
      {...props}
    />
  );
}

export const OptimizedFlatList = memo(OptimizedFlatListComponent) as typeof OptimizedFlatListComponent;
```

---

## 📚 Usage Examples

### Initialize in App.tsx

```typescript
import { CacheManager } from './src/utils/cacheManager';
import { ThumbnailGenerator } from './src/utils/thumbnailGenerator';

useEffect(() => {
  async function initializeApp() {
    await CacheManager.initialize();
    await ThumbnailGenerator.initialize();
  }
  initializeApp();
}, []);
```

### Feed Screen with All Features

```typescript
import React, { useState, useEffect } from 'react';
import { View } from 'react-native';
import { OptimizedFlatList } from '../components/OptimizedFlatList';
import { VideoThumbnail } from '../components/VideoThumbnail';
import { ShareButton } from '../components/ShareButton';
import { CacheManager } from '../utils/cacheManager';
import { api } from '../utils/api';

export const EnhancedFeedScreen = () => {
  const [feed, setFeed] = useState([]);
  const [refreshing, setRefreshing] = useState(false);

  useEffect(() => {
    loadFeed();
  }, []);

  const loadFeed = async () => {
    // Try cache first
    const cached = await CacheManager.getCachedFeed();
    if (cached) {
      setFeed(cached);
    }

    // Fetch fresh data
    const response = await api.getFeed(20, 0);
    setFeed(response.feed);
    await CacheManager.cacheFeed(response.feed);
  };

  const renderItem = ({ item }: any) => (
    <View>
      <VideoThumbnail videoUri={item.video_url} />
      <ShareButton nft={item} />
    </View>
  );

  return (
    <OptimizedFlatList
      data={feed}
      renderItem={renderItem}
      keyExtractor={(item) => item.mint_address}
      onRefresh={loadFeed}
      refreshing={refreshing}
    />
  );
};
```

---

## ✅ Implementation Checklist

- [ ] Install required dependencies
- [ ] Implement thumbnail generation
- [ ] Add pull-to-refresh to all screens
- [ ] Implement map search component
- [ ] Add map filters (date, radius, users)
- [ ] Implement cache manager
- [ ] Add share functionality
- [ ] Optimize FlatList performance
- [ ] Test on iOS and Android
- [ ] Measure performance improvements

---

**Created**: November 5, 2025  
**Version**: 1.0.0  
**Status**: Implementation Ready
