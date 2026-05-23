package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/zhaozeguang/timecapsule-api/internal/config"
	"github.com/zhaozeguang/timecapsule-api/internal/dto"
	"github.com/zhaozeguang/timecapsule-api/internal/repository"
)

type AIService struct {
	memoryRepo *repository.MemoryRepo
	cfg        *config.Config
	httpClient *http.Client
}

func NewAIService(memoryRepo *repository.MemoryRepo, cfg *config.Config) *AIService {
	return &AIService{
		memoryRepo: memoryRepo,
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *AIService) GenerateMemory(req *dto.GenerateMemoryRequest) (*dto.GenerateMemoryResponse, error) {
	return s.callOpenAI(req.Text, req.Style)
}

func (s *AIService) RegenerateForMemory(memoryID, userID string, style string) (*dto.GenerateMemoryResponse, error) {
	memory, err := s.memoryRepo.FindByID(memoryID, userID)
	if err != nil {
		return nil, errors.New("memory not found")
	}

	result, err := s.callOpenAI(memory.Title, style)
	if err != nil {
		return nil, err
	}

	memory.AIContent = result.Content
	memory.AIModel = result.Model
	memory.AIPrompt = memory.Title
	_ = s.memoryRepo.Update(memory)

	return result, nil
}

func (s *AIService) callOpenAI(text, style string) (*dto.GenerateMemoryResponse, error) {
	if s.cfg.OpenAIKey == "" {
		return nil, errors.New("OPENAI_API_KEY not configured")
	}

	systemPrompt := `You are an emotional memory storyteller.

Your task is to transform short memory descriptions into warm, cinematic, emotionally resonant reflections.

Style:
- poetic
- concise
- warm
- realistic
- cinematic

Avoid:
- exaggerated romance
- cringe internet language
- overly dramatic wording`

	if style != "" {
		systemPrompt += "\n\nPrefer " + style + " tone."
	}

	body := map[string]interface{}{
		"model": s.cfg.OpenAIModel,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": "Memory:\n" + text},
		},
		"temperature":    0.8,
		"max_tokens":     150,
		"top_p":          0.9,
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.OpenAIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Error != nil {
		return nil, errors.New("openai error: " + result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return nil, errors.New("no response from AI")
	}

	return &dto.GenerateMemoryResponse{
		Content:    result.Choices[0].Message.Content,
		Model:      s.cfg.OpenAIModel,
		TokensUsed: result.Usage.TotalTokens,
	}, nil
}
