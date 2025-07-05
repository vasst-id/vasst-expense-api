package processors

import (
	"fmt"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
)

// InstagramProcessor handles Instagram Graph API webhook payloads
type InstagramProcessor struct{}

func (p *InstagramProcessor) GetMediumID() int {
	return entities.MediumInstagram
}

func (p *InstagramProcessor) ValidatePayload(payload map[string]interface{}) error {
	entry, ok := payload["entry"].([]interface{})
	if !ok || len(entry) == 0 {
		return fmt.Errorf("invalid Instagram payload: missing entry")
	}
	return nil
}

func (p *InstagramProcessor) ExtractMessages(payload map[string]interface{}) ([]*MessageInfo, error) {
	if err := p.ValidatePayload(payload); err != nil {
		return nil, err
	}

	var messages []*MessageInfo

	// Parse Instagram Graph API webhook structure
	entry, _ := payload["entry"].([]interface{})
	
	for _, entryItem := range entry {
		entryMap, ok := entryItem.(map[string]interface{})
		if !ok {
			continue
		}

		messaging, ok := entryMap["messaging"].([]interface{})
		if !ok {
			continue
		}

		for _, msg := range messaging {
			msgMap, ok := msg.(map[string]interface{})
			if !ok {
				continue
			}

			messageInfo := p.extractSingleMessage(msgMap)
			if messageInfo != nil {
				messages = append(messages, messageInfo)
			}
		}
	}

	return messages, nil
}

func (p *InstagramProcessor) extractSingleMessage(messaging map[string]interface{}) *MessageInfo {
	sender, ok := messaging["sender"].(map[string]interface{})
	if !ok {
		return nil
	}

	senderID, ok := sender["id"].(string)
	if !ok {
		return nil
	}

	message, exists := messaging["message"]
	if !exists {
		return nil
	}

	messageMap, ok := message.(map[string]interface{})
	if !ok {
		return nil
	}

	messageType := entities.MessageTypeText
	content := ""
	mediaURL := ""
	metadata := make(map[string]interface{})

	// Extract content based on message type
	if text, exists := messageMap["text"]; exists {
		if textStr, ok := text.(string); ok {
			content = textStr
		}
	} else if attachments, exists := messageMap["attachments"]; exists {
		if attachArray, ok := attachments.([]interface{}); ok && len(attachArray) > 0 {
			if attach, ok := attachArray[0].(map[string]interface{}); ok {
				if attachType, ok := attach["type"].(string); ok {
					switch attachType {
					case "image":
						messageType = entities.MessageTypeImage
					case "video":
						messageType = entities.MessageTypeVideo
					case "audio":
						messageType = entities.MessageTypeAudio
					}
				}
				if payload, ok := attach["payload"].(map[string]interface{}); ok {
					if url, ok := payload["url"].(string); ok {
						mediaURL = url
					}
				}
			}
		}
	}

	// Add Instagram-specific metadata
	metadata["instagram_sender_id"] = senderID
	if msgId, ok := messageMap["mid"].(string); ok {
		metadata["instagram_message_id"] = msgId
	}

	return &MessageInfo{
		PhoneNumber: senderID, // Instagram uses sender ID instead of phone
		Content:     content,
		MediaURL:    mediaURL,
		MessageType: messageType,
		Metadata:    metadata,
	}
}
