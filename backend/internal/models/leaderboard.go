package models

import (
	"time"

	"github.com/google/uuid"
)

// PeriodType represents the leaderboard period
type PeriodType string

const (
	PeriodWeekly  PeriodType = "WEEKLY"
	PeriodMonthly PeriodType = "MONTHLY"
	PeriodYearly  PeriodType = "YEARLY"
	PeriodAllTime PeriodType = "ALLTIME"
)

// AwardType represents the type of award
type AwardType string

const (
	AwardEcoChampion   AwardType = "ECO_CHAMPION"
	AwardEcoRunnerUp   AwardType = "ECO_RUNNER_UP"
	AwardEcoThird      AwardType = "ECO_THIRD"
	AwardNewEntry      AwardType = "NEW_ENTRY"
	AwardRecordBreaker AwardType = "RECORD_BREAKER"
	AwardTopWeek       AwardType = "TOP_WEEK"
	AwardEcoLegend     AwardType = "ECO_LEGEND"
	AwardMilestone     AwardType = "MILESTONE"
)

// InterviewStatus represents the status of interview for an award
type InterviewStatus string

const (
	InterviewNotRequired InterviewStatus = "NOT_REQUIRED"
	InterviewPending     InterviewStatus = "PENDING"
	InterviewContacted   InterviewStatus = "CONTACTED"
	InterviewScheduled   InterviewStatus = "SCHEDULED"
	InterviewRecorded    InterviewStatus = "RECORDED"
	InterviewPublished   InterviewStatus = "PUBLISHED"
	InterviewDeclined    InterviewStatus = "DECLINED"
)

// TaskStatus represents the status of an admin task
type TaskStatus string

const (
	TaskPending    TaskStatus = "PENDING"
	TaskInProgress TaskStatus = "IN_PROGRESS"
	TaskCompleted  TaskStatus = "COMPLETED"
	TaskSkipped    TaskStatus = "SKIPPED"
	TaskCancelled  TaskStatus = "CANCELLED"
)

// TaskPriority represents the priority of an admin task
type TaskPriority string

const (
	TaskPriorityUrgent TaskPriority = "URGENT"
	TaskPriorityHigh   TaskPriority = "HIGH"
	TaskPriorityNormal TaskPriority = "NORMAL"
	TaskPriorityLow    TaskPriority = "LOW"
)

// TaskType represents the type of admin task
type TaskType string

const (
	TaskInterviewContact   TaskType = "INTERVIEW_CONTACT"
	TaskInterviewSchedule  TaskType = "INTERVIEW_SCHEDULE"
	TaskInterviewRecord    TaskType = "INTERVIEW_RECORD"
	TaskYouTubeEdit        TaskType = "YOUTUBE_EDIT"
	TaskYouTubePublish     TaskType = "YOUTUBE_PUBLISH"
	TaskSocialPost         TaskType = "SOCIAL_POST"
	TaskHallOfFameUpdate   TaskType = "HALL_OF_FAME_UPDATE"
	TaskEmailWinner        TaskType = "EMAIL_WINNER"
	TaskBadgeAssign        TaskType = "BADGE_ASSIGN"
	TaskOther              TaskType = "OTHER"
)

// LeaderboardEntry represents a single entry in the leaderboard
type LeaderboardEntry struct {
	Rank              int       `json:"rank"`
	UserID            uuid.UUID `json:"user_id"`
	BusinessName      string    `json:"business_name,omitempty"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	AccountType       string    `json:"account_type"`
	City              string    `json:"city,omitempty"`
	TotalCO2Saved     float64   `json:"total_co2_saved"`
	TotalWaterSaved   float64   `json:"total_water_saved,omitempty"`
	TotalProductsSold int       `json:"total_products_sold"`
	TotalOrders       int       `json:"total_orders,omitempty"`
	AvatarURL         string    `json:"avatar_url,omitempty"`
}

// LeaderboardSnapshot represents a saved snapshot of rankings
type LeaderboardSnapshot struct {
	ID                uuid.UUID  `json:"id"`
	UserID            uuid.UUID  `json:"user_id"`
	PeriodType        PeriodType `json:"period_type"`
	PeriodStart       time.Time  `json:"period_start"`
	PeriodEnd         time.Time  `json:"period_end"`
	TotalCO2Saved     float64    `json:"total_co2_saved"`
	TotalWaterSaved   float64    `json:"total_water_saved"`
	TotalProductsSold int        `json:"total_products_sold"`
	TotalOrders       int        `json:"total_orders"`
	Rank              int        `json:"rank"`
	CreatedAt         time.Time  `json:"created_at"`
}

// Award represents an award given to a user
type Award struct {
	ID                   uuid.UUID       `json:"id"`
	UserID               uuid.UUID       `json:"user_id"`
	AwardType            AwardType       `json:"award_type"`
	PeriodType           PeriodType      `json:"period_type,omitempty"`
	PeriodStart          *time.Time      `json:"period_start,omitempty"`
	PeriodEnd            *time.Time      `json:"period_end,omitempty"`
	Title                string          `json:"title"`
	Description          string          `json:"description,omitempty"`
	BadgeURL             string          `json:"badge_url,omitempty"`
	CO2Saved             float64         `json:"co2_saved,omitempty"`
	ProductsCount        int             `json:"products_count,omitempty"`
	YouTubeURL           string          `json:"youtube_url,omitempty"`
	InterviewStatus      InterviewStatus `json:"interview_status"`
	InterviewScheduledAt *time.Time      `json:"interview_scheduled_at,omitempty"`
	InterviewNotes       string          `json:"interview_notes,omitempty"`
	IsFeatured           bool            `json:"is_featured"`
	IsPublic             bool            `json:"is_public"`
	CreatedAt            time.Time       `json:"created_at"`
	PublishedAt          *time.Time      `json:"published_at,omitempty"`

	// Joined data
	User *UserPublicMinimal `json:"user,omitempty"`
}

// AdminContentTask represents a task for admin to complete
type AdminContentTask struct {
	ID          uuid.UUID    `json:"id"`
	AwardID     *uuid.UUID   `json:"award_id,omitempty"`
	UserID      *uuid.UUID   `json:"user_id,omitempty"`
	TaskType    TaskType     `json:"task_type"`
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	DueDate     *time.Time   `json:"due_date,omitempty"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
	CompletedBy *uuid.UUID   `json:"completed_by,omitempty"`
	Notes       string       `json:"notes,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`

	// Joined data
	Award *Award             `json:"award,omitempty"`
	User  *UserPublicMinimal `json:"user,omitempty"`
}

// ImpactActionType represents the type of impact action
type ImpactActionType string

const (
	ImpactPurchase        ImpactActionType = "PURCHASE"
	ImpactSale            ImpactActionType = "SALE"
	ImpactGiftGiven       ImpactActionType = "GIFT_GIVEN"
	ImpactGiftReceived    ImpactActionType = "GIFT_RECEIVED"
	ImpactFirstPurchase   ImpactActionType = "FIRST_PURCHASE"
	ImpactFirstSale       ImpactActionType = "FIRST_SALE"
	ImpactPickupBonus     ImpactActionType = "PICKUP_BONUS"
	ImpactLastChanceBonus ImpactActionType = "LAST_CHANCE_BONUS"
	ImpactReviewBonus     ImpactActionType = "REVIEW_BONUS"
	ImpactReferralBonus   ImpactActionType = "REFERRAL_BONUS"
	ImpactMilestoneBonus  ImpactActionType = "MILESTONE_BONUS"
	ImpactSocialShare     ImpactActionType = "SOCIAL_SHARE"
	ImpactProfileComplete ImpactActionType = "PROFILE_COMPLETE"
	ImpactRedeemBoost     ImpactActionType = "REDEEM_BOOST"
	ImpactRedeemTree      ImpactActionType = "REDEEM_TREE"
	ImpactRedeemBadge     ImpactActionType = "REDEEM_BADGE"
	ImpactPointsExpired   ImpactActionType = "POINTS_EXPIRED"
	ImpactAdminAdjustment ImpactActionType = "ADMIN_ADJUSTMENT"
)

// ImpactLog represents a log entry for impact/credits
type ImpactLog struct {
	ID                uuid.UUID        `json:"id"`
	UserID            uuid.UUID        `json:"user_id"`
	OrderID           *uuid.UUID       `json:"order_id,omitempty"`
	ActionType        ImpactActionType `json:"action_type"`
	CO2Saved          float64          `json:"co2_saved"`
	WaterSaved        float64          `json:"water_saved"`
	EcoCreditsEarned  int              `json:"eco_credits_earned"`
	EcoCreditsSpent   int              `json:"eco_credits_spent"`
	EcoCreditsBalance int              `json:"eco_credits_balance"`
	Description       string           `json:"description,omitempty"`
	CreatedAt         time.Time        `json:"created_at"`
}

// HallOfFame represents the hall of fame data
type HallOfFame struct {
	EcoChampion *Award   `json:"eco_champion,omitempty"`
	RunnerUp    *Award   `json:"runner_up,omitempty"`
	Third       *Award   `json:"third,omitempty"`
	NewEntry    *Award   `json:"new_entry,omitempty"`
	RecordOf    *Award   `json:"record,omitempty"`
	Period      string   `json:"period"` // e.g., "Dicembre 2024"
	PeriodStart string   `json:"period_start"`
	PeriodEnd   string   `json:"period_end"`
}

// CommunityStats represents global community statistics
type CommunityStats struct {
	TotalCO2Saved    float64 `json:"total_co2_saved"`
	TotalWaterSaved  float64 `json:"total_water_saved"`
	TotalTreesPlanted int    `json:"total_trees_planted"`
	TotalProducts    int     `json:"total_products_saved"`
	TotalUsers       int     `json:"total_users"`
}

// --- Request/Response types ---

type LeaderboardResponse struct {
	Period     PeriodType         `json:"period"`
	PeriodName string             `json:"period_name"` // e.g., "Questa settimana"
	Entries    []LeaderboardEntry `json:"entries"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

type CreateAwardRequest struct {
	UserID          uuid.UUID   `json:"user_id" validate:"required"`
	AwardType       AwardType   `json:"award_type" validate:"required"`
	PeriodType      PeriodType  `json:"period_type,omitempty"`
	PeriodStart     *time.Time  `json:"period_start,omitempty"`
	PeriodEnd       *time.Time  `json:"period_end,omitempty"`
	Title           string      `json:"title" validate:"required"`
	Description     string      `json:"description,omitempty"`
	CO2Saved        float64     `json:"co2_saved,omitempty"`
	ProductsCount   int         `json:"products_count,omitempty"`
	InterviewStatus InterviewStatus `json:"interview_status,omitempty"`
}

type UpdateAwardRequest struct {
	YouTubeURL           *string          `json:"youtube_url,omitempty"`
	InterviewStatus      *InterviewStatus `json:"interview_status,omitempty"`
	InterviewScheduledAt *time.Time       `json:"interview_scheduled_at,omitempty"`
	InterviewNotes       *string          `json:"interview_notes,omitempty"`
	IsFeatured           *bool            `json:"is_featured,omitempty"`
	IsPublic             *bool            `json:"is_public,omitempty"`
}

type UpdateTaskRequest struct {
	Status   *TaskStatus   `json:"status,omitempty"`
	Priority *TaskPriority `json:"priority,omitempty"`
	DueDate  *time.Time    `json:"due_date,omitempty"`
	Notes    *string       `json:"notes,omitempty"`
}

type CompleteTaskRequest struct {
	Notes string `json:"notes,omitempty"`
}

type AdminTasksResponse struct {
	Urgent    []AdminContentTask `json:"urgent"`
	ThisWeek  []AdminContentTask `json:"this_week"`
	Later     []AdminContentTask `json:"later"`
	Completed []AdminContentTask `json:"completed"`
	Stats     TaskStats          `json:"stats"`
}

type TaskStats struct {
	TotalPending   int `json:"total_pending"`
	TotalCompleted int `json:"total_completed"`
	OverdueCount   int `json:"overdue_count"`
}

type RedeemRequest struct {
	RewardType string `json:"reward_type" validate:"required"` // "BOOST_24H", "BOOST_7D", "TREE", "TOP_CATEGORY", "BADGE"
	ProductID  *uuid.UUID `json:"product_id,omitempty"` // For boosts
}

type RedeemResponse struct {
	Success         bool   `json:"success"`
	Message         string `json:"message"`
	CreditsSpent    int    `json:"credits_spent"`
	NewBalance      int    `json:"new_balance"`
	TreeCertURL     string `json:"tree_certificate_url,omitempty"`
}
