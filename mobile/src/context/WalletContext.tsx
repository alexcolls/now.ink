import React, { createContext, useContext, useState, useCallback, ReactNode } from 'react';
import { transact } from '@solana-mobile/mobile-wallet-adapter-protocol';
import { Connection, PublicKey } from '@solana/web3.js';

const SOLANA_RPC_URL = process.env.SOLANA_RPC_URL || 'https://api.devnet.solana.com';

interface WalletContextType {
  publicKey: PublicKey | null;
  isConnected: boolean;
  isConnecting: boolean;
  connect: () => Promise<void>;
  disconnect: () => void;
  signMessage: (message: string) => Promise<string>;
}

const WalletContext = createContext<WalletContextType | undefined>(undefined);

export const useWallet = () => {
  const context = useContext(WalletContext);
  if (!context) {
    throw new Error('useWallet must be used within WalletProvider');
  }
  return context;
};

interface WalletProviderProps {
  children: ReactNode;
}

export const WalletProvider: React.FC<WalletProviderProps> = ({ children }) => {
  const [publicKey, setPublicKey] = useState<PublicKey | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [isConnecting, setIsConnecting] = useState(false);

  const connect = useCallback(async () => {
    try {
      setIsConnecting(true);
      
      const result = await transact(async (wallet) => {
        // Request authorization from wallet
        const authResult = await wallet.authorize({
          cluster: 'devnet',
          identity: {
            name: 'now.ink',
            uri: 'https://now.ink',
            icon: 'favicon.ico',
          },
        });

        const address = authResult.accounts[0].address;
        return new PublicKey(address);
      });

      setPublicKey(result);
      setIsConnected(true);
      console.log('✅ Wallet connected:', result.toBase58());
    } catch (error) {
      console.error('❌ Wallet connection failed:', error);
      throw error;
    } finally {
      setIsConnecting(false);
    }
  }, []);

  const disconnect = useCallback(() => {
    setPublicKey(null);
    setIsConnected(false);
    console.log('🔌 Wallet disconnected');
  }, []);

  const signMessage = useCallback(async (message: string): Promise<string> => {
    if (!publicKey) {
      throw new Error('Wallet not connected');
    }

    try {
      const encodedMessage = new TextEncoder().encode(message);
      
      const signature = await transact(async (wallet) => {
        // Re-authorize if needed
        const authResult = await wallet.authorize({
          cluster: 'devnet',
          identity: {
            name: 'now.ink',
          },
        });

        // Sign the message
        const signResult = await wallet.signMessages({
          addresses: [publicKey.toBase58()],
          payloads: [encodedMessage],
        });

        return Buffer.from(signResult.signedPayloads[0]).toString('base64');
      });

      return signature;
    } catch (error) {
      console.error('❌ Message signing failed:', error);
      throw error;
    }
  }, [publicKey]);

  const value: WalletContextType = {
    publicKey,
    isConnected,
    isConnecting,
    connect,
    disconnect,
    signMessage,
  };

  return (
    <WalletContext.Provider value={value}>
      {children}
    </WalletContext.Provider>
  );
};
