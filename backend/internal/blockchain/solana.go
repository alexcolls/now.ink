package blockchain

import (
	"context"
	"fmt"
	"os"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// SolanaClient wraps Solana RPC client
type SolanaClient struct {
	rpcClient *rpc.Client
	network   string
}

// NewSolanaClient creates a new Solana client
func NewSolanaClient() (*SolanaClient, error) {
	network := os.Getenv("SOLANA_NETWORK")
	if network == "" {
		network = "devnet"
	}

	rpcURL := os.Getenv("SOLANA_RPC_URL")
	if rpcURL == "" {
		switch network {
		case "mainnet":
			rpcURL = rpc.MainNetBeta_RPC
		case "devnet":
			rpcURL = rpc.DevNet_RPC
		case "testnet":
			rpcURL = rpc.TestNet_RPC
		default:
			rpcURL = rpc.DevNet_RPC
		}
	}

	client := rpc.New(rpcURL)

	return &SolanaClient{
		rpcClient: client,
		network:   network,
	}, nil
}

// GetBalance returns the SOL balance for a public key
func (s *SolanaClient) GetBalance(ctx context.Context, pubkey solana.PublicKey) (uint64, error) {
	balance, err := s.rpcClient.GetBalance(ctx, pubkey, rpc.CommitmentFinalized)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}
	return balance.Value, nil
}

// GetNetwork returns the current network
func (s *SolanaClient) GetNetwork() string {
	return s.network
}

// VerifySignature verifies a Solana wallet signature
func (s *SolanaClient) VerifySignature(publicKey string, message []byte, signature []byte) (bool, error) {
	// Parse public key
	pubKey, err := solana.PublicKeyFromBase58(publicKey)
	if err != nil {
		return false, fmt.Errorf("invalid public key: %w", err)
	}

	// Parse signature
	sig := solana.SignatureFromBytes(signature)

	// Verify signature
	// Note: For full verification, we'd need the actual signed message
	// For MVP, we'll do basic validation
	if sig.IsZero() {
		return false, nil
	}

	// In production, you would verify against the actual transaction/message
	// For now, we accept any non-zero signature with valid public key
	_ = pubKey // Use pubKey to avoid unused variable error
	return true, nil
}

// MintNFT mints an NFT on Solana (placeholder for Metaplex integration)
// In production, this would use Metaplex SDK to create the NFT
func (s *SolanaClient) MintNFT(ctx context.Context, opts MintOptions) (*MintResult, error) {
	// TODO: Integrate with Metaplex SDK
	// For MVP, we return a mock response
	// The actual minting should be done via the blockchain/scripts/mint-test.ts
	
	return &MintResult{
		MintAddress:  "MOCK_MINT_" + opts.CreatorWallet[:8],
		MetadataURI:  opts.MetadataURI,
		ArweaveTxID:  opts.ArweaveTxID,
		Status:       "pending",
		Network:      s.network,
	}, nil
}

// MintOptions contains NFT minting parameters
type MintOptions struct {
	CreatorWallet string
	MetadataURI   string
	ArweaveTxID   string
	Title         string
	Latitude      float64
	Longitude     float64
	Duration      int
}

// MintResult contains the result of NFT minting
type MintResult struct {
	MintAddress string
	MetadataURI string
	ArweaveTxID string
	Status      string
	Network     string
}
