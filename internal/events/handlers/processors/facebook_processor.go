package processors

import (
	"fmt"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
)

// FacebookProcessor handles Facebook Messenger webhook payloads
type FacebookProcessor struct{}

func (p *FacebookProcessor) GetMediumID() int {
	return entities.MediumFacebook
}

func (p *FacebookProcessor) ValidatePayload(payload map[string]interface{}) error {
	entry, ok := payload["entry"].([]interface{})
	if !ok || len(entry) == 0 {
		return fmt.Errorf("invalid Facebook payload: missing entry")
	}
	return nil
}

func (p *FacebookProcessor) ExtractMessages(payload map[string]interface{}) ([]*MessageInfo, error) {
	if err := p.ValidatePayload(payload); err != nil {
		return nil, err
	}

	var messages []*MessageInfo

	// Parse Facebook Messenger webhook structure (similar to Instagram)
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

func (p *FacebookProcessor) extractSingleMessage(messaging map[string]interface{}) *MessageInfo {
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
					case "file":
						messageType = entities.MessageTypeDocument
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

	// Add Facebook-specific metadata
	metadata["facebook_sender_id"] = senderID
	if msgId, ok := messageMap["mid"].(string); ok {
		metadata["facebook_message_id"] = msgId
	}
	if timestamp, ok := messaging["timestamp"].(float64); ok {
		metadata["timestamp"] = timestamp
	}

	return &MessageInfo{
		PhoneNumber: senderID, // Facebook uses sender ID instead of phone
		Content:     content,
		MediaURL:    mediaURL,
		MessageType: messageType,
		Metadata:    metadata,
	}
}
