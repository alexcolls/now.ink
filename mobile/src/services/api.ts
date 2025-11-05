import axios from 'axios';

const API_URL = process.env.API_URL || 'http://localhost:8080/api/v1';

// Create axios instance
const api = axios.create({
  baseURL: API_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Types
interface StartStreamRequest {
  title: string;
  latitude: number;
  longitude: number;
  is_public: boolean;
}

interface Stream {
  id: string;
  user_id: string;
  title: string;
  is_live: boolean;
  started_at: string;
  latitude: number;
  longitude: number;
}

interface MintResponse {
  stream_id: string;
  mint: {
    mint_address: string;
    metadata_uri: string;
    arweave_hash: string;
    status: string;
  };
  message: string;
}

// API Methods
export const apiClient = {
  // Auth
  getNonce: async (walletAddress: string) => {
    const { data } = await api.post('/auth/nonce', { wallet_address: walletAddress });
    return data;
  },

  verifyWallet: async (walletAddress: string, signature: string, nonce: string) => {
    const { data } = await api.post('/auth/verify', {
      wallet_address: walletAddress,
      signature,
      nonce,
    });
    return data;
  },

  // Streams
  startStream: async (req: StartStreamRequest): Promise<Stream> => {
    const { data } = await api.post('/streams/start', req);
    return data;
  },

  endStream: async (streamId: string): Promise<Stream> => {
    const { data } = await api.post(`/streams/${streamId}/end`);
    return data;
  },

  saveStream: async (streamId: string, videoUri: string): Promise<MintResponse> => {
    // TODO: Upload video file
    const { data } = await api.post(`/streams/${streamId}/save`, {
      video_uri: videoUri,
    });
    return data;
  },

  // NFTs
  listNFTs: async (params?: {
    latitude?: number;
    longitude?: number;
    radius_km?: number;
  }) => {
    const { data } = await api.get('/nfts', { params });
    return data;
  },

  getNFT: async (mintAddress: string) => {
    const { data } = await api.get(`/nfts/${mintAddress}`);
    return data;
  },

  // Health check
  health: async () => {
    const { data } = await api.get('/health');
    return data;
  },
};

// Set auth token
export const setAuthToken = (token: string) => {
  api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
};

// Remove auth token
export const removeAuthToken = () => {
  delete api.defaults.headers.common['Authorization'];
};

export default api;
