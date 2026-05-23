package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/zhaozeguang/timecapsule-api/internal/config"
	"github.com/zhaozeguang/timecapsule-api/internal/dto"
	"github.com/zhaozeguang/timecapsule-api/internal/middleware"
	"github.com/zhaozeguang/timecapsule-api/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	cfg         *config.Config
}

func NewAuthHandler(authService *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{authService: authService, cfg: cfg}
}

func (h *AuthHandler) Register(r *gin.Engine) {
	api := r.Group("/api/v1/auth")
	{
		api.POST("/login", h.Login)
		api.POST("/refresh", h.Refresh)
		api.GET("/me", middleware.AuthRequired(h.cfg.JWTSecret), h.Me)
		api.GET("/google", h.GoogleLogin)
		api.GET("/github", h.GitHubLogin)
		api.GET("/callback/google", h.GoogleCallback)
		api.GET("/callback/github", h.GitHubCallback)
		api.POST("/dev-login", h.DevLogin)
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

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	redirectURI := h.cfg.GoogleRedirectURL
	authURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid+email+profile",
		url.QueryEscape(h.cfg.GoogleClientID),
		url.QueryEscape(redirectURI),
	)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *AuthHandler) GitHubLogin(c *gin.Context) {
	redirectURI := h.cfg.GitHubRedirectURL
	authURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email",
		url.QueryEscape(h.cfg.GitHubClientID),
		url.QueryEscape(redirectURI),
	)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 40001, Message: "missing code"})
		return
	}

	result, err := h.authService.Login("google", code)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/auth/callback?error=%s", h.cfg.FrontendURL, url.QueryEscape(err.Error())))
		return
	}

	c.Redirect(http.StatusTemporaryRedirect,
		fmt.Sprintf("%s/auth/callback?access_token=%s&refresh_token=%s&expires_in=%d",
			h.cfg.FrontendURL, result.AccessToken, result.RefreshToken, result.ExpiresIn))
}

func (h *AuthHandler) GitHubCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 40001, Message: "missing code"})
		return
	}

	result, err := h.authService.Login("github", code)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s/auth/callback?error=%s", h.cfg.FrontendURL, url.QueryEscape(err.Error())))
		return
	}

	c.Redirect(http.StatusTemporaryRedirect,
		fmt.Sprintf("%s/auth/callback?access_token=%s&refresh_token=%s&expires_in=%d",
			h.cfg.FrontendURL, result.AccessToken, result.RefreshToken, result.ExpiresIn))
}

func (h *AuthHandler) DevLogin(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
		Name  string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 40001, Message: err.Error()})
		return
	}

	result, err := h.authService.DevLogin(req.Email, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 50000, Message: "dev login failed"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 0, Message: "ok", Data: result})
}
