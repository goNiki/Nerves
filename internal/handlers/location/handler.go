package location

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/goNiki/Nerves/internal/converter"
	"github.com/goNiki/Nerves/internal/dto"
	"github.com/goNiki/Nerves/internal/service"
)

type Handler struct {
	locationService service.LocationService
}

func NewHandler(locationService service.LocationService) *Handler {
	return &Handler{
		locationService: locationService,
	}
}

func (h *Handler) Check(c *gin.Context) {
	var req dto.CheckLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_REQUEST", err.Error()))
		return
	}

	checkReq := converter.CheckLocationRequestToModel(&req)

	response, err := h.locationService.Check(c.Request.Context(), checkReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "Failed to check location"))
		return
	}

	c.JSON(http.StatusOK, converter.LocationCheckToResponce(response))
}
