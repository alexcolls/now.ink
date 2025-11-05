package stream

import (
	"context"
	"fmt"

	"github.com/alexcolls/now.ink/backend/internal/db"
	"github.com/google/uuid"
)

// UpdateStreamMintInfo updates a stream with NFT minting information
func (s *Service) UpdateStreamMintInfo(ctx context.Context, streamID, mintAddress, arweaveTxID string) error {
	id, err := uuid.Parse(streamID)
	if err != nil {
		return fmt.Errorf("invalid stream_id: %w", err)
	}

	query := `
		UPDATE streams
		SET nft_mint_address = $1, arweave_tx_id = $2
		WHERE id = $3
	`

	result, err := db.DB.ExecContext(ctx, query, mintAddress, arweaveTxID, id)
	if err != nil {
		return fmt.Errorf("failed to update stream: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stream not found: %s", streamID)
	}

	return nil
}
