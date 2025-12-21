package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v76"

	"github.com/gecogreen/backend/internal/models"
	"github.com/gecogreen/backend/internal/repository"
	"github.com/gecogreen/backend/internal/services"
)

type OrderHandler struct {
	orderRepo     *repository.OrderRepository
	productRepo   *repository.ProductRepository
	userRepo      *repository.UserRepository
	stripeService *services.StripeService
	emailService  *services.EmailService
	frontendURL   string
}

func NewOrderHandler(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository, userRepo *repository.UserRepository, stripeService *services.StripeService, emailService *services.EmailService, frontendURL string) *OrderHandler {
	return &OrderHandler{
		orderRepo:     orderRepo,
		productRepo:   productRepo,
		userRepo:      userRepo,
		stripeService: stripeService,
		emailService:  emailService,
		frontendURL:   frontendURL,
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
	if req.Quantity > product.QuantityAvail {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":     "Quantità non disponibile",
			"available": product.QuantityAvail,
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
		fmt.Printf("❌ Order creation error: %v\n", err)
		if err == repository.ErrCannotOrder {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Account sospeso per troppi strike. Contatta il supporto.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore creazione ordine"})
	}
	fmt.Printf("✅ Order created: %s, totalAmount: %.2f\n", order.ID, totalAmount)

	// Handle free products (gifts) - no Stripe needed
	if totalAmount == 0 {
		// Mark as paid immediately for free items
		_ = h.orderRepo.MarkAsPaid(ctx, order.ID, "FREE_GIFT")
		successURL := fmt.Sprintf("%s/orders/%s/success", h.frontendURL, order.ID.String())
		return c.Status(fiber.StatusCreated).JSON(models.CheckoutResponse{
			OrderID:           order.ID,
			StripeCheckoutURL: successURL, // Direct to success page
			TotalAmount:       0,
			ExpiresAt:         time.Now().Add(30 * time.Minute),
		})
	}

	// Create Stripe Checkout Session for paid items
	successURL := fmt.Sprintf("%s/orders/%s/success?session_id={CHECKOUT_SESSION_ID}", h.frontendURL, order.ID.String())
	cancelURL := fmt.Sprintf("%s/orders/%s/cancel", h.frontendURL, order.ID.String())

	checkoutSession, err := h.stripeService.CreateCheckoutSession(order, product, successURL, cancelURL)
	if err != nil {
		// Log error but don't fail - return order ID so they can retry
		fmt.Printf("Stripe checkout error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":    "Errore creazione pagamento",
			"order_id": order.ID,
		})
	}

	// Save Stripe session ID to order
	order.StripeCheckoutSessionID = checkoutSession.ID
	_ = h.orderRepo.UpdateStripeSession(ctx, order.ID, checkoutSession.ID)

	return c.Status(fiber.StatusCreated).JSON(models.CheckoutResponse{
		OrderID:           order.ID,
		StripeCheckoutURL: checkoutSession.URL,
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
		fmt.Printf("ListMyOrders error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore recupero ordini: " + err.Error()})
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

// HandleStripeWebhook processes Stripe webhook events
// POST /api/v1/webhooks/stripe
func (h *OrderHandler) HandleStripeWebhook(c *fiber.Ctx) error {
	fmt.Printf("🔔 Stripe webhook received!\n")

	payload, err := io.ReadAll(c.Request().BodyStream())
	if err != nil {
		fmt.Printf("❌ Cannot read webhook body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot read body"})
	}

	fmt.Printf("📦 Webhook payload size: %d bytes\n", len(payload))

	signature := c.Get("Stripe-Signature")
	event, err := h.stripeService.VerifyWebhookSignature(payload, signature)
	if err != nil {
		fmt.Printf("❌ Webhook signature verification failed: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid signature"})
	}

	fmt.Printf("✅ Webhook verified, event type: %s\n", event.Type)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			fmt.Printf("Error parsing checkout session: %v\n", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event data"})
		}

		fmt.Printf("🔔 Webhook received: checkout.session.completed, session_id=%s\n", session.ID)

		orderID, err := uuid.Parse(session.Metadata["order_id"])
		if err != nil {
			fmt.Printf("Invalid order_id in metadata: %v\n", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order_id"})
		}

		// Get PaymentIntent ID safely (it may be nil or just an ID string in webhook)
		paymentIntentID := "webhook"
		if session.PaymentIntent != nil {
			paymentIntentID = session.PaymentIntent.ID
		}

		// Update order status to PAID
		err = h.orderRepo.MarkAsPaid(ctx, orderID, paymentIntentID)
		if err != nil {
			fmt.Printf("Error marking order as paid: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update order"})
		}

		fmt.Printf("✅ Order %s marked as PAID (PaymentIntent: %s)\n", orderID, session.PaymentIntent.ID)

		// Get order to decrement product quantity and send emails
		order, err := h.orderRepo.GetByID(ctx, orderID)
		if err != nil {
			fmt.Printf("Error getting order for quantity update: %v\n", err)
		} else {
			// Decrement product quantity
			if err := h.productRepo.DecrementQuantity(ctx, order.ProductID, order.Quantity); err != nil {
				fmt.Printf("Error decrementing product quantity: %v\n", err)
			} else {
				fmt.Printf("✅ Product %s quantity decremented by %d\n", order.ProductID, order.Quantity)
			}

			// Send email notifications
			go h.sendOrderEmails(order)
		}

	case "checkout.session.expired":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event data"})
		}

		orderID, err := uuid.Parse(session.Metadata["order_id"])
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order_id"})
		}

		// Cancel the order
		err = h.orderRepo.UpdateStatus(ctx, orderID, models.OrderCancelled)
		if err != nil {
			fmt.Printf("Error cancelling expired order: %v\n", err)
		}

		fmt.Printf("⏰ Order %s cancelled (checkout expired)\n", orderID)

	case "payment_intent.payment_failed":
		fmt.Printf("❌ Payment failed: %s\n", event.ID)
		// Could notify user here

	default:
		fmt.Printf("Unhandled event type: %s\n", event.Type)
	}

	return c.JSON(fiber.Map{"received": true})
}

// sendOrderEmails sends email notifications to buyer and seller
func (h *OrderHandler) sendOrderEmails(order *models.Order) {
	if h.emailService == nil {
		fmt.Println("Email service not configured, skipping notifications")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get product details
	product, err := h.productRepo.GetByID(ctx, order.ProductID)
	if err != nil {
		fmt.Printf("Error getting product for email: %v\n", err)
		return
	}

	// Get buyer info
	buyer, err := h.userRepo.GetByID(ctx, order.BuyerID)
	if err != nil {
		fmt.Printf("Error getting buyer for email: %v\n", err)
		return
	}

	// Get seller info
	seller, err := h.userRepo.GetByID(ctx, order.SellerID)
	if err != nil {
		fmt.Printf("Error getting seller for email: %v\n", err)
		return
	}

	// Send to buyer
	buyerName := buyer.FirstName
	if buyerName == "" {
		buyerName = "Cliente"
	}
	if err := h.emailService.SendOrderConfirmationToBuyer(order, product, buyer.Email, buyerName); err != nil {
		fmt.Printf("Error sending email to buyer: %v\n", err)
	} else {
		fmt.Printf("✅ Email sent to buyer: %s\n", buyer.Email)
	}

	// Send to seller
	sellerName := seller.FirstName
	if sellerName == "" {
		sellerName = "Venditore"
	}
	if err := h.emailService.SendNewOrderToSeller(order, product, seller.Email, sellerName, buyerName); err != nil {
		fmt.Printf("Error sending email to seller: %v\n", err)
	} else {
		fmt.Printf("✅ Email sent to seller: %s\n", seller.Email)
	}
}
