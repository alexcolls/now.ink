import { Connection, Keypair, clusterApiUrl } from '@solana/web3.js';
import { Metaplex, keypairIdentity, bundlrStorage } from '@metaplex-foundation/js';
import fs from 'fs';
import path from 'path';

/**
 * Test NFT minting script for now.ink
 * 
 * Usage:
 *   npx tsx scripts/mint-test.ts
 * 
 * Prerequisites:
 *   1. Have a wallet at blockchain/wallets/platform-wallet.json
 *   2. Airdrop SOL: solana airdrop 2 -k wallets/platform-wallet.json
 *   3. Run on devnet only!
 */

async function main() {
  console.log('🚀 now.ink Test NFT Minting\n');

  // Connect to devnet
  const connection = new Connection(clusterApiUrl('devnet'), 'confirmed');
  console.log('📡 Connected to Solana devnet');

  // Load platform wallet
  const walletPath = path.join(__dirname, '../wallets/platform-wallet.json');
  
  if (!fs.existsSync(walletPath)) {
    console.error('❌ Wallet not found at:', walletPath);
    console.log('💡 Run: solana-keygen new --outfile blockchain/wallets/platform-wallet.json');
    process.exit(1);
  }

  const walletKeypair = Keypair.fromSecretKey(
    new Uint8Array(JSON.parse(fs.readFileSync(walletPath, 'utf8')))
  );
  
  console.log('🔑 Platform wallet:', walletKeypair.publicKey.toBase58());

  // Check balance
  const balance = await connection.getBalance(walletKeypair.publicKey);
  console.log(`💰 Balance: ${balance / 1e9} SOL`);
  
  if (balance < 0.1 * 1e9) {
    console.log('⚠️  Low balance. Run: solana airdrop 2 -k wallets/platform-wallet.json');
  }

  // Initialize Metaplex
  const metaplex = Metaplex.make(connection)
    .use(keypairIdentity(walletKeypair))
    .use(bundlrStorage({
      address: 'https://devnet.bundlr.network',
      providerUrl: clusterApiUrl('devnet'),
      timeout: 60000,
    }));

  console.log('\n🪙 Minting test NFT...\n');

  // Test metadata (simulating a now.ink moment)
  const metadata = {
    name: 'now.ink Test Moment #1',
    symbol: 'NOWINK',
    description: 'Test NFT minted from now.ink on Solana devnet',
    image: 'https://arweave.net/test-image-hash', // Placeholder
    animation_url: 'https://arweave.net/test-video-hash', // Placeholder
    external_url: 'https://now.ink/nft/test-1',
    attributes: [
      { trait_type: 'Latitude', value: '40.7128' },
      { trait_type: 'Longitude', value: '-74.0060' },
      { trait_type: 'Timestamp', value: new Date().toISOString() },
      { trait_type: 'Creator', value: walletKeypair.publicKey.toBase58().slice(0, 8) + '...' },
      { trait_type: 'Duration (seconds)', value: '42' },
      { trait_type: 'Location', value: 'Times Square, New York' },
      { trait_type: 'App', value: 'now.ink' },
    ],
    properties: {
      files: [
        {
          uri: 'https://arweave.net/test-video-hash',
          type: 'video/mp4',
        },
      ],
      category: 'video',
    },
  };

  try {
    // Upload metadata to Arweave (via Bundlr)
    console.log('📤 Uploading metadata to Arweave...');
    const { uri } = await metaplex.nfts().uploadMetadata(metadata);
    console.log('✅ Metadata URI:', uri);

    // Mint NFT
    console.log('\n⚡ Minting NFT on Solana...');
    const { nft } = await metaplex.nfts().create({
      uri,
      name: metadata.name,
      symbol: metadata.symbol,
      sellerFeeBasisPoints: 0, // No resale royalties for MVP
      creators: [
        {
          address: walletKeypair.publicKey,
          share: 100,
          verified: true,
        },
      ],
    });

    console.log('\n🎉 NFT Minted Successfully!\n');
    console.log('📋 Details:');
    console.log('  Mint Address:', nft.address.toBase58());
    console.log('  Metadata URI:', uri);
    console.log('  Name:', nft.name);
    console.log('  Symbol:', nft.symbol);
    console.log('\n🔍 View on Solscan:');
    console.log(`  https://solscan.io/token/${nft.address.toBase58()}?cluster=devnet`);
    console.log('\n✅ Test minting completed!\n');
  } catch (error) {
    console.error('\n❌ Minting failed:', error);
    process.exit(1);
  }
}

main().catch(console.error);
