package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/alexcolls/now.ink/backend/internal/db"
	"github.com/alexcolls/now.ink/backend/internal/models"
	"github.com/google/uuid"
)

// UserProfile represents a user profile with stats
type UserProfile struct {
	models.User
	FollowersCount int `json:"followers_count"`
	FollowingCount int `json:"following_count"`
	NFTCount       int `json:"nft_count"`
}

// FollowUser creates a follow relationship
func (s *Service) FollowUser(ctx context.Context, followerID, followingID string) error {
	query := `
		INSERT INTO follows (follower_id, following_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (follower_id, following_id) DO NOTHING
	`

	_, err := db.DB.ExecContext(ctx, query, followerID, followingID, time.Now())
	return err
}

// UnfollowUser removes a follow relationship
func (s *Service) UnfollowUser(ctx context.Context, followerID, followingID string) error {
	query := `
		DELETE FROM follows
		WHERE follower_id = $1 AND following_id = $2
	`

	result, err := db.DB.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("follow relationship not found")
	}

	return nil
}

// IsFollowing checks if follower follows following
func (s *Service) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM follows
			WHERE follower_id = $1 AND following_id = $2
		)
	`

	var exists bool
	err := db.DB.QueryRowContext(ctx, query, followerID, followingID).Scan(&exists)
	return exists, err
}

// GetFollowers gets a user's followers
func (s *Service) GetFollowers(ctx context.Context, userID string, limit, offset int) ([]*models.User, error) {
	query := `
		SELECT 
			u.id, u.wallet_address, u.username, u.bio, u.avatar_url,
			u.is_premium, u.created_at, u.updated_at
		FROM users u
		INNER JOIN follows f ON u.id = f.follower_id
		WHERE f.following_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.WalletAddress,
			&user.Username,
			&user.Bio,
			&user.AvatarURL,
			&user.IsPremium,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if users == nil {
		users = []*models.User{}
	}

	return users, rows.Err()
}

// GetFollowing gets users a user is following
func (s *Service) GetFollowing(ctx context.Context, userID string, limit, offset int) ([]*models.User, error) {
	query := `
		SELECT 
			u.id, u.wallet_address, u.username, u.bio, u.avatar_url,
			u.is_premium, u.created_at, u.updated_at
		FROM users u
		INNER JOIN follows f ON u.id = f.following_id
		WHERE f.follower_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.WalletAddress,
			&user.Username,
			&user.Bio,
			&user.AvatarURL,
			&user.IsPremium,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if users == nil {
		users = []*models.User{}
	}

	return users, rows.Err()
}

// SearchUsers searches for users by username or wallet address
func (s *Service) SearchUsers(ctx context.Context, query string, limit, offset int) ([]*models.User, error) {
	searchQuery := `
		SELECT 
			id, wallet_address, username, bio, avatar_url,
			is_premium, created_at, updated_at
		FROM users
		WHERE 
			username ILIKE '%' || $1 || '%' OR
			wallet_address ILIKE '%' || $1 || '%'
		ORDER BY 
			CASE 
				WHEN username = $1 THEN 0
				WHEN wallet_address = $1 THEN 1
				WHEN username ILIKE $1 || '%' THEN 2
				ELSE 3
			END,
			created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.DB.QueryContext(ctx, searchQuery, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.WalletAddress,
			&user.Username,
			&user.Bio,
			&user.AvatarURL,
			&user.IsPremium,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if users == nil {
		users = []*models.User{}
	}

	return users, rows.Err()
}

// GetUserProfile gets detailed user profile with stats
func (s *Service) GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
	// Parse UUID
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get user basic info
	user := &models.User{}
	userQuery := `
		SELECT id, wallet_address, username, bio, avatar_url,
			   is_premium, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err = db.DB.QueryRowContext(ctx, userQuery, uid).Scan(
		&user.ID,
		&user.WalletAddress,
		&user.Username,
		&user.Bio,
		&user.AvatarURL,
		&user.IsPremium,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Get stats
	statsQuery := `
		SELECT 
			(SELECT COUNT(*) FROM follows WHERE following_id = $1) as followers_count,
			(SELECT COUNT(*) FROM follows WHERE follower_id = $1) as following_count,
			(SELECT COUNT(*) FROM nfts WHERE owner_wallet = $2) as nft_count
	`

	profile := &UserProfile{
		User: *user,
	}

	err = db.DB.QueryRowContext(ctx, statsQuery, uid, user.WalletAddress).Scan(
		&profile.FollowersCount,
		&profile.FollowingCount,
		&profile.NFTCount,
	)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
