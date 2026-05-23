package repository

import (
	"github.com/zhaozeguang/timecapsule-api/internal/model"
	"github.com/zhaozeguang/timecapsule-api/pkg/database"
)

type MemoryRepo struct{}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{}
}

func (r *MemoryRepo) Create(memory *model.Memory) error {
	return database.DB.Create(memory).Error
}

func (r *MemoryRepo) FindByID(id, userID string) (*model.Memory, error) {
	var memory model.Memory
	err := database.DB.First(&memory, "id = ? AND user_id = ?", id, userID).Error
	if err != nil {
		return nil, err
	}
	return &memory, nil
}

func (r *MemoryRepo) FindBySpaceID(spaceID, userID string, page, pageSize int) ([]model.Memory, int64, error) {
	var memories []model.Memory
	var total int64

	query := database.DB.Model(&model.Memory{}).Where("space_id = ? AND user_id = ?", spaceID, userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("memory_date DESC").Offset(offset).Limit(pageSize).Find(&memories).Error
	return memories, total, err
}

func (r *MemoryRepo) Update(memory *model.Memory) error {
	return database.DB.Save(memory).Error
}

func (r *MemoryRepo) SoftDelete(id, userID string) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Memory{}).Error
}

func (r *MemoryRepo) CountBySpaceID(spaceID string) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Memory{}).Where("space_id = ?", spaceID).Count(&count).Error
	return count, err
}
