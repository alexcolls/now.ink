import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { useWallet } from '../context/WalletContext';

interface ProfileScreenProps {
  onBack: () => void;
}

export default function ProfileScreen({ onBack }: ProfileScreenProps) {
  const { publicKey, isConnected } = useWallet();

  return (
    <View style={styles.container}>
      <TouchableOpacity style={styles.backButton} onPress={onBack}>
        <Text style={styles.backButtonText}>← Back</Text>
      </TouchableOpacity>
      
      <View style={styles.content}>
        <Text style={styles.emoji}>👤</Text>
        <Text style={styles.title}>Profile</Text>
        
        {isConnected ? (
          <>
            <Text style={styles.wallet}>
              {publicKey?.toBase58().slice(0, 16)}...
            </Text>
            <Text style={styles.subtitle}>Coming Soon</Text>
            <Text style={styles.description}>
              View your minted NFTs{'\n'}
              See followers & following{'\n'}
              Edit profile settings
            </Text>
          </>
        ) : (
          <Text style={styles.description}>
            Connect your wallet to view profile
          </Text>
        )}
      </View>
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
  },
  backButtonText: {
    color: '#007AFF',
    fontSize: 18,
  },
  content: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 40,
  },
  emoji: {
    fontSize: 80,
    marginBottom: 20,
  },
  title: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#fff',
    marginBottom: 8,
  },
  wallet: {
    fontSize: 14,
    color: '#888',
    fontFamily: 'monospace',
    marginBottom: 20,
  },
  subtitle: {
    fontSize: 18,
    color: '#888',
    marginBottom: 20,
  },
  description: {
    fontSize: 14,
    color: '#666',
    textAlign: 'center',
    lineHeight: 20,
  },
});
