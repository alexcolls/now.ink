package nft

import (
	"context"
	"database/sql"

	"github.com/alexcolls/now.ink/backend/internal/db"
)

// FeedItem represents an NFT in a user's feed
type FeedItem struct {
	NFTDetails
	CreatorUsername *string `json:"creator_username"`
	CreatorAvatar   *string `json:"creator_avatar"`
}

// GetUserFeed gets chronological feed from users that userID follows
func (s *Service) GetUserFeed(ctx context.Context, userID string, limit, offset int) ([]*FeedItem, error) {
	query := `
		SELECT 
			n.mint_address,
			n.metadata_uri,
			n.title,
			n.creator_wallet,
			n.latitude,
			n.longitude,
			n.timestamp,
			n.duration_seconds,
			n.video_url,
			n.thumbnail_url,
			n.views,
			n.created_at,
			u.username,
			u.avatar_url
		FROM nfts n
		INNER JOIN users u ON n.creator_wallet = u.wallet_address
		INNER JOIN follows f ON u.id = f.following_id
		WHERE f.follower_id = $1
		ORDER BY n.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []*FeedItem
	for rows.Next() {
		item := &FeedItem{}
		var title, videoURL, thumbnailURL sql.NullString
		var durationSeconds, views sql.NullInt64
		var username, avatar sql.NullString

		err := rows.Scan(
			&item.MintAddress,
			&item.MetadataURI,
			&title,
			&item.Creator,
			&item.Latitude,
			&item.Longitude,
			&item.Timestamp,
			&durationSeconds,
			&videoURL,
			&thumbnailURL,
			&views,
			&item.CreatedAt,
			&username,
			&avatar,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		if title.Valid {
			item.Name = title.String
		}
		if videoURL.Valid {
			item.VideoURL = videoURL.String
		}
		if thumbnailURL.Valid {
			item.ThumbnailURL = &thumbnailURL.String
		}
		if durationSeconds.Valid {
			item.DurationSeconds = int(durationSeconds.Int64)
		}
		if views.Valid {
			item.Views = int(views.Int64)
		}
		if username.Valid {
			item.CreatorUsername = &username.String
		}
		if avatar.Valid {
			item.CreatorAvatar = &avatar.String
		}

		feed = append(feed, item)
	}

	if feed == nil {
		feed = []*FeedItem{}
	}

	return feed, rows.Err()
}
