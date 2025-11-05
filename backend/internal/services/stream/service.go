package stream

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/alexcolls/now.ink/backend/internal/db"
	"github.com/google/uuid"
)

// Service handles stream operations
type Service struct{}

// NewService creates a new stream service
func NewService() *Service {
	return &Service{}
}

// Stream represents a live stream or recording
type Stream struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	IsLive      bool      `json:"is_live"`
	IsPublic    bool      `json:"is_public"`
	StartedAt   time.Time `json:"started_at"`
	EndedAt     *time.Time `json:"ended_at,omitempty"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	ViewerCount int       `json:"viewer_count"`
	MintAddress string    `json:"mint_address,omitempty"`
	ArweaveHash string    `json:"arweave_hash,omitempty"`
}

// StartStreamRequest represents stream start data
type StartStreamRequest struct {
	UserID    string  `json:"user_id"`
	Title     string  `json:"title"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	IsPublic  bool    `json:"is_public"`
}

// StartStream initiates a new live stream
func (s *Service) StartStream(ctx context.Context, req *StartStreamRequest) (*Stream, error) {
	streamID := uuid.New()
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	query := `
		INSERT INTO streams (id, user_id, title, is_live, is_public, started_at, location, viewer_count, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, ST_SetSRID(ST_MakePoint($7, $8), 4326), $9, $10)
		RETURNING id, user_id, title, is_live, is_public, started_at, ended_at, 
		          ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude,
		          viewer_count, nft_mint_address, arweave_tx_id
	`

	now := time.Now()
	stream := &Stream{}
	var endedAt sql.NullTime
	var mintAddress, arweaveTxID sql.NullString
	var dbUserID uuid.UUID

	err = db.DB.QueryRowContext(ctx, query,
		streamID, userID, req.Title, true, req.IsPublic, now,
		req.Longitude, req.Latitude, 0, now,
	).Scan(
		&stream.ID, &dbUserID, &stream.Title, &stream.IsLive, &stream.IsPublic,
		&stream.StartedAt, &endedAt, &stream.Longitude, &stream.Latitude,
		&stream.ViewerCount, &mintAddress, &arweaveTxID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	stream.UserID = dbUserID.String()
	if endedAt.Valid {
		stream.EndedAt = &endedAt.Time
	}
	if mintAddress.Valid {
		stream.MintAddress = mintAddress.String
	}
	if arweaveTxID.Valid {
		stream.ArweaveHash = arweaveTxID.String
	}

	return stream, nil
}

// EndStream stops a live stream
func (s *Service) EndStream(ctx context.Context, streamID string) (*Stream, error) {
	id, err := uuid.Parse(streamID)
	if err != nil {
		return nil, fmt.Errorf("invalid stream_id: %w", err)
	}

	// First get the stream to calculate duration
	var startedAt time.Time
	query := `SELECT started_at FROM streams WHERE id = $1 AND is_live = true`
	err = db.DB.QueryRowContext(ctx, query, id).Scan(&startedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("stream not found or already ended")
		}
		return nil, err
	}

	now := time.Now()
	duration := int(now.Sub(startedAt).Seconds())

	// Update stream
	updateQuery := `
		UPDATE streams 
		SET is_live = false, ended_at = $1, duration_seconds = $2
		WHERE id = $3
		RETURNING id, user_id, title, is_live, is_public, started_at, ended_at,
		          ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude,
		          viewer_count, duration_seconds, nft_mint_address, arweave_tx_id
	`

	stream := &Stream{}
	var endedAt sql.NullTime
	var durationSeconds sql.NullInt64
	var mintAddress, arweaveTxID sql.NullString
	var dbUserID uuid.UUID

	err = db.DB.QueryRowContext(ctx, updateQuery, now, duration, id).Scan(
		&stream.ID, &dbUserID, &stream.Title, &stream.IsLive, &stream.IsPublic,
		&stream.StartedAt, &endedAt, &stream.Longitude, &stream.Latitude,
		&stream.ViewerCount, &durationSeconds, &mintAddress, &arweaveTxID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to end stream: %w", err)
	}

	stream.UserID = dbUserID.String()
	if endedAt.Valid {
		stream.EndedAt = &endedAt.Time
	}
	if mintAddress.Valid {
		stream.MintAddress = mintAddress.String
	}
	if arweaveTxID.Valid {
		stream.ArweaveHash = arweaveTxID.String
	}

	return stream, nil
}

// GetStream retrieves a stream by ID
func (s *Service) GetStream(ctx context.Context, streamID string) (*Stream, error) {
	id, err := uuid.Parse(streamID)
	if err != nil {
		return nil, fmt.Errorf("invalid stream_id: %w", err)
	}

	query := `
		SELECT id, user_id, title, is_live, is_public, started_at, ended_at,
		       ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude,
		       viewer_count, duration_seconds, nft_mint_address, arweave_tx_id
		FROM streams
		WHERE id = $1
	`

	stream := &Stream{}
	var endedAt sql.NullTime
	var durationSeconds sql.NullInt64
	var mintAddress, arweaveTxID sql.NullString
	var dbUserID uuid.UUID

	err = db.DB.QueryRowContext(ctx, query, id).Scan(
		&stream.ID, &dbUserID, &stream.Title, &stream.IsLive, &stream.IsPublic,
		&stream.StartedAt, &endedAt, &stream.Longitude, &stream.Latitude,
		&stream.ViewerCount, &durationSeconds, &mintAddress, &arweaveTxID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("stream not found")
		}
		return nil, err
	}

	stream.UserID = dbUserID.String()
	if endedAt.Valid {
		stream.EndedAt = &endedAt.Time
	}
	if mintAddress.Valid {
		stream.MintAddress = mintAddress.String
	}
	if arweaveTxID.Valid {
		stream.ArweaveHash = arweaveTxID.String
	}

	return stream, nil
}

// ListLiveStreams returns all currently live streams
func (s *Service) ListLiveStreams(ctx context.Context, limit, offset int) ([]*Stream, error) {
	query := `
		SELECT id, user_id, title, is_live, is_public, started_at, ended_at,
		       ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude,
		       viewer_count, duration_seconds, nft_mint_address, arweave_tx_id
		FROM streams
		WHERE is_live = true AND is_public = true
		ORDER BY started_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := db.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	streams := []*Stream{}
	for rows.Next() {
		stream := &Stream{}
		var endedAt sql.NullTime
		var durationSeconds sql.NullInt64
		var mintAddress, arweaveTxID sql.NullString
		var dbUserID uuid.UUID

		err := rows.Scan(
			&stream.ID, &dbUserID, &stream.Title, &stream.IsLive, &stream.IsPublic,
			&stream.StartedAt, &endedAt, &stream.Longitude, &stream.Latitude,
			&stream.ViewerCount, &durationSeconds, &mintAddress, &arweaveTxID,
		)
		if err != nil {
			return nil, err
		}

		stream.UserID = dbUserID.String()
		if endedAt.Valid {
			stream.EndedAt = &endedAt.Time
		}
		if mintAddress.Valid {
			stream.MintAddress = mintAddress.String
		}
		if arweaveTxID.Valid {
			stream.ArweaveHash = arweaveTxID.String
		}

		streams = append(streams, stream)
	}

	return streams, rows.Err()
}
