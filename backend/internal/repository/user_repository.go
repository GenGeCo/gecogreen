package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gecogreen/backend/internal/models"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrLocationNotFound   = errors.New("location not found")
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (
			id, email, password_hash, first_name, last_name, phone,
			city, province, postal_code, account_type, business_name, vat_number,
			has_multiple_locations, fiscal_code, sdi_code, pec_email, eu_vat_id, billing_country,
			status, email_verified, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10::account_type, $11, $12, $13, $14, $15, $16, $17, $18, $19::user_status, $20, $21, $22
		)
	`

	user.ID = uuid.New()
	user.Status = models.UserStatusActive // Auto-activate for now
	user.EmailVerified = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if user.AccountType == "" {
		user.AccountType = models.AccountPrivate
	}

	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.Phone,
		user.City, user.Province, user.PostalCode, string(user.AccountType), user.BusinessName, user.VATNumber,
		user.HasMultipleLocations, user.FiscalCode, user.SDICode, user.PECEmail, user.EUVatID, user.BillingCountry,
		string(user.Status), user.EmailVerified, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil && strings.Contains(err.Error(), "users_email_key") {
		return ErrEmailAlreadyExists
	}
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name,
		       COALESCE(phone, ''), COALESCE(city, ''), COALESCE(province, ''), COALESCE(postal_code, ''),
		       COALESCE(account_type::text, 'PRIVATE'), COALESCE(business_name, ''), COALESCE(vat_number, ''),
		       COALESCE(has_multiple_locations, false),
		       COALESCE(fiscal_code, ''), COALESCE(sdi_code, ''), COALESCE(pec_email, ''),
		       COALESCE(eu_vat_id, ''), COALESCE(billing_address, ''), COALESCE(billing_city, ''),
		       COALESCE(billing_province, ''), COALESCE(billing_postal_code, ''), COALESCE(billing_country, 'IT'),
		       status::text, email_verified, COALESCE(is_admin, false),
		       COALESCE(avatar_url, ''), COALESCE(social_links, '{}'), COALESCE(business_photos, '[]'),
		       stripe_customer_id, stripe_account_id,
		       COALESCE(total_co2_saved, 0), COALESCE(total_water_saved, 0), COALESCE(eco_credits, 0), COALESCE(eco_level, 'Germoglio'),
		       COALESCE(rating_avg, 0), COALESCE(rating_count, 0),
		       last_login_at, created_at, updated_at
		FROM users WHERE id = $1 AND deleted_at IS NULL
	`

	user := &models.User{}
	var status string
	var accountType string
	var socialLinksJSON []byte
	var businessPhotosJSON []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Phone, &user.City, &user.Province, &user.PostalCode,
		&accountType, &user.BusinessName, &user.VATNumber,
		&user.HasMultipleLocations,
		&user.FiscalCode, &user.SDICode, &user.PECEmail,
		&user.EUVatID, &user.BillingAddress, &user.BillingCity,
		&user.BillingProvince, &user.BillingPostalCode, &user.BillingCountry,
		&status, &user.EmailVerified, &user.IsAdmin,
		&user.AvatarURL, &socialLinksJSON, &businessPhotosJSON,
		&user.StripeCustomerID, &user.StripeAccountID,
		&user.TotalCO2Saved, &user.TotalWaterSaved, &user.EcoCredits, &user.EcoLevel,
		&user.RatingAvg, &user.RatingCount,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user.Status = models.UserStatus(status)
	user.AccountType = models.AccountType(accountType)

	// Parse JSON fields
	json.Unmarshal(socialLinksJSON, &user.SocialLinks)
	json.Unmarshal(businessPhotosJSON, &user.BusinessPhotos)

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name,
		       COALESCE(phone, ''), COALESCE(city, ''), COALESCE(province, ''), COALESCE(postal_code, ''),
		       COALESCE(account_type::text, 'PRIVATE'), COALESCE(business_name, ''), COALESCE(vat_number, ''),
		       COALESCE(has_multiple_locations, false),
		       COALESCE(fiscal_code, ''), COALESCE(sdi_code, ''), COALESCE(pec_email, ''),
		       COALESCE(eu_vat_id, ''), COALESCE(billing_address, ''), COALESCE(billing_city, ''),
		       COALESCE(billing_province, ''), COALESCE(billing_postal_code, ''), COALESCE(billing_country, 'IT'),
		       status::text, email_verified, COALESCE(is_admin, false),
		       COALESCE(avatar_url, ''), COALESCE(social_links, '{}'), COALESCE(business_photos, '[]'),
		       stripe_customer_id, stripe_account_id,
		       COALESCE(total_co2_saved, 0), COALESCE(total_water_saved, 0), COALESCE(eco_credits, 0), COALESCE(eco_level, 'Germoglio'),
		       COALESCE(rating_avg, 0), COALESCE(rating_count, 0),
		       last_login_at, created_at, updated_at
		FROM users WHERE email = $1 AND deleted_at IS NULL
	`

	user := &models.User{}
	var status string
	var accountType string
	var socialLinksJSON []byte
	var businessPhotosJSON []byte

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Phone, &user.City, &user.Province, &user.PostalCode,
		&accountType, &user.BusinessName, &user.VATNumber,
		&user.HasMultipleLocations,
		&user.FiscalCode, &user.SDICode, &user.PECEmail,
		&user.EUVatID, &user.BillingAddress, &user.BillingCity,
		&user.BillingProvince, &user.BillingPostalCode, &user.BillingCountry,
		&status, &user.EmailVerified, &user.IsAdmin,
		&user.AvatarURL, &socialLinksJSON, &businessPhotosJSON,
		&user.StripeCustomerID, &user.StripeAccountID,
		&user.TotalCO2Saved, &user.TotalWaterSaved, &user.EcoCredits, &user.EcoLevel,
		&user.RatingAvg, &user.RatingCount,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user.Status = models.UserStatus(status)
	user.AccountType = models.AccountType(accountType)

	// Parse JSON fields
	json.Unmarshal(socialLinksJSON, &user.SocialLinks)
	json.Unmarshal(businessPhotosJSON, &user.BusinessPhotos)

	return user, nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET last_login_at = $1, updated_at = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, time.Now(), id)
	return err
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) UpdateProfile(ctx context.Context, id uuid.UUID, req *models.UpdateProfileRequest) error {
	// Build dynamic update query
	updates := []string{}
	args := []interface{}{}
	argNum := 1

	addField := func(column string, value interface{}) {
		updates = append(updates, fmt.Sprintf("%s = $%d", column, argNum))
		args = append(args, value)
		argNum++
	}

	if req.FirstName != nil {
		addField("first_name", *req.FirstName)
	}
	if req.LastName != nil {
		addField("last_name", *req.LastName)
	}
	if req.Phone != nil {
		addField("phone", *req.Phone)
	}
	if req.City != nil {
		addField("city", *req.City)
	}
	if req.Province != nil {
		addField("province", *req.Province)
	}
	if req.PostalCode != nil {
		addField("postal_code", *req.PostalCode)
	}
	if req.AccountType != nil {
		addField("account_type", string(*req.AccountType))
	}
	if req.BusinessName != nil {
		addField("business_name", *req.BusinessName)
	}
	if req.VATNumber != nil {
		addField("vat_number", *req.VATNumber)
	}
	if req.SocialLinks != nil {
		socialJSON, _ := json.Marshal(req.SocialLinks)
		addField("social_links", socialJSON)
	}
	if req.HasMultipleLocations != nil {
		addField("has_multiple_locations", *req.HasMultipleLocations)
	}
	// Billing fields
	if req.FiscalCode != nil {
		addField("fiscal_code", *req.FiscalCode)
	}
	if req.SDICode != nil {
		addField("sdi_code", *req.SDICode)
	}
	if req.PECEmail != nil {
		addField("pec_email", *req.PECEmail)
	}
	if req.EUVatID != nil {
		addField("eu_vat_id", *req.EUVatID)
	}
	if req.BillingAddress != nil {
		addField("billing_address", *req.BillingAddress)
	}
	if req.BillingCity != nil {
		addField("billing_city", *req.BillingCity)
	}
	if req.BillingProvince != nil {
		addField("billing_province", *req.BillingProvince)
	}
	if req.BillingPostalCode != nil {
		addField("billing_postal_code", *req.BillingPostalCode)
	}
	if req.BillingCountry != nil {
		addField("billing_country", *req.BillingCountry)
	}

	if len(updates) == 0 {
		return nil
	}

	updates = append(updates, "updated_at = NOW()")
	args = append(args, id)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(updates, ", "), argNum)
	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

func (r *UserRepository) UpdateAvatar(ctx context.Context, id uuid.UUID, avatarURL string) error {
	query := `UPDATE users SET avatar_url = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, avatarURL, id)
	return err
}

func (r *UserRepository) AddBusinessPhoto(ctx context.Context, id uuid.UUID, photoURL string) error {
	query := `UPDATE users SET business_photos = business_photos || $1::jsonb, updated_at = NOW() WHERE id = $2`
	photoJSON, _ := json.Marshal([]string{photoURL})
	_, err := r.pool.Exec(ctx, query, photoJSON, id)
	return err
}

// Location methods

func (r *UserRepository) CreateLocation(ctx context.Context, loc *models.Location) error {
	query := `
		INSERT INTO seller_locations (
			id, user_id, name, is_primary, is_active,
			address_street, address_city, address_province, address_postal_code,
			phone, email, pickup_hours, pickup_instructions, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	loc.ID = uuid.New()
	loc.IsActive = true
	loc.CreatedAt = time.Now()

	// Handle empty pickup_hours (JSONB column requires valid JSON or NULL)
	var pickupHours interface{}
	if loc.PickupHours == "" {
		pickupHours = nil
	} else {
		pickupHours = loc.PickupHours
	}

	_, err := r.pool.Exec(ctx, query,
		loc.ID, loc.UserID, loc.Name, loc.IsPrimary, loc.IsActive,
		loc.AddressStreet, loc.AddressCity, loc.AddressProvince, loc.AddressPostal,
		loc.Phone, loc.Email, pickupHours, loc.PickupInstructions, loc.CreatedAt, loc.CreatedAt,
	)

	// If this is primary, unset other primaries
	if loc.IsPrimary {
		r.pool.Exec(ctx, `UPDATE seller_locations SET is_primary = false WHERE user_id = $1 AND id != $2`, loc.UserID, loc.ID)
	}

	return err
}

func (r *UserRepository) GetLocationsByUser(ctx context.Context, userID uuid.UUID) ([]models.Location, error) {
	query := `
		SELECT id, user_id, name, is_primary, is_active,
		       address_street, address_city, COALESCE(address_province, ''), address_postal_code,
		       COALESCE(phone, ''), COALESCE(email, ''), COALESCE(pickup_hours::text, ''), COALESCE(pickup_instructions, ''),
		       created_at
		FROM seller_locations
		WHERE user_id = $1 AND is_active = true
		ORDER BY is_primary DESC, created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var loc models.Location
		err := rows.Scan(
			&loc.ID, &loc.UserID, &loc.Name, &loc.IsPrimary, &loc.IsActive,
			&loc.AddressStreet, &loc.AddressCity, &loc.AddressProvince, &loc.AddressPostal,
			&loc.Phone, &loc.Email, &loc.PickupHours, &loc.PickupInstructions,
			&loc.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}

	return locations, nil
}

func (r *UserRepository) GetLocationByID(ctx context.Context, id uuid.UUID) (*models.Location, error) {
	query := `
		SELECT id, user_id, name, is_primary, is_active,
		       address_street, address_city, COALESCE(address_province, ''), address_postal_code,
		       COALESCE(phone, ''), COALESCE(email, ''), COALESCE(pickup_hours::text, ''), COALESCE(pickup_instructions, ''),
		       created_at
		FROM seller_locations
		WHERE id = $1
	`

	loc := &models.Location{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&loc.ID, &loc.UserID, &loc.Name, &loc.IsPrimary, &loc.IsActive,
		&loc.AddressStreet, &loc.AddressCity, &loc.AddressProvince, &loc.AddressPostal,
		&loc.Phone, &loc.Email, &loc.PickupHours, &loc.PickupInstructions,
		&loc.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}

	return loc, nil
}

func (r *UserRepository) DeleteLocation(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE seller_locations SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// Admin methods

func (r *UserRepository) IsAdmin(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT COALESCE(is_admin, false) FROM users WHERE id = $1`
	var isAdmin bool
	err := r.pool.QueryRow(ctx, query, id).Scan(&isAdmin)
	return isAdmin, err
}
