package repository

import (
	"github.com/zhaozeguang/timecapsule-api/internal/model"
	"github.com/zhaozeguang/timecapsule-api/pkg/database"
)

type SpaceRepo struct{}

func NewSpaceRepo() *SpaceRepo {
	return &SpaceRepo{}
}

func (r *SpaceRepo) Create(space *model.MemorySpace) error {
	return database.DB.Create(space).Error
}

func (r *SpaceRepo) FindByUserID(userID string, page, pageSize int) ([]model.MemorySpace, int64, error) {
	var spaces []model.MemorySpace
	var total int64

	query := database.DB.Model(&model.MemorySpace{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&spaces).Error
	return spaces, total, err
}

func (r *SpaceRepo) FindByID(id, userID string) (*model.MemorySpace, error) {
	var space model.MemorySpace
	err := database.DB.First(&space, "id = ? AND user_id = ?", id, userID).Error
	if err != nil {
		return nil, err
	}
	return &space, nil
}

func (r *SpaceRepo) Update(space *model.MemorySpace) error {
	return database.DB.Save(space).Error
}

func (r *SpaceRepo) SoftDelete(id, userID string) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.MemorySpace{}).Error
}
