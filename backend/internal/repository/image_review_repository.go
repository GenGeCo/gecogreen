package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gecogreen/backend/internal/models"
)

type ImageReviewRepository struct {
	pool *pgxpool.Pool
}

func NewImageReviewRepository(pool *pgxpool.Pool) *ImageReviewRepository {
	return &ImageReviewRepository{pool: pool}
}

// Create adds a new image review to the moderation queue
func (r *ImageReviewRepository) Create(ctx context.Context, review *models.ImageReview) error {
	query := `
		INSERT INTO image_reviews (
			id, user_id, product_id, image_url, image_type,
			detected_text, detected_phone, detected_email, detected_url, confidence_score,
			status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	review.ID = uuid.New()
	review.Status = "PENDING"
	review.CreatedAt = time.Now()

	_, err := r.pool.Exec(ctx, query,
		review.ID, review.UserID, review.ProductID, review.ImageURL, review.ImageType,
		review.DetectedText, review.DetectedPhone, review.DetectedEmail, review.DetectedURL, review.Confidence,
		review.Status, review.CreatedAt,
	)
	return err
}

// GetPending returns all pending image reviews
func (r *ImageReviewRepository) GetPending(ctx context.Context, limit, offset int) ([]models.ImageReview, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM image_reviews WHERE status = 'PENDING'`
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, user_id, product_id, image_url, image_type,
		       COALESCE(detected_text, ''), detected_phone, detected_email, detected_url, COALESCE(confidence_score, 0),
		       status, reviewed_by, reviewed_at, COALESCE(rejection_reason, ''), created_at
		FROM image_reviews
		WHERE status = 'PENDING'
		ORDER BY created_at ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reviews []models.ImageReview
	for rows.Next() {
		var rev models.ImageReview
		err := rows.Scan(
			&rev.ID, &rev.UserID, &rev.ProductID, &rev.ImageURL, &rev.ImageType,
			&rev.DetectedText, &rev.DetectedPhone, &rev.DetectedEmail, &rev.DetectedURL, &rev.Confidence,
			&rev.Status, &rev.ReviewedBy, &rev.ReviewedAt, &rev.RejectionReason, &rev.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		reviews = append(reviews, rev)
	}

	return reviews, total, nil
}

// GetByID returns an image review by ID
func (r *ImageReviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ImageReview, error) {
	query := `
		SELECT id, user_id, product_id, image_url, image_type,
		       COALESCE(detected_text, ''), detected_phone, detected_email, detected_url, COALESCE(confidence_score, 0),
		       status, reviewed_by, reviewed_at, COALESCE(rejection_reason, ''), created_at
		FROM image_reviews
		WHERE id = $1
	`

	rev := &models.ImageReview{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&rev.ID, &rev.UserID, &rev.ProductID, &rev.ImageURL, &rev.ImageType,
		&rev.DetectedText, &rev.DetectedPhone, &rev.DetectedEmail, &rev.DetectedURL, &rev.Confidence,
		&rev.Status, &rev.ReviewedBy, &rev.ReviewedAt, &rev.RejectionReason, &rev.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return rev, nil
}

// Approve marks an image review as approved
func (r *ImageReviewRepository) Approve(ctx context.Context, reviewID, adminID uuid.UUID) error {
	query := `
		UPDATE image_reviews
		SET status = 'APPROVED', reviewed_by = $1, reviewed_at = $2
		WHERE id = $3
	`
	_, err := r.pool.Exec(ctx, query, adminID, time.Now(), reviewID)
	return err
}

// Reject marks an image review as rejected
func (r *ImageReviewRepository) Reject(ctx context.Context, reviewID, adminID uuid.UUID, reason string) error {
	query := `
		UPDATE image_reviews
		SET status = 'REJECTED', reviewed_by = $1, reviewed_at = $2, rejection_reason = $3
		WHERE id = $4
	`
	_, err := r.pool.Exec(ctx, query, adminID, time.Now(), reason, reviewID)
	return err
}

// GetStats returns moderation statistics
func (r *ImageReviewRepository) GetStats(ctx context.Context) (map[string]int, error) {
	query := `
		SELECT status, COUNT(*) FROM image_reviews
		GROUP BY status
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := map[string]int{
		"PENDING":  0,
		"APPROVED": 0,
		"REJECTED": 0,
	}

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats[status] = count
	}

	return stats, nil
}
