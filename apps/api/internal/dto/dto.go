package dto

import "time"

// Unified response
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedData struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// Auth
type LoginRequest struct {
	Provider string `json:"provider" binding:"required,oneof=google github"`
	Code     string `json:"code" binding:"required"`
}

// Memory Space
type CreateSpaceRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
	Theme       string `json:"theme" binding:"oneof=default warm ocean forest sunset"`
}

// Memory
type CreateMemoryRequest struct {
	Title      string    `json:"title" binding:"required,min=1,max=300"`
	Content    string    `json:"content"`
	Location   string    `json:"location"`
	MemoryDate time.Time `json:"memory_date" binding:"required"`
	SpaceID    string    `json:"space_id" binding:"required"`
	GenerateAI bool      `json:"generate_ai"`
}

// AI
type GenerateMemoryRequest struct {
	Text     string `json:"text" binding:"required,min=1"`
	Style    string `json:"style" binding:"oneof=warm poetic cinematic"`
	MemoryID string `json:"memory_id"`
}

type GenerateMemoryResponse struct {
	Content    string `json:"content"`
	Model      string `json:"model"`
	TokensUsed int    `json:"tokens_used"`
}
