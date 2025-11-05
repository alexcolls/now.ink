package user

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/alexcolls/now.ink/backend/internal/db"
	"github.com/alexcolls/now.ink/backend/internal/models"
	"github.com/google/uuid"
)

// Service handles user-related operations
type Service struct{}

// NewService creates a new user service
func NewService() *Service {
	return &Service{}
}

// GetOrCreateUser gets an existing user or creates a new one
func (s *Service) GetOrCreateUser(walletAddress string) (*models.User, error) {
	user, err := s.GetUserByWallet(walletAddress)
	if err == nil {
		return user, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Create new user
	newUser := &models.User{
		ID:            uuid.New(),
		WalletAddress: walletAddress,
		IsPremium:     false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	query := `
		INSERT INTO users (id, wallet_address, is_premium, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, wallet_address, username, bio, avatar_url, is_premium, created_at, updated_at
	`

	err = db.DB.QueryRow(
		query,
		newUser.ID,
		newUser.WalletAddress,
		newUser.IsPremium,
		newUser.CreatedAt,
		newUser.UpdatedAt,
	).Scan(
		&newUser.ID,
		&newUser.WalletAddress,
		&newUser.Username,
		&newUser.Bio,
		&newUser.AvatarURL,
		&newUser.IsPremium,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)

	return newUser, err
}

// GetUserByWallet retrieves a user by wallet address
func (s *Service) GetUserByWallet(walletAddress string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, wallet_address, username, bio, avatar_url, is_premium, created_at, updated_at
		FROM users
		WHERE wallet_address = $1
	`

	err := db.DB.QueryRow(query, walletAddress).Scan(
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

	return user, nil
}

// GenerateNonce creates a new authentication nonce
func (s *Service) GenerateNonce(walletAddress string) (string, error) {
	// Generate random nonce
	nonceBytes := make([]byte, 32)
	if _, err := rand.Read(nonceBytes); err != nil {
		return "", err
	}
	nonce := hex.EncodeToString(nonceBytes)

	// Store in sessions table
	sessionID := uuid.New()
	expiresAt := time.Now().Add(5 * time.Minute) // 5 minute expiry

	query := `
		INSERT INTO sessions (id, wallet_address, nonce, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := db.DB.Exec(query, sessionID, walletAddress, nonce, time.Now(), expiresAt)
	if err != nil {
		return "", err
	}

	return nonce, nil
}

// ValidateNonce checks if a nonce is valid for a wallet
func (s *Service) ValidateNonce(walletAddress, nonce string) (bool, error) {
	query := `
		SELECT id FROM sessions
		WHERE wallet_address = $1 AND nonce = $2 AND expires_at > NOW()
	`

	var sessionID uuid.UUID
	err := db.DB.QueryRow(query, walletAddress, nonce).Scan(&sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	// Delete used nonce
	deleteQuery := `DELETE FROM sessions WHERE id = $1`
	_, _ = db.DB.Exec(deleteQuery, sessionID)

	return true, nil
}

// CleanExpiredSessions removes expired sessions
func (s *Service) CleanExpiredSessions() error {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`
	_, err := db.DB.Exec(query)
	return err
}
