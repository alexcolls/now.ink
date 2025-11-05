package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a now.ink user identified by wallet
type User struct {
	ID            uuid.UUID  `json:"id"`
	WalletAddress string     `json:"wallet_address"`
	Username      *string    `json:"username,omitempty"`
	Bio           *string    `json:"bio,omitempty"`
	AvatarURL     *string    `json:"avatar_url,omitempty"`
	IsPremium     bool       `json:"is_premium"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Stream represents a live or recorded stream
type Stream struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	Title           *string    `json:"title,omitempty"`
	IsLive          bool       `json:"is_live"`
	IsPublic        bool       `json:"is_public"`
	StartedAt       time.Time  `json:"started_at"`
	EndedAt         *time.Time `json:"ended_at,omitempty"`
	Latitude        float64    `json:"latitude"`
	Longitude       float64    `json:"longitude"`
	NFTMintAddress  *string    `json:"nft_mint_address,omitempty"`
	ArweaveTxID     *string    `json:"arweave_tx_id,omitempty"`
	ViewerCount     int        `json:"viewer_count"`
	DurationSeconds *int       `json:"duration_seconds,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// NFT represents a minted NFT on Solana
type NFT struct {
	ID              uuid.UUID  `json:"id"`
	StreamID        *uuid.UUID `json:"stream_id,omitempty"`
	MintAddress     string     `json:"mint_address"`
	MetadataURI     string     `json:"metadata_uri"`
	CreatorWallet   string     `json:"creator_wallet"`
	Title           *string    `json:"title,omitempty"`
	Latitude        float64    `json:"latitude"`
	Longitude       float64    `json:"longitude"`
	Timestamp       time.Time  `json:"timestamp"`
	DurationSeconds *int       `json:"duration_seconds,omitempty"`
	ThumbnailURL    *string    `json:"thumbnail_url,omitempty"`
	VideoURL        *string    `json:"video_url,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// Follow represents a social connection
type Follow struct {
	FollowerID  uuid.UUID `json:"follower_id"`
	FollowingID uuid.UUID `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// Session represents an authentication session
type Session struct {
	ID            uuid.UUID `json:"id"`
	WalletAddress string    `json:"wallet_address"`
	Nonce         string    `json:"nonce"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}

// StartStreamRequest is the request body for starting a stream
type StartStreamRequest struct {
	Title     *string  `json:"title,omitempty"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	IsPublic  bool     `json:"is_public"`
}

// StartStreamResponse is the response for starting a stream
type StartStreamResponse struct {
	StreamID string `json:"stream_id"`
	Message  string `json:"message"`
}

// GetNonceRequest is the request body for getting an auth nonce
type GetNonceRequest struct {
	WalletAddress string `json:"wallet_address"`
}

// GetNonceResponse is the response for getting an auth nonce
type GetNonceResponse struct {
	Nonce string `json:"nonce"`
}

// VerifyWalletRequest is the request body for wallet verification
type VerifyWalletRequest struct {
	WalletAddress string `json:"wallet_address"`
	Signature     string `json:"signature"`
	Nonce         string `json:"nonce"`
}

// VerifyWalletResponse is the response for wallet verification
type VerifyWalletResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
