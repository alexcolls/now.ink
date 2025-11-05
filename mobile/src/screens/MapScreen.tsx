import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';

interface MapScreenProps {
  onBack: () => void;
}

export default function MapScreen({ onBack }: MapScreenProps) {
  return (
    <View style={styles.container}>
      <TouchableOpacity style={styles.backButton} onPress={onBack}>
        <Text style={styles.backButtonText}>← Back</Text>
      </TouchableOpacity>
      
      <View style={styles.content}>
        <Text style={styles.emoji}>🗺️</Text>
        <Text style={styles.title}>Map View</Text>
        <Text style={styles.subtitle}>Coming Soon</Text>
        <Text style={styles.description}>
          Discover geo-tagged NFT moments{'\n'}
          on an interactive map with{'\n'}
          time-slider filtering
        </Text>
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
