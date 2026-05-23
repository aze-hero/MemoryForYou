package service

import (
	"github.com/zhaozeguang/timecapsule-api/internal/dto"
	"github.com/zhaozeguang/timecapsule-api/internal/model"
	"github.com/zhaozeguang/timecapsule-api/internal/repository"
)

type SpaceService struct {
	spaceRepo  *repository.SpaceRepo
	memoryRepo *repository.MemoryRepo
}

func NewSpaceService(spaceRepo *repository.SpaceRepo, memoryRepo *repository.MemoryRepo) *SpaceService {
	return &SpaceService{spaceRepo: spaceRepo, memoryRepo: memoryRepo}
}

func (s *SpaceService) Create(userID string, req *dto.CreateSpaceRequest) (*model.MemorySpace, error) {
	space := &model.MemorySpace{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		CoverImage:  req.CoverImage,
		Theme:       req.Theme,
	}
	if space.Theme == "" {
		space.Theme = "default"
	}
	if err := s.spaceRepo.Create(space); err != nil {
		return nil, err
	}
	return space, nil
}

func (s *SpaceService) List(userID string, page, pageSize int) (*dto.PaginatedData, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	spaces, total, err := s.spaceRepo.FindByUserID(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, len(spaces))
	for i, space := range spaces {
		count, _ := s.memoryRepo.CountBySpaceID(space.ID)
		items[i] = map[string]interface{}{
			"id":           space.ID,
			"user_id":      space.UserID,
			"title":        space.Title,
			"description":  space.Description,
			"cover_image":  space.CoverImage,
			"theme":        space.Theme,
			"memory_count": count,
			"created_at":   space.CreatedAt,
			"updated_at":   space.UpdatedAt,
		}
	}

	return &dto.PaginatedData{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *SpaceService) GetByID(id, userID string) (*model.MemorySpace, error) {
	return s.spaceRepo.FindByID(id, userID)
}

func (s *SpaceService) Update(id, userID string, req *dto.CreateSpaceRequest) (*model.MemorySpace, error) {
	space, err := s.spaceRepo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		space.Title = req.Title
	}
	if req.Description != "" {
		space.Description = req.Description
	}
	if req.CoverImage != "" {
		space.CoverImage = req.CoverImage
	}
	if req.Theme != "" {
		space.Theme = req.Theme
	}

	if err := s.spaceRepo.Update(space); err != nil {
		return nil, err
	}
	return space, nil
}

func (s *SpaceService) Delete(id, userID string) error {
	return s.spaceRepo.SoftDelete(id, userID)
}
