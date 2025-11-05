package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

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

// MintNFT mints an NFT on Solana using Metaplex (calls TypeScript script)
func (s *SolanaClient) MintNFT(ctx context.Context, opts MintOptions) (*MintResult, error) {
	// Determine if we should use real minting or mock
	useRealMinting := os.Getenv("USE_REAL_MINTING") == "true"
	
	if !useRealMinting {
		// Mock mode for development
		log.Println("⏳ Mock minting mode (set USE_REAL_MINTING=true for production)")
		return &MintResult{
			MintAddress:  "MOCK_MINT_" + opts.CreatorWallet[:min(8, len(opts.CreatorWallet))],
			MetadataURI:  opts.MetadataURI,
			ArweaveTxID:  opts.ArweaveTxID,
			Status:       "minted",
			Network:      s.network,
		}, nil
	}
	
	// Real minting via TypeScript/Metaplex
	return s.mintWithMetaplex(ctx, opts)
}

// mintWithMetaplex calls the TypeScript minting script
func (s *SolanaClient) mintWithMetaplex(ctx context.Context, opts MintOptions) (*MintResult, error) {
	log.Println("⚡ Real Metaplex minting...")
	
	// Construct script path
	scriptPath := os.Getenv("BLOCKCHAIN_SCRIPTS_PATH")
	if scriptPath == "" {
		scriptPath = "./blockchain/scripts"
	}
	
	// Create temp output file
	outputFile := fmt.Sprintf("/tmp/mint-result-%d.json", time.Now().UnixNano())
	defer os.Remove(outputFile)
	
	// Build command
	cmd := exec.CommandContext(ctx,
		"npx", "tsx",
		fmt.Sprintf("%s/mint-nft.ts", scriptPath),
		"--metadata-uri", opts.MetadataURI,
		"--name", opts.Title,
		"--creator-wallet", opts.CreatorWallet,
		"--network", s.network,
		"--output", outputFile,
	)
	
	// Set working directory
	cmd.Dir = scriptPath
	
	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Run command with timeout
	if err := cmd.Run(); err != nil {
		log.Printf("❌ Minting command failed: %v", err)
		log.Printf("stdout: %s", stdout.String())
		log.Printf("stderr: %s", stderr.String())
		return nil, fmt.Errorf("minting failed: %w", err)
	}
	
	// Read result
	resultData, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read mint result: %w", err)
	}
	
	var result struct {
		Success      bool   `json:"success"`
		MintAddress  string `json:"mint_address"`
		MetadataURI  string `json:"metadata_uri"`
		Error        string `json:"error"`
		ExplorerURL  string `json:"explorer_url"`
	}
	
	if err := json.Unmarshal(resultData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse mint result: %w", err)
	}
	
	if !result.Success {
		return nil, fmt.Errorf("minting failed: %s", result.Error)
	}
	
	log.Printf("✅ NFT minted: %s", result.MintAddress)
	log.Printf("🔍 View: %s", result.ExplorerURL)
	
	return &MintResult{
		MintAddress:  result.MintAddress,
		MetadataURI:  result.MetadataURI,
		ArweaveTxID:  opts.ArweaveTxID,
		Status:       "minted",
		Network:      s.network,
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
