package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhaozeguang/timecapsule-api/internal/dto"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Register(r *gin.Engine) {
	r.GET("/health", h.Health)
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, dto.Response{
		Code:    0,
		Message: "ok",
		Data:    gin.H{"status": "healthy"},
	})
}
