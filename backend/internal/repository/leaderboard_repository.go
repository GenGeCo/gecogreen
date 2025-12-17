package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gecogreen/backend/internal/models"
)

var (
	ErrAwardNotFound = errors.New("award not found")
	ErrTaskNotFound  = errors.New("task not found")
)

type LeaderboardRepository struct {
	pool *pgxpool.Pool
}

func NewLeaderboardRepository(pool *pgxpool.Pool) *LeaderboardRepository {
	return &LeaderboardRepository{pool: pool}
}

// =====================
// LEADERBOARD QUERIES
// =====================

// GetLeaderboard returns the leaderboard for a given period
func (r *LeaderboardRepository) GetLeaderboard(ctx context.Context, period models.PeriodType, limit int) ([]models.LeaderboardEntry, error) {
	var periodStart, periodEnd time.Time
	now := time.Now()

	switch period {
	case models.PeriodWeekly:
		// Start of current week (Monday)
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		periodStart = now.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
		periodEnd = periodStart.AddDate(0, 0, 7)
	case models.PeriodMonthly:
		periodStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		periodEnd = periodStart.AddDate(0, 1, 0)
	case models.PeriodYearly:
		periodStart = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		periodEnd = periodStart.AddDate(1, 0, 0)
	case models.PeriodAllTime:
		periodStart = time.Date(2020, 1, 1, 0, 0, 0, 0, now.Location())
		periodEnd = now.AddDate(10, 0, 0)
	}

	query := `
		WITH user_stats AS (
			SELECT
				u.id as user_id,
				COALESCE(u.business_name, '') as business_name,
				u.first_name,
				u.last_name,
				COALESCE(u.account_type::text, 'PRIVATE') as account_type,
				COALESCE(u.city, '') as city,
				COALESCE(u.avatar_url, '') as avatar_url,
				COALESCE(SUM(il.co2_saved), 0) as total_co2_saved,
				COALESCE(SUM(il.water_saved), 0) as total_water_saved,
				COUNT(DISTINCT CASE WHEN il.action_type IN ('SALE', 'PURCHASE') THEN il.order_id END) as total_orders,
				COALESCE(SUM(CASE WHEN il.action_type = 'SALE' THEN 1 ELSE 0 END), 0) as total_products_sold
			FROM users u
			LEFT JOIN impact_logs il ON u.id = il.user_id
				AND il.created_at >= $1
				AND il.created_at < $2
			WHERE u.deleted_at IS NULL
				AND u.status = 'ACTIVE'
			GROUP BY u.id, u.business_name, u.first_name, u.last_name, u.account_type, u.city, u.avatar_url
			HAVING COALESCE(SUM(il.co2_saved), 0) > 0
		)
		SELECT
			ROW_NUMBER() OVER (ORDER BY total_co2_saved DESC) as rank,
			user_id, business_name, first_name, last_name, account_type, city, avatar_url,
			total_co2_saved, total_water_saved, total_products_sold, total_orders
		FROM user_stats
		ORDER BY total_co2_saved DESC
		LIMIT $3
	`

	rows, err := r.pool.Query(ctx, query, periodStart, periodEnd, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.LeaderboardEntry
	for rows.Next() {
		var entry models.LeaderboardEntry
		err := rows.Scan(
			&entry.Rank,
			&entry.UserID,
			&entry.BusinessName,
			&entry.FirstName,
			&entry.LastName,
			&entry.AccountType,
			&entry.City,
			&entry.AvatarURL,
			&entry.TotalCO2Saved,
			&entry.TotalWaterSaved,
			&entry.TotalProductsSold,
			&entry.TotalOrders,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// GetUserRank returns the rank of a specific user for a given period
func (r *LeaderboardRepository) GetUserRank(ctx context.Context, userID uuid.UUID, period models.PeriodType) (int, float64, error) {
	var periodStart, periodEnd time.Time
	now := time.Now()

	switch period {
	case models.PeriodWeekly:
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		periodStart = now.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
		periodEnd = periodStart.AddDate(0, 0, 7)
	case models.PeriodMonthly:
		periodStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		periodEnd = periodStart.AddDate(0, 1, 0)
	case models.PeriodYearly:
		periodStart = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		periodEnd = periodStart.AddDate(1, 0, 0)
	case models.PeriodAllTime:
		periodStart = time.Date(2020, 1, 1, 0, 0, 0, 0, now.Location())
		periodEnd = now.AddDate(10, 0, 0)
	}

	query := `
		WITH user_stats AS (
			SELECT
				u.id as user_id,
				COALESCE(SUM(il.co2_saved), 0) as total_co2_saved
			FROM users u
			LEFT JOIN impact_logs il ON u.id = il.user_id
				AND il.created_at >= $1
				AND il.created_at < $2
			WHERE u.deleted_at IS NULL AND u.status = 'ACTIVE'
			GROUP BY u.id
		),
		ranked AS (
			SELECT user_id, total_co2_saved,
				ROW_NUMBER() OVER (ORDER BY total_co2_saved DESC) as rank
			FROM user_stats
			WHERE total_co2_saved > 0
		)
		SELECT rank, total_co2_saved FROM ranked WHERE user_id = $3
	`

	var rank int
	var co2Saved float64
	err := r.pool.QueryRow(ctx, query, periodStart, periodEnd, userID).Scan(&rank, &co2Saved)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, nil // Not ranked
		}
		return 0, 0, err
	}

	return rank, co2Saved, nil
}

// SaveSnapshot saves a leaderboard snapshot for historical tracking
func (r *LeaderboardRepository) SaveSnapshot(ctx context.Context, snapshot *models.LeaderboardSnapshot) error {
	query := `
		INSERT INTO leaderboard_snapshots (
			id, user_id, period_type, period_start, period_end,
			total_co2_saved, total_water_saved, total_products_sold, total_orders, rank, created_at
		) VALUES ($1, $2, $3::period_type, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	snapshot.ID = uuid.New()
	snapshot.CreatedAt = time.Now()

	_, err := r.pool.Exec(ctx, query,
		snapshot.ID, snapshot.UserID, string(snapshot.PeriodType),
		snapshot.PeriodStart, snapshot.PeriodEnd,
		snapshot.TotalCO2Saved, snapshot.TotalWaterSaved,
		snapshot.TotalProductsSold, snapshot.TotalOrders,
		snapshot.Rank, snapshot.CreatedAt,
	)

	return err
}

// =====================
// AWARDS
// =====================

// CreateAward creates a new award
func (r *LeaderboardRepository) CreateAward(ctx context.Context, award *models.Award) error {
	query := `
		INSERT INTO awards (
			id, user_id, award_type, period_type, period_start, period_end,
			title, description, badge_url, co2_saved, products_count,
			interview_status, is_featured, is_public, created_at
		) VALUES ($1, $2, $3::award_type, $4::period_type, $5, $6, $7, $8, $9, $10, $11, $12::interview_status, $13, $14, $15)
	`

	award.ID = uuid.New()
	award.CreatedAt = time.Now()
	if award.InterviewStatus == "" {
		award.InterviewStatus = models.InterviewNotRequired
	}

	var periodType *string
	if award.PeriodType != "" {
		pt := string(award.PeriodType)
		periodType = &pt
	}

	_, err := r.pool.Exec(ctx, query,
		award.ID, award.UserID, string(award.AwardType), periodType,
		award.PeriodStart, award.PeriodEnd,
		award.Title, award.Description, award.BadgeURL,
		award.CO2Saved, award.ProductsCount,
		string(award.InterviewStatus), award.IsFeatured, award.IsPublic, award.CreatedAt,
	)

	return err
}

// GetAwardByID retrieves an award by ID
func (r *LeaderboardRepository) GetAwardByID(ctx context.Context, id uuid.UUID) (*models.Award, error) {
	query := `
		SELECT a.id, a.user_id, a.award_type::text, COALESCE(a.period_type::text, ''),
			a.period_start, a.period_end, a.title, COALESCE(a.description, ''),
			COALESCE(a.badge_url, ''), COALESCE(a.co2_saved, 0), COALESCE(a.products_count, 0),
			COALESCE(a.youtube_url, ''), a.interview_status::text,
			a.interview_scheduled_at, COALESCE(a.interview_notes, ''),
			a.is_featured, a.is_public, a.created_at, a.published_at,
			u.id, COALESCE(u.business_name, ''), u.first_name, u.last_name,
			COALESCE(u.avatar_url, ''), COALESCE(u.city, '')
		FROM awards a
		JOIN users u ON a.user_id = u.id
		WHERE a.id = $1
	`

	award := &models.Award{User: &models.UserPublicMinimal{}}
	var awardType, periodType, interviewStatus string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&award.ID, &award.UserID, &awardType, &periodType,
		&award.PeriodStart, &award.PeriodEnd, &award.Title, &award.Description,
		&award.BadgeURL, &award.CO2Saved, &award.ProductsCount,
		&award.YouTubeURL, &interviewStatus,
		&award.InterviewScheduledAt, &award.InterviewNotes,
		&award.IsFeatured, &award.IsPublic, &award.CreatedAt, &award.PublishedAt,
		&award.User.ID, &award.User.BusinessName, &award.User.FirstName, &award.User.LastName,
		&award.User.AvatarURL, &award.User.City,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAwardNotFound
		}
		return nil, err
	}

	award.AwardType = models.AwardType(awardType)
	award.PeriodType = models.PeriodType(periodType)
	award.InterviewStatus = models.InterviewStatus(interviewStatus)

	return award, nil
}

// GetUserAwards returns all awards for a user
func (r *LeaderboardRepository) GetUserAwards(ctx context.Context, userID uuid.UUID) ([]models.Award, error) {
	query := `
		SELECT id, user_id, award_type::text, COALESCE(period_type::text, ''),
			period_start, period_end, title, COALESCE(description, ''),
			COALESCE(badge_url, ''), COALESCE(co2_saved, 0), COALESCE(products_count, 0),
			COALESCE(youtube_url, ''), interview_status::text,
			interview_scheduled_at, COALESCE(interview_notes, ''),
			is_featured, is_public, created_at, published_at
		FROM awards
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var awards []models.Award
	for rows.Next() {
		var award models.Award
		var awardType, periodType, interviewStatus string

		err := rows.Scan(
			&award.ID, &award.UserID, &awardType, &periodType,
			&award.PeriodStart, &award.PeriodEnd, &award.Title, &award.Description,
			&award.BadgeURL, &award.CO2Saved, &award.ProductsCount,
			&award.YouTubeURL, &interviewStatus,
			&award.InterviewScheduledAt, &award.InterviewNotes,
			&award.IsFeatured, &award.IsPublic, &award.CreatedAt, &award.PublishedAt,
		)
		if err != nil {
			return nil, err
		}

		award.AwardType = models.AwardType(awardType)
		award.PeriodType = models.PeriodType(periodType)
		award.InterviewStatus = models.InterviewStatus(interviewStatus)
		awards = append(awards, award)
	}

	return awards, nil
}

// GetFeaturedAwards returns featured/public awards for Hall of Fame
func (r *LeaderboardRepository) GetFeaturedAwards(ctx context.Context, periodType models.PeriodType, limit int) ([]models.Award, error) {
	query := `
		SELECT a.id, a.user_id, a.award_type::text, COALESCE(a.period_type::text, ''),
			a.period_start, a.period_end, a.title, COALESCE(a.description, ''),
			COALESCE(a.badge_url, ''), COALESCE(a.co2_saved, 0), COALESCE(a.products_count, 0),
			COALESCE(a.youtube_url, ''), a.interview_status::text,
			a.interview_scheduled_at, COALESCE(a.interview_notes, ''),
			a.is_featured, a.is_public, a.created_at, a.published_at,
			u.id, COALESCE(u.business_name, ''), u.first_name, u.last_name,
			COALESCE(u.avatar_url, ''), COALESCE(u.city, '')
		FROM awards a
		JOIN users u ON a.user_id = u.id
		WHERE a.is_public = true
			AND ($1 = '' OR a.period_type::text = $1)
		ORDER BY a.is_featured DESC, a.created_at DESC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, string(periodType), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var awards []models.Award
	for rows.Next() {
		award := models.Award{User: &models.UserPublicMinimal{}}
		var awardType, periodTypeStr, interviewStatus string

		err := rows.Scan(
			&award.ID, &award.UserID, &awardType, &periodTypeStr,
			&award.PeriodStart, &award.PeriodEnd, &award.Title, &award.Description,
			&award.BadgeURL, &award.CO2Saved, &award.ProductsCount,
			&award.YouTubeURL, &interviewStatus,
			&award.InterviewScheduledAt, &award.InterviewNotes,
			&award.IsFeatured, &award.IsPublic, &award.CreatedAt, &award.PublishedAt,
			&award.User.ID, &award.User.BusinessName, &award.User.FirstName, &award.User.LastName,
			&award.User.AvatarURL, &award.User.City,
		)
		if err != nil {
			return nil, err
		}

		award.AwardType = models.AwardType(awardType)
		award.PeriodType = models.PeriodType(periodTypeStr)
		award.InterviewStatus = models.InterviewStatus(interviewStatus)
		awards = append(awards, award)
	}

	return awards, nil
}

// UpdateAward updates an existing award
func (r *LeaderboardRepository) UpdateAward(ctx context.Context, id uuid.UUID, req *models.UpdateAwardRequest) error {
	updates := []string{}
	args := []interface{}{}
	argNum := 1

	if req.YouTubeURL != nil {
		updates = append(updates, fmt.Sprintf("youtube_url = $%d", argNum))
		args = append(args, *req.YouTubeURL)
		argNum++
	}
	if req.InterviewStatus != nil {
		updates = append(updates, fmt.Sprintf("interview_status = $%d::interview_status", argNum))
		args = append(args, string(*req.InterviewStatus))
		argNum++
	}
	if req.InterviewScheduledAt != nil {
		updates = append(updates, fmt.Sprintf("interview_scheduled_at = $%d", argNum))
		args = append(args, *req.InterviewScheduledAt)
		argNum++
	}
	if req.InterviewNotes != nil {
		updates = append(updates, fmt.Sprintf("interview_notes = $%d", argNum))
		args = append(args, *req.InterviewNotes)
		argNum++
	}
	if req.IsFeatured != nil {
		updates = append(updates, fmt.Sprintf("is_featured = $%d", argNum))
		args = append(args, *req.IsFeatured)
		argNum++
	}
	if req.IsPublic != nil {
		updates = append(updates, fmt.Sprintf("is_public = $%d", argNum))
		args = append(args, *req.IsPublic)
		argNum++

		if *req.IsPublic {
			updates = append(updates, fmt.Sprintf("published_at = $%d", argNum))
			args = append(args, time.Now())
			argNum++
		}
	}

	if len(updates) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE awards SET %s WHERE id = $%d",
		joinStrings(updates, ", "), argNum)

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// =====================
// ADMIN TASKS
// =====================

// GetAdminTasks returns admin tasks organized by priority/due date
func (r *LeaderboardRepository) GetAdminTasks(ctx context.Context) (*models.AdminTasksResponse, error) {
	query := `
		SELECT t.id, t.award_id, t.user_id, t.task_type::text, t.title,
			COALESCE(t.description, ''), t.status::text, t.priority::text,
			t.due_date, t.completed_at, t.completed_by, COALESCE(t.notes, ''),
			t.created_at, t.updated_at
		FROM admin_content_tasks t
		ORDER BY
			CASE t.priority
				WHEN 'URGENT' THEN 1
				WHEN 'HIGH' THEN 2
				WHEN 'NORMAL' THEN 3
				WHEN 'LOW' THEN 4
			END,
			t.due_date NULLS LAST,
			t.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	response := &models.AdminTasksResponse{
		Urgent:    []models.AdminContentTask{},
		ThisWeek:  []models.AdminContentTask{},
		Later:     []models.AdminContentTask{},
		Completed: []models.AdminContentTask{},
	}

	now := time.Now()
	weekEnd := now.AddDate(0, 0, 7)

	for rows.Next() {
		var task models.AdminContentTask
		var taskType, status, priority string

		err := rows.Scan(
			&task.ID, &task.AwardID, &task.UserID, &taskType, &task.Title,
			&task.Description, &status, &priority,
			&task.DueDate, &task.CompletedAt, &task.CompletedBy, &task.Notes,
			&task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		task.TaskType = models.TaskType(taskType)
		task.Status = models.TaskStatus(status)
		task.Priority = models.TaskPriority(priority)

		// Categorize tasks
		if task.Status == models.TaskCompleted || task.Status == models.TaskSkipped || task.Status == models.TaskCancelled {
			response.Completed = append(response.Completed, task)
		} else if task.Priority == models.TaskPriorityUrgent {
			response.Urgent = append(response.Urgent, task)
		} else if task.DueDate != nil && task.DueDate.Before(weekEnd) {
			response.ThisWeek = append(response.ThisWeek, task)
		} else {
			response.Later = append(response.Later, task)
		}
	}

	// Get stats
	statsQuery := `
		SELECT
			COUNT(*) FILTER (WHERE status NOT IN ('COMPLETED', 'SKIPPED', 'CANCELLED')) as pending,
			COUNT(*) FILTER (WHERE status = 'COMPLETED') as completed,
			COUNT(*) FILTER (WHERE due_date < NOW() AND status NOT IN ('COMPLETED', 'SKIPPED', 'CANCELLED')) as overdue
		FROM admin_content_tasks
	`

	err = r.pool.QueryRow(ctx, statsQuery).Scan(
		&response.Stats.TotalPending,
		&response.Stats.TotalCompleted,
		&response.Stats.OverdueCount,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetTaskByID retrieves a task by ID
func (r *LeaderboardRepository) GetTaskByID(ctx context.Context, id uuid.UUID) (*models.AdminContentTask, error) {
	query := `
		SELECT t.id, t.award_id, t.user_id, t.task_type::text, t.title,
			COALESCE(t.description, ''), t.status::text, t.priority::text,
			t.due_date, t.completed_at, t.completed_by, COALESCE(t.notes, ''),
			t.created_at, t.updated_at
		FROM admin_content_tasks t
		WHERE t.id = $1
	`

	task := &models.AdminContentTask{}
	var taskType, status, priority string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&task.ID, &task.AwardID, &task.UserID, &taskType, &task.Title,
		&task.Description, &status, &priority,
		&task.DueDate, &task.CompletedAt, &task.CompletedBy, &task.Notes,
		&task.CreatedAt, &task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	task.TaskType = models.TaskType(taskType)
	task.Status = models.TaskStatus(status)
	task.Priority = models.TaskPriority(priority)

	return task, nil
}

// UpdateTask updates an existing task
func (r *LeaderboardRepository) UpdateTask(ctx context.Context, id uuid.UUID, req *models.UpdateTaskRequest) error {
	updates := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argNum := 1

	if req.Status != nil {
		updates = append(updates, fmt.Sprintf("status = $%d::task_status", argNum))
		args = append(args, string(*req.Status))
		argNum++
	}
	if req.Priority != nil {
		updates = append(updates, fmt.Sprintf("priority = $%d::task_priority", argNum))
		args = append(args, string(*req.Priority))
		argNum++
	}
	if req.DueDate != nil {
		updates = append(updates, fmt.Sprintf("due_date = $%d", argNum))
		args = append(args, *req.DueDate)
		argNum++
	}
	if req.Notes != nil {
		updates = append(updates, fmt.Sprintf("notes = $%d", argNum))
		args = append(args, *req.Notes)
		argNum++
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE admin_content_tasks SET %s WHERE id = $%d",
		joinStrings(updates, ", "), argNum)

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// CompleteTask marks a task as completed
func (r *LeaderboardRepository) CompleteTask(ctx context.Context, id uuid.UUID, completedBy uuid.UUID, notes string) error {
	query := `
		UPDATE admin_content_tasks
		SET status = 'COMPLETED'::task_status,
			completed_at = NOW(),
			completed_by = $1,
			notes = CASE WHEN $2 = '' THEN notes ELSE $2 END,
			updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.pool.Exec(ctx, query, completedBy, notes, id)
	return err
}

// CreateTask creates a new admin task
func (r *LeaderboardRepository) CreateTask(ctx context.Context, task *models.AdminContentTask) error {
	query := `
		INSERT INTO admin_content_tasks (
			id, award_id, user_id, task_type, title, description,
			status, priority, due_date, created_at, updated_at
		) VALUES ($1, $2, $3, $4::task_type, $5, $6, $7::task_status, $8::task_priority, $9, $10, $10)
	`

	task.ID = uuid.New()
	task.CreatedAt = time.Now()
	task.UpdatedAt = task.CreatedAt
	if task.Status == "" {
		task.Status = models.TaskPending
	}
	if task.Priority == "" {
		task.Priority = models.TaskPriorityNormal
	}

	_, err := r.pool.Exec(ctx, query,
		task.ID, task.AwardID, task.UserID, string(task.TaskType),
		task.Title, task.Description,
		string(task.Status), string(task.Priority), task.DueDate,
		task.CreatedAt,
	)

	return err
}

// =====================
// IMPACT LOG
// =====================

// AddImpactLog adds an impact log entry
func (r *LeaderboardRepository) AddImpactLog(ctx context.Context, log *models.ImpactLog) error {
	// First get current balance
	var currentBalance int
	err := r.pool.QueryRow(ctx,
		"SELECT COALESCE(eco_credits, 0) FROM users WHERE id = $1",
		log.UserID,
	).Scan(&currentBalance)
	if err != nil {
		return err
	}

	log.ID = uuid.New()
	log.EcoCreditsBalance = currentBalance + log.EcoCreditsEarned - log.EcoCreditsSpent
	log.CreatedAt = time.Now()

	// Start transaction
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert impact log
	query := `
		INSERT INTO impact_logs (
			id, user_id, order_id, action_type, co2_saved, water_saved,
			eco_credits_earned, eco_credits_spent, eco_credits_balance,
			description, created_at
		) VALUES ($1, $2, $3, $4::impact_action_type, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = tx.Exec(ctx, query,
		log.ID, log.UserID, log.OrderID, string(log.ActionType),
		log.CO2Saved, log.WaterSaved,
		log.EcoCreditsEarned, log.EcoCreditsSpent, log.EcoCreditsBalance,
		log.Description, log.CreatedAt,
	)
	if err != nil {
		return err
	}

	// Update user totals
	updateQuery := `
		UPDATE users SET
			total_co2_saved = COALESCE(total_co2_saved, 0) + $1,
			total_water_saved = COALESCE(total_water_saved, 0) + $2,
			eco_credits = $3,
			updated_at = NOW()
		WHERE id = $4
	`

	_, err = tx.Exec(ctx, updateQuery,
		log.CO2Saved, log.WaterSaved, log.EcoCreditsBalance, log.UserID,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetUserImpactHistory returns impact history for a user
func (r *LeaderboardRepository) GetUserImpactHistory(ctx context.Context, userID uuid.UUID, limit int) ([]models.ImpactLog, error) {
	query := `
		SELECT id, user_id, order_id, action_type::text, co2_saved, water_saved,
			eco_credits_earned, eco_credits_spent, eco_credits_balance,
			COALESCE(description, ''), created_at
		FROM impact_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.ImpactLog
	for rows.Next() {
		var log models.ImpactLog
		var actionType string

		err := rows.Scan(
			&log.ID, &log.UserID, &log.OrderID, &actionType,
			&log.CO2Saved, &log.WaterSaved,
			&log.EcoCreditsEarned, &log.EcoCreditsSpent, &log.EcoCreditsBalance,
			&log.Description, &log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		log.ActionType = models.ImpactActionType(actionType)
		logs = append(logs, log)
	}

	return logs, nil
}

// =====================
// COMMUNITY STATS
// =====================

// GetCommunityStats returns global community statistics
func (r *LeaderboardRepository) GetCommunityStats(ctx context.Context) (*models.CommunityStats, error) {
	query := `
		SELECT
			COALESCE(SUM(total_co2_saved), 0),
			COALESCE(SUM(total_water_saved), 0),
			COALESCE((SELECT SUM(trees_count) FROM trees_planted), 0),
			COALESCE((SELECT COUNT(*) FROM products WHERE deleted_at IS NULL), 0),
			COUNT(*)
		FROM users
		WHERE deleted_at IS NULL AND status = 'ACTIVE'
	`

	stats := &models.CommunityStats{}
	err := r.pool.QueryRow(ctx, query).Scan(
		&stats.TotalCO2Saved,
		&stats.TotalWaterSaved,
		&stats.TotalTreesPlanted,
		&stats.TotalProducts,
		&stats.TotalUsers,
	)

	return stats, err
}

// =====================
// HALL OF FAME
// =====================

// GetHallOfFame returns the current Hall of Fame
func (r *LeaderboardRepository) GetHallOfFame(ctx context.Context, periodType models.PeriodType) (*models.HallOfFame, error) {
	now := time.Now()
	var periodStart, periodEnd time.Time
	var periodName string

	switch periodType {
	case models.PeriodMonthly:
		periodStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		periodEnd = periodStart.AddDate(0, 1, 0)
		periodName = periodStart.Format("January 2006")
	default:
		periodStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		periodEnd = periodStart.AddDate(0, 1, 0)
		periodName = periodStart.Format("January 2006")
	}

	hof := &models.HallOfFame{
		Period:      periodName,
		PeriodStart: periodStart.Format("2006-01-02"),
		PeriodEnd:   periodEnd.Format("2006-01-02"),
	}

	// Get awards for this period
	query := `
		SELECT a.id, a.user_id, a.award_type::text, COALESCE(a.period_type::text, ''),
			a.period_start, a.period_end, a.title, COALESCE(a.description, ''),
			COALESCE(a.badge_url, ''), COALESCE(a.co2_saved, 0), COALESCE(a.products_count, 0),
			COALESCE(a.youtube_url, ''), a.interview_status::text,
			a.interview_scheduled_at, COALESCE(a.interview_notes, ''),
			a.is_featured, a.is_public, a.created_at, a.published_at,
			u.id, COALESCE(u.business_name, ''), u.first_name, u.last_name,
			COALESCE(u.avatar_url, ''), COALESCE(u.city, '')
		FROM awards a
		JOIN users u ON a.user_id = u.id
		WHERE a.is_public = true
			AND a.period_start >= $1
			AND a.period_end <= $2
		ORDER BY a.award_type, a.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, periodStart, periodEnd)
	if err != nil {
		return hof, err
	}
	defer rows.Close()

	for rows.Next() {
		award := &models.Award{User: &models.UserPublicMinimal{}}
		var awardType, periodTypeStr, interviewStatus string

		err := rows.Scan(
			&award.ID, &award.UserID, &awardType, &periodTypeStr,
			&award.PeriodStart, &award.PeriodEnd, &award.Title, &award.Description,
			&award.BadgeURL, &award.CO2Saved, &award.ProductsCount,
			&award.YouTubeURL, &interviewStatus,
			&award.InterviewScheduledAt, &award.InterviewNotes,
			&award.IsFeatured, &award.IsPublic, &award.CreatedAt, &award.PublishedAt,
			&award.User.ID, &award.User.BusinessName, &award.User.FirstName, &award.User.LastName,
			&award.User.AvatarURL, &award.User.City,
		)
		if err != nil {
			continue
		}

		award.AwardType = models.AwardType(awardType)
		award.PeriodType = models.PeriodType(periodTypeStr)
		award.InterviewStatus = models.InterviewStatus(interviewStatus)

		switch award.AwardType {
		case models.AwardEcoChampion:
			hof.EcoChampion = award
		case models.AwardEcoRunnerUp:
			hof.RunnerUp = award
		case models.AwardEcoThird:
			hof.Third = award
		case models.AwardNewEntry:
			hof.NewEntry = award
		case models.AwardRecordBreaker:
			hof.RecordOf = award
		}
	}

	return hof, nil
}

// Helper function
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
