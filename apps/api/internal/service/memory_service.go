package service

import (
	"errors"

	"github.com/zhaozeguang/timecapsule-api/internal/dto"
	"github.com/zhaozeguang/timecapsule-api/internal/model"
	"github.com/zhaozeguang/timecapsule-api/internal/repository"
)

type MemoryService struct {
	memoryRepo *repository.MemoryRepo
}

func NewMemoryService(memoryRepo *repository.MemoryRepo) *MemoryService {
	return &MemoryService{memoryRepo: memoryRepo}
}

func (s *MemoryService) Create(userID string, req *dto.CreateMemoryRequest, imageURL, thumbnailURL string) (*model.Memory, error) {
	memory := &model.Memory{
		UserID:       userID,
		SpaceID:      req.SpaceID,
		ImageURL:     imageURL,
		ThumbnailURL: thumbnailURL,
		Title:        req.Title,
		Content:      req.Content,
		Location:     req.Location,
		MemoryDate:   req.MemoryDate,
	}

	if err := s.memoryRepo.Create(memory); err != nil {
		return nil, err
	}
	return memory, nil
}

func (s *MemoryService) GetByID(id, userID string) (*model.Memory, error) {
	memory, err := s.memoryRepo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}
	return memory, nil
}

func (s *MemoryService) List(spaceID, userID string, page, pageSize int) (*dto.PaginatedData, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	if spaceID == "" {
		return nil, errors.New("space_id is required")
	}

	memories, total, err := s.memoryRepo.FindBySpaceID(spaceID, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &dto.PaginatedData{
		Items:    memories,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *MemoryService) Update(id, userID, aiContent, aiModel, aiPrompt string) (*model.Memory, error) {
	memory, err := s.memoryRepo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}

	if aiContent != "" {
		memory.AIContent = aiContent
	}
	if aiModel != "" {
		memory.AIModel = aiModel
	}
	if aiPrompt != "" {
		memory.AIPrompt = aiPrompt
	}

	if err := s.memoryRepo.Update(memory); err != nil {
		return nil, err
	}
	return memory, nil
}

func (s *MemoryService) Delete(id, userID string) error {
	return s.memoryRepo.SoftDelete(id, userID)
}
