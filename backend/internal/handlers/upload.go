package handlers

import (
	"context"
	"io"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/gecogreen/backend/internal/models"
	"github.com/gecogreen/backend/internal/moderation"
	"github.com/gecogreen/backend/internal/repository"
	"github.com/gecogreen/backend/internal/storage"
)

// UploadHandler handles file upload endpoints
type UploadHandler struct {
	storage         *storage.R2Storage
	productRepo     *repository.ProductRepository
	imageReviewRepo *repository.ImageReviewRepository
	ocrService      *moderation.OCRService
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(storage *storage.R2Storage, productRepo *repository.ProductRepository) *UploadHandler {
	return &UploadHandler{
		storage:     storage,
		productRepo: productRepo,
	}
}

// SetModerationServices sets the moderation services (optional)
func (h *UploadHandler) SetModerationServices(ocrService *moderation.OCRService, imageReviewRepo *repository.ImageReviewRepository) {
	h.ocrService = ocrService
	h.imageReviewRepo = imageReviewRepo
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

	if product.SellerID != user.ID && !user.IsAdmin {
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

	// Run OCR moderation if service is available
	if h.ocrService != nil && h.imageReviewRepo != nil {
		go h.moderateImage(user.ID, &productID, imageURL, "PRODUCT")
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

	if product.SellerID != user.ID && !user.IsAdmin {
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

// moderateImage runs OCR on the uploaded image and queues it for review if suspicious
func (h *UploadHandler) moderateImage(userID uuid.UUID, productID *uuid.UUID, imageURL string, imageType string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := h.ocrService.AnalyzeImageFromURL(ctx, imageURL)
	if err != nil {
		// Log error but don't fail - moderation is optional
		return
	}

	// If suspicious, add to review queue
	if result.IsSuspicious {
		review := &models.ImageReview{
			UserID:        userID,
			ProductID:     productID,
			ImageURL:      imageURL,
			ImageType:     imageType,
			DetectedText:  result.DetectedText,
			DetectedPhone: result.DetectedPhone,
			DetectedEmail: result.DetectedEmail,
			DetectedURL:   result.DetectedURL,
			Confidence:    result.Confidence,
		}
		_ = h.imageReviewRepo.Create(ctx, review)
	}
}

// UploadWithModeration uploads a file and runs moderation
func (h *UploadHandler) UploadWithModeration(ctx context.Context, reader io.Reader, filename, contentType, folder string, userID uuid.UUID, productID *uuid.UUID, imageType string) (string, error) {
	imageURL, err := h.storage.Upload(ctx, reader, filename, contentType, folder)
	if err != nil {
		return "", err
	}

	// Run moderation asynchronously
	if h.ocrService != nil && h.imageReviewRepo != nil {
		go h.moderateImage(userID, productID, imageURL, imageType)
	}

	return imageURL, nil
}
