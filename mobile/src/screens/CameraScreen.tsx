import React, { useState, useRef, useEffect } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Alert, ActivityIndicator } from 'react-native';
import { CameraView, CameraType, useCameraPermissions } from 'expo-camera';
import * as Location from 'expo-location';
import * as FileSystem from 'expo-file-system';
import { useWallet } from '../context/WalletContext';
import apiClient from '../services/api';

export default function CameraScreen() {
  const [permission, requestPermission] = useCameraPermissions();
  const [isRecording, setIsRecording] = useState(false);
  const [recordingDuration, setRecordingDuration] = useState(0);
  const [location, setLocation] = useState<Location.LocationObject | null>(null);
  const [cameraType, setCameraType] = useState<CameraType>('back');
  const [isMinting, setIsMinting] = useState(false);
  const [mintingProgress, setMintingProgress] = useState('');
  const [streamId, setStreamId] = useState<string | null>(null);
  const cameraRef = useRef<CameraView>(null);
  const { isConnected, publicKey } = useWallet();

  useEffect(() => {
    (async () => {
      const { status } = await Location.requestForegroundPermissionsAsync();
      if (status !== 'granted') {
        Alert.alert('Permission denied', 'Location permission is required to geo-tag moments');
      }
    })();
  }, []);

  useEffect(() => {
    let interval: NodeJS.Timeout;
    if (isRecording) {
      interval = setInterval(() => {
        setRecordingDuration(prev => prev + 1);
      }, 1000);
    }
    return () => clearInterval(interval);
  }, [isRecording]);

  const startRecording = async () => {
    if (!cameraRef.current) return;
    
    // Get current location
    try {
      const loc = await Location.getCurrentPositionAsync({});
      setLocation(loc);
      console.log('📍 Location:', loc.coords.latitude, loc.coords.longitude);
    } catch (error) {
      console.error('❌ Failed to get location:', error);
    }

    try {
      setIsRecording(true);
      setRecordingDuration(0);
      
      const video = await cameraRef.current.recordAsync();
      console.log('🎥 Video recorded:', video?.uri);
      
      // After recording stops, show save/discard dialog
      if (video) {
        handleRecordingComplete(video.uri);
      }
    } catch (error) {
      console.error('❌ Recording failed:', error);
      setIsRecording(false);
    }
  };

  const stopRecording = () => {
    if (cameraRef.current && isRecording) {
      cameraRef.current.stopRecording();
      setIsRecording(false);
    }
  };

  const handleRecordingComplete = (videoUri: string) => {
    Alert.alert(
      'Save as NFT?',
      `Do you want to mint this ${recordingDuration}s video as an NFT?`,
      [
        {
          text: 'Discard',
          style: 'cancel',
          onPress: () => {
            console.log('❌ Video discarded');
            // Delete local video
          }
        },
        {
          text: 'Save & Mint',
          onPress: () => handleMintNFT(videoUri)
        }
      ]
    );
  };

  const handleMintNFT = async (videoUri: string) => {
    if (!isConnected || !publicKey) {
      Alert.alert('Connect Wallet', 'Please connect your Solana wallet first');
      return;
    }

    if (!location) {
      Alert.alert('Location Required', 'Unable to get location for geo-tagging');
      return;
    }

    setIsMinting(true);
    setMintingProgress('Starting stream...');

    try {
      // 1. Start stream on backend
      console.log('🎬 Starting stream...');
      const streamResponse = await apiClient.startStream({
        title: `Moment at ${new Date().toLocaleTimeString()}`,
        latitude: location.coords.latitude,
        longitude: location.coords.longitude,
        is_public: true,
      });
      
      const newStreamId = streamResponse.id;
      setStreamId(newStreamId);
      console.log('✅ Stream started:', newStreamId);

      // 2. Upload video file
      setMintingProgress('Uploading video...');
      console.log('📤 Uploading video...');
      
      const videoInfo = await FileSystem.getInfoAsync(videoUri);
      if (!videoInfo.exists) {
        throw new Error('Video file not found');
      }

      const formData = new FormData();
      formData.append('video', {
        uri: videoUri,
        name: 'moment.mp4',
        type: 'video/mp4',
      } as any);

      // Upload to backend (backend will handle Arweave upload)
      await apiClient.saveStream(newStreamId, formData);
      console.log('✅ Video uploaded');

      // 3. Wait for NFT minting (backend does this)
      setMintingProgress('Minting NFT on Solana...');
      console.log('🪙 Minting NFT...');
      
      // Poll for stream status to get mint address
      let attempts = 0;
      const maxAttempts = 30; // 30 seconds
      let mintAddress = null;

      while (attempts < maxAttempts && !mintAddress) {
        await new Promise(resolve => setTimeout(resolve, 1000));
        const stream = await apiClient.getStream(newStreamId);
        if (stream.mint_address) {
          mintAddress = stream.mint_address;
          break;
        }
        attempts++;
      }

      setIsMinting(false);
      setMintingProgress('');
      
      if (mintAddress) {
        Alert.alert(
          'Success! 🎉',
          `Your moment has been minted!\n\nNFT: ${mintAddress.slice(0, 8)}...${mintAddress.slice(-8)}`,
          [{ text: 'OK', onPress: () => setStreamId(null) }]
        );
      } else {
        Alert.alert(
          'Processing...',
          'Your moment is being minted. Check back in your profile later!',
          [{ text: 'OK', onPress: () => setStreamId(null) }]
        );
      }
    } catch (error) {
      console.error('❌ Minting failed:', error);
      setIsMinting(false);
      setMintingProgress('');
      Alert.alert('Error', `Failed to mint NFT: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  };

  if (!permission) {
    return <View style={styles.container}><Text>Requesting permissions...</Text></View>;
  }

  if (!permission.granted) {
    return (
      <View style={styles.container}>
        <Text style={styles.message}>Camera permission is required</Text>
        <TouchableOpacity onPress={requestPermission} style={styles.button}>
          <Text style={styles.buttonText}>Grant Permission</Text>
        </TouchableOpacity>
      </View>
    );
  }

  const formatDuration = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  return (
    <View style={styles.container}>
      <CameraView
        ref={cameraRef}
        style={styles.camera}
        facing={cameraType}
      >
        {/* Recording indicator */}
        {isRecording && (
          <View style={styles.recordingIndicator}>
            <View style={styles.recordingDot} />
            <Text style={styles.recordingText}>
              REC {formatDuration(recordingDuration)}
            </Text>
          </View>
        )}

        {/* Minting progress overlay */}
        {isMinting && (
          <View style={styles.mintingOverlay}>
            <View style={styles.mintingModal}>
              <ActivityIndicator size="large" color="#fff" />
              <Text style={styles.mintingText}>{mintingProgress}</Text>
              <Text style={styles.mintingSubtext}>Please wait...</Text>
            </View>
          </View>
        )}

        {/* Controls */}
        <View style={styles.controls}>
          {/* Flip camera button */}
          <TouchableOpacity
            style={styles.controlButton}
            onPress={() => setCameraType(current => current === 'back' ? 'front' : 'back')}
          >
            <Text style={styles.controlButtonText}>🔄</Text>
          </TouchableOpacity>

          {/* Record button */}
          <TouchableOpacity
            style={[styles.recordButton, isRecording && styles.recordButtonActive]}
            onPress={isRecording ? stopRecording : startRecording}
          >
            <View style={[styles.recordButtonInner, isRecording && styles.recordButtonInnerActive]} />
          </TouchableOpacity>

          {/* Wallet status */}
          <View style={styles.controlButton}>
            <Text style={styles.controlButtonText}>
              {isConnected ? '🟢' : '🔴'}
            </Text>
          </View>
        </View>
      </CameraView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    backgroundColor: '#000',
  },
  message: {
    textAlign: 'center',
    color: '#fff',
    fontSize: 16,
    marginBottom: 20,
  },
  camera: {
    flex: 1,
  },
  recordingIndicator: {
    position: 'absolute',
    top: 60,
    left: 20,
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: 'rgba(255, 0, 0, 0.8)',
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 20,
  },
  recordingDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    backgroundColor: '#fff',
    marginRight: 8,
  },
  recordingText: {
    color: '#fff',
    fontWeight: 'bold',
    fontSize: 14,
  },
  controls: {
    position: 'absolute',
    bottom: 40,
    left: 0,
    right: 0,
    flexDirection: 'row',
    justifyContent: 'space-around',
    alignItems: 'center',
    paddingHorizontal: 40,
  },
  controlButton: {
    width: 50,
    height: 50,
    borderRadius: 25,
    backgroundColor: 'rgba(255, 255, 255, 0.3)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  controlButtonText: {
    fontSize: 24,
  },
  recordButton: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: 'rgba(255, 255, 255, 0.8)',
    justifyContent: 'center',
    alignItems: 'center',
    borderWidth: 4,
    borderColor: '#fff',
  },
  recordButtonActive: {
    backgroundColor: 'rgba(255, 0, 0, 0.8)',
  },
  recordButtonInner: {
    width: 60,
    height: 60,
    borderRadius: 30,
    backgroundColor: '#ff0000',
  },
  recordButtonInnerActive: {
    borderRadius: 8,
    width: 30,
    height: 30,
  },
  mintingOverlay: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: 'rgba(0, 0, 0, 0.9)',
    justifyContent: 'center',
    alignItems: 'center',
    zIndex: 999,
  },
  mintingModal: {
    backgroundColor: 'rgba(255, 255, 255, 0.1)',
    padding: 40,
    borderRadius: 20,
    alignItems: 'center',
    minWidth: 200,
  },
  mintingText: {
    color: '#fff',
    fontSize: 18,
    fontWeight: 'bold',
    marginTop: 20,
    textAlign: 'center',
  },
  mintingSubtext: {
    color: '#aaa',
    fontSize: 14,
    marginTop: 8,
    textAlign: 'center',
  },
  button: {
    backgroundColor: '#007AFF',
    paddingHorizontal: 20,
    paddingVertical: 12,
    borderRadius: 8,
    alignSelf: 'center',
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
});
