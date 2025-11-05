package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// HandleFollowUser follows a user
func (h *Handlers) HandleFollowUser(c *fiber.Ctx) error {
	// Get wallet address from JWT
	followerWallet, ok := c.Locals("wallet_address").(string)
	if !ok || followerWallet == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Get follower user
	follower, err := h.UserService.GetUserByWallet(followerWallet)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "user not found"})
	}

	// Get target user ID from params
	targetUserID := c.Params("user_id")
	if targetUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id required"})
	}

	// Prevent self-follow
	if follower.ID.String() == targetUserID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot follow yourself"})
	}

	// Create follow relationship
	err = h.UserService.FollowUser(c.Context(), follower.ID.String(), targetUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":   "followed successfully",
		"following": true,
	})
}

// HandleUnfollowUser unfollows a user
func (h *Handlers) HandleUnfollowUser(c *fiber.Ctx) error {
	// Get wallet address from JWT
	followerWallet, ok := c.Locals("wallet_address").(string)
	if !ok || followerWallet == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Get follower user
	follower, err := h.UserService.GetUserByWallet(followerWallet)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "user not found"})
	}

	// Get target user ID from params
	targetUserID := c.Params("user_id")
	if targetUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id required"})
	}

	// Unfollow
	err = h.UserService.UnfollowUser(c.Context(), follower.ID.String(), targetUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":   "unfollowed successfully",
		"following": false,
	})
}

// HandleGetFollowers gets a user's followers
func (h *Handlers) HandleGetFollowers(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id required"})
	}

	// Parse pagination params
	limit := parseInt(c.Query("limit", "50"), 50)
	offset := parseInt(c.Query("offset", "0"), 0)

	followers, err := h.UserService.GetFollowers(c.Context(), userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"followers": followers,
		"count":     len(followers),
	})
}

// HandleGetFollowing gets users a user is following
func (h *Handlers) HandleGetFollowing(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id required"})
	}

	// Parse pagination params
	limit := parseInt(c.Query("limit", "50"), 50)
	offset := parseInt(c.Query("offset", "0"), 0)

	following, err := h.UserService.GetFollowing(c.Context(), userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"following": following,
		"count":     len(following),
	})
}

// HandleGetFeed gets chronological feed from followed users
func (h *Handlers) HandleGetFeed(c *fiber.Ctx) error {
	// Get wallet address from JWT
	walletAddress, ok := c.Locals("wallet_address").(string)
	if !ok || walletAddress == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Get user
	user, err := h.UserService.GetUserByWallet(walletAddress)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "user not found"})
	}

	// Parse pagination params
	limit := parseInt(c.Query("limit", "20"), 20)
	offset := parseInt(c.Query("offset", "0"), 0)

	// Get feed (NFTs from users I follow, chronologically)
	feed, err := h.NFTService.GetUserFeed(c.Context(), user.ID.String(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"feed":  feed,
		"count": len(feed),
	})
}

// HandleSearchUsers searches for users by username or wallet
func (h *Handlers) HandleSearchUsers(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "query parameter 'q' required"})
	}

	// Parse pagination params
	limit := parseInt(c.Query("limit", "20"), 20)
	offset := parseInt(c.Query("offset", "0"), 0)

	users, err := h.UserService.SearchUsers(c.Context(), query, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"users": users,
		"count": len(users),
		"query": query,
	})
}

// HandleCheckFollowing checks if current user follows target user
func (h *Handlers) HandleCheckFollowing(c *fiber.Ctx) error {
	// Get wallet address from JWT
	walletAddress, ok := c.Locals("wallet_address").(string)
	if !ok || walletAddress == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Get user
	user, err := h.UserService.GetUserByWallet(walletAddress)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "user not found"})
	}

	// Get target user ID
	targetUserID := c.Params("user_id")
	if targetUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id required"})
	}

	isFollowing, err := h.UserService.IsFollowing(c.Context(), user.ID.String(), targetUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"following": isFollowing,
	})
}

// HandleGetUserProfile gets detailed user profile with stats
func (h *Handlers) HandleGetUserProfile(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id required"})
	}

	// Get user profile with stats
	profile, err := h.UserService.GetUserProfile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(profile)
}

// parseInt safely parses string to int with default
func parseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil || val < 0 {
		return defaultVal
	}
	return val
}
