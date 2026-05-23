package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhaozeguang/timecapsule-api/internal/dto"
	"github.com/zhaozeguang/timecapsule-api/internal/middleware"
	"github.com/zhaozeguang/timecapsule-api/internal/service"
)

type MemoryHandler struct {
	memoryService *service.MemoryService
	aiService     *service.AIService
	jwtSecret     string
}

func NewMemoryHandler(memoryService *service.MemoryService, aiService *service.AIService, jwtSecret string) *MemoryHandler {
	return &MemoryHandler{memoryService: memoryService, aiService: aiService, jwtSecret: jwtSecret}
}

func (h *MemoryHandler) Register(r *gin.Engine) {
	memories := r.Group("/api/v1/memories", middleware.AuthRequired(h.jwtSecret))
	{
		memories.POST("", h.Create)
		memories.GET("", h.List)
		memories.GET("/:id", h.Get)
		memories.PATCH("/:id", h.Update)
		memories.DELETE("/:id", h.Delete)
	}

	ai := r.Group("/api/v1/ai", middleware.AuthRequired(h.jwtSecret))
	{
		ai.POST("/generate-memory", h.GenerateMemory)
		ai.POST("/regenerate", h.Regenerate)
	}
}

func (h *MemoryHandler) Create(c *gin.Context) {
	var req dto.CreateMemoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 40001, Message: err.Error()})
		return
	}

	userID := c.GetString("user_id")

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 40001, Message: "image file is required"})
		return
	}

	if file.Size > 20*1024*1024 {
		c.JSON(http.StatusRequestEntityTooLarge, dto.Response{Code: 41300, Message: "image must be under 20MB"})
		return
	}

	imageURL := "temp://" + file.Filename
	thumbnailURL := "temp://thumb_" + file.Filename

	memory, err := h.memoryService.Create(userID, &req, imageURL, thumbnailURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 50000, Message: err.Error()})
		return
	}

	if req.GenerateAI {
		go func() {
			result, aiErr := h.aiService.RegenerateForMemory(memory.ID, userID, "")
			if aiErr == nil {
				memory.AIContent = result.Content
				memory.AIModel = result.Model
			}
		}()
	}

	c.JSON(http.StatusCreated, dto.Response{Code: 0, Message: "ok", Data: memory})
}

func (h *MemoryHandler) List(c *gin.Context) {
	spaceID := c.Query("space_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	userID := c.GetString("user_id")
	result, err := h.memoryService.List(spaceID, userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 40001, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 0, Message: "ok", Data: result})
}

func (h *MemoryHandler) Get(c *gin.Context) {
	userID := c.GetString("user_id")
	memory, err := h.memoryService.GetByID(c.Param("id"), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 40400, Message: "memory not found"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 0, Message: "ok", Data: memory})
}

func (h *MemoryHandler) Update(c *gin.Context) {
	var req struct {
		AIContent string `json:"ai_content"`
		AIModel   string `json:"ai_model"`
		AIPrompt  string `json:"ai_prompt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 40001, Message: err.Error()})
		return
	}

	userID := c.GetString("user_id")
	memory, err := h.memoryService.Update(c.Param("id"), userID, req.AIContent, req.AIModel, req.AIPrompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 50000, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 0, Message: "ok", Data: memory})
}

func (h *MemoryHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	if err := h.memoryService.Delete(c.Param("id"), userID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 50000, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 0, Message: "ok"})
}

func (h *MemoryHandler) GenerateMemory(c *gin.Context) {
	var req dto.GenerateMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 40001, Message: err.Error()})
		return
	}

	userID := c.GetString("user_id")

	if req.MemoryID != "" {
		result, err := h.aiService.RegenerateForMemory(req.MemoryID, userID, req.Style)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.Response{Code: 50001, Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, dto.Response{Code: 0, Message: "ok", Data: result})
		return
	}

	result, err := h.aiService.GenerateMemory(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 50001, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 0, Message: "ok", Data: result})
}

func (h *MemoryHandler) Regenerate(c *gin.Context) {
	var req struct {
		MemoryID string `json:"memory_id" binding:"required"`
		Style    string `json:"style"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 40001, Message: err.Error()})
		return
	}

	userID := c.GetString("user_id")
	result, err := h.aiService.RegenerateForMemory(req.MemoryID, userID, req.Style)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 50001, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 0, Message: "ok", Data: result})
}
