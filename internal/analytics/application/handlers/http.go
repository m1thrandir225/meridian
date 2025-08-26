package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/analytics/application/services"
	"github.com/m1thrandir225/meridian/internal/analytics/domain"
	"github.com/m1thrandir225/meridian/pkg/logging"
	"go.uber.org/zap"
)

type HTTPHandler struct {
	analyticsService *services.AnalyticsService
	logger           *logging.Logger
}

func NewHTTPHandler(analyticsService *services.AnalyticsService, logger *logging.Logger) *HTTPHandler {
	return &HTTPHandler{
		analyticsService: analyticsService,
		logger:           logger,
	}
}

func (h *HTTPHandler) handleGetDashboard(c *gin.Context) {
	logger := h.logger.WithMethod("handleGetDashboard")
	logger.Info("Getting dashboard data")

	timeRangeStr := c.DefaultQuery("timeRange", "7d")
	timeRange, err := parseTimeRange(timeRangeStr)
	if err != nil {
		logger.Error("Invalid time range", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time range"})
		return
	}

	query := domain.GetDashboardDataQuery{
		TimeRange: timeRange,
	}

	dashboardData, err := h.analyticsService.GetDashboardData(c.Request.Context(), query)
	if err != nil {
		logger.Error("Failed to get dashboard data", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard data"})
		return
	}

	c.JSON(http.StatusOK, dashboardData)
}

func (h *HTTPHandler) handleGetUserGrowth(c *gin.Context) {
	logger := h.logger.WithMethod("handleGetUserGrowth")
	logger.Info("Getting user growth data")

	startDate, endDate, err := parseDateRange(c)
	if err != nil {
		logger.Error("Invalid date range", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date range"})
		return
	}

	interval := c.DefaultQuery("interval", "daily")

	query := domain.GetUserGrowthQuery{
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  interval,
	}

	growthData, err := h.analyticsService.GetUserGrowth(c.Request.Context(), query)
	if err != nil {
		logger.Error("Failed to get user growth data", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user growth data"})
		return
	}

	c.JSON(http.StatusOK, growthData)
}

func (h *HTTPHandler) handleGetMessageVolume(c *gin.Context) {
	logger := h.logger.WithMethod("handleGetMessageVolume")
	logger.Info("Getting message volume data")

	startDate, endDate, err := parseDateRange(c)
	if err != nil {
		logger.Error("Invalid date range", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date range"})
		return
	}

	var channelID *string
	if channelIDStr := c.Query("channelId"); channelIDStr != "" {
		channelID = &channelIDStr
	}

	query := domain.GetMessageVolumeQuery{
		StartDate: startDate,
		EndDate:   endDate,
		ChannelID: channelID,
	}

	volumeData, err := h.analyticsService.GetMessageVolume(c.Request.Context(), query)
	if err != nil {
		logger.Error("Failed to get message volume data", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get message volume data"})
		return
	}

	c.JSON(http.StatusOK, volumeData)
}

func (h *HTTPHandler) handleGetChannelActivity(c *gin.Context) {
	logger := h.logger.WithMethod("handleGetChannelActivity")
	logger.Info("Getting channel activity data")

	startDate, endDate, err := parseDateRange(c)
	if err != nil {
		logger.Error("Invalid date range", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date range"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		logger.Error("Invalid limit", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	query := domain.GetChannelActivityQuery{
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     limit,
	}

	activityData, err := h.analyticsService.GetChannelActivity(c.Request.Context(), query)
	if err != nil {
		logger.Error("Failed to get channel activity data", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get channel activity data"})
		return
	}

	c.JSON(http.StatusOK, activityData)
}

func (h *HTTPHandler) handleGetTopUsers(c *gin.Context) {
	logger := h.logger.WithMethod("handleGetTopUsers")
	logger.Info("Getting top users data")

	startDate, endDate, err := parseDateRange(c)
	if err != nil {
		logger.Error("Invalid date range", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date range"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		logger.Error("Invalid limit", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	query := domain.GetTopUsersQuery{
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     limit,
	}

	topUsers, err := h.analyticsService.GetTopUsers(c.Request.Context(), query)
	if err != nil {
		logger.Error("Failed to get top users data", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top users data"})
		return
	}

	c.JSON(http.StatusOK, topUsers)
}

func (h *HTTPHandler) handleGetReactionUsage(c *gin.Context) {
	logger := h.logger.WithMethod("handleGetReactionUsage")
	logger.Info("Getting reaction usage data")

	startDate, endDate, err := parseDateRange(c)
	if err != nil {
		logger.Error("Invalid date range", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date range"})
		return
	}

	query := domain.GetReactionUsageQuery{
		StartDate: startDate,
		EndDate:   endDate,
	}

	reactionData, err := h.analyticsService.GetReactionUsage(c.Request.Context(), query)
	if err != nil {
		logger.Error("Failed to get reaction usage data", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reaction usage data"})
		return
	}

	c.JSON(http.StatusOK, reactionData)
}

func (h *HTTPHandler) handleGetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "analytics",
	})
}

func parseTimeRange(timeRangeStr string) (time.Duration, error) {
	switch timeRangeStr {
	case "1h":
		return time.Hour, nil
	case "24h", "1d":
		return 24 * time.Hour, nil
	case "7d":
		return 7 * 24 * time.Hour, nil
	case "30d":
		return 30 * 24 * time.Hour, nil
	default:
		return 7 * 24 * time.Hour, nil
	}
}

func parseDateRange(c *gin.Context) (time.Time, time.Time, error) {
	startDateStr := c.DefaultQuery("startDate", time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("endDate", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return startDate, endDate, nil
}
