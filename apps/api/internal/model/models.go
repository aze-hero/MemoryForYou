package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email      string         `gorm:"uniqueIndex;not null" json:"email"`
	Name       string         `gorm:"not null" json:"name"`
	AvatarURL  string         `json:"avatar_url"`
	Provider   string         `gorm:"not null" json:"provider"`
	ProviderID string         `gorm:"not null" json:"provider_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type MemorySpace struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID      string         `gorm:"not null;index" json:"user_id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	CoverImage  string         `json:"cover_image"`
	Theme       string         `gorm:"default:default" json:"theme"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Memory struct {
	ID           string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID       string         `gorm:"not null;index" json:"user_id"`
	SpaceID      string         `gorm:"not null;index" json:"space_id"`
	ImageURL     string         `gorm:"not null" json:"image_url"`
	ThumbnailURL string         `gorm:"not null" json:"thumbnail_url"`
	Title        string         `gorm:"not null" json:"title"`
	Content      string         `json:"content"`
	AIContent    string         `json:"ai_content"`
	AIModel      string         `json:"ai_model"`
	AIPrompt     string         `json:"ai_prompt"`
	Location     string         `json:"location"`
	MemoryDate   time.Time      `gorm:"not null" json:"memory_date"`
	SortOrder    int            `gorm:"default:0" json:"sort_order"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
