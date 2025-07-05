package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/vasst-id/vasst-expense-api/config"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=openai_service.go -package=mock -destination=mock/openai_service_mock.go
type (
	OpenAIService interface {
		ProcessCustomerMessage(ctx context.Context, message string) (string, error)
	}

	openAIService struct {
		apiKey         string
		baseURL        string
		httpClient     *http.Client
		messageService MessageService
	}

	// OpenAIRequest represents the request body for OpenAI API
	OpenAIRequest struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
	}

	// Message represents a message in the OpenAI conversation
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	// OpenAIResponse represents the response from OpenAI API
	OpenAIResponse struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
)

// NewOpenAIService creates a new OpenAI service
func NewOpenAIService(cfg *config.Config, messageService MessageService) (OpenAIService, error) {
	apiKey := cfg.OpenAIApiKey
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY environment variable is required")
	}

	return &openAIService{
		apiKey:         apiKey,
		baseURL:        "https://api.openai.com/v1/chat/completions",
		httpClient:     &http.Client{},
		messageService: messageService,
	}, nil
}

// ProcessCustomerMessage processes a customer message and generates a response
func (s *openAIService) ProcessCustomerMessage(ctx context.Context, message string) (string, error) {
	// Load system prompt from docs/tge_agent_prompt.md
	systemPromptPath := filepath.Join("internal", "docs", "tge_agent_prompt.md")
	systemPromptBytes, err := os.ReadFile(systemPromptPath)
	if err != nil {
		return "", fmt.Errorf("failed to read system prompt: %w", err)
	}
	systemPrompt := string(systemPromptBytes)

	// Load knowledge base from docs/knowledge.md
	knowledgePath := filepath.Join("internal", "docs", "knowledge.md")
	knowledgeBytes, err := os.ReadFile(knowledgePath)
	if err != nil {
		return "", fmt.Errorf("failed to read knowledge base: %w", err)
	}
	knowledge := string(knowledgeBytes)

	// Combine system prompt and knowledge base
	fullSystemPrompt := systemPrompt + "\n\n# Knowledge Base\n" + knowledge

	// Create the conversation messages
	messages := []Message{
		{
			Role:    "system",
			Content: fullSystemPrompt,
		},
		{
			Role:    "user",
			Content: message,
		},
	}

	// Create OpenAI request
	request := OpenAIRequest{
		Model:    "gpt-4.1-nano",
		Messages: messages,
	}

	// Send request to OpenAI
	response, err := s.sendToOpenAI(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to get OpenAI response: %w", err)
	}

	// Extract the assistant's response
	if len(response.Choices) == 0 {
		return "", errors.New("no response from OpenAI")
	}
	assistantResponse := response.Choices[0].Message.Content

	return assistantResponse, nil
}

// sendToOpenAI sends a request to the OpenAI API
func (s *openAIService) sendToOpenAI(ctx context.Context, request OpenAIRequest) (*OpenAIResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errorsutil.New(resp.StatusCode, fmt.Sprintf("OpenAI API error: %s, response: %s", resp.Status, string(body)))
	}

	var response OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
