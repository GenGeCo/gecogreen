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

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email non valida"})
	}

	if !auth.ValidatePassword(req.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password deve essere almeno 8 caratteri"})
	}

	if req.Role != models.RoleBuyer && req.Role != models.RoleSeller {
		req.Role = models.RoleBuyer
	}

	if req.FirstName == "" || req.LastName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nome e cognome sono obbligatori"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	exists, err := h.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore del server"})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email già registrata"})
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore del server"})
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Roles:        []models.UserRole{req.Role},
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	if err := h.userRepo.Create(ctx, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nella creazione account"})
	}

	accessToken, _ := h.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Roles[0]))
	refreshToken, _ := h.jwtManager.GenerateRefreshToken(user.ID, user.Email, string(user.Roles[0]))

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

	role := string(models.RoleBuyer)
	if len(user.Roles) > 0 {
		role = string(user.Roles[0])
	}

	accessToken, _ := h.jwtManager.GenerateAccessToken(user.ID, user.Email, role)
	refreshToken, _ := h.jwtManager.GenerateRefreshToken(user.ID, user.Email, role)

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

	role := string(models.RoleBuyer)
	if len(user.Roles) > 0 {
		role = string(user.Roles[0])
	}

	accessToken, _ := h.jwtManager.GenerateAccessToken(user.ID, user.Email, role)
	refreshToken, _ := h.jwtManager.GenerateRefreshToken(user.ID, user.Email, role)

	return c.JSON(models.AuthResponse{
		User: *user, AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: h.jwtManager.GetAccessTokenExpiry(),
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	return c.JSON(user)
}
