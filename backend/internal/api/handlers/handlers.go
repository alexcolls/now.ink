package handlers

import (
	"fmt"

	"github.com/alexcolls/now.ink/backend/internal/api/middleware"
	"github.com/alexcolls/now.ink/backend/internal/models"
	"github.com/alexcolls/now.ink/backend/internal/services/nft"
	"github.com/alexcolls/now.ink/backend/internal/services/stream"
	"github.com/alexcolls/now.ink/backend/internal/services/user"
	"github.com/gofiber/fiber/v2"
)

// Handlers holds all service dependencies
type Handlers struct {
	StreamService *stream.Service
	NFTService    *nft.Service
	UserService   *user.Service
}

// NewHandlers creates new handlers with services
func NewHandlers() *Handlers {
	return &Handlers{
		StreamService: stream.NewService(),
		NFTService:    nft.NewService(),
		UserService:   user.NewService(),
	}
}

// RegisterRoutes registers all API routes
func (h *Handlers) RegisterRoutes(api fiber.Router) {
	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/nonce", h.HandleNonce)
	auth.Post("/verify", h.HandleVerify)

	// Stream routes (authenticated)
	streams := api.Group("/streams", middleware.AuthRequired())
	streams.Post("/start", h.HandleStartStream)
	streams.Post("/:id/end", h.HandleEndStream)
	streams.Post("/:id/save", h.HandleSaveStream)
	streams.Get("/live", h.HandleListLiveStreams)
	streams.Get("/:id", h.HandleGetStream)

	// NFT routes
	nfts := api.Group("/nfts")
	nfts.Get("/", h.HandleListNFTs)
	nfts.Get("/:mint_address", h.HandleGetNFT)
	nfts.Get("/:mint_address/playback", h.HandleGetPlayback)
}

// HandleNonce generates a nonce for wallet signature
func (h *Handlers) HandleNonce(c *fiber.Ctx) error {
	var req models.GetNonceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.WalletAddress == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_address required"})
	}

	nonce, err := h.UserService.GenerateNonce(req.WalletAddress)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate nonce"})
	}

	return c.JSON(models.GetNonceResponse{
		Nonce: nonce,
	})
}

// HandleVerify verifies wallet signature and issues JWT
func (h *Handlers) HandleVerify(c *fiber.Ctx) error {
	var req models.VerifyWalletRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.WalletAddress == "" || req.Signature == "" || req.Nonce == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_address, signature, and nonce required"})
	}

	// Validate nonce
	valid, err := h.UserService.ValidateNonce(req.WalletAddress, req.Nonce)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to validate nonce"})
	}

	if !valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired nonce"})
	}

	// TODO: Verify Solana signature (req.Signature)
	// For MVP, we trust the nonce validation

	// Get or create user
	user, err := h.UserService.GetOrCreateUser(req.WalletAddress)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(req.WalletAddress)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.JSON(models.VerifyWalletResponse{
		Token: token,
		User:  *user,
	})
}

// HandleStartStream starts a new live stream
func (h *Handlers) HandleStartStream(c *fiber.Ctx) error {
	var req stream.StartStreamRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// Get wallet address from JWT context
	walletAddress, ok := c.Locals("wallet_address").(string)
	if !ok || walletAddress == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Get user from wallet address
	user, err := h.UserService.GetUserByWallet(walletAddress)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "user not found"})
	}

	req.UserID = user.ID.String()

	stream, err := h.StreamService.StartStream(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(stream)
}

// HandleEndStream ends a live stream
func (h *Handlers) HandleEndStream(c *fiber.Ctx) error {
	streamID := c.Params("id")
	
	stream, err := h.StreamService.EndStream(c.Context(), streamID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(stream)
}

// HandleSaveStream saves stream as NFT (triggers minting)
func (h *Handlers) HandleSaveStream(c *fiber.Ctx) error {
	streamID := c.Params("id")

	// Get wallet address from JWT
	walletAddress, ok := c.Locals("wallet_address").(string)
	if !ok || walletAddress == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Parse multipart form
	file, err := c.FormFile("video")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "video file required"})
	}

	// Validate file size (max 100MB for MVP)
	if file.Size > 100*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "video file too large (max 100MB)"})
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if contentType != "video/mp4" && contentType != "video/quicktime" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "only mp4 and mov videos supported"})
	}

	// Fetch stream details from database
	stream, err := h.StreamService.GetStream(c.Context(), streamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "stream not found"})
	}

	// Verify ownership
	user, err := h.UserService.GetUserByWallet(walletAddress)
	if err != nil || stream.UserID != user.ID.String() {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "not authorized to save this stream"})
	}

	// Save video file to temporary storage
	// In production, this would go to S3 or similar
	videoPath := fmt.Sprintf("/tmp/nowink-videos/%s.mp4", streamID)
	if err := c.SaveFile(file, videoPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save video"})
	}

	videoURL := videoPath // Local path, will be uploaded to Arweave during minting

	// End the stream
	_, err = h.StreamService.EndStream(c.Context(), streamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to end stream"})
	}

	// Prepare minting request
	mintReq := &nft.MintRequest{
		VideoURL:   videoURL,
		Title:      stream.Title,
		UserWallet: walletAddress,
		Latitude:   stream.Latitude,
		Longitude:  stream.Longitude,
		Duration:   calculateDuration(stream),
		Timestamp:  stream.StartedAt,
	}

	// Mint NFT
	mintResp, err := h.NFTService.Mint(c.Context(), mintReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Update stream with mint info
	err = h.StreamService.UpdateStreamMintInfo(c.Context(), streamID, mintResp.MintAddress, mintResp.ArweaveHash)
	if err != nil {
		// Log error but don't fail - NFT was already minted
		fmt.Printf("⚠️  Failed to update stream with mint info: %v\n", err)
	}

	return c.JSON(fiber.Map{
		"stream_id":    streamID,
		"mint_address": mintResp.MintAddress,
		"arweave_tx":   mintResp.ArweaveHash,
		"metadata_uri":  mintResp.MetadataURI,
		"status":       "minted",
		"message":      "NFT minted successfully!",
	})
}

func calculateDuration(s *stream.Stream) int {
	if s.EndedAt == nil {
		return 0
	}
	return int(s.EndedAt.Sub(s.StartedAt).Seconds())
}

// HandleListLiveStreams lists currently live streams
func (h *Handlers) HandleListLiveStreams(c *fiber.Ctx) error {
	streams, err := h.StreamService.ListLiveStreams(c.Context(), 50, 0)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"streams": streams,
		"total":   len(streams),
	})
}

// HandleGetStream gets a specific stream
func (h *Handlers) HandleGetStream(c *fiber.Ctx) error {
	streamID := c.Params("id")
	
	stream, err := h.StreamService.GetStream(c.Context(), streamID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "stream not found"})
	}

	return c.JSON(stream)
}

// HandleListNFTs lists NFTs with filters
func (h *Handlers) HandleListNFTs(c *fiber.Ctx) error {
	// TODO: Parse query parameters for filters
	filters := &nft.NFTFilters{
		Limit:  50,
		Offset: 0,
	}

	nfts, err := h.NFTService.ListNFTs(c.Context(), filters)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"nfts":  nfts,
		"total": len(nfts),
	})
}

// HandleGetNFT gets a specific NFT
func (h *Handlers) HandleGetNFT(c *fiber.Ctx) error {
	mintAddress := c.Params("mint_address")
	
	nft, err := h.NFTService.GetNFT(c.Context(), mintAddress)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "nft not found"})
	}

	return c.JSON(nft)
}

// HandleGetPlayback gets playback URL for NFT
func (h *Handlers) HandleGetPlayback(c *fiber.Ctx) error {
	mintAddress := c.Params("mint_address")
	
	// TODO: Check proximity or premium status
	// TODO: Return Arweave URL if authorized

	return c.JSON(fiber.Map{
		"mint_address": mintAddress,
		"playback_url": "ar://placeholder-video-hash",
		"access_reason": "testing",
	})
}
