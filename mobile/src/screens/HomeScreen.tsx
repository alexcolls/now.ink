import React, { useState, useEffect } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Alert } from 'react-native';
import { useWallet } from '../context/WalletContext';
import { apiClient } from '../services/api';

interface HomeScreenProps {
  onNavigateToCamera: () => void;
  onNavigateToMap: () => void;
  onNavigateToProfile: () => void;
}

export default function HomeScreen({ 
  onNavigateToCamera, 
  onNavigateToMap,
  onNavigateToProfile 
}: HomeScreenProps) {
  const { publicKey, isConnected, isConnecting, connect, disconnect } = useWallet();
  const [apiStatus, setApiStatus] = useState<'checking' | 'online' | 'offline'>('checking');

  useEffect(() => {
    checkApiStatus();
  }, []);

  const checkApiStatus = async () => {
    try {
      await apiClient.health();
      setApiStatus('online');
    } catch (error) {
      setApiStatus('offline');
      console.log('⚠️  Backend API offline - continuing in offline mode');
    }
  };

  const handleConnect = async () => {
    try {
      await connect();
      Alert.alert('Success', 'Wallet connected! 🎉');
    } catch (error) {
      Alert.alert('Error', 'Failed to connect wallet');
    }
  };

  return (
    <View style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.logo}>now.ink</Text>
        <Text style={styles.tagline}>Your life, minted</Text>
      </View>

      {/* Wallet Status */}
      <View style={styles.walletSection}>
        {isConnected ? (
          <>
            <Text style={styles.walletLabel}>Connected</Text>
            <Text style={styles.walletAddress}>
              {publicKey?.toBase58().slice(0, 8)}...{publicKey?.toBase58().slice(-8)}
            </Text>
            <TouchableOpacity style={styles.disconnectButton} onPress={disconnect}>
              <Text style={styles.disconnectButtonText}>Disconnect</Text>
            </TouchableOpacity>
          </>
        ) : (
          <>
            <Text style={styles.walletLabel}>Connect your Solana wallet</Text>
            <TouchableOpacity 
              style={styles.connectButton} 
              onPress={handleConnect}
              disabled={isConnecting}
            >
              <Text style={styles.connectButtonText}>
                {isConnecting ? 'Connecting...' : 'Connect Phantom Wallet'}
              </Text>
            </TouchableOpacity>
          </>
        )}
      </View>

      {/* Main Actions */}
      <View style={styles.actions}>
        <TouchableOpacity 
          style={[styles.actionButton, styles.primaryAction]}
          onPress={onNavigateToCamera}
        >
          <Text style={styles.actionIcon}>📷</Text>
          <Text style={styles.actionTitle}>Record Moment</Text>
          <Text style={styles.actionSubtitle}>Capture & mint as NFT</Text>
        </TouchableOpacity>

        <TouchableOpacity 
          style={styles.actionButton}
          onPress={onNavigateToMap}
        >
          <Text style={styles.actionIcon}>🗺️</Text>
          <Text style={styles.actionTitle}>Explore Map</Text>
          <Text style={styles.actionSubtitle}>Discover moments</Text>
        </TouchableOpacity>

        <TouchableOpacity 
          style={styles.actionButton}
          onPress={onNavigateToProfile}
        >
          <Text style={styles.actionIcon}>👤</Text>
          <Text style={styles.actionTitle}>Profile</Text>
          <Text style={styles.actionSubtitle}>Your NFTs</Text>
        </TouchableOpacity>
      </View>

      {/* API Status */}
      <View style={styles.footer}>
        <Text style={styles.apiStatus}>
          API: {apiStatus === 'online' ? '🟢 Online' : apiStatus === 'checking' ? '🟡 Checking...' : '🔴 Offline'}
        </Text>
        <Text style={styles.version}>v0.1.0 • Solana Devnet</Text>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#000',
    padding: 20,
  },
  header: {
    marginTop: 60,
    marginBottom: 40,
    alignItems: 'center',
  },
  logo: {
    fontSize: 48,
    fontWeight: 'bold',
    color: '#fff',
    marginBottom: 8,
  },
  tagline: {
    fontSize: 16,
    color: '#888',
  },
  walletSection: {
    backgroundColor: '#111',
    borderRadius: 12,
    padding: 20,
    marginBottom: 30,
    alignItems: 'center',
  },
  walletLabel: {
    fontSize: 14,
    color: '#888',
    marginBottom: 12,
  },
  walletAddress: {
    fontSize: 16,
    color: '#fff',
    fontFamily: 'monospace',
    marginBottom: 16,
  },
  connectButton: {
    backgroundColor: '#007AFF',
    paddingHorizontal: 24,
    paddingVertical: 12,
    borderRadius: 8,
  },
  connectButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
  disconnectButton: {
    paddingHorizontal: 24,
    paddingVertical: 12,
  },
  disconnectButtonText: {
    color: '#ff4444',
    fontSize: 14,
  },
  actions: {
    flex: 1,
  },
  actionButton: {
    backgroundColor: '#111',
    borderRadius: 12,
    padding: 20,
    marginBottom: 16,
    alignItems: 'center',
  },
  primaryAction: {
    backgroundColor: '#007AFF',
  },
  actionIcon: {
    fontSize: 48,
    marginBottom: 12,
  },
  actionTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#fff',
    marginBottom: 4,
  },
  actionSubtitle: {
    fontSize: 14,
    color: '#888',
  },
  footer: {
    alignItems: 'center',
    paddingVertical: 20,
  },
  apiStatus: {
    fontSize: 12,
    color: '#888',
    marginBottom: 4,
  },
  version: {
    fontSize: 12,
    color: '#555',
  },
});
