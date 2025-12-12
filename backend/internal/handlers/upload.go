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

// UploadHandler handles file upload endpoints
type UploadHandler struct {
	storage     *storage.R2Storage
	productRepo *repository.ProductRepository
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(storage *storage.R2Storage, productRepo *repository.ProductRepository) *UploadHandler {
	return &UploadHandler{
		storage:     storage,
		productRepo: productRepo,
	}
}

// UploadProductImage handles POST /api/v1/upload/product/:product_id/image
func (h *UploadHandler) UploadProductImage(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	// Parse product ID
	productIDStr := c.Params("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID prodotto non valido",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Check product ownership
	product, err := h.productRepo.GetByID(ctx, productID)
	if err != nil {
		if err == repository.ErrProductNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Prodotto non trovato",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Errore nel recupero del prodotto",
		})
	}

	if product.SellerID != user.ID && !user.IsAdmin() {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Non sei autorizzato a modificare questo prodotto",
		})
	}

	// Get file from form
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nessun file caricato",
		})
	}

	// Validate file size
	if file.Size > storage.MaxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File troppo grande. Massimo 5MB",
		})
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if !storage.IsValidImageType(contentType) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tipo di file non valido. Usa JPEG, PNG, GIF o WebP",
		})
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Errore nell'apertura del file",
		})
	}
	defer src.Close()

	// Upload to R2
	folder := "products/" + productID.String()
	imageURL, err := h.storage.Upload(ctx, src, file.Filename, contentType, folder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Errore nel caricamento del file",
		})
	}

	// Save to database
	if err := h.productRepo.AddImage(ctx, productID, imageURL); err != nil {
		// Try to delete the uploaded file
		_ = h.storage.Delete(ctx, imageURL)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Errore nel salvataggio dell'immagine",
		})
	}

	return c.JSON(fiber.Map{
		"url": imageURL,
	})
}

// UploadExpiryPhoto handles POST /api/v1/upload/product/:product_id/expiry-photo
func (h *UploadHandler) UploadExpiryPhoto(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	productIDStr := c.Params("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID prodotto non valido",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Check product ownership
	product, err := h.productRepo.GetByID(ctx, productID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Prodotto non trovato",
		})
	}

	if product.SellerID != user.ID && !user.IsAdmin() {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Non autorizzato",
		})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nessun file caricato",
		})
	}

	if file.Size > storage.MaxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File troppo grande. Massimo 5MB",
		})
	}

	contentType := file.Header.Get("Content-Type")
	if !storage.IsValidImageType(contentType) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tipo di file non valido",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Errore nell'apertura del file",
		})
	}
	defer src.Close()

	folder := "expiry/" + productID.String()
	imageURL, err := h.storage.Upload(ctx, src, file.Filename, contentType, folder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Errore nel caricamento",
		})
	}

	// TODO: Update product with expiry_photo_url

	return c.JSON(fiber.Map{
		"url": imageURL,
	})
}

// GetPresignedURL handles POST /api/v1/upload/presign
// Returns a presigned URL for direct client-side upload
func (h *UploadHandler) GetPresignedURL(c *fiber.Ctx) error {
	_ = c.Locals("user").(*models.User)

	var req struct {
		Filename    string `json:"filename"`
		ContentType string `json:"content_type"`
		Folder      string `json:"folder"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dati non validi",
		})
	}

	if !storage.IsValidImageType(req.ContentType) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tipo di file non valido",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	presignedURL, publicURL, err := h.storage.GeneratePresignedURL(ctx, req.Filename, req.ContentType, req.Folder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Errore nella generazione dell'URL",
		})
	}

	return c.JSON(fiber.Map{
		"upload_url": presignedURL,
		"public_url": publicURL,
	})
}
