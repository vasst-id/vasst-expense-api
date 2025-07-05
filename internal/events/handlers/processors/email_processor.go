package processors

import (
	"fmt"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
)

// EmailProcessor handles email webhook payloads (e.g., from SendGrid, Mailgun)
type EmailProcessor struct{}

func (p *EmailProcessor) GetMediumID() int {
	return entities.MediumEmail
}

func (p *EmailProcessor) ValidatePayload(payload map[string]interface{}) error {
	if _, ok := payload["email"]; !ok {
		if _, ok := payload["from"]; !ok {
			return fmt.Errorf("invalid email payload: missing email or from field")
		}
	}
	return nil
}

func (p *EmailProcessor) ExtractMessages(payload map[string]interface{}) ([]*MessageInfo, error) {
	if err := p.ValidatePayload(payload); err != nil {
		return nil, err
	}

	// Handle different email webhook formats
	var messages []*MessageInfo

	// Check if it's a single email or array of emails
	if emailArray, ok := payload["emails"].([]interface{}); ok {
		// Multiple emails in array
		for _, email := range emailArray {
			if emailMap, ok := email.(map[string]interface{}); ok {
				messageInfo := p.extractSingleEmail(emailMap)
				if messageInfo != nil {
					messages = append(messages, messageInfo)
				}
			}
		}
	} else {
		// Single email
		messageInfo := p.extractSingleEmail(payload)
		if messageInfo != nil {
			messages = append(messages, messageInfo)
		}
	}

	return messages, nil
}

func (p *EmailProcessor) extractSingleEmail(email map[string]interface{}) *MessageInfo {
	var fromEmail string
	var content string
	var subject string
	metadata := make(map[string]interface{})

	// Extract sender email (try different field names)
	if emailAddr, ok := email["email"].(string); ok {
		fromEmail = emailAddr
	} else if from, ok := email["from"].(string); ok {
		fromEmail = from
	} else if sender, ok := email["sender"].(string); ok {
		fromEmail = sender
	} else {
		return nil // No sender found
	}

	// Extract subject
	if subj, ok := email["subject"].(string); ok {
		subject = subj
		metadata["subject"] = subject
	}

	// Extract content (try different field names)
	if text, ok := email["text"].(string); ok {
		content = text
	} else if body, ok := email["body"].(string); ok {
		content = body
	} else if message, ok := email["message"].(string); ok {
		content = message
	} else if html, ok := email["html"].(string); ok {
		content = html
		metadata["content_type"] = "html"
	}

	// Combine subject and content if both exist
	if subject != "" && content != "" {
		content = fmt.Sprintf("Subject: %s\n\n%s", subject, content)
	} else if subject != "" && content == "" {
		content = subject
	}

	// Add email-specific metadata
	if msgId, ok := email["message_id"].(string); ok {
		metadata["email_message_id"] = msgId
	}
	if timestamp, ok := email["timestamp"].(string); ok {
		metadata["timestamp"] = timestamp
	}
	if to, ok := email["to"].(string); ok {
		metadata["to"] = to
	}

	return &MessageInfo{
		PhoneNumber: fromEmail, // Use email as identifier
		Content:     content,
		MediaURL:    "", // Emails typically don't have direct media URLs
		MessageType: entities.MessageTypeText,
		Metadata:    metadata,
	}
}
