package stats

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/goNiki/Nerves/internal/dto"
	"github.com/goNiki/Nerves/internal/service"
)

type Handler struct {
	statsService service.StatsService
}

func NewHandler(statsService service.StatsService) *Handler {
	return &Handler{
		statsService: statsService,
	}
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.statsService.GetIncidentStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "Failed to get stats"))
		return
	}

	items := make([]dto.IncidentStatsResponse, len(stats.Items))
	for i, stat := range stats.Items {
		items[i] = dto.IncidentStatsResponse{
			IncidentID: stat.IncidentID.String(),
			Title:      stat.Title,
			UserCount:  stat.UserCount,
		}
	}

	response := dto.IncidentStatsListResponse{
		Items:             items,
		TimeWindowMinutes: stats.TimeWindowMinutes,
	}

	c.JSON(http.StatusOK, response)
}
