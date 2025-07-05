package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vasst-id/vasst-expense-api/config"
	"google.golang.org/genai"
)

//go:generate mockgen -source=gemini_service.go -package=mock -destination=mock/gemini_service_mock.go
type (
	GeminiService interface {
		ProcessCustomerMessage(ctx context.Context, message string) (string, error)
		ProcessCustomerMessageWithPrompt(ctx context.Context, message, systemPrompt string) (string, error)
		ProcessMediaMessage(ctx context.Context, message string, messageType int, mediaURL string, systemPrompt string) (string, error)
	}

	geminiService struct {
		apiKey         string
		messageService MessageService
		client         *genai.Client
	}
)

// NewGeminiService creates a new Gemini service
func NewGeminiService(cfg *config.Config, messageService MessageService) (GeminiService, error) {
	apiKey := cfg.GeminiApiKey
	if apiKey == "" {
		return nil, errors.New("GEMINI_API_KEY environment variable is required")
	}

	// Set the API key as environment variable for the genai client
	os.Setenv("GEMINI_API_KEY", apiKey)

	// Create the Gemini client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &geminiService{
		apiKey:         apiKey,
		messageService: messageService,
		client:         client,
	}, nil
}

// ProcessCustomerMessage processes a customer message and generates a response
func (s *geminiService) ProcessCustomerMessage(ctx context.Context, message string) (string, error) {
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

	// Create the complete prompt with system instructions and user message
	completePrompt := fmt.Sprintf("%s\n\nUser: %s\nAssistant:", fullSystemPrompt, message)

	// Call Gemini API to generate response
	result, err := s.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(completePrompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content with Gemini: %w", err)
	}

	assistantResponse := result.Text()

	return assistantResponse, nil
}

// ProcessCustomerMessageWithPrompt processes a customer message with a custom system prompt
func (s *geminiService) ProcessCustomerMessageWithPrompt(ctx context.Context, message, systemPrompt string) (string, error) {
	if message == "" {
		return "", errors.New("message cannot be empty")
	}

	if systemPrompt == "" {
		return "", errors.New("system prompt cannot be empty")
	}

	// Create the complete prompt with system instructions and user message
	completePrompt := fmt.Sprintf("%s\n\nUser: %s\nAssistant:", systemPrompt, message)

	// Call Gemini API to generate response
	result, err := s.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(completePrompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content with Gemini: %w", err)
	}

	assistantResponse := result.Text()

	return assistantResponse, nil
}

// ProcessMediaMessage processes a customer message with media content
func (s *geminiService) ProcessMediaMessage(ctx context.Context, message string, messageType int, mediaURL string, systemPrompt string) (string, error) {
	if systemPrompt == "" {
		return "", errors.New("system prompt cannot be empty")
	}

	// Create media-aware prompt based on message type
	var mediaPrompt string
	switch messageType {
	case 2: // MessageTypeImage
		if mediaURL != "" {
			mediaPrompt = fmt.Sprintf("The user has shared an image (URL: %s). %s", mediaURL, message)
		} else {
			mediaPrompt = fmt.Sprintf("The user has shared an image. %s", message)
		}
	case 3: // MessageTypeVideo
		if mediaURL != "" {
			mediaPrompt = fmt.Sprintf("The user has shared a video (URL: %s). %s", mediaURL, message)
		} else {
			mediaPrompt = fmt.Sprintf("The user has shared a video. %s", message)
		}
	case 4: // MessageTypeAudio
		if mediaURL != "" {
			mediaPrompt = fmt.Sprintf("The user has shared an audio message (URL: %s). %s", mediaURL, message)
		} else {
			mediaPrompt = fmt.Sprintf("The user has shared an audio message. %s", message)
		}
	case 5: // MessageTypeDocument
		if mediaURL != "" {
			mediaPrompt = fmt.Sprintf("The user has shared a document (URL: %s). %s", mediaURL, message)
		} else {
			mediaPrompt = fmt.Sprintf("The user has shared a document. %s", message)
		}
	case 6: // MessageTypeLocation
		mediaPrompt = fmt.Sprintf("The user has shared their location. %s", message)
	case 7: // MessageTypeContact
		mediaPrompt = fmt.Sprintf("The user has shared a contact. %s", message)
	case 8: // MessageTypeSticker
		if mediaURL != "" {
			mediaPrompt = fmt.Sprintf("The user has sent a sticker (URL: %s). %s", mediaURL, message)
		} else {
			mediaPrompt = fmt.Sprintf("The user has sent a sticker. %s", message)
		}
	default:
		mediaPrompt = message
	}

	// Add context about how to handle media
	mediaContext := `

When responding to media messages:
- For videos: Acknowledge the video content and respond to any context provided
- For audio: Acknowledge the voice message and respond to the content
- For documents: Acknowledge the document and ask if they need help with it
- For locations: Acknowledge the location sharing and respond appropriately
- For contacts: Acknowledge the contact sharing and respond appropriately
- For stickers: Respond naturally to the sticker context

Always be helpful and acknowledge the media type shared.`

	// Create the complete prompt with media context
	completePrompt := fmt.Sprintf("%s%s\n\nUser: %s\nAssistant:", systemPrompt, mediaContext, mediaPrompt)

	// Call Gemini API to generate response
	result, err := s.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(completePrompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content with Gemini: %w", err)
	}

	assistantResponse := result.Text()

	return assistantResponse, nil
}
