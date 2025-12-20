package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/gecogreen/backend/internal/models"
	"github.com/gecogreen/backend/internal/repository"
)

type CategoryHandler struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryHandler(categoryRepo *repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{
		categoryRepo: categoryRepo,
	}
}

// List returns all categories
// GET /api/v1/categories
func (h *CategoryHandler) List(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	categories, err := h.categoryRepo.List(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore recupero categorie"})
	}

	if categories == nil {
		categories = []models.Category{}
	}

	return c.JSON(categories)
}

// GetByID returns a category by ID
// GET /api/v1/categories/:id
func (h *CategoryHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	category, err := h.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Categoria non trovata"})
	}

	return c.JSON(category)
}
