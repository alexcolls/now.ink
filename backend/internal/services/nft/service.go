package nft

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/alexcolls/now.ink/backend/internal/blockchain"
	"github.com/alexcolls/now.ink/backend/internal/db"
)

// Service handles NFT minting operations
type Service struct {
	solanaClient *blockchain.SolanaClient
}

// NewService creates a new NFT service
func NewService() *Service {
	client, _ := blockchain.NewSolanaClient()
	return &Service{
		solanaClient: client,
	}
}

// MintRequest represents the data needed to mint an NFT
type MintRequest struct {
	VideoURL    string    `json:"video_url"`
	Title       string    `json:"title"`
	UserWallet  string    `json:"user_wallet"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Duration    int       `json:"duration_seconds"`
	Timestamp   time.Time `json:"timestamp"`
}

// MintResponse represents the minting result
type MintResponse struct {
	MintAddress  string    `json:"mint_address"`
	MetadataURI  string    `json:"metadata_uri"`
	ArweaveHash  string    `json:"arweave_hash"`
	Transaction  string    `json:"transaction"`
	Status       string    `json:"status"`
	MintedAt     time.Time `json:"minted_at"`
}

// Mint mints a new NFT on Solana
func (s *Service) Mint(ctx context.Context, req *MintRequest) (*MintResponse, error) {
	// For MVP, use blockchain client to get mock response
	// In production, this will call Metaplex SDK
	mintOpts := blockchain.MintOptions{
		CreatorWallet: req.UserWallet,
		MetadataURI:   "ar://pending",
		ArweaveTxID:   "pending",
		Title:         req.Title,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		Duration:      req.Duration,
	}

	result, err := s.solanaClient.MintNFT(ctx, mintOpts)
	if err != nil {
		return nil, err
	}

	// TODO: Save to database

	return &MintResponse{
		MintAddress: result.MintAddress,
		MetadataURI: result.MetadataURI,
		ArweaveHash: result.ArweaveTxID,
		Transaction: "pending",
		Status:      result.Status,
		MintedAt:    time.Now(),
	}, nil
}

// GetNFT retrieves NFT details by mint address
func (s *Service) GetNFT(ctx context.Context, mintAddress string) (*NFTDetails, error) {
	query := `
		SELECT mint_address, metadata_uri, title, creator_wallet, latitude, longitude, 
		       timestamp, duration_seconds, video_url, thumbnail_url
		FROM nfts
		WHERE mint_address = $1
	`

	details := &NFTDetails{}
	var title, videoURL, thumbnailURL sql.NullString
	var durationSeconds sql.NullInt64

	err := db.DB.QueryRowContext(ctx, query, mintAddress).Scan(
		&details.MintAddress,
		&details.MetadataURI,
		&title,
		&details.Creator,
		&details.Latitude,
		&details.Longitude,
		&details.Timestamp,
		&durationSeconds,
		&videoURL,
		&thumbnailURL,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("NFT not found")
		}
		return nil, err
	}

	if title.Valid {
		details.Name = title.String
	}
	if videoURL.Valid {
		details.VideoURL = videoURL.String
	}
	if thumbnailURL.Valid {
		details.ThumbnailURL = thumbnailURL.String
	}
	if durationSeconds.Valid {
		details.Duration = int(durationSeconds.Int64)
	}

	details.Symbol = "NOWINK"

	return details, nil
}

// ListNFTs lists NFTs with filters
func (s *Service) ListNFTs(ctx context.Context, filters *NFTFilters) ([]*NFTDetails, error) {
	query := `
		SELECT mint_address, metadata_uri, title, creator_wallet, latitude, longitude,
		       timestamp, duration_seconds, video_url, thumbnail_url
		FROM nfts
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	// Add filters
	if filters.Creator != "" {
		query += fmt.Sprintf(" AND creator_wallet = $%d", argCount)
		args = append(args, filters.Creator)
		argCount++
	}

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND timestamp >= $%d", argCount)
		args = append(args, *filters.StartDate)
		argCount++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND timestamp <= $%d", argCount)
		args = append(args, *filters.EndDate)
		argCount++
	}

	// TODO: Add geographic radius filter using PostGIS

	query += " ORDER BY timestamp DESC"

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.Limit)
		argCount++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filters.Offset)
	}

	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nfts := []*NFTDetails{}
	for rows.Next() {
		details := &NFTDetails{}
		var title, videoURL, thumbnailURL sql.NullString
		var durationSeconds sql.NullInt64

		err := rows.Scan(
			&details.MintAddress,
			&details.MetadataURI,
			&title,
			&details.Creator,
			&details.Latitude,
			&details.Longitude,
			&details.Timestamp,
			&durationSeconds,
			&videoURL,
			&thumbnailURL,
		)
		if err != nil {
			return nil, err
		}

		if title.Valid {
			details.Name = title.String
		}
		if videoURL.Valid {
			details.VideoURL = videoURL.String
		}
		if thumbnailURL.Valid {
			details.ThumbnailURL = thumbnailURL.String
		}
		if durationSeconds.Valid {
			details.Duration = int(durationSeconds.Int64)
		}

		details.Symbol = "NOWINK"
		nfts = append(nfts, details)
	}

	return nfts, rows.Err()
}

// NFTDetails represents complete NFT information
type NFTDetails struct {
	MintAddress  string    `json:"mint_address"`
	MetadataURI  string    `json:"metadata_uri"`
	Name         string    `json:"name"`
	Symbol       string    `json:"symbol"`
	Creator      string    `json:"creator"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	Timestamp    time.Time `json:"timestamp"`
	VideoURL     string    `json:"video_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Duration     int       `json:"duration_seconds"`
}

// NFTFilters represents query filters for NFTs
type NFTFilters struct {
	Latitude  *float64  `json:"latitude"`
	Longitude *float64  `json:"longitude"`
	RadiusKm  float64   `json:"radius_km"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Creator   string    `json:"creator"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}
