package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vasst-id/vasst-expense-api/config"
)

//go:generate mockgen -source=whatsapp_media_service.go -package=mock -destination=mock/whatsapp_media_service_mock.go
type (
	WhatsAppMediaService interface {
		GetMediaURL(ctx context.Context, mediaID string) (string, error)
		DownloadMedia(ctx context.Context, mediaURL string) ([]byte, string, error)
		GetMediaInfo(ctx context.Context, mediaID string) (*MediaInfo, error)
	}

	whatsAppMediaService struct {
		accessToken string
		client      *http.Client
	}

	MediaInfo struct {
		ID       string `json:"id"`
		MimeType string `json:"mime_type"`
		Sha256   string `json:"sha256"`
		FileSize int64  `json:"file_size"`
		URL      string `json:"url"`
	}

	WhatsAppMediaResponse struct {
		URL          string `json:"url"`
		MimeType     string `json:"mime_type"`
		Sha256       string `json:"sha256"`
		FileSize     int64  `json:"file_size"`
		ID           string `json:"id"`
		MessagingProduct string `json:"messaging_product"`
	}
)

// NewWhatsAppMediaService creates a new WhatsApp Media service
func NewWhatsAppMediaService(cfg *config.Config) WhatsAppMediaService {
	return &whatsAppMediaService{
		accessToken: cfg.WhatsAppAccessToken,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetMediaURL retrieves the download URL for a WhatsApp media file
func (s *whatsAppMediaService) GetMediaURL(ctx context.Context, mediaID string) (string, error) {
	if mediaID == "" {
		return "", fmt.Errorf("media ID cannot be empty")
	}

	url := fmt.Sprintf("https://graph.facebook.com/v21.0/%s", mediaID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get media URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("WhatsApp API error: %d - %s", resp.StatusCode, string(body))
	}

	var mediaResponse WhatsAppMediaResponse
	if err := json.NewDecoder(resp.Body).Decode(&mediaResponse); err != nil {
		return "", fmt.Errorf("failed to decode media response: %w", err)
	}

	return mediaResponse.URL, nil
}

// DownloadMedia downloads media content from WhatsApp's media URL
func (s *whatsAppMediaService) DownloadMedia(ctx context.Context, mediaURL string) ([]byte, string, error) {
	if mediaURL == "" {
		return nil, "", fmt.Errorf("media URL cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", mediaURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create download request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.accessToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download media: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("failed to download media: %d - %s", resp.StatusCode, string(body))
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read media data: %w", err)
	}

	return data, contentType, nil
}

// GetMediaInfo retrieves detailed information about a WhatsApp media file
func (s *whatsAppMediaService) GetMediaInfo(ctx context.Context, mediaID string) (*MediaInfo, error) {
	if mediaID == "" {
		return nil, fmt.Errorf("media ID cannot be empty")
	}

	url := fmt.Sprintf("https://graph.facebook.com/v21.0/%s", mediaID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get media info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("WhatsApp API error: %d - %s", resp.StatusCode, string(body))
	}

	var mediaResponse WhatsAppMediaResponse
	if err := json.NewDecoder(resp.Body).Decode(&mediaResponse); err != nil {
		return nil, fmt.Errorf("failed to decode media response: %w", err)
	}

	return &MediaInfo{
		ID:       mediaResponse.ID,
		MimeType: mediaResponse.MimeType,
		Sha256:   mediaResponse.Sha256,
		FileSize: mediaResponse.FileSize,
		URL:      mediaResponse.URL,
	}, nil
}