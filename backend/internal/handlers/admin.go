package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/gecogreen/backend/internal/models"
	"github.com/gecogreen/backend/internal/repository"
)

type AdminHandler struct {
	userRepo        *repository.UserRepository
	imageReviewRepo *repository.ImageReviewRepository
}

func NewAdminHandler(userRepo *repository.UserRepository, imageReviewRepo *repository.ImageReviewRepository) *AdminHandler {
	return &AdminHandler{
		userRepo:        userRepo,
		imageReviewRepo: imageReviewRepo,
	}
}

// GetPendingReviews returns all pending image reviews
func (h *AdminHandler) GetPendingReviews(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	reviews, total, err := h.imageReviewRepo.GetPending(ctx, perPage, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore caricamento"})
	}

	totalPages := (total + perPage - 1) / perPage

	return c.JSON(fiber.Map{
		"reviews":     reviews,
		"total":       total,
		"page":        page,
		"per_page":    perPage,
		"total_pages": totalPages,
	})
}

// GetReviewStats returns moderation statistics
func (h *AdminHandler) GetReviewStats(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	stats, err := h.imageReviewRepo.GetStats(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore statistiche"})
	}

	return c.JSON(stats)
}

// ApproveReview approves an image review
func (h *AdminHandler) ApproveReview(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)

	reviewID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	if err := h.imageReviewRepo.Approve(ctx, reviewID, admin.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore approvazione"})
	}

	return c.JSON(fiber.Map{"success": true, "status": "APPROVED"})
}

// RejectReview rejects an image review
func (h *AdminHandler) RejectReview(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)

	reviewID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	if err := h.imageReviewRepo.Reject(ctx, reviewID, admin.ID, req.Reason); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore rifiuto"})
	}

	return c.JSON(fiber.Map{"success": true, "status": "REJECTED"})
}

// GetReviewDetail returns details of a single review
func (h *AdminHandler) GetReviewDetail(c *fiber.Ctx) error {
	reviewID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	review, err := h.imageReviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Review non trovata"})
	}

	// Get user info
	user, _ := h.userRepo.GetByID(ctx, review.UserID)

	return c.JSON(fiber.Map{
		"review": review,
		"user":   user,
	})
}
