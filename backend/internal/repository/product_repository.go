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
	ErrProductNotFound = errors.New("product not found")
)

// ProductRepository handles product database operations
type ProductRepository struct {
	pool *pgxpool.Pool
}

// NewProductRepository creates a new product repository
func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

// Create creates a new product
func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (
			id, seller_id, category_id, title, description, price, original_price,
			listing_type, shipping_method, shipping_cost, pickup_location_ids, quantity, quantity_available,
			expiry_date, is_dutch_auction, dutch_start_price, dutch_decrease_amount,
			dutch_decrease_hours, dutch_min_price, dutch_started_at,
			city, province, postal_code, latitude, longitude,
			images, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29
		)
	`

	product.ID = uuid.New()
	product.Status = models.ProductStatusActive
	product.QuantityAvail = product.Quantity
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	if product.IsDutchAuction && product.DutchStartPrice != nil {
		now := time.Now()
		product.DutchStartedAt = &now
		product.Price = *product.DutchStartPrice
	}

	// Default values
	if product.ListingType == "" {
		product.ListingType = models.ListingSale
	}
	if product.ShippingMethod == "" {
		product.ShippingMethod = models.ShippingPickup
	}

	// Convert images to JSON
	imagesJSON, _ := json.Marshal(product.Images)
	if product.Images == nil {
		imagesJSON = []byte("[]")
	}

	// Convert pickup_location_ids to JSON
	pickupLocationIDsJSON, _ := json.Marshal(product.PickupLocationIDs)
	if product.PickupLocationIDs == nil {
		pickupLocationIDsJSON = []byte("[]")
	}

	_, err := r.pool.Exec(ctx, query,
		product.ID,
		product.SellerID,
		product.CategoryID,
		product.Title,
		product.Description,
		product.Price,
		product.OriginalPrice,
		product.ListingType,
		product.ShippingMethod,
		product.ShippingCost,
		pickupLocationIDsJSON,
		product.Quantity,
		product.QuantityAvail,
		product.ExpiryDate,
		product.IsDutchAuction,
		product.DutchStartPrice,
		product.DutchDecreaseAmount,
		product.DutchDecreaseHours,
		product.DutchMinPrice,
		product.DutchStartedAt,
		product.City,
		product.Province,
		product.PostalCode,
		product.Latitude,
		product.Longitude,
		imagesJSON,
		product.Status,
		product.CreatedAt,
		product.UpdatedAt,
	)

	return err
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	query := `
		SELECT p.id, p.seller_id, p.category_id, p.title, p.description, p.price, p.original_price,
			   p.listing_type::text, p.shipping_method::text, p.shipping_cost, COALESCE(p.pickup_location_ids, '[]'),
			   p.quantity, p.quantity_available,
			   p.expiry_date, p.expiry_photo_url, p.is_dutch_auction, p.dutch_start_price,
			   p.dutch_decrease_amount, p.dutch_decrease_hours, p.dutch_min_price, p.dutch_started_at,
			   COALESCE(p.city, ''), COALESCE(p.province, ''), COALESCE(p.postal_code, ''), p.latitude, p.longitude,
			   p.images, p.status::text, p.view_count, p.favorite_count, p.created_at, p.updated_at,
			   u.id, u.first_name, u.last_name, COALESCE(u.avatar_url, ''), u.created_at
		FROM products p
		JOIN users u ON p.seller_id = u.id
		WHERE p.id = $1
	`

	product := &models.Product{}
	seller := &models.UserProfile{}
	var imagesJSON, pickupLocationIDsJSON []byte
	var listingType, shippingMethod, status string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&product.ID, &product.SellerID, &product.CategoryID, &product.Title, &product.Description,
		&product.Price, &product.OriginalPrice, &listingType, &shippingMethod,
		&product.ShippingCost, &pickupLocationIDsJSON, &product.Quantity, &product.QuantityAvail,
		&product.ExpiryDate, &product.ExpiryPhotoURL, &product.IsDutchAuction, &product.DutchStartPrice,
		&product.DutchDecreaseAmount, &product.DutchDecreaseHours, &product.DutchMinPrice, &product.DutchStartedAt,
		&product.City, &product.Province, &product.PostalCode, &product.Latitude, &product.Longitude,
		&imagesJSON, &status, &product.ViewCount, &product.FavoriteCount, &product.CreatedAt, &product.UpdatedAt,
		&seller.ID, &seller.FirstName, &seller.LastName, &seller.AvatarURL, &seller.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	// Convert string types
	product.ListingType = models.ListingType(listingType)
	product.ShippingMethod = models.ShippingMethod(shippingMethod)
	product.Status = models.ProductStatus(status)

	// Parse images JSON
	if imagesJSON != nil {
		_ = json.Unmarshal(imagesJSON, &product.Images)
	}

	// Parse pickup_location_ids JSON
	if pickupLocationIDsJSON != nil {
		_ = json.Unmarshal(pickupLocationIDsJSON, &product.PickupLocationIDs)
	}

	product.Seller = seller
	return product, nil
}

// List retrieves products with filters
func (r *ProductRepository) List(ctx context.Context, filters models.ProductFilters) (*models.ProductListResponse, error) {
	var conditions []string
	var args []interface{}
	argNum := 1

	// Default to active products
	if filters.Status != nil {
		conditions = append(conditions, fmt.Sprintf("p.status = $%d", argNum))
		args = append(args, *filters.Status)
		argNum++
	} else {
		conditions = append(conditions, fmt.Sprintf("p.status = $%d", argNum))
		args = append(args, models.ProductStatusActive)
		argNum++
	}

	if filters.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("p.category_id = $%d", argNum))
		args = append(args, *filters.CategoryID)
		argNum++
	}

	if filters.SellerID != nil {
		conditions = append(conditions, fmt.Sprintf("p.seller_id = $%d", argNum))
		args = append(args, *filters.SellerID)
		argNum++
	}

	if filters.ListingType != nil {
		conditions = append(conditions, fmt.Sprintf("p.listing_type = $%d", argNum))
		args = append(args, *filters.ListingType)
		argNum++
	}

	if filters.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("p.price >= $%d", argNum))
		args = append(args, *filters.MinPrice)
		argNum++
	}

	if filters.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("p.price <= $%d", argNum))
		args = append(args, *filters.MaxPrice)
		argNum++
	}

	if filters.City != nil && *filters.City != "" {
		conditions = append(conditions, fmt.Sprintf("p.city ILIKE $%d", argNum))
		args = append(args, "%"+*filters.City+"%")
		argNum++
	}

	if filters.Search != nil && *filters.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(p.title ILIKE $%d OR p.description ILIKE $%d)", argNum, argNum))
		args = append(args, "%"+*filters.Search+"%")
		argNum++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products p %s", whereClause)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	// Pagination
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PerPage < 1 || filters.PerPage > 100 {
		filters.PerPage = 20
	}
	offset := (filters.Page - 1) * filters.PerPage

	// Sort
	sortBy := "created_at"
	if filters.SortBy != "" {
		allowedSorts := map[string]bool{"created_at": true, "price": true, "view_count": true}
		if allowedSorts[filters.SortBy] {
			sortBy = filters.SortBy
		}
	}
	sortOrder := "DESC"
	if filters.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT p.id, p.seller_id, p.category_id, p.title, p.description, p.price, p.original_price,
			   p.listing_type::text, p.shipping_method::text, p.shipping_cost, p.quantity, p.quantity_available,
			   p.expiry_date, p.is_dutch_auction, p.dutch_start_price, p.dutch_min_price,
			   COALESCE(p.city, ''), COALESCE(p.province, ''), p.images, p.status::text,
			   p.view_count, p.favorite_count, p.created_at,
			   u.id, u.first_name, u.last_name, COALESCE(u.avatar_url, '')
		FROM products p
		JOIN users u ON p.seller_id = u.id
		%s
		ORDER BY p.%s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortBy, sortOrder, argNum, argNum+1)

	args = append(args, filters.PerPage, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var seller models.UserProfile
		var imagesJSON []byte
		var listingType, shippingMethod, status string

		err := rows.Scan(
			&p.ID, &p.SellerID, &p.CategoryID, &p.Title, &p.Description, &p.Price, &p.OriginalPrice,
			&listingType, &shippingMethod, &p.ShippingCost, &p.Quantity, &p.QuantityAvail,
			&p.ExpiryDate, &p.IsDutchAuction, &p.DutchStartPrice, &p.DutchMinPrice,
			&p.City, &p.Province, &imagesJSON, &status,
			&p.ViewCount, &p.FavoriteCount, &p.CreatedAt,
			&seller.ID, &seller.FirstName, &seller.LastName, &seller.AvatarURL,
		)
		if err != nil {
			return nil, err
		}

		// Convert string types
		p.ListingType = models.ListingType(listingType)
		p.ShippingMethod = models.ShippingMethod(shippingMethod)
		p.Status = models.ProductStatus(status)

		if imagesJSON != nil {
			_ = json.Unmarshal(imagesJSON, &p.Images)
		}

		p.Seller = &seller
		products = append(products, p)
	}

	if products == nil {
		products = []models.Product{}
	}

	totalPages := (total + filters.PerPage - 1) / filters.PerPage

	return &models.ProductListResponse{
		Products:   products,
		Total:      total,
		Page:       filters.Page,
		PerPage:    filters.PerPage,
		TotalPages: totalPages,
	}, nil
}

// Update updates a product
func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE products SET
			title = $1, description = $2, price = $3, quantity = $4, quantity_available = $5,
			status = $6, shipping_method = $7, shipping_cost = $8, updated_at = $9
		WHERE id = $10 AND seller_id = $11
	`
	product.UpdatedAt = time.Now()

	result, err := r.pool.Exec(ctx, query,
		product.Title, product.Description, product.Price, product.Quantity, product.QuantityAvail,
		product.Status, product.ShippingMethod, product.ShippingCost, product.UpdatedAt,
		product.ID, product.SellerID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrProductNotFound
	}
	return nil
}

// Delete soft-deletes a product
func (r *ProductRepository) Delete(ctx context.Context, id, sellerID uuid.UUID) error {
	query := `UPDATE products SET status = $1, updated_at = $2 WHERE id = $3 AND seller_id = $4`
	result, err := r.pool.Exec(ctx, query, models.ProductStatusDeleted, time.Now(), id, sellerID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrProductNotFound
	}
	return nil
}

// IncrementViewCount increments the view count
func (r *ProductRepository) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE products SET view_count = view_count + 1 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// AddImage adds an image URL to a product's images array
func (r *ProductRepository) AddImage(ctx context.Context, productID uuid.UUID, imageURL string) error {
	query := `
		UPDATE products
		SET images = images || $1::jsonb, updated_at = $2
		WHERE id = $3
	`
	imageJSON, _ := json.Marshal([]string{imageURL})
	_, err := r.pool.Exec(ctx, query, imageJSON, time.Now(), productID)
	return err
}

// GetImages gets all images for a product
func (r *ProductRepository) GetImages(ctx context.Context, productID uuid.UUID) ([]string, error) {
	query := `SELECT images FROM products WHERE id = $1`
	var imagesJSON []byte
	err := r.pool.QueryRow(ctx, query, productID).Scan(&imagesJSON)
	if err != nil {
		return nil, err
	}

	var images []string
	if imagesJSON != nil {
		_ = json.Unmarshal(imagesJSON, &images)
	}
	return images, nil
}
