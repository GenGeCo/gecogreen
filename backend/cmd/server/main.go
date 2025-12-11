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

	"github.com/gecogreen/backend/internal/config"
	"github.com/gecogreen/backend/internal/database"
	"github.com/gecogreen/backend/internal/handlers"
)

// @title GecoGreen API
// @version 0.1.0
// @description API per la piattaforma antispreco GecoGreen
// @host localhost:8080
// @BasePath /api/v1

func main() {
	// Load configuration
	cfg := config.Load()

	log.Printf("🌿 GecoGreen API Starting...")
	log.Printf("   Environment: %s", cfg.AppEnv)
	log.Printf("   Port: %s", cfg.Port)

	// Connect to databases
	db, err := database.New(cfg.DatabaseURL, cfg.RedisURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Connected to PostgreSQL")
	log.Println("✅ Connected to Redis")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "GecoGreen API v0.1.0",
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	// CORS - più permissivo in development
	app.Use(cors.New(cors.Config{
		AllowOrigins: func() string {
			if cfg.IsDevelopment() {
				return "*"
			}
			return "https://www.gecogreen.com,https://gecogreen.com"
		}(),
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Setup routes
	setupRoutes(app, db)

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("🛑 Shutting down server...")
		_ = app.Shutdown()
	}()

	// Start server
	log.Printf("🚀 Server running on http://localhost:%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}

func setupRoutes(app *fiber.App, db *database.DB) {
	// Health handlers
	healthHandler := handlers.NewHealthHandler(db)

	// Root routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":    "GecoGreen API",
			"version": "0.1.0",
			"docs":    "/api/v1/docs",
		})
	})

	app.Get("/ping", healthHandler.Ping)
	app.Get("/health", healthHandler.Check)

	// API v1
	v1 := app.Group("/api/v1")
	{
		// Health
		v1.Get("/health", healthHandler.Check)

		// Auth routes (TODO)
		// auth := v1.Group("/auth")
		// {
		// 	auth.Post("/register", authHandler.Register)
		// 	auth.Post("/login", authHandler.Login)
		// 	auth.Post("/refresh", authHandler.Refresh)
		// 	auth.Post("/logout", authHandler.Logout)
		// }

		// Products routes (TODO)
		// products := v1.Group("/products")
		// {
		// 	products.Get("/", productHandler.List)
		// 	products.Get("/:id", productHandler.Get)
		// 	products.Post("/", authMiddleware, productHandler.Create)
		// 	products.Put("/:id", authMiddleware, productHandler.Update)
		// 	products.Delete("/:id", authMiddleware, productHandler.Delete)
		// }

		// Orders routes (TODO)
		// Users routes (TODO)
		// Chat routes (TODO)
	}

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Not Found",
			"message": "La risorsa richiesta non esiste",
		})
	})
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   message,
		"code":    code,
	})
}
