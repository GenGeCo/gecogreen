package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/gecogreen/backend/internal/auth"
	"github.com/gecogreen/backend/internal/models"
	"github.com/gecogreen/backend/internal/repository"
)

type AuthHandler struct {
	userRepo   *repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewAuthHandler(userRepo *repository.UserRepository, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, jwtManager: jwtManager}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	// Validate email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email non valida"})
	}

	// Validate password
	if !auth.ValidatePassword(req.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password deve essere almeno 8 caratteri"})
	}

	// Validate name
	if req.FirstName == "" || req.LastName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nome e cognome sono obbligatori"})
	}

	// Validate city (required for all)
	if req.City == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Città obbligatoria"})
	}

	// Validate account type - default to PRIVATE
	if req.AccountType == "" {
		req.AccountType = models.AccountPrivate
	}

	// Validate business fields
	if req.AccountType == models.AccountBusiness {
		if req.BusinessName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ragione sociale obbligatoria per account aziendali"})
		}
		if req.VATNumber == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Partita IVA obbligatoria per account aziendali"})
		}
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	// Check if email exists
	exists, err := h.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore del server"})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email già registrata"})
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore del server"})
	}

	// Create user
	user := &models.User{
		Email:                req.Email,
		PasswordHash:         hashedPassword,
		FirstName:            req.FirstName,
		LastName:             req.LastName,
		AccountType:          req.AccountType,
		BusinessName:         req.BusinessName,
		VATNumber:            req.VATNumber,
		HasMultipleLocations: req.HasMultipleLocations,
		FiscalCode:           req.FiscalCode,
		SDICode:              req.SDICode,
		PECEmail:             req.PECEmail,
		EUVatID:              req.EUVatID,
		BillingCountry:       req.BillingCountry,
		City:                 req.City,
		Province:             req.Province,
		PostalCode:           req.PostalCode,
	}

	if err := h.userRepo.Create(ctx, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nella creazione account"})
	}

	// Create primary location if address provided
	if req.AddressStreet != "" {
		loc := &models.Location{
			UserID:          user.ID,
			Name:            "Sede principale",
			IsPrimary:       true,
			AddressStreet:   req.AddressStreet,
			AddressCity:     req.City,
			AddressProvince: req.Province,
			AddressPostal:   req.PostalCode,
		}
		// Don't fail registration if location creation fails
		_ = h.userRepo.CreateLocation(ctx, loc)
	}

	// Generate tokens
	accessToken, _ := h.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.AccountType))
	refreshToken, _ := h.jwtManager.GenerateRefreshToken(user.ID, user.Email, string(user.AccountType))

	return c.Status(fiber.StatusCreated).JSON(models.AuthResponse{
		User: *user, AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: h.jwtManager.GetAccessTokenExpiry(),
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	user, err := h.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email o password non corretti"})
	}

	if !user.IsActive() {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Account non attivo"})
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email o password non corretti"})
	}

	_ = h.userRepo.UpdateLastLogin(ctx, user.ID)

	accessToken, _ := h.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.AccountType))
	refreshToken, _ := h.jwtManager.GenerateRefreshToken(user.ID, user.Email, string(user.AccountType))

	return c.JSON(models.AuthResponse{
		User: *user, AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: h.jwtManager.GetAccessTokenExpiry(),
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req models.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	claims, err := h.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token non valido"})
	}

	if claims.TokenType != auth.RefreshToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token type non valido"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	user, err := h.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Utente non trovato"})
	}

	if !user.IsActive() {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Account non attivo"})
	}

	accessToken, _ := h.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.AccountType))
	refreshToken, _ := h.jwtManager.GenerateRefreshToken(user.ID, user.Email, string(user.AccountType))

	return c.JSON(models.AuthResponse{
		User: *user, AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: h.jwtManager.GetAccessTokenExpiry(),
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	return c.JSON(user)
}
