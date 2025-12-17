package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/gecogreen/backend/internal/models"
	"github.com/gecogreen/backend/internal/repository"
)

type LeaderboardHandler struct {
	leaderboardRepo *repository.LeaderboardRepository
	userRepo        *repository.UserRepository
}

func NewLeaderboardHandler(leaderboardRepo *repository.LeaderboardRepository, userRepo *repository.UserRepository) *LeaderboardHandler {
	return &LeaderboardHandler{
		leaderboardRepo: leaderboardRepo,
		userRepo:        userRepo,
	}
}

// GetLeaderboard returns the leaderboard for a given period
// GET /api/leaderboard?period=WEEKLY|MONTHLY|YEARLY|ALLTIME
func (h *LeaderboardHandler) GetLeaderboard(c *fiber.Ctx) error {
	periodStr := c.Query("period", "MONTHLY")
	period := models.PeriodType(periodStr)

	// Validate period
	validPeriods := map[models.PeriodType]bool{
		models.PeriodWeekly:  true,
		models.PeriodMonthly: true,
		models.PeriodYearly:  true,
		models.PeriodAllTime: true,
	}
	if !validPeriods[period] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Periodo non valido"})
	}

	limit := c.QueryInt("limit", 50)
	if limit > 100 {
		limit = 100
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	entries, err := h.leaderboardRepo.GetLeaderboard(ctx, period, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nel recupero classifica"})
	}

	// Period name for display
	periodNames := map[models.PeriodType]string{
		models.PeriodWeekly:  "Questa settimana",
		models.PeriodMonthly: "Questo mese",
		models.PeriodYearly:  "Quest'anno",
		models.PeriodAllTime: "Sempre",
	}

	return c.JSON(models.LeaderboardResponse{
		Period:     period,
		PeriodName: periodNames[period],
		Entries:    entries,
		UpdatedAt:  time.Now(),
	})
}

// GetMyRank returns the current user's rank
// GET /api/leaderboard/my-rank?period=MONTHLY
func (h *LeaderboardHandler) GetMyRank(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	periodStr := c.Query("period", "MONTHLY")
	period := models.PeriodType(periodStr)

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	rank, co2Saved, err := h.leaderboardRepo.GetUserRank(ctx, user.ID, period)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nel recupero rank"})
	}

	return c.JSON(fiber.Map{
		"rank":          rank,
		"co2_saved":     co2Saved,
		"period":        period,
		"eco_credits":   user.EcoCredits,
		"total_co2":     user.TotalCO2Saved,
		"total_water":   user.TotalWaterSaved,
	})
}

// GetHallOfFame returns the Hall of Fame
// GET /api/leaderboard/hall-of-fame?period=MONTHLY
func (h *LeaderboardHandler) GetHallOfFame(c *fiber.Ctx) error {
	periodStr := c.Query("period", "MONTHLY")
	period := models.PeriodType(periodStr)

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	hof, err := h.leaderboardRepo.GetHallOfFame(ctx, period)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nel recupero Hall of Fame"})
	}

	return c.JSON(hof)
}

// GetCommunityStats returns global community statistics
// GET /api/leaderboard/community-stats
func (h *LeaderboardHandler) GetCommunityStats(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	stats, err := h.leaderboardRepo.GetCommunityStats(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nel recupero statistiche"})
	}

	return c.JSON(stats)
}

// GetMyAwards returns the current user's awards
// GET /api/leaderboard/my-awards
func (h *LeaderboardHandler) GetMyAwards(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	awards, err := h.leaderboardRepo.GetUserAwards(ctx, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nel recupero premi"})
	}

	return c.JSON(fiber.Map{"awards": awards})
}

// GetMyImpactHistory returns the current user's impact history
// GET /api/leaderboard/my-impact?limit=50
func (h *LeaderboardHandler) GetMyImpactHistory(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	limit := c.QueryInt("limit", 50)
	if limit > 100 {
		limit = 100
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	logs, err := h.leaderboardRepo.GetUserImpactHistory(ctx, user.ID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nel recupero storico"})
	}

	return c.JSON(fiber.Map{
		"history":     logs,
		"eco_credits": user.EcoCredits,
		"total_co2":   user.TotalCO2Saved,
	})
}

// RedeemReward allows users to redeem eco-credits for rewards
// POST /api/leaderboard/redeem
func (h *LeaderboardHandler) RedeemReward(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	var req models.RedeemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	// Define reward costs
	rewardCosts := map[string]int{
		"BOOST_24H":    100,
		"BOOST_7D":     500,
		"TOP_CATEGORY": 200,
		"TREE":         300,
		"BADGE":        150,
	}

	cost, valid := rewardCosts[req.RewardType]
	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tipo premio non valido"})
	}

	// Check if user has enough credits
	if user.EcoCredits < cost {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":           "EcoCredits insufficienti",
			"required":        cost,
			"current_balance": user.EcoCredits,
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	// Create impact log for redemption
	impactLog := &models.ImpactLog{
		UserID:          user.ID,
		ActionType:      getRedeemActionType(req.RewardType),
		EcoCreditsSpent: cost,
		Description:     "Riscatto: " + req.RewardType,
	}

	err := h.leaderboardRepo.AddImpactLog(ctx, impactLog)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nel riscatto"})
	}

	response := models.RedeemResponse{
		Success:      true,
		Message:      getRedeemMessage(req.RewardType),
		CreditsSpent: cost,
		NewBalance:   impactLog.EcoCreditsBalance,
	}

	// TODO: Handle specific reward logic (boost product, plant tree, etc.)
	// For trees, we would call Tree-Nation API and set TreeCertURL

	return c.JSON(response)
}

// GetFeaturedAwards returns featured awards for public display
// GET /api/leaderboard/featured?limit=10
func (h *LeaderboardHandler) GetFeaturedAwards(c *fiber.Ctx) error {
	periodStr := c.Query("period", "")
	period := models.PeriodType(periodStr)

	limit := c.QueryInt("limit", 10)
	if limit > 50 {
		limit = 50
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	awards, err := h.leaderboardRepo.GetFeaturedAwards(ctx, period, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nel recupero premi"})
	}

	return c.JSON(fiber.Map{"awards": awards})
}

// =====================
// ADMIN ENDPOINTS
// =====================

// AdminGetTasks returns all admin tasks
// GET /api/admin/awards/tasks
func (h *LeaderboardHandler) AdminGetTasks(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	if !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Accesso negato"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	tasks, err := h.leaderboardRepo.GetAdminTasks(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nel recupero tasks"})
	}

	return c.JSON(tasks)
}

// AdminCreateAward creates a new award
// POST /api/admin/awards
func (h *LeaderboardHandler) AdminCreateAward(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	if !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Accesso negato"})
	}

	var req models.CreateAwardRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	// Validate required fields
	if req.UserID == uuid.Nil || req.AwardType == "" || req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Campi obbligatori mancanti"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	// Verify user exists
	_, err := h.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Utente non trovato"})
	}

	award := &models.Award{
		UserID:          req.UserID,
		AwardType:       req.AwardType,
		PeriodType:      req.PeriodType,
		PeriodStart:     req.PeriodStart,
		PeriodEnd:       req.PeriodEnd,
		Title:           req.Title,
		Description:     req.Description,
		CO2Saved:        req.CO2Saved,
		ProductsCount:   req.ProductsCount,
		InterviewStatus: req.InterviewStatus,
	}

	// Set default interview status based on award type
	if award.InterviewStatus == "" {
		switch award.AwardType {
		case models.AwardEcoChampion, models.AwardNewEntry, models.AwardRecordBreaker:
			award.InterviewStatus = models.InterviewPending
		default:
			award.InterviewStatus = models.InterviewNotRequired
		}
	}

	err = h.leaderboardRepo.CreateAward(ctx, award)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nella creazione premio"})
	}

	return c.Status(fiber.StatusCreated).JSON(award)
}

// AdminUpdateAward updates an existing award
// PUT /api/admin/awards/:id
func (h *LeaderboardHandler) AdminUpdateAward(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	if !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Accesso negato"})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	var req models.UpdateAwardRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	err = h.leaderboardRepo.UpdateAward(ctx, id, &req)
	if err != nil {
		if err == repository.ErrAwardNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Premio non trovato"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nell'aggiornamento"})
	}

	// Return updated award
	award, _ := h.leaderboardRepo.GetAwardByID(ctx, id)
	return c.JSON(award)
}

// AdminGetAward returns a single award
// GET /api/admin/awards/:id
func (h *LeaderboardHandler) AdminGetAward(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	if !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Accesso negato"})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	award, err := h.leaderboardRepo.GetAwardByID(ctx, id)
	if err != nil {
		if err == repository.ErrAwardNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Premio non trovato"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nel recupero"})
	}

	return c.JSON(award)
}

// AdminUpdateTask updates a task
// PUT /api/admin/awards/tasks/:id
func (h *LeaderboardHandler) AdminUpdateTask(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	if !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Accesso negato"})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	var req models.UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	err = h.leaderboardRepo.UpdateTask(ctx, id, &req)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task non trovato"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nell'aggiornamento"})
	}

	task, _ := h.leaderboardRepo.GetTaskByID(ctx, id)
	return c.JSON(task)
}

// AdminCompleteTask marks a task as completed
// POST /api/admin/awards/tasks/:id/complete
func (h *LeaderboardHandler) AdminCompleteTask(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	if !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Accesso negato"})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID non valido"})
	}

	var req models.CompleteTaskRequest
	c.BodyParser(&req) // Optional body

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	err = h.leaderboardRepo.CompleteTask(ctx, id, user.ID, req.Notes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nel completamento"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Task completato"})
}

// AdminCreateTask creates a new manual task
// POST /api/admin/awards/tasks
func (h *LeaderboardHandler) AdminCreateTask(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	if !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Accesso negato"})
	}

	var task models.AdminContentTask
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dati non validi"})
	}

	if task.Title == "" || task.TaskType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Titolo e tipo sono obbligatori"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	err := h.leaderboardRepo.CreateTask(ctx, &task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Errore nella creazione"})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

// Helper functions
func getRedeemActionType(rewardType string) models.ImpactActionType {
	switch rewardType {
	case "BOOST_24H", "BOOST_7D":
		return models.ImpactRedeemBoost
	case "TREE":
		return models.ImpactRedeemTree
	case "BADGE":
		return models.ImpactRedeemBadge
	default:
		return models.ImpactRedeemBoost
	}
}

func getRedeemMessage(rewardType string) string {
	switch rewardType {
	case "BOOST_24H":
		return "Boost 24h attivato! Il tuo prodotto avrà più visibilità."
	case "BOOST_7D":
		return "Boost 7 giorni attivato! Massima visibilità per una settimana."
	case "TOP_CATEGORY":
		return "Il tuo prodotto è ora in evidenza nella categoria!"
	case "TREE":
		return "Grazie! Stiamo piantando un albero a tuo nome. Riceverai il certificato via email."
	case "BADGE":
		return "Badge esclusivo sbloccato! Visibile sul tuo profilo."
	default:
		return "Premio riscattato con successo!"
	}
}
