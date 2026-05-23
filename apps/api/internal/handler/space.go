package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhaozeguang/timecapsule-api/internal/dto"
	"github.com/zhaozeguang/timecapsule-api/internal/middleware"
	"github.com/zhaozeguang/timecapsule-api/internal/service"
)

type SpaceHandler struct {
	spaceService *service.SpaceService
	jwtSecret    string
}

func NewSpaceHandler(spaceService *service.SpaceService, jwtSecret string) *SpaceHandler {
	return &SpaceHandler{spaceService: spaceService, jwtSecret: jwtSecret}
}

func (h *SpaceHandler) Register(r *gin.Engine) {
	api := r.Group("/api/v1/spaces", middleware.AuthRequired(h.jwtSecret))
	{
		api.POST("", h.Create)
		api.GET("", h.List)
		api.GET("/:id", h.Get)
		api.PATCH("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}
}

func (h *SpaceHandler) Create(c *gin.Context) {
	var req dto.CreateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 40001, Message: err.Error()})
		return
	}

	userID := c.GetString("user_id")
	space, err := h.spaceService.Create(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 50000, Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{Code: 0, Message: "ok", Data: space})
}

func (h *SpaceHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	userID := c.GetString("user_id")
	result, err := h.spaceService.List(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 50000, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 0, Message: "ok", Data: result})
}

func (h *SpaceHandler) Get(c *gin.Context) {
	userID := c.GetString("user_id")
	space, err := h.spaceService.GetByID(c.Param("id"), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 40400, Message: "space not found"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 0, Message: "ok", Data: space})
}

func (h *SpaceHandler) Update(c *gin.Context) {
	var req dto.CreateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 40001, Message: err.Error()})
		return
	}

	userID := c.GetString("user_id")
	space, err := h.spaceService.Update(c.Param("id"), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 50000, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 0, Message: "ok", Data: space})
}

func (h *SpaceHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	if err := h.spaceService.Delete(c.Param("id"), userID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 50000, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 0, Message: "ok"})
}
