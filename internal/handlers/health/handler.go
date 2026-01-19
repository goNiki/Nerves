package health

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/goNiki/Nerves/internal/dto"
	"github.com/goNiki/Nerves/internal/service"
)

type Handler struct {
	healthService service.HealthService
}

func NewHandler(healthService service.HealthService) *Handler {
	return &Handler{
		healthService: healthService,
	}
}

func (h *Handler) Check(c *gin.Context) {
	status, err := h.healthService.Check(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	httpStatus := http.StatusOK
	if status.Status != "healthy" {
		httpStatus = http.StatusServiceUnavailable
	}

	response := dto.HealthResponse{
		Status:     status.Status,
		Components: status.Components,
	}

	c.JSON(httpStatus, response)
}
