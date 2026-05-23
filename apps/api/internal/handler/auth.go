package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhaozeguang/timecapsule-api/internal/dto"
	"github.com/zhaozeguang/timecapsule-api/internal/middleware"
	"github.com/zhaozeguang/timecapsule-api/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	jwtSecret   string
}

func NewAuthHandler(authService *service.AuthService, jwtSecret string) *AuthHandler {
	return &AuthHandler{authService: authService, jwtSecret: jwtSecret}
}

func (h *AuthHandler) Register(r *gin.Engine) {
	api := r.Group("/api/v1/auth")
	{
		api.POST("/login", h.Login)
		api.POST("/refresh", h.Refresh)
		api.GET("/me", middleware.AuthRequired(h.jwtSecret), h.Me)
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 40001, Message: err.Error()})
		return
	}

	result, err := h.authService.Login(req.Provider, req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 50000, Message: "login failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 0, Message: "ok", Data: result})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 40001, Message: err.Error()})
		return
	}

	result, err := h.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{Code: 40100, Message: "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 0, Message: "ok", Data: result})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.authService.GetMe(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 40400, Message: "user not found"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 0, Message: "ok", Data: user})
}
