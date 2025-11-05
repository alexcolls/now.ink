package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
)

// ArweaveClient handles uploads to Arweave permanent storage
type ArweaveClient struct {
	wallet *goar.Wallet
	client *goar.Client
}

// NewArweaveClient creates a new Arweave client
func NewArweaveClient() (*ArweaveClient, error) {
	// Load wallet from environment or create ephemeral one for testing
	walletPath := os.Getenv("ARWEAVE_WALLET_PATH")
	
	var wallet *goar.Wallet
	var err error

	if walletPath != "" && fileExists(walletPath) {
		// Load existing wallet
		wallet, err = goar.NewWalletFromPath(walletPath, "https://arweave.net")
		if err != nil {
			return nil, fmt.Errorf("failed to load Arweave wallet: %w", err)
		}
		log.Println("✅ Loaded Arweave wallet from", walletPath)
	} else {
		// For development, use a mock/test setup
		log.Println("⚠️  No Arweave wallet found, using test mode")
		// In production, you must configure a real wallet
		return &ArweaveClient{}, nil
	}

	client := goar.NewClient("https://arweave.net")

	return &ArweaveClient{
		wallet: wallet,
		client: client,
	}, nil
}

// UploadVideo uploads a video file to Arweave and returns the transaction ID
func (a *ArweaveClient) UploadVideo(ctx context.Context, videoPath string, metadata VideoMetadata) (string, error) {
	// For MVP without wallet, return mock transaction ID
	if a.client == nil {
		log.Println("📦 Mock Arweave upload:", videoPath)
		// Generate a fake transaction ID for testing
		mockTxID := fmt.Sprintf("MOCK_AR_%d", time.Now().Unix())
		return mockTxID, nil
	}

	// Read video file
	file, err := os.Open(videoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open video file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read video file: %w", err)
	}

	// Create tags for the upload
	tags := []types.Tag{
		{Name: "Content-Type", Value: "video/mp4"},
		{Name: "App-Name", Value: "now.ink"},
		{Name: "App-Version", Value: "0.1.0"},
		{Name: "Type", Value: "video"},
		{Name: "Title", Value: metadata.Title},
		{Name: "Creator", Value: metadata.Creator},
		{Name: "Latitude", Value: fmt.Sprintf("%.6f", metadata.Latitude)},
		{Name: "Longitude", Value: fmt.Sprintf("%.6f", metadata.Longitude)},
		{Name: "Timestamp", Value: metadata.Timestamp.Format(time.RFC3339)},
		{Name: "Duration", Value: fmt.Sprintf("%d", metadata.Duration)},
	}

	// Create transaction
	tx, err := a.wallet.SendData(data, tags)
	if err != nil {
		return "", fmt.Errorf("failed to upload to Arweave: %w", err)
	}

	log.Printf("✅ Uploaded to Arweave: %s (%.2f MB)", tx.ID, float64(len(data))/(1024*1024))
	return tx.ID, nil
}

// UploadMetadata uploads NFT metadata JSON to Arweave
func (a *ArweaveClient) UploadMetadata(ctx context.Context, metadata NFTMetadata) (string, error) {
	// For MVP without wallet, return mock transaction ID
	if a.client == nil {
		log.Println("📦 Mock Arweave metadata upload")
		mockTxID := fmt.Sprintf("MOCK_AR_META_%d", time.Now().Unix())
		return mockTxID, nil
	}

	// Convert metadata to JSON
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Create tags
	tags := []types.Tag{
		{Name: "Content-Type", Value: "application/json"},
		{Name: "App-Name", Value: "now.ink"},
		{Name: "Type", Value: "nft-metadata"},
	}

	// Create transaction
	tx, err := a.wallet.SendData(jsonData, tags)
	if err != nil {
		return "", fmt.Errorf("failed to upload metadata to Arweave: %w", err)
	}

	log.Printf("✅ Uploaded metadata to Arweave: %s", tx.ID)
	return tx.ID, nil
}

// GetTransactionStatus checks if a transaction is confirmed
func (a *ArweaveClient) GetTransactionStatus(txID string) (string, error) {
	if a.client == nil {
		return "confirmed", nil // Mock always confirmed
	}

	status, err := a.client.GetTransactionStatus(txID)
	if err != nil {
		return "pending", nil
	}

	if status.NumberOfConfirmations > 0 {
		return "confirmed", nil
	}

	return "pending", nil
}

// VideoMetadata represents video upload metadata
type VideoMetadata struct {
	Title     string
	Creator   string
	Latitude  float64
	Longitude float64
	Timestamp time.Time
	Duration  int
}

// NFTMetadata represents the Metaplex-compatible metadata
type NFTMetadata struct {
	Name                 string              `json:"name"`
	Symbol               string              `json:"symbol"`
	Description          string              `json:"description"`
	SellerFeeBasisPoints int                 `json:"seller_fee_basis_points"`
	Image                string              `json:"image,omitempty"`
	AnimationURL         string              `json:"animation_url"`
	ExternalURL          string              `json:"external_url,omitempty"`
	Attributes           []MetadataAttribute `json:"attributes"`
	Properties           MetadataProperties  `json:"properties"`
}

// MetadataAttribute represents an NFT attribute
type MetadataAttribute struct {
	TraitType string      `json:"trait_type"`
	Value     interface{} `json:"value"`
}

// MetadataProperties represents NFT properties
type MetadataProperties struct {
	Files    []MetadataFile    `json:"files"`
	Category string            `json:"category"`
	Creators []MetadataCreator `json:"creators"`
}

// MetadataFile represents a file reference
type MetadataFile struct {
	URI  string `json:"uri"`
	Type string `json:"type"`
}

// MetadataCreator represents a creator with share
type MetadataCreator struct {
	Address string `json:"address"`
	Share   int    `json:"share"`
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
