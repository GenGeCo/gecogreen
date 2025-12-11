package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gecogreen/backend/internal/database"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db *database.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status   string            `json:"status"`
	Version  string            `json:"version"`
	Services map[string]string `json:"services"`
}

// Check handles GET /health
func (h *HealthHandler) Check(c *fiber.Ctx) error {
	services := h.db.Health(c.Context())

	// Determine overall status
	status := "healthy"
	for _, v := range services {
		if v != "ok" {
			status = "unhealthy"
			break
		}
	}

	response := HealthResponse{
		Status:   status,
		Version:  "0.1.0",
		Services: services,
	}

	if status == "unhealthy" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(response)
	}

	return c.JSON(response)
}

// Ping handles GET /ping (simple liveness probe)
func (h *HealthHandler) Ping(c *fiber.Ctx) error {
	return c.SendString("pong")
}
