package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/gecogreen/backend/internal/auth"
	"github.com/gecogreen/backend/internal/config"
	"github.com/gecogreen/backend/internal/database"
	"github.com/gecogreen/backend/internal/handlers"
	"github.com/gecogreen/backend/internal/middleware"
	"github.com/gecogreen/backend/internal/moderation"
	"github.com/gecogreen/backend/internal/repository"
	"github.com/gecogreen/backend/internal/storage"
)

func main() {
	cfg := config.Load()

	log.Printf("🌿 GecoGreen API Starting...")
	log.Printf("   Environment: %s", cfg.AppEnv)
	log.Printf("   Port: %s", cfg.Port)

	db, err := database.New(cfg.DatabaseURL, cfg.RedisURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Connected to PostgreSQL")
	log.Println("✅ Connected to Redis")

	// Run database migrations
	ctx := context.Background()
	if err := db.Migrate(ctx); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	// Repositories
	userRepo := repository.NewUserRepository(db.Pool)
	productRepo := repository.NewProductRepository(db.Pool)
	imageReviewRepo := repository.NewImageReviewRepository(db.Pool)
	leaderboardRepo := repository.NewLeaderboardRepository(db.Pool)
	orderRepo := repository.NewOrderRepository(db.Pool)

	// JWT
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	// Handlers
	healthHandler := handlers.NewHealthHandler(db)
	authHandler := handlers.NewAuthHandler(userRepo, jwtManager)
	productHandler := handlers.NewProductHandler(productRepo)
	adminHandler := handlers.NewAdminHandler(userRepo, imageReviewRepo)
	leaderboardHandler := handlers.NewLeaderboardHandler(leaderboardRepo, userRepo)
	orderHandler := handlers.NewOrderHandler(orderRepo, productRepo, userRepo)

	// Optional: R2 Storage (for image uploads)
	var uploadHandler *handlers.UploadHandler
	var profileHandler *handlers.ProfileHandler
	var r2Storage *storage.R2Storage

	if cfg.R2AccountID != "" && cfg.R2AccessKeyID != "" {
		r2Storage, err = storage.NewR2Storage(cfg.R2AccountID, cfg.R2AccessKeyID, cfg.R2SecretKey, cfg.R2BucketName)
		if err != nil {
			log.Printf("⚠️  R2 Storage not configured: %v", err)
		} else {
			uploadHandler = handlers.NewUploadHandler(r2Storage, productRepo)
			profileHandler = handlers.NewProfileHandler(userRepo, r2Storage)
			log.Println("✅ Connected to Cloudflare R2")
		}
	} else {
		log.Println("⚠️  R2 Storage not configured (uploads disabled)")
		profileHandler = handlers.NewProfileHandler(userRepo, nil)
	}

	// Optional: Google Cloud Vision OCR
	var ocrService *moderation.OCRService
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" {
		ocrService, err = moderation.NewOCRService(ctx)
		if err != nil {
			log.Printf("⚠️  OCR Service not configured: %v", err)
		} else {
			defer ocrService.Close()
			log.Println("✅ Connected to Google Cloud Vision")
			// Connect OCR to upload handler
			if uploadHandler != nil {
				uploadHandler.SetModerationServices(ocrService, imageReviewRepo)
			}
		}
	} else {
		log.Println("⚠️  OCR Service not configured (GOOGLE_APPLICATION_CREDENTIALS not set)")
	}

	// Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "GecoGreen API v0.2.0",
		ErrorHandler: customErrorHandler,
		BodyLimit:    10 * 1024 * 1024,
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: func() string {
			if cfg.IsDevelopment() {
				return "*"
			}
			// Allow custom frontend origin from env, plus default domains
			frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
			origins := "https://www.gecogreen.com,https://gecogreen.com"
			if frontendOrigin != "" {
				origins = frontendOrigin + "," + origins
			}
			return origins
		}(),
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	authMiddleware := middleware.AuthMiddleware(jwtManager, userRepo)
	adminMiddleware := middleware.AdminOnly(userRepo)

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"name": "GecoGreen API", "version": "0.2.0"})
	})
	app.Get("/ping", healthHandler.Ping)
	app.Get("/health", healthHandler.Check)

	v1 := app.Group("/api/v1")
	v1.Get("/health", healthHandler.Check)

	// Auth
	authRoutes := v1.Group("/auth")
	authRoutes.Post("/register", authHandler.Register)
	authRoutes.Post("/login", authHandler.Login)
	authRoutes.Post("/refresh", authHandler.Refresh)
	authRoutes.Get("/me", authMiddleware, authHandler.Me)

	// Profile
	profile := v1.Group("/profile", authMiddleware)
	profile.Get("/", profileHandler.GetProfile)
	profile.Put("/", profileHandler.UpdateProfile)
	profile.Get("/locations", profileHandler.GetLocations)
	profile.Post("/locations", profileHandler.CreateLocation)
	profile.Delete("/locations/:id", profileHandler.DeleteLocation)
	if r2Storage != nil {
		profile.Post("/avatar", profileHandler.UploadAvatar)
		profile.Post("/business-photos", profileHandler.UploadBusinessPhoto)
	}

	// Products
	products := v1.Group("/products")
	products.Get("/", productHandler.List)
	products.Get("/:id", productHandler.Get)
	products.Post("/", authMiddleware, productHandler.Create)
	products.Put("/:id", authMiddleware, productHandler.Update)
	products.Delete("/:id", authMiddleware, productHandler.Delete)
	products.Get("/seller/my", authMiddleware, productHandler.MyProducts)

	// Upload (only if R2 configured)
	if uploadHandler != nil {
		upload := v1.Group("/upload")
		upload.Post("/product/:product_id/image", authMiddleware, uploadHandler.UploadProductImage)
		upload.Post("/product/:product_id/expiry-photo", authMiddleware, uploadHandler.UploadExpiryPhoto)
		upload.Post("/presign", authMiddleware, uploadHandler.GetPresignedURL)
	}

	// Admin routes
	admin := v1.Group("/admin", authMiddleware, adminMiddleware)
	admin.Get("/reviews", adminHandler.GetPendingReviews)
	admin.Get("/reviews/stats", adminHandler.GetReviewStats)
	admin.Get("/reviews/:id", adminHandler.GetReviewDetail)
	admin.Post("/reviews/:id/approve", adminHandler.ApproveReview)
	admin.Post("/reviews/:id/reject", adminHandler.RejectReview)

	// Leaderboard routes (public)
	leaderboard := v1.Group("/leaderboard")
	leaderboard.Get("/", leaderboardHandler.GetLeaderboard)
	leaderboard.Get("/hall-of-fame", leaderboardHandler.GetHallOfFame)
	leaderboard.Get("/community-stats", leaderboardHandler.GetCommunityStats)
	leaderboard.Get("/featured", leaderboardHandler.GetFeaturedAwards)
	// Leaderboard routes (authenticated)
	leaderboard.Get("/my-rank", authMiddleware, leaderboardHandler.GetMyRank)
	leaderboard.Get("/my-awards", authMiddleware, leaderboardHandler.GetMyAwards)
	leaderboard.Get("/my-impact", authMiddleware, leaderboardHandler.GetMyImpactHistory)
	leaderboard.Post("/redeem", authMiddleware, leaderboardHandler.RedeemReward)

	// Admin Awards/Tasks routes
	adminAwards := admin.Group("/awards")
	adminAwards.Get("/tasks", leaderboardHandler.AdminGetTasks)
	adminAwards.Post("/tasks", leaderboardHandler.AdminCreateTask)
	adminAwards.Put("/tasks/:id", leaderboardHandler.AdminUpdateTask)
	adminAwards.Post("/tasks/:id/complete", leaderboardHandler.AdminCompleteTask)
	adminAwards.Post("/", leaderboardHandler.AdminCreateAward)
	adminAwards.Get("/:id", leaderboardHandler.AdminGetAward)
	adminAwards.Put("/:id", leaderboardHandler.AdminUpdateAward)

	// Orders routes
	orders := v1.Group("/orders", authMiddleware)
	orders.Get("/", orderHandler.ListMyOrders)           // List my orders (as buyer)
	orders.Get("/seller", orderHandler.ListSellerOrders) // List orders (as seller)
	orders.Post("/", orderHandler.CreateOrder)           // Create new order
	orders.Get("/:id", orderHandler.GetOrder)            // Get order details
	orders.Put("/:id/status", orderHandler.UpdateOrderStatus) // Update status (seller)
	orders.Post("/:id/cancel", orderHandler.CancelOrder) // Cancel order
	orders.Get("/:id/qr", orderHandler.GetQRCode)        // Get QR code (buyer)
	orders.Post("/confirm-pickup", orderHandler.ConfirmPickup) // Confirm pickup (seller scans QR)
	orders.Post("/:id/dispute", orderHandler.OpenDispute) // Open dispute
	orders.Get("/:id/dispute", orderHandler.GetDispute)   // Get dispute
	orders.Post("/:id/review", orderHandler.CreateReview) // Create review

	// 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not Found"})
	})

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("🛑 Shutting down server...")
		_ = app.Shutdown()
	}()

	log.Printf("🚀 Server running on http://localhost:%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}
	return c.Status(code).JSON(fiber.Map{"error": message, "code": code})
}
