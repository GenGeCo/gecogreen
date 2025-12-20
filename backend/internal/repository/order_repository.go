package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gecogreen/backend/internal/models"
)

var (
	ErrOrderNotFound   = errors.New("order not found")
	ErrDisputeNotFound = errors.New("dispute not found")
	ErrCannotOrder     = errors.New("user cannot place orders (too many strikes)")
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

// =====================
// ORDER CRUD
// =====================

// Create creates a new order
func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	// Check if user can order
	var canOrder bool
	err := r.pool.QueryRow(ctx, "SELECT can_user_order($1)", order.BuyerID).Scan(&canOrder)
	if err == nil && !canOrder {
		return ErrCannotOrder
	}

	// Generate QR token for pickup orders
	if order.DeliveryType == models.DeliveryPickup {
		token := make([]byte, 32)
		rand.Read(token)
		order.QRCodeToken = hex.EncodeToString(token)
		expires := time.Now().AddDate(0, 0, 7) // 7 days
		order.QRCodeExpiresAt = &expires

		// Set pickup deadline
		deadline := time.Now().AddDate(0, 0, 7)
		order.PickupDeadline = &deadline
	}

	// Calculate fees
	order.PlatformFee = round(order.TotalAmount * 0.10)
	order.StripeFee = round(order.TotalAmount*0.014 + 0.25)
	order.SellerPayout = order.TotalAmount - order.PlatformFee - order.StripeFee

	order.ID = uuid.New()
	order.Status = models.OrderPending
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	query := `
		INSERT INTO orders (
			id, buyer_id, seller_id, product_id,
			quantity, unit_price, shipping_cost, total_amount,
			platform_fee, stripe_fee, seller_payout,
			status, delivery_type,
			pickup_location_id, pickup_address, pickup_instructions, pickup_deadline,
			shipping_address, shipping_city, shipping_province, shipping_postal_code, shipping_country,
			qr_code_token, qr_code_expires_at,
			co2_saved, water_saved, eco_credits_buyer, eco_credits_seller,
			buyer_notes,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8,
			$9, $10, $11,
			$12::order_status, $13::delivery_type,
			$14, $15, $16, $17,
			$18, $19, $20, $21, $22,
			$23, $24,
			$25, $26, $27, $28,
			$29,
			$30, $31
		)
	`

	_, err = r.pool.Exec(ctx, query,
		order.ID, order.BuyerID, order.SellerID, order.ProductID,
		order.Quantity, order.UnitPrice, order.ShippingCost, order.TotalAmount,
		order.PlatformFee, order.StripeFee, order.SellerPayout,
		string(order.Status), string(order.DeliveryType),
		order.PickupLocationID, order.PickupAddress, order.PickupInstructions, order.PickupDeadline,
		order.ShippingAddress, order.ShippingCity, order.ShippingProvince, order.ShippingPostalCode, order.ShippingCountry,
		order.QRCodeToken, order.QRCodeExpiresAt,
		order.CO2Saved, order.WaterSaved, order.EcoCreditsBuyer, order.EcoCreditsSeller,
		order.BuyerNotes,
		order.CreatedAt, order.UpdatedAt,
	)

	return err
}

// GetByID retrieves an order by ID with all related data
func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	query := `
		SELECT
			o.id, o.buyer_id, o.seller_id, o.product_id,
			o.quantity, o.unit_price, COALESCE(o.shipping_cost, 0), o.total_amount,
			COALESCE(o.platform_fee, 0), COALESCE(o.stripe_fee, 0), COALESCE(o.seller_payout, 0),
			o.status::text, o.delivery_type::text,
			o.pickup_location_id, COALESCE(o.pickup_address, ''), COALESCE(o.pickup_instructions, ''), o.pickup_deadline,
			COALESCE(o.shipping_address, ''), COALESCE(o.shipping_city, ''), COALESCE(o.shipping_province, ''),
			COALESCE(o.shipping_postal_code, ''), COALESCE(o.shipping_country, 'IT'),
			COALESCE(o.tracking_number, ''), COALESCE(o.tracking_url, ''), COALESCE(o.shipping_carrier, ''), o.shipped_at,
			COALESCE(o.qr_code_token, ''), o.qr_code_expires_at, o.qr_scanned_at,
			COALESCE(o.stripe_payment_intent_id, ''), COALESCE(o.stripe_checkout_session_id, ''), o.paid_at,
			o.payout_scheduled_at, o.payout_completed_at,
			COALESCE(o.co2_saved, 0), COALESCE(o.water_saved, 0), COALESCE(o.eco_credits_buyer, 0), COALESCE(o.eco_credits_seller, 0),
			COALESCE(o.buyer_notes, ''), COALESCE(o.seller_notes, ''),
			o.created_at, o.updated_at, o.completed_at, o.cancelled_at, COALESCE(o.cancellation_reason, ''),
			-- Buyer
			b.id, COALESCE(b.business_name, ''), b.first_name, b.last_name, COALESCE(b.avatar_url, ''), COALESCE(b.city, ''),
			-- Seller
			s.id, COALESCE(s.business_name, ''), s.first_name, s.last_name, COALESCE(s.avatar_url, ''), COALESCE(s.city, ''),
			-- Product
			p.id, p.title, COALESCE(p.images->>0, ''), p.price
		FROM orders o
		JOIN users b ON o.buyer_id = b.id
		JOIN users s ON o.seller_id = s.id
		JOIN products p ON o.product_id = p.id
		WHERE o.id = $1
	`

	order := &models.Order{
		Buyer:   &models.UserPublicMinimal{},
		Seller:  &models.UserPublicMinimal{},
		Product: &models.ProductMinimal{},
	}
	var status, deliveryType string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&order.ID, &order.BuyerID, &order.SellerID, &order.ProductID,
		&order.Quantity, &order.UnitPrice, &order.ShippingCost, &order.TotalAmount,
		&order.PlatformFee, &order.StripeFee, &order.SellerPayout,
		&status, &deliveryType,
		&order.PickupLocationID, &order.PickupAddress, &order.PickupInstructions, &order.PickupDeadline,
		&order.ShippingAddress, &order.ShippingCity, &order.ShippingProvince,
		&order.ShippingPostalCode, &order.ShippingCountry,
		&order.TrackingNumber, &order.TrackingURL, &order.ShippingCarrier, &order.ShippedAt,
		&order.QRCodeToken, &order.QRCodeExpiresAt, &order.QRScannedAt,
		&order.StripePaymentIntentID, &order.StripeCheckoutSessionID, &order.PaidAt,
		&order.PayoutScheduledAt, &order.PayoutCompletedAt,
		&order.CO2Saved, &order.WaterSaved, &order.EcoCreditsBuyer, &order.EcoCreditsSeller,
		&order.BuyerNotes, &order.SellerNotes,
		&order.CreatedAt, &order.UpdatedAt, &order.CompletedAt, &order.CancelledAt, &order.CancellationReason,
		// Buyer
		&order.Buyer.ID, &order.Buyer.BusinessName, &order.Buyer.FirstName, &order.Buyer.LastName, &order.Buyer.AvatarURL, &order.Buyer.City,
		// Seller
		&order.Seller.ID, &order.Seller.BusinessName, &order.Seller.FirstName, &order.Seller.LastName, &order.Seller.AvatarURL, &order.Seller.City,
		// Product
		&order.Product.ID, &order.Product.Title, &order.Product.MainImageURL, &order.Product.Price,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	order.Status = models.OrderStatus(status)
	order.DeliveryType = models.DeliveryType(deliveryType)

	return order, nil
}

// ListByBuyer returns orders for a buyer
func (r *OrderRepository) ListByBuyer(ctx context.Context, buyerID uuid.UUID, filters models.OrderFilters) (*models.OrdersListResponse, error) {
	return r.listOrders(ctx, &buyerID, nil, filters)
}

// ListBySeller returns orders for a seller
func (r *OrderRepository) ListBySeller(ctx context.Context, sellerID uuid.UUID, filters models.OrderFilters) (*models.OrdersListResponse, error) {
	return r.listOrders(ctx, nil, &sellerID, filters)
}

func (r *OrderRepository) listOrders(ctx context.Context, buyerID, sellerID *uuid.UUID, filters models.OrderFilters) (*models.OrdersListResponse, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PerPage < 1 || filters.PerPage > 50 {
		filters.PerPage = 20
	}
	offset := (filters.Page - 1) * filters.PerPage

	// Build WHERE clause
	where := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if buyerID != nil {
		where += fmt.Sprintf(" AND o.buyer_id = $%d", argNum)
		args = append(args, *buyerID)
		argNum++
	}
	if sellerID != nil {
		where += fmt.Sprintf(" AND o.seller_id = $%d", argNum)
		args = append(args, *sellerID)
		argNum++
	}
	if filters.Status != nil {
		where += fmt.Sprintf(" AND o.status = $%d::order_status", argNum)
		args = append(args, string(*filters.Status))
		argNum++
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM orders o %s", where)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Get orders
	query := fmt.Sprintf(`
		SELECT
			o.id, o.buyer_id, o.seller_id, o.product_id,
			o.quantity, o.unit_price, COALESCE(o.shipping_cost, 0), o.total_amount,
			o.status::text, o.delivery_type::text,
			COALESCE(o.tracking_number, ''),
			o.created_at, o.updated_at,
			-- Buyer
			b.id, COALESCE(b.business_name, ''), b.first_name, b.last_name, COALESCE(b.avatar_url, ''),
			-- Seller
			s.id, COALESCE(s.business_name, ''), s.first_name, s.last_name, COALESCE(s.avatar_url, ''),
			-- Product
			p.id, p.title, COALESCE(p.images->>0, ''), p.price
		FROM orders o
		JOIN users b ON o.buyer_id = b.id
		JOIN users s ON o.seller_id = s.id
		JOIN products p ON o.product_id = p.id
		%s
		ORDER BY o.created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argNum, argNum+1)

	args = append(args, filters.PerPage, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		order := models.Order{
			Buyer:   &models.UserPublicMinimal{},
			Seller:  &models.UserPublicMinimal{},
			Product: &models.ProductMinimal{},
		}
		var status, deliveryType string

		err := rows.Scan(
			&order.ID, &order.BuyerID, &order.SellerID, &order.ProductID,
			&order.Quantity, &order.UnitPrice, &order.ShippingCost, &order.TotalAmount,
			&status, &deliveryType,
			&order.TrackingNumber,
			&order.CreatedAt, &order.UpdatedAt,
			&order.Buyer.ID, &order.Buyer.BusinessName, &order.Buyer.FirstName, &order.Buyer.LastName, &order.Buyer.AvatarURL,
			&order.Seller.ID, &order.Seller.BusinessName, &order.Seller.FirstName, &order.Seller.LastName, &order.Seller.AvatarURL,
			&order.Product.ID, &order.Product.Title, &order.Product.MainImageURL, &order.Product.Price,
		)
		if err != nil {
			return nil, err
		}

		order.Status = models.OrderStatus(status)
		order.DeliveryType = models.DeliveryType(deliveryType)
		orders = append(orders, order)
	}

	totalPages := (total + filters.PerPage - 1) / filters.PerPage

	return &models.OrdersListResponse{
		Orders:     orders,
		Total:      total,
		Page:       filters.Page,
		TotalPages: totalPages,
	}, nil
}

// UpdateStatus updates order status
func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.OrderStatus) error {
	query := `UPDATE orders SET status = $1::order_status, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, string(status), id)
	return err
}

// MarkAsPaid marks order as paid
func (r *OrderRepository) MarkAsPaid(ctx context.Context, id uuid.UUID, paymentIntentID string) error {
	query := `
		UPDATE orders SET
			status = 'PAID'::order_status,
			stripe_payment_intent_id = $1,
			paid_at = NOW(),
			updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.pool.Exec(ctx, query, paymentIntentID, id)
	return err
}

// UpdateTracking updates shipping tracking info
func (r *OrderRepository) UpdateTracking(ctx context.Context, id uuid.UUID, trackingNumber, trackingURL, carrier string) error {
	query := `
		UPDATE orders SET
			tracking_number = $1,
			tracking_url = $2,
			shipping_carrier = $3,
			status = 'SHIPPED'::order_status,
			shipped_at = NOW(),
			updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.pool.Exec(ctx, query, trackingNumber, trackingURL, carrier, id)
	return err
}

// ConfirmPickup confirms pickup via QR code
func (r *OrderRepository) ConfirmPickup(ctx context.Context, qrToken string) (*models.Order, error) {
	// Find and validate order
	query := `
		SELECT id FROM orders
		WHERE qr_code_token = $1
			AND status IN ('PAID', 'READY_FOR_PICKUP')
			AND (qr_code_expires_at IS NULL OR qr_code_expires_at > NOW())
	`

	var orderID uuid.UUID
	err := r.pool.QueryRow(ctx, query, qrToken).Scan(&orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("QR code non valido o scaduto")
		}
		return nil, err
	}

	// Update order
	updateQuery := `
		UPDATE orders SET
			status = 'DELIVERED'::order_status,
			qr_scanned_at = NOW(),
			payout_scheduled_at = NOW() + INTERVAL '48 hours',
			updated_at = NOW()
		WHERE id = $1
	`
	_, err = r.pool.Exec(ctx, updateQuery, orderID)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, orderID)
}

// CancelOrder cancels an order
func (r *OrderRepository) CancelOrder(ctx context.Context, id uuid.UUID, cancelledBy uuid.UUID, reason string) error {
	query := `
		UPDATE orders SET
			status = 'CANCELLED'::order_status,
			cancelled_at = NOW(),
			cancelled_by = $1,
			cancellation_reason = $2,
			updated_at = NOW()
		WHERE id = $3 AND status IN ('PENDING', 'PAID')
	`
	result, err := r.pool.Exec(ctx, query, cancelledBy, reason, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("ordine non cancellabile")
	}
	return nil
}

// =====================
// DISPUTES
// =====================

// CreateDispute opens a new dispute
func (r *OrderRepository) CreateDispute(ctx context.Context, dispute *models.Dispute) error {
	dispute.ID = uuid.New()
	dispute.Status = models.DisputeOpen
	dispute.CreatedAt = time.Now()
	dispute.UpdatedAt = time.Now()

	// Seller has 48h to respond
	deadline := time.Now().Add(48 * time.Hour)
	dispute.SellerResponseDeadline = &deadline

	query := `
		INSERT INTO disputes (
			id, order_id, opened_by, reason, description, evidence_urls,
			status, seller_response_deadline, created_at, updated_at
		) VALUES ($1, $2, $3, $4::dispute_reason, $5, $6, $7::dispute_status, $8, $9, $10)
	`

	_, err := r.pool.Exec(ctx, query,
		dispute.ID, dispute.OrderID, dispute.OpenedBy,
		string(dispute.Reason), dispute.Description, dispute.EvidenceURLs,
		string(dispute.Status), dispute.SellerResponseDeadline,
		dispute.CreatedAt, dispute.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// Update order status
	r.pool.Exec(ctx, "UPDATE orders SET status = 'DISPUTED'::order_status WHERE id = $1", dispute.OrderID)

	return nil
}

// GetDisputeByOrderID gets dispute for an order
func (r *OrderRepository) GetDisputeByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Dispute, error) {
	query := `
		SELECT id, order_id, opened_by, reason::text, description, evidence_urls,
			status::text, COALESCE(seller_response, ''), seller_response_at, seller_evidence_urls,
			resolved_by, COALESCE(resolution_notes, ''), COALESCE(refund_amount, 0), COALESCE(seller_payout_amount, 0),
			created_at, updated_at, resolved_at, seller_response_deadline
		FROM disputes
		WHERE order_id = $1
	`

	dispute := &models.Dispute{}
	var reason, status string

	err := r.pool.QueryRow(ctx, query, orderID).Scan(
		&dispute.ID, &dispute.OrderID, &dispute.OpenedBy, &reason, &dispute.Description, &dispute.EvidenceURLs,
		&status, &dispute.SellerResponse, &dispute.SellerResponseAt, &dispute.SellerEvidenceURLs,
		&dispute.ResolvedBy, &dispute.ResolutionNotes, &dispute.RefundAmount, &dispute.SellerPayoutAmount,
		&dispute.CreatedAt, &dispute.UpdatedAt, &dispute.ResolvedAt, &dispute.SellerResponseDeadline,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDisputeNotFound
		}
		return nil, err
	}

	dispute.Reason = models.DisputeReason(reason)
	dispute.Status = models.DisputeStatus(status)

	return dispute, nil
}

// =====================
// REVIEWS
// =====================

// CreateReview creates a new review
func (r *OrderRepository) CreateReview(ctx context.Context, review *models.OrderReview) error {
	review.ID = uuid.New()
	review.CreatedAt = time.Now()

	query := `
		INSERT INTO order_reviews (id, order_id, reviewer_id, reviewed_id, rating, comment, is_anonymous, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.pool.Exec(ctx, query,
		review.ID, review.OrderID, review.ReviewerID, review.ReviewedID,
		review.Rating, review.Comment, review.IsAnonymous, review.CreatedAt,
	)

	return err
}

// GetReviewsForUser gets reviews for a user
func (r *OrderRepository) GetReviewsForUser(ctx context.Context, userID uuid.UUID, limit int) ([]models.OrderReview, error) {
	query := `
		SELECT r.id, r.order_id, r.reviewer_id, r.reviewed_id, r.rating, COALESCE(r.comment, ''), r.is_anonymous, r.created_at,
			u.id, COALESCE(u.business_name, ''), u.first_name, u.last_name, COALESCE(u.avatar_url, '')
		FROM order_reviews r
		LEFT JOIN users u ON r.reviewer_id = u.id
		WHERE r.reviewed_id = $1 AND r.is_approved = TRUE
		ORDER BY r.created_at DESC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.OrderReview
	for rows.Next() {
		review := models.OrderReview{Reviewer: &models.UserPublicMinimal{}}
		err := rows.Scan(
			&review.ID, &review.OrderID, &review.ReviewerID, &review.ReviewedID,
			&review.Rating, &review.Comment, &review.IsAnonymous, &review.CreatedAt,
			&review.Reviewer.ID, &review.Reviewer.BusinessName, &review.Reviewer.FirstName, &review.Reviewer.LastName, &review.Reviewer.AvatarURL,
		)
		if err != nil {
			return nil, err
		}

		// Hide reviewer info if anonymous
		if review.IsAnonymous {
			review.Reviewer = nil
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}

// =====================
// STRIKES
// =====================

// AddStrike adds a strike to a user
func (r *OrderRepository) AddStrike(ctx context.Context, strike *models.UserStrike) error {
	strike.ID = uuid.New()
	strike.IsActive = true
	strike.CreatedAt = time.Now()

	query := `
		INSERT INTO user_strikes (id, user_id, order_id, strike_type, description, expires_at, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.pool.Exec(ctx, query,
		strike.ID, strike.UserID, strike.OrderID, strike.StrikeType,
		strike.Description, strike.ExpiresAt, strike.IsActive, strike.CreatedAt,
	)

	return err
}

// GetUserStrikes gets active strikes for a user
func (r *OrderRepository) GetUserStrikes(ctx context.Context, userID uuid.UUID) ([]models.UserStrike, error) {
	query := `
		SELECT id, user_id, order_id, strike_type, COALESCE(description, ''), expires_at, is_active, created_at
		FROM user_strikes
		WHERE user_id = $1 AND is_active = TRUE
			AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var strikes []models.UserStrike
	for rows.Next() {
		var strike models.UserStrike
		err := rows.Scan(
			&strike.ID, &strike.UserID, &strike.OrderID, &strike.StrikeType,
			&strike.Description, &strike.ExpiresAt, &strike.IsActive, &strike.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		strikes = append(strikes, strike)
	}

	return strikes, nil
}

// Helper function
func round(f float64) float64 {
	return float64(int(f*100+0.5)) / 100
}
