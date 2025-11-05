import { StatusBar } from 'expo-status-bar';
import { StyleSheet, View } from 'react-native';
import { useState } from 'react';
import { WalletProvider } from './src/context/WalletContext';
import HomeScreen from './src/screens/HomeScreen';
import CameraScreen from './src/screens/CameraScreen';
import MapScreen from './src/screens/MapScreen';
import ProfileScreen from './src/screens/ProfileScreen';

type Screen = 'home' | 'camera' | 'map' | 'profile';

export default function App() {
  const [currentScreen, setCurrentScreen] = useState<Screen>('home');

  const renderScreen = () => {
    switch (currentScreen) {
      case 'home':
        return (
          <HomeScreen
            onNavigateToCamera={() => setCurrentScreen('camera')}
            onNavigateToMap={() => setCurrentScreen('map')}
            onNavigateToProfile={() => setCurrentScreen('profile')}
          />
        );
      case 'camera':
        return <CameraScreen />;
      case 'map':
        return <MapScreen onBack={() => setCurrentScreen('home')} />;
      case 'profile':
        return <ProfileScreen onBack={() => setCurrentScreen('home')} />;
    }
  };

  return (
    <WalletProvider>
      <View style={styles.container}>
        {renderScreen()}
        <StatusBar style="light" />
      </View>
    </WalletProvider>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#000',
  },
});
