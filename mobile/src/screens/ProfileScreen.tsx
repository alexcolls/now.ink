import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  ScrollView,
  ActivityIndicator,
  Image,
  Alert,
} from 'react-native';
import { useWallet } from '../context/WalletContext';
import { apiClient } from '/home/quantium/labs/now.ink/mobile/src/services/api';
import VideoPlayer from '/home/quantium/labs/now.ink/mobile/src/components/VideoPlayer';

interface ProfileScreenProps {
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

export default function ProfileScreen({ onBack }: ProfileScreenProps) {
  const { publicKey, isConnected } = useWallet();
  const [nfts, setNfts] = useState<NFT[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedNFT, setSelectedNFT] = useState<NFT | null>(null);

  useEffect(() => {
    if (isConnected) {
      loadUserNFTs();
    }
  }, [isConnected]);

  const loadUserNFTs = async () => {
    try {
      setLoading(true);
      const response = await apiClient.listNFTs();
      // Filter NFTs by current user's wallet
      const userWallet = publicKey?.toBase58();
      const userNFTs = response.nfts?.filter(
        (nft: NFT) => nft.creator_wallet === userWallet
      ) || [];
      setNfts(userNFTs);
    } catch (error: any) {
      Alert.alert('Error', error.message || 'Failed to load NFTs');
    } finally {
      setLoading(false);
    }
  };

  const getArweaveVideoUrl = (arweaveTxId: string) => {
    if (arweaveTxId.startsWith('MOCK_')) {
      return 'https://www.w3schools.com/html/mov_bbb.mp4';
    }
    return `https://arweave.net/${arweaveTxId}`;
  };

  const formatDate = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
    });
  };

  return (
    <View style={styles.container}>
      <TouchableOpacity style={styles.backButton} onPress={onBack}>
        <Text style={styles.backButtonText}>← Back</Text>
      </TouchableOpacity>

      <ScrollView style={styles.scrollView}>
        {/* Header */}
        <View style={styles.header}>
          <View style={styles.avatarContainer}>
            <Text style={styles.avatarEmoji}>👤</Text>
          </View>
          
          {isConnected ? (
            <>
              <Text style={styles.walletText}>
                {publicKey?.toBase58().slice(0, 8)}...{publicKey?.toBase58().slice(-8)}
              </Text>
              <View style={styles.stats}>
                <View style={styles.statItem}>
                  <Text style={styles.statValue}>{nfts.length}</Text>
                  <Text style={styles.statLabel}>Moments</Text>
                </View>
                <View style={styles.statDivider} />
                <View style={styles.statItem}>
                  <Text style={styles.statValue}>0</Text>
                  <Text style={styles.statLabel}>Followers</Text>
                </View>
                <View style={styles.statDivider} />
                <View style={styles.statItem}>
                  <Text style={styles.statValue}>0</Text>
                  <Text style={styles.statLabel}>Following</Text>
                </View>
              </View>
            </>
          ) : (
            <Text style={styles.notConnectedText}>
              Connect your wallet to view profile
            </Text>
          )}
        </View>

        {/* NFT Grid */}
        {isConnected && (
          <View style={styles.nftSection}>
            <Text style={styles.sectionTitle}>Your Moments</Text>
            
            {loading ? (
              <View style={styles.loadingContainer}>
                <ActivityIndicator size="large" color="#007AFF" />
              </View>
            ) : nfts.length === 0 ? (
              <View style={styles.emptyContainer}>
                <Text style={styles.emptyEmoji}>🎥</Text>
                <Text style={styles.emptyText}>
                  No moments yet{' \n'}
                  Record your first moment!
                </Text>
              </View>
            ) : (
              <View style={styles.nftGrid}>
                {nfts.map((nft) => (
                  <TouchableOpacity
                    key={nft.mint_address}
                    style={styles.nftCard}
                    onPress={() => setSelectedNFT(nft)}
                  >
                    <View style={styles.nftThumbnail}>
                      <Text style={styles.thumbnailIcon}>📹</Text>
                    </View>
                    <View style={styles.nftInfo}>
                      <Text style={styles.nftTitle} numberOfLines={1}>
                        {nft.title}
                      </Text>
                      <Text style={styles.nftDate}>
                        {formatDate(nft.timestamp)}
                      </Text>
                    </View>
                  </TouchableOpacity>
                ))}
              </View>
            )}
          </View>
        )}
      </ScrollView>

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
    marginTop: 60,
    marginLeft: 20,
    marginBottom: 20,
  },
  backButtonText: {
    color: '#007AFF',
    fontSize: 18,
    fontWeight: '600',
  },
  scrollView: {
    flex: 1,
  },
  header: {
    alignItems: 'center',
    padding: 20,
    borderBottomWidth: 1,
    borderBottomColor: '#222',
  },
  avatarContainer: {
    width: 100,
    height: 100,
    borderRadius: 50,
    backgroundColor: '#1a1a1a',
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 16,
  },
  avatarEmoji: {
    fontSize: 50,
  },
  walletText: {
    fontSize: 14,
    color: '#888',
    fontFamily: 'monospace',
    marginBottom: 20,
  },
  notConnectedText: {
    fontSize: 14,
    color: '#666',
    textAlign: 'center',
  },
  stats: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    width: '100%',
    paddingVertical: 16,
  },
  statItem: {
    alignItems: 'center',
    flex: 1,
  },
  statValue: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#fff',
    marginBottom: 4,
  },
  statLabel: {
    fontSize: 12,
    color: '#666',
  },
  statDivider: {
    width: 1,
    backgroundColor: '#222',
  },
  nftSection: {
    padding: 20,
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#fff',
    marginBottom: 16,
  },
  loadingContainer: {
    padding: 40,
    alignItems: 'center',
  },
  emptyContainer: {
    padding: 40,
    alignItems: 'center',
  },
  emptyEmoji: {
    fontSize: 60,
    marginBottom: 16,
  },
  emptyText: {
    fontSize: 14,
    color: '#666',
    textAlign: 'center',
    lineHeight: 20,
  },
  nftGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    marginHorizontal: -8,
  },
  nftCard: {
    width: '50%',
    padding: 8,
  },
  nftThumbnail: {
    aspectRatio: 1,
    backgroundColor: '#1a1a1a',
    borderRadius: 12,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 8,
  },
  thumbnailIcon: {
    fontSize: 40,
  },
  nftInfo: {
    paddingHorizontal: 4,
  },
  nftTitle: {
    fontSize: 14,
    fontWeight: '600',
    color: '#fff',
    marginBottom: 4,
  },
  nftDate: {
    fontSize: 12,
    color: '#666',
  },
});
