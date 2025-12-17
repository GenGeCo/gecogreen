package models

import (
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	AccountPrivate  AccountType = "PRIVATE"
	AccountBusiness AccountType = "BUSINESS"
)

type UserStatus string

const (
	UserStatusPending   UserStatus = "PENDING"
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
	UserStatusBanned    UserStatus = "BANNED"
)

// SocialLinks holds user's social media links
type SocialLinks struct {
	Instagram string `json:"instagram,omitempty"`
	Facebook  string `json:"facebook,omitempty"`
	Website   string `json:"website,omitempty"`
	LinkedIn  string `json:"linkedin,omitempty"`
}

// User represents a user in the system
type User struct {
	ID           uuid.UUID   `json:"id"`
	Email        string      `json:"email"`
	PasswordHash string      `json:"-"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	Phone        string      `json:"phone,omitempty"`
	City         string      `json:"city,omitempty"`
	Province     string      `json:"province,omitempty"`
	PostalCode   string      `json:"postal_code,omitempty"`

	// Account type
	AccountType          AccountType `json:"account_type"`
	BusinessName         string      `json:"business_name,omitempty"`
	VATNumber            string      `json:"vat_number,omitempty"`
	HasMultipleLocations bool        `json:"has_multiple_locations"`

	// Billing info
	FiscalCode        string `json:"fiscal_code,omitempty"`         // Codice Fiscale (Italia)
	SDICode           string `json:"sdi_code,omitempty"`            // Codice Univoco SDI
	PECEmail          string `json:"pec_email,omitempty"`           // PEC per fatturazione elettronica
	EUVatID           string `json:"eu_vat_id,omitempty"`           // VAT ID europeo
	BillingAddress    string `json:"billing_address,omitempty"`     // Indirizzo fatturazione
	BillingCity       string `json:"billing_city,omitempty"`        // Città fatturazione
	BillingProvince   string `json:"billing_province,omitempty"`    // Provincia fatturazione
	BillingPostalCode string `json:"billing_postal_code,omitempty"` // CAP fatturazione
	BillingCountry    string `json:"billing_country,omitempty"`     // Paese fatturazione (ISO 3166-1)

	// Profile
	AvatarURL      string      `json:"avatar_url,omitempty"`
	SocialLinks    SocialLinks `json:"social_links,omitempty"`
	BusinessPhotos []string    `json:"business_photos,omitempty"`

	// Status & Admin
	Status        UserStatus `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	IsAdmin       bool       `json:"is_admin"`

	// Stripe
	StripeCustomerID *string `json:"stripe_customer_id,omitempty"`
	StripeAccountID  *string `json:"stripe_account_id,omitempty"`

	// Eco stats
	TotalCO2Saved   float64 `json:"total_co2_saved"`
	TotalWaterSaved float64 `json:"total_water_saved"`
	EcoCredits      int     `json:"eco_credits"`
	EcoLevel        string  `json:"eco_level"`

	// Rating
	RatingAvg   float64 `json:"rating_avg"`
	RatingCount int     `json:"rating_count"`

	// Timestamps
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (u *User) IsBusiness() bool { return u.AccountType == AccountBusiness }
func (u *User) IsActive() bool   { return u.Status == UserStatusActive }

// UserProfile is the public-facing profile (shown after purchase)
type UserProfile struct {
	ID           uuid.UUID   `json:"id"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	AvatarURL    string      `json:"avatar_url,omitempty"`
	AccountType  AccountType `json:"account_type"`
	BusinessName string      `json:"business_name,omitempty"`
	City         string      `json:"city,omitempty"`
	RatingAvg    float64     `json:"rating_avg"`
	RatingCount  int         `json:"rating_count"`
	CreatedAt    time.Time   `json:"created_at"`
}

// UserPublicMinimal is shown in product listings and leaderboards (minimal public info)
type UserPublicMinimal struct {
	ID           uuid.UUID   `json:"id"`
	AccountType  AccountType `json:"account_type"`
	BusinessName string      `json:"business_name,omitempty"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	AvatarURL    string      `json:"avatar_url,omitempty"`
	City         string      `json:"city,omitempty"`
	RatingAvg    float64     `json:"rating_avg"`
}

// UserFullProfile is shown after purchase (includes contact info)
type UserFullProfile struct {
	UserProfile
	Phone          string      `json:"phone,omitempty"`
	SocialLinks    SocialLinks `json:"social_links,omitempty"`
	BusinessPhotos []string    `json:"business_photos,omitempty"`
	VATNumber      string      `json:"vat_number,omitempty"`
}

func (u *User) ToProfile() UserProfile {
	return UserProfile{
		ID:           u.ID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		AvatarURL:    u.AvatarURL,
		AccountType:  u.AccountType,
		BusinessName: u.BusinessName,
		City:         u.City,
		RatingAvg:    u.RatingAvg,
		RatingCount:  u.RatingCount,
		CreatedAt:    u.CreatedAt,
	}
}

func (u *User) ToMinimal() UserPublicMinimal {
	return UserPublicMinimal{
		ID:           u.ID,
		AccountType:  u.AccountType,
		BusinessName: u.BusinessName,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		AvatarURL:    u.AvatarURL,
		City:         u.City,
		RatingAvg:    u.RatingAvg,
	}
}

func (u *User) ToFullProfile() UserFullProfile {
	return UserFullProfile{
		UserProfile:    u.ToProfile(),
		Phone:          u.Phone,
		SocialLinks:    u.SocialLinks,
		BusinessPhotos: u.BusinessPhotos,
		VATNumber:      u.VATNumber,
	}
}

// Location represents a pickup location
type Location struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	Name            string     `json:"name"`
	IsPrimary       bool       `json:"is_primary"`
	IsActive        bool       `json:"is_active"`
	AddressStreet   string     `json:"address_street"`
	AddressCity     string     `json:"address_city"`
	AddressProvince string     `json:"address_province,omitempty"`
	AddressPostal   string     `json:"address_postal_code"`
	Phone           string     `json:"phone,omitempty"`
	Email           string     `json:"email,omitempty"`
	PickupHours     string     `json:"pickup_hours,omitempty"` // JSON string
	PickupInstructions string  `json:"pickup_instructions,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// RegisterRequest for new user registration
type RegisterRequest struct {
	Email       string      `json:"email" validate:"required,email"`
	Password    string      `json:"password" validate:"required,min=8"`
	FirstName   string      `json:"first_name" validate:"required"`
	LastName    string      `json:"last_name" validate:"required"`
	AccountType AccountType `json:"account_type" validate:"required,oneof=PRIVATE BUSINESS"`

	// Business fields (required if AccountType is BUSINESS)
	BusinessName         string `json:"business_name,omitempty"`
	VATNumber            string `json:"vat_number,omitempty"`
	HasMultipleLocations bool   `json:"has_multiple_locations"`

	// Billing info (for BUSINESS accounts)
	FiscalCode     string `json:"fiscal_code,omitempty"`      // Codice Fiscale
	SDICode        string `json:"sdi_code,omitempty"`         // Codice Univoco SDI (7 chars)
	PECEmail       string `json:"pec_email,omitempty"`        // PEC for e-invoicing
	EUVatID        string `json:"eu_vat_id,omitempty"`        // EU VAT ID for non-IT businesses
	BillingCountry string `json:"billing_country,omitempty"`  // ISO 3166-1 alpha-2 (IT, DE, FR, etc.)

	// Primary location (required for all)
	City          string `json:"city" validate:"required"`
	Province      string `json:"province,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
	AddressStreet string `json:"address_street,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// UpdateProfileRequest for updating user profile
type UpdateProfileRequest struct {
	FirstName            *string      `json:"first_name,omitempty"`
	LastName             *string      `json:"last_name,omitempty"`
	Phone                *string      `json:"phone,omitempty"`
	City                 *string      `json:"city,omitempty"`
	Province             *string      `json:"province,omitempty"`
	PostalCode           *string      `json:"postal_code,omitempty"`
	AccountType          *AccountType `json:"account_type,omitempty"` // Allow switching between PRIVATE/BUSINESS
	BusinessName         *string      `json:"business_name,omitempty"`
	VATNumber            *string      `json:"vat_number,omitempty"`
	SocialLinks          *SocialLinks `json:"social_links,omitempty"`
	HasMultipleLocations *bool        `json:"has_multiple_locations,omitempty"`

	// Billing info
	FiscalCode        *string `json:"fiscal_code,omitempty"`
	SDICode           *string `json:"sdi_code,omitempty"`
	PECEmail          *string `json:"pec_email,omitempty"`
	EUVatID           *string `json:"eu_vat_id,omitempty"`
	BillingAddress    *string `json:"billing_address,omitempty"`
	BillingCity       *string `json:"billing_city,omitempty"`
	BillingProvince   *string `json:"billing_province,omitempty"`
	BillingPostalCode *string `json:"billing_postal_code,omitempty"`
	BillingCountry    *string `json:"billing_country,omitempty"`
}

// CreateLocationRequest for adding a new location
type CreateLocationRequest struct {
	Name               string `json:"name" validate:"required"`
	AddressStreet      string `json:"address_street" validate:"required"`
	AddressCity        string `json:"address_city" validate:"required"`
	AddressProvince    string `json:"address_province,omitempty"`
	AddressPostal      string `json:"address_postal_code" validate:"required"`
	Phone              string `json:"phone,omitempty"`
	Email              string `json:"email,omitempty"`
	PickupHours        string `json:"pickup_hours,omitempty"`
	PickupInstructions string `json:"pickup_instructions,omitempty"`
	IsPrimary          bool   `json:"is_primary"`
}

// ImageReview for moderation queue
type ImageReview struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	ProductID       *uuid.UUID `json:"product_id,omitempty"`
	ImageURL        string     `json:"image_url"`
	ImageType       string     `json:"image_type"` // PRODUCT, PROFILE, BUSINESS

	// OCR results
	DetectedText  string  `json:"detected_text,omitempty"`
	DetectedPhone bool    `json:"detected_phone"`
	DetectedEmail bool    `json:"detected_email"`
	DetectedURL   bool    `json:"detected_url"`
	Confidence    float64 `json:"confidence_score"`

	// Moderation
	Status          string     `json:"status"` // PENDING, APPROVED, REJECTED
	ReviewedBy      *uuid.UUID `json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	RejectionReason string     `json:"rejection_reason,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}
