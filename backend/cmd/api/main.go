package main

import (
	"log"
	"os"

	"github.com/alexcolls/now.ink/backend/internal/api/handlers"
	"github.com/alexcolls/now.ink/backend/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}

	// Connect to database
	if err := db.Connect(); err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer db.Close()

	// Ensure video storage directory exists
	videoDir := "/tmp/nowink-videos"
	if err := os.MkdirAll(videoDir, 0755); err != nil {
		log.Fatal("❌ Failed to create video storage directory:", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "now.ink API v0.1.0",
		ServerHeader: "now.ink",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"version": "0.1.0",
			"service": "now.ink API",
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// Root endpoint
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "now.ink API - Your life, minted",
			"version": "0.1.0",
			"docs":    "https://github.com/alexcolls/now.ink/tree/main/docs",
		})
	})

	// Initialize handlers
	handlers := handlers.NewHandlers()
	handlers.RegisterRoutes(api)

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("🚀 now.ink API starting on port %s", port)
	log.Printf("📖 Docs: https://github.com/alexcolls/now.ink/tree/main/docs")
	log.Printf("🔗 Health: http://localhost:%s/health", port)
	
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("❌ Server failed to start:", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
