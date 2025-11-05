package handlers

import (
	"github.com/alexcolls/now.ink/backend/internal/services/nft"
	"github.com/alexcolls/now.ink/backend/internal/services/stream"
	"github.com/gofiber/fiber/v2"
)

// Handlers holds all service dependencies
type Handlers struct {
	StreamService *stream.Service
	NFTService    *nft.Service
}

// NewHandlers creates new handlers with services
func NewHandlers() *Handlers {
	return &Handlers{
		StreamService: stream.NewService(),
		NFTService:    nft.NewService(),
	}
}

// RegisterRoutes registers all API routes
func (h *Handlers) RegisterRoutes(api fiber.Router) {
	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/nonce", h.HandleNonce)
	auth.Post("/verify", h.HandleVerify)

	// Stream routes
	streams := api.Group("/streams")
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
	// TODO: Generate and store nonce
	return c.JSON(fiber.Map{
		"nonce":   "test-nonce-placeholder",
		"message": "Sign this message to authenticate with now.ink",
	})
}

// HandleVerify verifies wallet signature and issues JWT
func (h *Handlers) HandleVerify(c *fiber.Ctx) error {
	// TODO: Verify signature and issue JWT
	return c.JSON(fiber.Map{
		"token": "test-jwt-token-placeholder",
		"user": fiber.Map{
			"wallet_address": "test-wallet",
		},
	})
}

// HandleStartStream starts a new live stream
func (h *Handlers) HandleStartStream(c *fiber.Ctx) error {
	var req stream.StartStreamRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	// TODO: Get user ID from JWT
	req.UserID = "test-user-id"

	stream, err := h.StreamService.StartStream(c.Context(), &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
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

	// TODO: Fetch stream details
	// TODO: Upload video to backend storage
	// TODO: Call NFT service to mint

	mintReq := &nft.MintRequest{
		VideoURL:   "placeholder-video-url",
		Title:      "Test Moment",
		UserWallet: "test-wallet",
		Latitude:   40.7128,
		Longitude:  -74.0060,
		Duration:   42,
	}

	mintResp, err := h.NFTService.Mint(c.Context(), mintReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"stream_id": streamID,
		"mint":      mintResp,
		"message":   "NFT minting in progress",
	})
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
