import { StatusBar } from 'expo-status-bar';
import { StyleSheet, View } from 'react-native';
import { WalletProvider } from './src/context/WalletContext';
import CameraScreen from './src/screens/CameraScreen';

export default function App() {
  return (
    <WalletProvider>
      <View style={styles.container}>
        <CameraScreen />
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
