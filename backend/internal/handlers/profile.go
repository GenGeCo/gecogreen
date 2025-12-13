package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/gecogreen/backend/internal/models"
	"github.com/gecogreen/backend/internal/repository"
	"github.com/gecogreen/backend/internal/storage"
)

var _ = storage.MaxFileSize // ensure storage import is used

type ProfileHandler struct {
	userRepo  *repository.UserRepository
	r2Storage *storage.R2Storage
}

func NewProfileHandler(userRepo *repository.UserRepository, r2Storage *storage.R2Storage) *ProfileHandler {
	return &ProfileHandler{userRepo: userRepo, r2Storage: r2Storage}
}

// GetProfile returns the current user's full profile
func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	// Get user's locations
	locations, err := h.userRepo.GetLocationsByUser(ctx, user.ID)
	if err != nil {
		locations = []models.Location{}
	}

	return c.JSON(fiber.Map{
		"user":      user,
		"locations": locations,
	})
}

// UpdateProfile updates the current user's profile
func (h *ProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	var req models.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	if err := h.userRepo.UpdateProfile(ctx, user.ID, &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nell'aggiornamento"})
	}

	// Return updated user
	updatedUser, err := h.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return c.JSON(fiber.Map{"success": true})
	}

	return c.JSON(updatedUser)
}

// UploadAvatar uploads user's avatar image
func (h *ProfileHandler) UploadAvatar(c *fiber.Ctx) error {
	if h.r2Storage == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Upload non disponibile"})
	}

	user := c.Locals("user").(*models.User)

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File richiesto"})
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Solo immagini JPG, PNG o WebP"})
	}

	// Max 2MB for avatars
	if file.Size > 2*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Immagine troppo grande (max 2MB)"})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore lettura file"})
	}
	defer src.Close()

	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Upload to R2
	folder := "avatars/" + user.ID.String()
	url, err := h.r2Storage.Upload(ctx, src, file.Header.Filename, contentType, folder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore upload"})
	}

	// Update user's avatar URL
	if err := h.userRepo.UpdateAvatar(ctx, user.ID, url); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore salvataggio"})
	}

	return c.JSON(fiber.Map{"avatar_url": url})
}

// UploadBusinessPhoto uploads a business photo
func (h *ProfileHandler) UploadBusinessPhoto(c *fiber.Ctx) error {
	if h.r2Storage == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Upload non disponibile"})
	}

	user := c.Locals("user").(*models.User)

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File richiesto"})
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Solo immagini JPG, PNG o WebP"})
	}

	// Max 5MB for business photos
	if file.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Immagine troppo grande (max 5MB)"})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore lettura file"})
	}
	defer src.Close()

	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Upload to R2
	folder := "business-photos/" + user.ID.String()
	url, err := h.r2Storage.Upload(ctx, src, file.Header.Filename, contentType, folder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore upload"})
	}

	// Add to user's business photos
	if err := h.userRepo.AddBusinessPhoto(ctx, user.ID, url); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore salvataggio"})
	}

	return c.JSON(fiber.Map{"photo_url": url})
}

// GetLocations returns user's pickup locations
func (h *ProfileHandler) GetLocations(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	locations, err := h.userRepo.GetLocationsByUser(ctx, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore caricamento sedi"})
	}

	return c.JSON(locations)
}

// CreateLocation adds a new pickup location
func (h *ProfileHandler) CreateLocation(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	var req models.CreateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	// Validate required fields
	if req.Name == "" || req.AddressStreet == "" || req.AddressCity == "" || req.AddressPostal == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nome, indirizzo, città e CAP obbligatori"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	loc := &models.Location{
		UserID:             user.ID,
		Name:               req.Name,
		IsPrimary:          req.IsPrimary,
		AddressStreet:      req.AddressStreet,
		AddressCity:        req.AddressCity,
		AddressProvince:    req.AddressProvince,
		AddressPostal:      req.AddressPostal,
		Phone:              req.Phone,
		Email:              req.Email,
		PickupHours:        req.PickupHours,
		PickupInstructions: req.PickupInstructions,
	}

	if err := h.userRepo.CreateLocation(ctx, loc); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore creazione sede"})
	}

	return c.Status(fiber.StatusCreated).JSON(loc)
}

// DeleteLocation removes a pickup location
func (h *ProfileHandler) DeleteLocation(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	locID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	// Verify ownership
	loc, err := h.userRepo.GetLocationByID(ctx, locID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sede non trovata"})
	}
	if loc.UserID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Non autorizzato"})
	}

	if err := h.userRepo.DeleteLocation(ctx, locID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore eliminazione"})
	}

	return c.JSON(fiber.Map{"success": true})
}
