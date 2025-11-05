import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  ActivityIndicator,
  Alert,
} from 'react-native';
import MapView, { Marker, PROVIDER_GOOGLE } from 'react-native-maps';
import { apiClient } from '/home/quantium/labs/now.ink/mobile/src/services/api';
import VideoPlayer from '/home/quantium/labs/now.ink/mobile/src/components/VideoPlayer';

interface MapScreenProps {
  onBack: () => void;
}

interface NFT {
  mint_address: string;
  metadata_uri: string;
  title: string;
  creator_wallet: string;
  latitude: number;
  longitude: number;
  timestamp: string;
  arweave_tx_id: string;
}

export default function MapScreen({ onBack }: MapScreenProps) {
  const [nfts, setNfts] = useState<NFT[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedNFT, setSelectedNFT] = useState<NFT | null>(null);
  const [region, setRegion] = useState({
    latitude: 40.7128,
    longitude: -74.006,
    latitudeDelta: 0.5,
    longitudeDelta: 0.5,
  });

  useEffect(() => {
    loadNFTs();
  }, []);

  const loadNFTs = async () => {
    try {
      setLoading(true);
      const response = await apiClient.listNFTs();
      if (response.nfts && response.nfts.length > 0) {
        setNfts(response.nfts);
        // Center on first NFT
        const firstNFT = response.nfts[0];
        setRegion({
          latitude: firstNFT.latitude,
          longitude: firstNFT.longitude,
          latitudeDelta: 0.5,
          longitudeDelta: 0.5,
        });
      }
    } catch (error: any) {
      Alert.alert('Error', error.message || 'Failed to load NFTs');
    } finally {
      setLoading(false);
    }
  };

  const handleMarkerPress = (nft: NFT) => {
    setSelectedNFT(nft);
  };

  const getArweaveVideoUrl = (arweaveTxId: string) => {
    // In mock mode, arweave_tx_id might be MOCK_AR_xxx
    if (arweaveTxId.startsWith('MOCK_')) {
      // Return a test video URL for mock data
      return 'https://www.w3schools.com/html/mov_bbb.mp4';
    }
    return `https://arweave.net/${arweaveTxId}`;
  };

  return (
    <View style={styles.container}>
      <TouchableOpacity style={styles.backButton} onPress={onBack}>
        <Text style={styles.backButtonText}>← Back</Text>
      </TouchableOpacity>

      {loading ? (
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color="#007AFF" />
          <Text style={styles.loadingText}>Loading NFT moments...</Text>
        </View>
      ) : nfts.length === 0 ? (
        <View style={styles.emptyContainer}>
          <Text style={styles.emoji}>🗺️</Text>
          <Text style={styles.emptyTitle}>No NFTs Yet</Text>
          <Text style={styles.emptyText}>
            Record your first moment to see it on the map!
          </Text>
        </View>
      ) : (
        <>
          <MapView
            provider={PROVIDER_GOOGLE}
            style={styles.map}
            initialRegion={region}
            showsUserLocation
            showsMyLocationButton
          >
            {nfts.map((nft) => (
              <Marker
                key={nft.mint_address}
                coordinate={{
                  latitude: nft.latitude,
                  longitude: nft.longitude,
                }}
                title={nft.title}
                description={`by ${nft.creator_wallet.slice(0, 8)}...`}
                onPress={() => handleMarkerPress(nft)}
              >
                <View style={styles.customMarker}>
                  <Text style={styles.markerText}>📹</Text>
                </View>
              </Marker>
            ))}
          </MapView>

          <View style={styles.statsBar}>
            <Text style={styles.statsText}>{nfts.length} moments on map</Text>
            <TouchableOpacity onPress={loadNFTs}>
              <Text style={styles.refreshText}>🔄 Refresh</Text>
            </TouchableOpacity>
          </View>
        </>
      )}

      {selectedNFT && (
        <VideoPlayer
          nft={selectedNFT}
          videoUrl={getArweaveVideoUrl(selectedNFT.arweave_tx_id)}
          onClose={() => setSelectedNFT(null)}
        />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#000',
  },
  backButton: {
    position: 'absolute',
    top: 60,
    left: 20,
    zIndex: 10,
    backgroundColor: 'rgba(0,0,0,0.6)',
    paddingHorizontal: 12,
    paddingVertical: 8,
    borderRadius: 8,
  },
  backButtonText: {
    color: '#fff',
    fontSize: 18,
    fontWeight: '600',
  },
  map: {
    flex: 1,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  loadingText: {
    color: '#888',
    marginTop: 12,
    fontSize: 14,
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 40,
  },
  emoji: {
    fontSize: 80,
    marginBottom: 20,
  },
  emptyTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#fff',
    marginBottom: 8,
  },
  emptyText: {
    fontSize: 14,
    color: '#666',
    textAlign: 'center',
    lineHeight: 20,
  },
  customMarker: {
    width: 40,
    height: 40,
    backgroundColor: 'rgba(0,122,255,0.8)',
    borderRadius: 20,
    justifyContent: 'center',
    alignItems: 'center',
    borderWidth: 3,
    borderColor: '#fff',
  },
  markerText: {
    fontSize: 20,
  },
  statsBar: {
    position: 'absolute',
    bottom: 40,
    left: 20,
    right: 20,
    backgroundColor: 'rgba(0,0,0,0.8)',
    padding: 16,
    borderRadius: 12,
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  statsText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
  refreshText: {
    color: '#007AFF',
    fontSize: 14,
    fontWeight: '600',
  },
});
