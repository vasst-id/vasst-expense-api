package processors

import (
	"fmt"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
)

// WhatsAppProcessor handles WhatsApp Business API webhook payloads
type WhatsAppProcessor struct{}

func (p *WhatsAppProcessor) GetMediumID() int {
	return entities.MediumWhatsApp
}

func (p *WhatsAppProcessor) ValidatePayload(payload map[string]interface{}) error {
	if len(payload) == 0 {
		return fmt.Errorf("empty WhatsApp payload")
	}

	entry, ok := payload["entry"]
	if !ok {
		return fmt.Errorf("invalid WhatsApp payload: missing 'entry' field")
	}

	entryArray, ok := entry.([]interface{})
	if !ok {
		return fmt.Errorf("invalid WhatsApp payload: 'entry' must be an array")
	}

	if len(entryArray) == 0 {
		return fmt.Errorf("invalid WhatsApp payload: 'entry' array is empty")
	}

	return nil
}

func (p *WhatsAppProcessor) ExtractMessages(payload map[string]interface{}) ([]*MessageInfo, error) {
	if err := p.ValidatePayload(payload); err != nil {
		return nil, err
	}

	var messages []*MessageInfo

	// Parse WhatsApp Business API webhook structure
	entry, _ := payload["entry"].([]interface{})

	for _, entryItem := range entry {
		entryMap, ok := entryItem.(map[string]interface{})
		if !ok {
			continue
		}

		changes, ok := entryMap["changes"].([]interface{})
		if !ok {
			continue
		}

		for _, change := range changes {
			changeMap, ok := change.(map[string]interface{})
			if !ok {
				continue
			}

			value, ok := changeMap["value"].(map[string]interface{})
			if !ok {
				continue
			}

			// Extract messages from value
			if msgs, exists := value["messages"]; exists {
				msgArray, ok := msgs.([]interface{})
				if !ok {
					continue
				}

				for _, msg := range msgArray {
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
		}
	}

	return messages, nil
}

func (p *WhatsAppProcessor) extractSingleMessage(message map[string]interface{}) *MessageInfo {
	phoneNumber, ok := message["from"].(string)
	originMessageID, ok := message["id"].(string)
	if !ok {
		return nil
	}

	messageType := entities.MessageTypeText
	content := ""
	mediaURL := ""
	metadata := make(map[string]interface{})

	// Extract content based on message type
	if textObj, exists := message["text"]; exists {
		if textMap, ok := textObj.(map[string]interface{}); ok {
			if body, ok := textMap["body"].(string); ok {
				content = body
			}
		}
	} else if imageObj, exists := message["image"]; exists {
		messageType = entities.MessageTypeImage
		if imageMap, ok := imageObj.(map[string]interface{}); ok {
			if id, ok := imageMap["id"].(string); ok {
				mediaURL = id // WhatsApp media ID
			}
			if caption, ok := imageMap["caption"].(string); ok {
				content = caption
			}
		}
	} else if videoObj, exists := message["video"]; exists {
		messageType = entities.MessageTypeVideo
		if videoMap, ok := videoObj.(map[string]interface{}); ok {
			if id, ok := videoMap["id"].(string); ok {
				mediaURL = id
			}
			if caption, ok := videoMap["caption"].(string); ok {
				content = caption
			}
		}
	} else if audioObj, exists := message["audio"]; exists {
		messageType = entities.MessageTypeAudio
		if audioMap, ok := audioObj.(map[string]interface{}); ok {
			if id, ok := audioMap["id"].(string); ok {
				mediaURL = id
			}
		}
	} else if docObj, exists := message["document"]; exists {
		messageType = entities.MessageTypeDocument
		if docMap, ok := docObj.(map[string]interface{}); ok {
			if id, ok := docMap["id"].(string); ok {
				mediaURL = id
			}
			if filename, ok := docMap["filename"].(string); ok {
				content = filename
			}
		}
	}

	// Add message metadata
	if msgId, ok := message["id"].(string); ok {
		fmt.Printf("DEBUG PROCESSOR: Extracted WhatsApp message ID: %s\n", msgId)
		metadata["whatsapp_message_id"] = msgId
	} else {
		fmt.Printf("DEBUG PROCESSOR: Failed to extract message ID from message: %+v\n", message)
	}
	if timestamp, ok := message["timestamp"].(string); ok {
		metadata["timestamp"] = timestamp
	}

	fmt.Printf("DEBUG PROCESSOR: Final metadata: %+v\n", metadata)

	return &MessageInfo{
		PhoneNumber:     phoneNumber,
		Content:         content,
		MediaURL:        mediaURL,
		MessageType:     messageType,
		Metadata:        metadata,
		MessageOriginId: originMessageID,
	}
}
