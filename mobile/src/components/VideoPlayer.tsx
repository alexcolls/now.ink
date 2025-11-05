import React, { useState, useRef } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  ActivityIndicator,
  Dimensions,
  Modal,
} from 'react-native';
import { Video, ResizeMode, AVPlaybackStatus } from 'expo-av';

const { width: SCREEN_WIDTH, height: SCREEN_HEIGHT } = Dimensions.get('window');

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

interface VideoPlayerProps {
  nft: NFT;
  videoUrl: string;
  onClose: () => void;
}

export default function VideoPlayer({ nft, videoUrl, onClose }: VideoPlayerProps) {
  const [isLoading, setIsLoading] = useState(true);
  const [isPlaying, setIsPlaying] = useState(true);
  const [showControls, setShowControls] = useState(true);
  const videoRef = useRef<Video>(null);

  const handlePlaybackStatusUpdate = (status: AVPlaybackStatus) => {
    if (status.isLoaded) {
      setIsLoading(false);
      setIsPlaying(status.isPlaying);
    }
  };

  const togglePlayPause = async () => {
    if (videoRef.current) {
      if (isPlaying) {
        await videoRef.current.pauseAsync();
      } else {
        await videoRef.current.playAsync();
      }
    }
  };

  const formatDate = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  const formatWallet = (wallet: string) => {
    return `${wallet.slice(0, 8)}...${wallet.slice(-8)}`;
  };

  return (
    <Modal visible={true} animationType="fade" onRequestClose={onClose}>
      <View style={styles.container}>
        {/* Video */}
        <TouchableOpacity
          activeOpacity={1}
          onPress={() => setShowControls(!showControls)}
          style={styles.videoContainer}
        >
          <Video
            ref={videoRef}
            source={{ uri: videoUrl }}
            style={styles.video}
            resizeMode={ResizeMode.CONTAIN}
            shouldPlay
            isLooping
            onPlaybackStatusUpdate={handlePlaybackStatusUpdate}
          />

          {isLoading && (
            <View style={styles.loadingOverlay}>
              <ActivityIndicator size="large" color="#fff" />
              <Text style={styles.loadingText}>Loading video...</Text>
            </View>
          )}

          {showControls && !isLoading && (
            <View style={styles.controlsOverlay}>
              <TouchableOpacity style={styles.playButton} onPress={togglePlayPause}>
                <Text style={styles.playButtonText}>
                  {isPlaying ? '⏸' : '▶'}
                </Text>
              </TouchableOpacity>
            </View>
          )}
        </TouchableOpacity>

        {/* NFT Info */}
        <View style={styles.infoContainer}>
          <View style={styles.infoHeader}>
            <View style={styles.infoLeft}>
              <Text style={styles.title}>{nft.title}</Text>
              <Text style={styles.creator}>by {formatWallet(nft.creator_wallet)}</Text>
              <Text style={styles.date}>📅 {formatDate(nft.timestamp)}</Text>
              <Text style={styles.location}>
                📍 {nft.latitude.toFixed(4)}, {nft.longitude.toFixed(4)}
              </Text>
            </View>
          </View>

          <View style={styles.metadata}>
            <View style={styles.metadataRow}>
              <Text style={styles.metadataLabel}>Mint Address</Text>
              <Text style={styles.metadataValue}>
                {formatWallet(nft.mint_address)}
              </Text>
            </View>
            <View style={styles.metadataRow}>
              <Text style={styles.metadataLabel}>Arweave TX</Text>
              <Text style={styles.metadataValue}>
                {nft.arweave_tx_id.slice(0, 16)}...
              </Text>
            </View>
          </View>

          <TouchableOpacity style={styles.closeButton} onPress={onClose}>
            <Text style={styles.closeButtonText}>Close</Text>
          </TouchableOpacity>
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#000',
  },
  videoContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  video: {
    width: SCREEN_WIDTH,
    height: SCREEN_HEIGHT * 0.6,
  },
  loadingOverlay: {
    ...StyleSheet.absoluteFillObject,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: 'rgba(0,0,0,0.8)',
  },
  loadingText: {
    color: '#fff',
    marginTop: 12,
    fontSize: 14,
  },
  controlsOverlay: {
    ...StyleSheet.absoluteFillObject,
    justifyContent: 'center',
    alignItems: 'center',
  },
  playButton: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: 'rgba(0,0,0,0.6)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  playButtonText: {
    fontSize: 36,
    color: '#fff',
  },
  infoContainer: {
    backgroundColor: '#111',
    padding: 20,
    borderTopLeftRadius: 20,
    borderTopRightRadius: 20,
  },
  infoHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 16,
  },
  infoLeft: {
    flex: 1,
  },
  title: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#fff',
    marginBottom: 8,
  },
  creator: {
    fontSize: 14,
    color: '#888',
    marginBottom: 8,
    fontFamily: 'monospace',
  },
  date: {
    fontSize: 12,
    color: '#666',
    marginBottom: 4,
  },
  location: {
    fontSize: 12,
    color: '#666',
  },
  metadata: {
    marginBottom: 20,
  },
  metadataRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingVertical: 8,
    borderBottomWidth: 1,
    borderBottomColor: '#222',
  },
  metadataLabel: {
    fontSize: 12,
    color: '#888',
  },
  metadataValue: {
    fontSize: 12,
    color: '#fff',
    fontFamily: 'monospace',
  },
  closeButton: {
    backgroundColor: '#007AFF',
    paddingVertical: 14,
    borderRadius: 12,
    alignItems: 'center',
  },
  closeButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
});
