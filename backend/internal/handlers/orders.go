package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/gecogreen/backend/internal/models"
	"github.com/gecogreen/backend/internal/repository"
)

type OrderHandler struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
	userRepo    *repository.UserRepository
}

func NewOrderHandler(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository, userRepo *repository.UserRepository) *OrderHandler {
	return &OrderHandler{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
	}
}

// CreateOrder creates a new order and returns checkout URL
// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	var req models.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	// Validate required fields
	if req.ProductID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Product ID richiesto"})
	}
	if req.Quantity < 1 {
		req.Quantity = 1
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	// Get product
	product, err := h.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Prodotto non trovato"})
	}

	// Can't buy your own product
	if product.SellerID == user.ID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Non puoi comprare il tuo prodotto"})
	}

	// Check quantity available
	if req.Quantity > product.QuantityAvailable {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":     "Quantità non disponibile",
			"available": product.QuantityAvailable,
		})
	}

	// Calculate shipping cost
	shippingCost := 0.0
	if req.DeliveryType == models.DeliverySellerShips {
		shippingCost = product.ShippingCost
		// Validate shipping address
		if req.ShippingAddress == "" || req.ShippingCity == "" || req.ShippingPostalCode == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Indirizzo di spedizione richiesto"})
		}
	}

	// Calculate totals
	unitPrice := product.Price
	if product.IsDutchAuction {
		unitPrice = product.GetCurrentPrice()
	}
	totalAmount := (unitPrice * float64(req.Quantity)) + shippingCost

	// Calculate eco impact (using category estimates)
	co2Saved := 2.0 * float64(req.Quantity)    // Default 2kg per item
	waterSaved := 500.0 * float64(req.Quantity) // Default 500L per item

	// EcoCredits: 10 per euro spent
	ecoCreditsBuyer := int(totalAmount * 10)
	ecoCreditsSeller := int(totalAmount * 15) // Seller gets more for selling

	// Create order
	order := &models.Order{
		BuyerID:     user.ID,
		SellerID:    product.SellerID,
		ProductID:   product.ID,
		Quantity:    req.Quantity,
		UnitPrice:   unitPrice,
		ShippingCost: shippingCost,
		TotalAmount: totalAmount,
		DeliveryType: req.DeliveryType,

		// Shipping
		ShippingAddress:    req.ShippingAddress,
		ShippingCity:       req.ShippingCity,
		ShippingProvince:   req.ShippingProvince,
		ShippingPostalCode: req.ShippingPostalCode,
		ShippingCountry:    req.ShippingCountry,

		// Pickup
		PickupLocationID: req.PickupLocationID,

		// Impact
		CO2Saved:         co2Saved,
		WaterSaved:       waterSaved,
		EcoCreditsBuyer:  ecoCreditsBuyer,
		EcoCreditsSeller: ecoCreditsSeller,

		BuyerNotes: req.BuyerNotes,
	}

	// If pickup, get location details
	if req.DeliveryType == models.DeliveryPickup && req.PickupLocationID != nil {
		location, err := h.userRepo.GetLocationByID(ctx, *req.PickupLocationID)
		if err == nil {
			order.PickupAddress = location.AddressStreet + ", " + location.AddressCity
			order.PickupInstructions = location.PickupInstructions
		}
	}

	err = h.orderRepo.Create(ctx, order)
	if err != nil {
		if err == repository.ErrCannotOrder {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Account sospeso per troppi strike. Contatta il supporto.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore creazione ordine"})
	}

	// TODO: Create Stripe Checkout Session
	// For now, return a mock checkout URL
	checkoutURL := "/checkout/" + order.ID.String()

	return c.Status(fiber.StatusCreated).JSON(models.CheckoutResponse{
		OrderID:           order.ID,
		StripeCheckoutURL: checkoutURL,
		TotalAmount:       totalAmount,
		ExpiresAt:         time.Now().Add(30 * time.Minute),
	})
}

// GetOrder returns order details
// GET /api/v1/orders/:id
func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	order, err := h.orderRepo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrOrderNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ordine non trovato"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore recupero ordine"})
	}

	// Only buyer, seller, or admin can see order
	if order.BuyerID != user.ID && order.SellerID != user.ID && !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Accesso negato"})
	}

	// Hide sensitive info based on status and role
	if order.Status == models.OrderPending && order.BuyerID == user.ID {
		// Buyer hasn't paid yet, hide pickup address
		order.PickupAddress = ""
		order.PickupInstructions = ""
	}

	// Hide QR code from seller (buyer shows it)
	if order.SellerID == user.ID {
		order.QRCodeToken = ""
	}

	return c.JSON(order)
}

// ListMyOrders returns orders for the current user (as buyer)
// GET /api/v1/orders
func (h *OrderHandler) ListMyOrders(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	filters := models.OrderFilters{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
	}

	if status := c.Query("status"); status != "" {
		s := models.OrderStatus(status)
		filters.Status = &s
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	response, err := h.orderRepo.ListByBuyer(ctx, user.ID, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore recupero ordini"})
	}

	return c.JSON(response)
}

// ListSellerOrders returns orders for the current user (as seller)
// GET /api/v1/orders/seller
func (h *OrderHandler) ListSellerOrders(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	filters := models.OrderFilters{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
	}

	if status := c.Query("status"); status != "" {
		s := models.OrderStatus(status)
		filters.Status = &s
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	response, err := h.orderRepo.ListBySeller(ctx, user.ID, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore recupero ordini"})
	}

	return c.JSON(response)
}

// UpdateOrderStatus updates order status (seller)
// PUT /api/v1/orders/:id/status
func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	var req models.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	// Get order
	order, err := h.orderRepo.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ordine non trovato"})
	}

	// Only seller can update
	if order.SellerID != user.ID && !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Accesso negato"})
	}

	// Update tracking if shipping
	if req.TrackingNumber != "" {
		err = h.orderRepo.UpdateTracking(ctx, id, req.TrackingNumber, req.TrackingURL, req.ShippingCarrier)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore aggiornamento tracking"})
		}
	}

	// Update status if provided
	if req.Status != "" {
		// Validate status transition
		allowedTransitions := map[models.OrderStatus][]models.OrderStatus{
			models.OrderPaid:           {models.OrderProcessing, models.OrderReadyForPickup, models.OrderShipped},
			models.OrderProcessing:     {models.OrderReadyForPickup, models.OrderShipped},
			models.OrderReadyForPickup: {models.OrderDelivered},
			models.OrderShipped:        {models.OrderInTransit, models.OrderDelivered},
			models.OrderInTransit:      {models.OrderDelivered},
			models.OrderDelivered:      {models.OrderCompleted},
		}

		allowed := false
		for _, s := range allowedTransitions[order.Status] {
			if s == req.Status {
				allowed = true
				break
			}
		}

		if !allowed && !user.IsAdmin {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":          "Transizione stato non permessa",
				"current_status": order.Status,
				"requested":      req.Status,
			})
		}

		err = h.orderRepo.UpdateStatus(ctx, id, req.Status)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore aggiornamento stato"})
		}
	}

	// Return updated order
	updatedOrder, _ := h.orderRepo.GetByID(ctx, id)
	return c.JSON(updatedOrder)
}

// ConfirmPickup confirms pickup via QR code (seller scans)
// POST /api/v1/orders/confirm-pickup
func (h *OrderHandler) ConfirmPickup(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	var req models.ConfirmPickupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	if req.QRCodeToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "QR code richiesto"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	order, err := h.orderRepo.ConfirmPickup(ctx, req.QRCodeToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Verify seller
	if order.SellerID != user.ID && !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Solo il seller può confermare"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Ritiro confermato! Payout in 48 ore.",
		"order":   order,
	})
}

// CancelOrder cancels an order
// POST /api/v1/orders/:id/cancel
func (h *OrderHandler) CancelOrder(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	var body struct {
		Reason string `json:"reason"`
	}
	c.BodyParser(&body)

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	// Get order first
	order, err := h.orderRepo.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ordine non trovato"})
	}

	// Only buyer, seller, or admin can cancel
	if order.BuyerID != user.ID && order.SellerID != user.ID && !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Accesso negato"})
	}

	err = h.orderRepo.CancelOrder(ctx, id, user.ID, body.Reason)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Ordine cancellato",
	})
}

// OpenDispute opens a dispute for an order
// POST /api/v1/orders/:id/dispute
func (h *OrderHandler) OpenDispute(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	var req models.OpenDisputeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	if len(req.Description) < 50 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Descrizione troppo breve (min 50 caratteri)"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	// Get order
	order, err := h.orderRepo.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ordine non trovato"})
	}

	// Only buyer can open dispute (for now)
	if order.BuyerID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Solo il buyer può aprire disputa"})
	}

	// Can only dispute delivered orders
	if order.Status != models.OrderDelivered && order.Status != models.OrderPaid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Non puoi aprire disputa per questo ordine"})
	}

	dispute := &models.Dispute{
		OrderID:      id,
		OpenedBy:     user.ID,
		Reason:       req.Reason,
		Description:  req.Description,
		EvidenceURLs: req.EvidenceURLs,
	}

	err = h.orderRepo.CreateDispute(ctx, dispute)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore apertura disputa"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Disputa aperta. Il seller ha 48h per rispondere.",
		"dispute": dispute,
	})
}

// GetDispute gets dispute for an order
// GET /api/v1/orders/:id/dispute
func (h *OrderHandler) GetDispute(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	// Get order first
	order, err := h.orderRepo.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ordine non trovato"})
	}

	// Only buyer, seller, or admin can see dispute
	if order.BuyerID != user.ID && order.SellerID != user.ID && !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Accesso negato"})
	}

	dispute, err := h.orderRepo.GetDisputeByOrderID(ctx, id)
	if err != nil {
		if err == repository.ErrDisputeNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Nessuna disputa"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore"})
	}

	return c.JSON(dispute)
}

// CreateReview creates a review for an order
// POST /api/v1/orders/:id/review
func (h *OrderHandler) CreateReview(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	var req models.CreateReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	if req.Rating < 1 || req.Rating > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Rating deve essere 1-5"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	// Get order
	order, err := h.orderRepo.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ordine non trovato"})
	}

	// Only completed orders
	if order.Status != models.OrderCompleted && order.Status != models.OrderDelivered {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Puoi recensire solo ordini completati"})
	}

	// Determine who is reviewing whom
	var reviewedID uuid.UUID
	if order.BuyerID == user.ID {
		reviewedID = order.SellerID // Buyer reviews seller
	} else if order.SellerID == user.ID {
		reviewedID = order.BuyerID // Seller reviews buyer
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Non puoi recensire"})
	}

	review := &models.OrderReview{
		OrderID:     id,
		ReviewerID:  user.ID,
		ReviewedID:  reviewedID,
		Rating:      req.Rating,
		Comment:     req.Comment,
		IsAnonymous: req.IsAnonymous,
	}

	err = h.orderRepo.CreateReview(ctx, review)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore creazione recensione"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Recensione pubblicata",
		"review":  review,
	})
}

// GetQRCode returns QR code data for pickup
// GET /api/v1/orders/:id/qr
func (h *OrderHandler) GetQRCode(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	order, err := h.orderRepo.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ordine non trovato"})
	}

	// Only buyer can see QR code
	if order.BuyerID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Accesso negato"})
	}

	// Only for pickup orders that are paid
	if order.DeliveryType != models.DeliveryPickup {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "QR code solo per ritiro"})
	}

	if order.Status != models.OrderPaid && order.Status != models.OrderReadyForPickup {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ordine non ancora pronto"})
	}

	return c.JSON(fiber.Map{
		"qr_code_token": order.QRCodeToken,
		"expires_at":    order.QRCodeExpiresAt,
		"pickup_address": order.PickupAddress,
		"pickup_instructions": order.PickupInstructions,
		"pickup_deadline": order.PickupDeadline,
	})
}
