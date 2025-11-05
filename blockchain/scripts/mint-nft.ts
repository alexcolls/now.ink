#!/usr/bin/env tsx
/**
 * Production NFT minting script for now.ink
 * Called by Go backend with command-line arguments
 * 
 * Usage:
 *   npx tsx scripts/mint-nft.ts \
 *     --metadata-uri ar://metadata-hash \
 *     --video-uri ar://video-hash \
 *     --name "Moment Title" \
 *     --creator-wallet WALLET_ADDRESS \
 *     --output /path/to/output.json
 */

import { Connection, Keypair, clusterApiUrl, PublicKey } from '@solana/web3.js';
import { Metaplex, keypairIdentity, bundlrStorage } from '@metaplex-foundation/js';
import fs from 'fs';
import path from 'path';

interface MintArgs {
  metadataUri: string;
  videoUri: string;
  name: string;
  creatorWallet: string;
  output?: string;
  network?: 'devnet' | 'mainnet-beta';
}

function parseArgs(): MintArgs {
  const args: any = {};
  
  for (let i = 2; i < process.argv.length; i += 2) {
    const key = process.argv[i].replace('--', '');
    const value = process.argv[i + 1];
    
    // Convert kebab-case to camelCase
    const camelKey = key.replace(/-([a-z])/g, (g) => g[1].toUpperCase());
    args[camelKey] = value;
  }
  
  if (!args.metadataUri || !args.name || !args.creatorWallet) {
    console.error('❌ Missing required arguments');
    console.log('Usage: npx tsx scripts/mint-nft.ts \\');
    console.log('  --metadata-uri ar://hash \\');
    console.log('  --name "Title" \\');
    console.log('  --creator-wallet ADDRESS \\');
    console.log('  [--network devnet|mainnet-beta] \\');
    console.log('  [--output /path/to/output.json]');
    process.exit(1);
  }
  
  return args as MintArgs;
}

async function main() {
  const args = parseArgs();
  const network = args.network || 'devnet';
  
  try {
    // Connect to Solana
    const connection = new Connection(clusterApiUrl(network), 'confirmed');
    
    // Load platform wallet
    const walletPath = path.join(__dirname, '../wallets/platform-wallet.json');
    
    if (!fs.existsSync(walletPath)) {
      throw new Error(`Wallet not found at: ${walletPath}`);
    }
    
    const walletKeypair = Keypair.fromSecretKey(
      new Uint8Array(JSON.parse(fs.readFileSync(walletPath, 'utf8')))
    );
    
    // Check balance
    const balance = await connection.getBalance(walletKeypair.publicKey);
    if (balance < 0.01 * 1e9) {
      throw new Error(`Insufficient balance: ${balance / 1e9} SOL`);
    }
    
    // Initialize Metaplex
    const bundlrUrl = network === 'devnet' 
      ? 'https://devnet.bundlr.network'
      : 'https://node1.bundlr.network';
    
    const metaplex = Metaplex.make(connection)
      .use(keypairIdentity(walletKeypair))
      .use(bundlrStorage({
        address: bundlrUrl,
        providerUrl: clusterApiUrl(network),
        timeout: 60000,
      }));
    
    // Parse creator wallet
    const creatorPubkey = new PublicKey(args.creatorWallet);
    
    // Mint NFT
    const { nft } = await metaplex.nfts().create({
      uri: args.metadataUri,
      name: args.name,
      symbol: 'NOWINK',
      sellerFeeBasisPoints: 500, // 5% platform commission
      creators: [
        {
          address: walletKeypair.publicKey, // Platform wallet (verified)
          share: 5,
          verified: true,
        },
        {
          address: creatorPubkey, // User wallet
          share: 95,
          verified: false, // User will verify later
        },
      ],
    });
    
    // Prepare result
    const result = {
      success: true,
      mint_address: nft.address.toBase58(),
      metadata_uri: args.metadataUri,
      name: nft.name,
      symbol: nft.symbol,
      update_authority: nft.updateAuthorityAddress.toBase58(),
      creators: [
        {
          address: walletKeypair.publicKey.toBase58(),
          share: 5,
          verified: true,
        },
        {
          address: args.creatorWallet,
          share: 95,
          verified: false,
        },
      ],
      network,
      timestamp: new Date().toISOString(),
      explorer_url: `https://solscan.io/token/${nft.address.toBase58()}${network === 'devnet' ? '?cluster=devnet' : ''}`,
    };
    
    // Output result
    if (args.output) {
      fs.writeFileSync(args.output, JSON.stringify(result, null, 2));
    } else {
      console.log(JSON.stringify(result));
    }
    
  } catch (error: any) {
    const errorResult = {
      success: false,
      error: error.message || 'Unknown error',
      timestamp: new Date().toISOString(),
    };
    
    if (args.output) {
      fs.writeFileSync(args.output, JSON.stringify(errorResult, null, 2));
    } else {
      console.error(JSON.stringify(errorResult));
    }
    
    process.exit(1);
  }
}

main().catch((error) => {
  console.error(JSON.stringify({
    success: false,
    error: error.message,
  }));
  process.exit(1);
});
