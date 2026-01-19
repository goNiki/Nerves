package incident

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/goNiki/Nerves/internal/converter"
	"github.com/goNiki/Nerves/internal/dto"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
	"github.com/goNiki/Nerves/internal/service"
)

type Handler struct {
	incidentService service.IncidentService
}

func NewHandler(incidentService service.IncidentService) *Handler {
	return &Handler{
		incidentService: incidentService,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateIncidentRequests
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_REQUEST", err.Error()))
		return
	}

	incident := converter.CreateIncidentRequestsToModel(&req)

	created, err := h.incidentService.Create(c.Request.Context(), incident)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "Failed to create incident"))
		return
	}

	c.JSON(http.StatusCreated, converter.IncidentToResponce(created))
}

func (h *Handler) List(c *gin.Context) {
	var query dto.ListIncidentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_REQUEST", err.Error()))
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 20
	}

	list, err := h.incidentService.List(c.Request.Context(), query.Page, query.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "Failed to list incidents"))
		return
	}

	c.JSON(http.StatusOK, converter.IncidentListToResponce(list))
}

func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_REQUEST", "Invalid incident ID"))
		return
	}

	incident, err := h.incidentService.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, errorapp.ErrIncidentNotFound) {
			c.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Incident not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "Failed to get incident"))
		return
	}

	c.JSON(http.StatusOK, converter.IncidentToResponce(incident))
}

func (h *Handler) Replace(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_REQUEST", "Invalid incident ID"))
		return
	}

	var req dto.CreateIncidentRequests
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_REQUEST", err.Error()))
		return
	}

	incident := converter.CreateIncidentRequestsToModel(&req)

	updated, err := h.incidentService.Replace(c.Request.Context(), id, incident)
	if err != nil {
		if errors.Is(err, errorapp.ErrIncidentNotFound) {
			c.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Incident not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "Failed to replace incident"))
		return
	}

	c.JSON(http.StatusOK, converter.IncidentToResponce(updated))
}

func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_REQUEST", "Invalid incident ID"))
		return
	}

	var req dto.UpdateIncidentRequests
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_REQUEST", err.Error()))
		return
	}

	updates := converter.UpdateIncidentRequestsToModel(&req)

	updated, err := h.incidentService.Update(c.Request.Context(), id, updates)
	if err != nil {
		if errors.Is(err, errorapp.ErrIncidentNotFound) {
			c.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Incident not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "Failed to update incident"))
		return
	}

	c.JSON(http.StatusOK, converter.IncidentToResponce(updated))
}

func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("INVALID_REQUEST", "Invalid incident ID"))
		return
	}

	if err := h.incidentService.Deactivate(c.Request.Context(), id); err != nil {
		if errors.Is(err, errorapp.ErrIncidentNotFound) {
			c.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Incident not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", "Failed to delete incident"))
		return
	}

	c.Status(http.StatusNoContent)
}
