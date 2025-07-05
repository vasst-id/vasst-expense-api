package entities

import (
	"time"

	"github.com/google/uuid"
)

// WhatsAppWebhook represents the incoming webhook from WhatsApp
type WhatsAppWebhook struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []WhatsAppMessage `json:"messages"`
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
	} `json:"entry"`
}

// WhatsAppMessage represents a message from WhatsApp webhook
type WhatsAppMessage struct {
	From      string `json:"from"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Text      *struct {
		Body string `json:"body"`
	} `json:"text,omitempty"`
	Image *struct {
		ID      string `json:"id"`
		Caption string `json:"caption"`
	} `json:"image,omitempty"`
	Video *struct {
		ID      string `json:"id"`
		Caption string `json:"caption"`
	} `json:"video,omitempty"`
	Audio *struct {
		ID string `json:"id"`
	} `json:"audio,omitempty"`
	Document *struct {
		ID       string `json:"id"`
		Caption  string `json:"caption"`
		Filename string `json:"filename"`
	} `json:"document,omitempty"`
	Sticker *struct {
		ID string `json:"id"`
	} `json:"sticker,omitempty"`
}

// WhatsAppTemplate represents a template message to be sent via WhatsApp
type WhatsAppTemplate struct {
	Name     string `json:"name"`
	Language struct {
		Code string `json:"code"`
	} `json:"language"`
}

// WhatsAppMessageRequest represents the request body for sending a WhatsApp message
type WhatsAppMessageRequest struct {
	MessagingProduct string            `json:"messaging_product"`
	To               string            `json:"to"`
	Type             string            `json:"type"`
	Template         WhatsAppTemplate  `json:"template,omitempty"`
	Text             *WhatsAppTextBody `json:"text,omitempty"`
}

// WhatsAppTextBody represents the text body of a WhatsApp message
type WhatsAppTextBody struct {
	Body string `json:"body"`
}

// WhatsAppOutgoingMessage represents a message to be sent via WhatsApp
type WhatsAppOutgoingMessage struct {
	ID        uuid.UUID
	To        string
	Type      string
	Content   string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewWhatsAppMessage creates a new WhatsApp outgoing message
func NewWhatsAppMessage(to, messageType, content string) *WhatsAppOutgoingMessage {
	return &WhatsAppOutgoingMessage{
		ID:        uuid.New(),
		To:        to,
		Type:      messageType,
		Content:   content,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
