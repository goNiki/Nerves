package incident

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	incidents := rg.Group("/incidents")
	incidents.POST("", h.Create)
	incidents.GET("", h.List)
	incidents.GET("/:id", h.GetByID)
	incidents.PUT("/:id", h.Replace)
	incidents.PATCH("/:id", h.Update)
	incidents.DELETE("/:id", h.Delete)
}
