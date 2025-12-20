package models

import (
	"time"

	"github.com/google/uuid"
)

// ProductStatus represents the status of a product (matches DB enum)
type ProductStatus string

const (
	ProductStatusDraft   ProductStatus = "DRAFT"
	ProductStatusActive  ProductStatus = "ACTIVE"
	ProductStatusSold    ProductStatus = "SOLD"
	ProductStatusExpired ProductStatus = "EXPIRED"
	ProductStatusDeleted ProductStatus = "DELETED"
)

// ListingType represents the type of listing (matches DB enum)
type ListingType string

const (
	ListingSale ListingType = "SALE"
	ListingGift ListingType = "GIFT"
)

// ShippingMethod represents shipping method (matches DB enum)
type ShippingMethod string

const (
	ShippingPickup            ShippingMethod = "PICKUP"
	ShippingSellerShips       ShippingMethod = "SELLER_SHIPS"
	ShippingBuyerArranges     ShippingMethod = "BUYER_ARRANGES"
	ShippingPlatformManaged   ShippingMethod = "PLATFORM_MANAGED"
	ShippingDigitalForwarders ShippingMethod = "DIGITAL_FORWARDERS" // Coming Soon
	ShippingBoth              ShippingMethod = "BOTH"               // Both pickup and shipping available
)

// QuantityUnit represents the unit of measurement for quantity
type QuantityUnit string

const (
	QuantityUnitPiece  QuantityUnit = "PIECE"
	QuantityUnitKG     QuantityUnit = "KG"
	QuantityUnitG      QuantityUnit = "G"
	QuantityUnitL      QuantityUnit = "L"
	QuantityUnitML     QuantityUnit = "ML"
	QuantityUnitCustom QuantityUnit = "CUSTOM"
)

// Product represents a product listing
type Product struct {
	ID              uuid.UUID      `json:"id"`
	SellerID        uuid.UUID      `json:"seller_id"`
	CategoryID      *uuid.UUID     `json:"category_id,omitempty"`
	LocationID      *uuid.UUID     `json:"location_id,omitempty"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	Price           float64        `json:"price"`
	OriginalPrice   *float64       `json:"original_price,omitempty"`
	ListingType     ListingType    `json:"listing_type"`
	ShippingMethod      ShippingMethod `json:"shipping_method"`
	ShippingCost        float64        `json:"shipping_cost"`
	PickupLocationIDs   []uuid.UUID    `json:"pickup_location_ids,omitempty"`
	Quantity            int            `json:"quantity"`
	QuantityAvail      int            `json:"quantity_available"`
	QuantityUnit       QuantityUnit   `json:"quantity_unit"`
	QuantityUnitCustom *string        `json:"quantity_unit_custom,omitempty"`
	ExpiryDate         *time.Time     `json:"expiry_date,omitempty"`
	ExpiryPhotoURL     *string        `json:"expiry_photo_url,omitempty"`

	// Dutch Auction
	IsDutchAuction      bool       `json:"is_dutch_auction"`
	DutchStartPrice     *float64   `json:"dutch_start_price,omitempty"`
	DutchDecreaseAmount *float64   `json:"dutch_decrease_amount,omitempty"`
	DutchDecreaseHours  *int       `json:"dutch_decrease_hours,omitempty"`
	DutchMinPrice       *float64   `json:"dutch_min_price,omitempty"`
	DutchStartedAt      *time.Time `json:"dutch_started_at,omitempty"`

	// Location
	City       string   `json:"city,omitempty"`
	Province   string   `json:"province,omitempty"`
	PostalCode string   `json:"postal_code,omitempty"`
	Latitude   *float64 `json:"latitude,omitempty"`
	Longitude  *float64 `json:"longitude,omitempty"`

	// Images (JSONB array in DB)
	Images []string `json:"images"`

	// Status
	Status ProductStatus `json:"status"`
	Slug   string        `json:"slug,omitempty"`

	// Stats
	ViewCount     int `json:"view_count"`
	FavoriteCount int `json:"favorite_count"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations (populated when needed)
	Seller   *UserProfile `json:"seller,omitempty"`
	Category *Category    `json:"category,omitempty"`
}

// Category represents a product category
type Category struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	Icon      string     `json:"icon,omitempty"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	SortOrder int        `json:"sort_order"`
}

// CreateProductRequest represents a request to create a product
type CreateProductRequest struct {
	Title               string         `json:"title"`
	Description         string         `json:"description"`
	CategoryID          *uuid.UUID     `json:"category_id,omitempty"`
	Price               float64        `json:"price"`
	OriginalPrice       *float64       `json:"original_price,omitempty"`
	Quantity            int            `json:"quantity"`
	QuantityUnit        QuantityUnit   `json:"quantity_unit"`
	QuantityUnitCustom  *string        `json:"quantity_unit_custom,omitempty"`
	ListingType         ListingType    `json:"listing_type"`
	ShippingMethod      ShippingMethod `json:"shipping_method"`
	ShippingCost        float64        `json:"shipping_cost"`
	PickupLocationIDs   []uuid.UUID    `json:"pickup_location_ids,omitempty"`
	ExpiryDate          *time.Time     `json:"expiry_date,omitempty"`
	IsDutchAuction      bool           `json:"is_dutch_auction"`
	DutchStartPrice     *float64       `json:"dutch_start_price,omitempty"`
	DutchDecreaseAmount *float64       `json:"dutch_decrease_amount,omitempty"`
	DutchDecreaseHours  *int           `json:"dutch_decrease_hours,omitempty"`
	DutchMinPrice       *float64       `json:"dutch_min_price,omitempty"`
	City                string         `json:"city"`
	Province            string         `json:"province,omitempty"`
	PostalCode          string         `json:"postal_code,omitempty"`
	Latitude            *float64       `json:"latitude,omitempty"`
	Longitude           *float64       `json:"longitude,omitempty"`
}

// UpdateProductRequest represents a request to update a product
type UpdateProductRequest struct {
	Title              *string         `json:"title,omitempty"`
	Description        *string         `json:"description,omitempty"`
	Price              *float64        `json:"price,omitempty"`
	Quantity           *int            `json:"quantity,omitempty"`
	QuantityUnit       *QuantityUnit   `json:"quantity_unit,omitempty"`
	QuantityUnitCustom *string         `json:"quantity_unit_custom,omitempty"`
	Status             *ProductStatus  `json:"status,omitempty"`
	ShippingMethod     *ShippingMethod `json:"shipping_method,omitempty"`
	ShippingCost       *float64        `json:"shipping_cost,omitempty"`
}

// ProductListResponse represents a paginated list of products
type ProductListResponse struct {
	Products   []Product `json:"products"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	PerPage    int       `json:"per_page"`
	TotalPages int       `json:"total_pages"`
}

// ProductFilters represents filters for product listing
type ProductFilters struct {
	CategoryID  *uuid.UUID
	SellerID    *uuid.UUID
	ListingType *ListingType
	MinPrice    *float64
	MaxPrice    *float64
	City        *string
	Status      *ProductStatus
	Search      *string
	Page        int
	PerPage     int
	SortBy      string
	SortOrder   string
}

// GetCurrentPrice calculates the current price considering Dutch Auction
func (p *Product) GetCurrentPrice() float64 {
	if !p.IsDutchAuction || p.DutchStartedAt == nil || p.DutchStartPrice == nil {
		return p.Price
	}

	hoursPassed := time.Since(*p.DutchStartedAt).Hours()
	decreaseHours := 24
	if p.DutchDecreaseHours != nil {
		decreaseHours = *p.DutchDecreaseHours
	}

	intervals := int(hoursPassed) / decreaseHours

	decreaseAmount := 1.0
	if p.DutchDecreaseAmount != nil {
		decreaseAmount = *p.DutchDecreaseAmount
	}

	currentPrice := *p.DutchStartPrice - (float64(intervals) * decreaseAmount)

	minPrice := 0.0
	if p.DutchMinPrice != nil {
		minPrice = *p.DutchMinPrice
	}

	if currentPrice < minPrice {
		return minPrice
	}
	return currentPrice
}
