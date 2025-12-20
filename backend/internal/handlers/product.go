package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/gecogreen/backend/internal/models"
	"github.com/gecogreen/backend/internal/repository"
)

type ProductHandler struct {
	productRepo *repository.ProductRepository
}

func NewProductHandler(productRepo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{productRepo: productRepo}
}

func (h *ProductHandler) List(c *fiber.Ctx) error {
	filters := models.ProductFilters{
		Page:      c.QueryInt("page", 1),
		PerPage:   c.QueryInt("per_page", 20),
		SortBy:    c.Query("sort_by", "created_at"),
		SortOrder: c.Query("sort_order", "desc"),
	}

	if search := c.Query("search"); search != "" {
		filters.Search = &search
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := uuid.Parse(categoryID); err == nil {
			filters.CategoryID = &id
		}
	}
	if listingType := c.Query("listing_type"); listingType != "" {
		lt := models.ListingType(listingType)
		filters.ListingType = &lt
	}
	if minPrice := c.Query("min_price"); minPrice != "" {
		if p, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filters.MinPrice = &p
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if p, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filters.MaxPrice = &p
		}
	}
	if city := c.Query("city"); city != "" {
		filters.City = &city
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	response, err := h.productRepo.List(ctx, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nel recupero prodotti"})
	}
	return c.JSON(response)
}

func (h *ProductHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	product, err := h.productRepo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrProductNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Prodotto non trovato"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nel recupero"})
	}

	if product.IsDutchAuction {
		product.Price = product.GetCurrentPrice()
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = h.productRepo.IncrementViewCount(bgCtx, id)
	}()

	return c.JSON(product)
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	// All authenticated users can create products now

	var req models.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Titolo obbligatorio"})
	}
	if req.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Descrizione obbligatoria"})
	}
	if req.Price < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Prezzo non valido"})
	}
	if req.Quantity < 1 {
		req.Quantity = 1
	}

	product := &models.Product{
		SellerID:            user.ID,
		Title:               req.Title,
		Description:         req.Description,
		CategoryID:          req.CategoryID,
		Price:               req.Price,
		OriginalPrice:       req.OriginalPrice,
		Quantity:            req.Quantity,
		ListingType:         req.ListingType,
		ShippingMethod:      req.ShippingMethod,
		ShippingCost:        req.ShippingCost,
		PickupLocationIDs:   req.PickupLocationIDs,
		ExpiryDate:          req.ExpiryDate,
		IsDutchAuction:      req.IsDutchAuction,
		DutchStartPrice:     req.DutchStartPrice,
		DutchDecreaseAmount: req.DutchDecreaseAmount,
		DutchDecreaseHours:  req.DutchDecreaseHours,
		DutchMinPrice:       req.DutchMinPrice,
		City:                req.City,
		Province:            req.Province,
		PostalCode:          req.PostalCode,
		Latitude:            req.Latitude,
		Longitude:           req.Longitude,
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	if err := h.productRepo.Create(ctx, product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nella creazione"})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	product, err := h.productRepo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrProductNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Prodotto non trovato"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore"})
	}

	if product.SellerID != user.ID && !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Non autorizzato"})
	}

	var req models.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	if req.Title != nil {
		product.Title = *req.Title
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Quantity != nil {
		product.Quantity = *req.Quantity
		product.QuantityAvail = *req.Quantity
	}
	if req.Status != nil {
		product.Status = *req.Status
	}
	if req.ShippingMethod != nil {
		product.ShippingMethod = *req.ShippingMethod
	}
	if req.ShippingCost != nil {
		product.ShippingCost = *req.ShippingCost
	}

	if err := h.productRepo.Update(ctx, product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nell'aggiornamento"})
	}

	return c.JSON(product)
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	product, err := h.productRepo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrProductNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Prodotto non trovato"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore"})
	}

	if product.SellerID != user.ID && !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Non autorizzato"})
	}

	if err := h.productRepo.Delete(ctx, id, product.SellerID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nell'eliminazione"})
	}

	return c.JSON(fiber.Map{"message": "Prodotto eliminato"})
}

func (h *ProductHandler) MyProducts(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	filters := models.ProductFilters{
		SellerID:  &user.ID,
		Page:      c.QueryInt("page", 1),
		PerPage:   c.QueryInt("per_page", 20),
		SortBy:    c.Query("sort_by", "created_at"),
		SortOrder: c.Query("sort_order", "desc"),
	}

	if status := c.Query("status"); status != "" {
		s := models.ProductStatus(status)
		filters.Status = &s
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	response, err := h.productRepo.List(ctx, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore"})
	}

	return c.JSON(response)
}
