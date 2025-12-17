package models

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderPending        OrderStatus = "PENDING"
	OrderPaid           OrderStatus = "PAID"
	OrderProcessing     OrderStatus = "PROCESSING"
	OrderShipped        OrderStatus = "SHIPPED"
	OrderReadyForPickup OrderStatus = "READY_FOR_PICKUP"
	OrderInTransit      OrderStatus = "IN_TRANSIT"
	OrderDelivered      OrderStatus = "DELIVERED"
	OrderCompleted      OrderStatus = "COMPLETED"
	OrderCancelled      OrderStatus = "CANCELLED"
	OrderRefunded       OrderStatus = "REFUNDED"
	OrderDisputed       OrderStatus = "DISPUTED"
)

// DeliveryType represents how the order will be delivered
type DeliveryType string

const (
	DeliveryPickup       DeliveryType = "PICKUP"
	DeliverySellerShips  DeliveryType = "SELLER_SHIPS"
	DeliveryBuyerArranges DeliveryType = "BUYER_ARRANGES"
)

// DisputeReason represents why a dispute was opened
type DisputeReason string

const (
	DisputeItemNotReceived    DisputeReason = "ITEM_NOT_RECEIVED"
	DisputeItemDamaged        DisputeReason = "ITEM_DAMAGED"
	DisputeItemNotAsDescribed DisputeReason = "ITEM_NOT_AS_DESCRIBED"
	DisputeSellerNoShow       DisputeReason = "SELLER_NO_SHOW"
	DisputeBuyerNoShow        DisputeReason = "BUYER_NO_SHOW"
	DisputeScamAttempt        DisputeReason = "SCAM_ATTEMPT"
	DisputeOther              DisputeReason = "OTHER"
)

// DisputeStatus represents the status of a dispute
type DisputeStatus string

const (
	DisputeOpen                  DisputeStatus = "OPEN"
	DisputeSellerResponse        DisputeStatus = "SELLER_RESPONSE"
	DisputeBuyerReview           DisputeStatus = "BUYER_REVIEW"
	DisputeAdminReview           DisputeStatus = "ADMIN_REVIEW"
	DisputeResolvedRefundFull    DisputeStatus = "RESOLVED_REFUND_FULL"
	DisputeResolvedRefundPartial DisputeStatus = "RESOLVED_REFUND_PARTIAL"
	DisputeResolvedPayoutSeller  DisputeStatus = "RESOLVED_PAYOUT_SELLER"
	DisputeResolvedSplit         DisputeStatus = "RESOLVED_SPLIT"
	DisputeClosed                DisputeStatus = "CLOSED"
)

// Order represents a complete order
type Order struct {
	ID        uuid.UUID `json:"id"`
	BuyerID   uuid.UUID `json:"buyer_id"`
	SellerID  uuid.UUID `json:"seller_id"`
	ProductID uuid.UUID `json:"product_id"`

	// Order details
	Quantity     int     `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	ShippingCost float64 `json:"shipping_cost"`
	TotalAmount  float64 `json:"total_amount"`

	// Fees
	PlatformFee  float64 `json:"platform_fee"`
	StripeFee    float64 `json:"stripe_fee"`
	SellerPayout float64 `json:"seller_payout"`

	// Status
	Status       OrderStatus  `json:"status"`
	DeliveryType DeliveryType `json:"delivery_type"`

	// Pickup details
	PickupLocationID   *uuid.UUID `json:"pickup_location_id,omitempty"`
	PickupAddress      string     `json:"pickup_address,omitempty"`
	PickupInstructions string     `json:"pickup_instructions,omitempty"`
	PickupDeadline     *time.Time `json:"pickup_deadline,omitempty"`

	// Shipping details
	ShippingAddress    string `json:"shipping_address,omitempty"`
	ShippingCity       string `json:"shipping_city,omitempty"`
	ShippingProvince   string `json:"shipping_province,omitempty"`
	ShippingPostalCode string `json:"shipping_postal_code,omitempty"`
	ShippingCountry    string `json:"shipping_country,omitempty"`
	TrackingNumber     string `json:"tracking_number,omitempty"`
	TrackingURL        string `json:"tracking_url,omitempty"`
	ShippingCarrier    string `json:"shipping_carrier,omitempty"`
	ShippedAt          *time.Time `json:"shipped_at,omitempty"`

	// QR Code
	QRCodeToken     string     `json:"qr_code_token,omitempty"`
	QRCodeExpiresAt *time.Time `json:"qr_code_expires_at,omitempty"`
	QRScannedAt     *time.Time `json:"qr_scanned_at,omitempty"`

	// Payment
	StripePaymentIntentID    string     `json:"stripe_payment_intent_id,omitempty"`
	StripeCheckoutSessionID  string     `json:"stripe_checkout_session_id,omitempty"`
	StripeTransferID         string     `json:"stripe_transfer_id,omitempty"`
	PaidAt                   *time.Time `json:"paid_at,omitempty"`

	// Payout
	PayoutScheduledAt *time.Time `json:"payout_scheduled_at,omitempty"`
	PayoutCompletedAt *time.Time `json:"payout_completed_at,omitempty"`
	PayoutHoldReason  string     `json:"payout_hold_reason,omitempty"`

	// Impact
	CO2Saved         float64 `json:"co2_saved"`
	WaterSaved       float64 `json:"water_saved"`
	EcoCreditsBuyer  int     `json:"eco_credits_buyer"`
	EcoCreditsSeller int     `json:"eco_credits_seller"`

	// Notes
	BuyerNotes    string `json:"buyer_notes,omitempty"`
	SellerNotes   string `json:"seller_notes,omitempty"`
	InternalNotes string `json:"internal_notes,omitempty"`

	// Timestamps
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
	CancelledAt        *time.Time `json:"cancelled_at,omitempty"`
	CancelledBy        *uuid.UUID `json:"cancelled_by,omitempty"`
	CancellationReason string     `json:"cancellation_reason,omitempty"`

	// Joined data
	Buyer   *UserPublicMinimal `json:"buyer,omitempty"`
	Seller  *UserPublicMinimal `json:"seller,omitempty"`
	Product *ProductMinimal    `json:"product,omitempty"`
}

// ProductMinimal for order display
type ProductMinimal struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	MainImageURL string    `json:"main_image_url"`
	Price        float64   `json:"price"`
}

// Dispute represents a dispute on an order
type Dispute struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	OpenedBy  uuid.UUID `json:"opened_by"`

	Reason      DisputeReason `json:"reason"`
	Description string        `json:"description"`
	EvidenceURLs []string     `json:"evidence_urls,omitempty"`

	Status DisputeStatus `json:"status"`

	SellerResponse    string     `json:"seller_response,omitempty"`
	SellerResponseAt  *time.Time `json:"seller_response_at,omitempty"`
	SellerEvidenceURLs []string  `json:"seller_evidence_urls,omitempty"`

	ResolvedBy         *uuid.UUID `json:"resolved_by,omitempty"`
	ResolutionNotes    string     `json:"resolution_notes,omitempty"`
	RefundAmount       float64    `json:"refund_amount,omitempty"`
	SellerPayoutAmount float64    `json:"seller_payout_amount,omitempty"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`

	SellerResponseDeadline *time.Time `json:"seller_response_deadline,omitempty"`
	AdminReviewDeadline    *time.Time `json:"admin_review_deadline,omitempty"`
}

// OrderReview represents a review for an order
type OrderReview struct {
	ID         uuid.UUID `json:"id"`
	OrderID    uuid.UUID `json:"order_id"`
	ReviewerID uuid.UUID `json:"reviewer_id"`
	ReviewedID uuid.UUID `json:"reviewed_id"`

	Rating      int    `json:"rating"`
	Comment     string `json:"comment,omitempty"`
	IsAnonymous bool   `json:"is_anonymous"`

	CreatedAt time.Time `json:"created_at"`

	// Joined
	Reviewer *UserPublicMinimal `json:"reviewer,omitempty"`
}

// UserStrike represents a strike against a user
type UserStrike struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	OrderID     *uuid.UUID `json:"order_id,omitempty"`
	StrikeType  string     `json:"strike_type"`
	Description string     `json:"description,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
}

// CartItem represents an item in the cart
type CartItem struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`

	Product *ProductMinimal `json:"product,omitempty"`
}

// --- Request/Response Types ---

// CreateOrderRequest is the request to create a new order
type CreateOrderRequest struct {
	ProductID    uuid.UUID    `json:"product_id" validate:"required"`
	Quantity     int          `json:"quantity" validate:"required,min=1"`
	DeliveryType DeliveryType `json:"delivery_type" validate:"required"`

	// For shipping
	ShippingAddress    string `json:"shipping_address,omitempty"`
	ShippingCity       string `json:"shipping_city,omitempty"`
	ShippingProvince   string `json:"shipping_province,omitempty"`
	ShippingPostalCode string `json:"shipping_postal_code,omitempty"`
	ShippingCountry    string `json:"shipping_country,omitempty"`

	// For pickup
	PickupLocationID *uuid.UUID `json:"pickup_location_id,omitempty"`

	BuyerNotes string `json:"buyer_notes,omitempty"`
}

// CheckoutResponse is returned after creating an order
type CheckoutResponse struct {
	OrderID            uuid.UUID `json:"order_id"`
	StripeCheckoutURL  string    `json:"stripe_checkout_url"`
	TotalAmount        float64   `json:"total_amount"`
	ExpiresAt          time.Time `json:"expires_at"`
}

// UpdateOrderStatusRequest for seller/admin to update status
type UpdateOrderStatusRequest struct {
	Status         OrderStatus `json:"status,omitempty"`
	TrackingNumber string      `json:"tracking_number,omitempty"`
	TrackingURL    string      `json:"tracking_url,omitempty"`
	ShippingCarrier string     `json:"shipping_carrier,omitempty"`
	SellerNotes    string      `json:"seller_notes,omitempty"`
}

// ConfirmPickupRequest for QR code scanning
type ConfirmPickupRequest struct {
	QRCodeToken string `json:"qr_code_token" validate:"required"`
}

// OpenDisputeRequest for opening a dispute
type OpenDisputeRequest struct {
	Reason      DisputeReason `json:"reason" validate:"required"`
	Description string        `json:"description" validate:"required,min=50"`
	EvidenceURLs []string     `json:"evidence_urls,omitempty"`
}

// RespondDisputeRequest for seller response
type RespondDisputeRequest struct {
	Response     string   `json:"response" validate:"required,min=50"`
	EvidenceURLs []string `json:"evidence_urls,omitempty"`
}

// CreateReviewRequest for leaving a review
type CreateReviewRequest struct {
	Rating      int    `json:"rating" validate:"required,min=1,max=5"`
	Comment     string `json:"comment,omitempty"`
	IsAnonymous bool   `json:"is_anonymous"`
}

// OrdersListResponse for paginated list
type OrdersListResponse struct {
	Orders      []Order `json:"orders"`
	Total       int     `json:"total"`
	Page        int     `json:"page"`
	TotalPages  int     `json:"total_pages"`
}

// OrderFilters for listing orders
type OrderFilters struct {
	Status       *OrderStatus  `json:"status,omitempty"`
	DeliveryType *DeliveryType `json:"delivery_type,omitempty"`
	DateFrom     *time.Time    `json:"date_from,omitempty"`
	DateTo       *time.Time    `json:"date_to,omitempty"`
	Page         int           `json:"page"`
	PerPage      int           `json:"per_page"`
}
