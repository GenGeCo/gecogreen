package main

import (
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

	// Repositories
	userRepo := repository.NewUserRepository(db.Pool)
	productRepo := repository.NewProductRepository(db.Pool)

	// JWT
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	// Handlers
	healthHandler := handlers.NewHealthHandler(db)
	authHandler := handlers.NewAuthHandler(userRepo, jwtManager)
	productHandler := handlers.NewProductHandler(productRepo)

	// Optional: R2 Storage (for image uploads)
	var uploadHandler *handlers.UploadHandler
	if cfg.R2AccountID != "" && cfg.R2AccessKeyID != "" {
		r2Storage, err := storage.NewR2Storage(cfg.R2AccountID, cfg.R2AccessKeyID, cfg.R2SecretKey, cfg.R2BucketName)
		if err != nil {
			log.Printf("⚠️  R2 Storage not configured: %v", err)
		} else {
			uploadHandler = handlers.NewUploadHandler(r2Storage, productRepo)
			log.Println("✅ Connected to Cloudflare R2")
		}
	} else {
		log.Println("⚠️  R2 Storage not configured (uploads disabled)")
	}

	// Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "GecoGreen API v0.1.0",
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

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"name": "GecoGreen API", "version": "0.1.0"})
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
