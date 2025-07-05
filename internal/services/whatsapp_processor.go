package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
)

//go:generate mockgen -source=whatsapp_processor.go -package=mock -destination=mock/whatsapp_processor_mock.go
type (
	WhatsAppProcessor interface {
		ProcessWebhookPayload(payload map[string]interface{}) ([]entities.WebhookMessage, error)
	}

	whatsAppProcessor struct{}
)

// NewWhatsAppProcessor creates a new WhatsApp processor
func NewWhatsAppProcessor() WhatsAppProcessor {
	return &whatsAppProcessor{}
}

// ProcessWebhookPayload processes WhatsApp webhook payload and extracts messages
func (p *whatsAppProcessor) ProcessWebhookPayload(payload map[string]interface{}) ([]entities.WebhookMessage, error) {
	var messages []entities.WebhookMessage

	entry, exists := payload["entry"]
	if !exists {
		return messages, fmt.Errorf("no entry field in payload")
	}

	entrySlice, ok := entry.([]interface{})
	if !ok {
		return messages, fmt.Errorf("entry field is not an array")
	}

	for _, entryItem := range entrySlice {
		entryMap, ok := entryItem.(map[string]interface{})
		if !ok {
			continue
		}

		changes, exists := entryMap["changes"]
		if !exists {
			continue
		}

		changesSlice, ok := changes.([]interface{})
		if !ok {
			continue
		}

		for _, change := range changesSlice {
			changeMap, ok := change.(map[string]interface{})
			if !ok {
				continue
			}

			value, exists := changeMap["value"]
			if !exists {
				continue
			}

			valueMap, ok := value.(map[string]interface{})
			if !ok {
				continue
			}

			messagesData, exists := valueMap["messages"]
			if !exists {
				continue
			}

			messagesSlice, ok := messagesData.([]interface{})
			if !ok {
				continue
			}

			for _, msgData := range messagesSlice {
				msgMap, ok := msgData.(map[string]interface{})
				if !ok {
					continue
				}

				message, err := p.extractMessage(msgMap)
				if err != nil {
					continue
				}

				messages = append(messages, message)
			}
		}
	}

	return messages, nil
}

// extractMessage extracts a WebhookMessage from WhatsApp message data
func (p *whatsAppProcessor) extractMessage(msgMap map[string]interface{}) (entities.WebhookMessage, error) {
	var message entities.WebhookMessage

	// Extract phone number (from field)
	if from, exists := msgMap["from"]; exists {
		if fromStr, ok := from.(string); ok {
			message.PhoneNumber = fromStr
		}
	}

	// Extract WhatsApp message ID (id field)
	if id, exists := msgMap["id"]; exists {
		if idStr, ok := id.(string); ok {
			message.WhatsAppMessageID = idStr
		}
	}

	// Extract message type
	if msgType, exists := msgMap["type"]; exists {
		if msgTypeStr, ok := msgType.(string); ok {
			message.MessageType = p.getMessageTypeID(msgTypeStr)
		}
	}

	// Extract message content based on type
	switch message.MessageType {
	case entities.MessageTypeText:
		if err := p.extractTextMessage(msgMap, &message); err != nil {
			return message, err
		}
	case entities.MessageTypeImage:
		if err := p.extractImageMessage(msgMap, &message); err != nil {
			return message, err
		}
	case entities.MessageTypeVideo:
		if err := p.extractVideoMessage(msgMap, &message); err != nil {
			return message, err
		}
	case entities.MessageTypeAudio:
		if err := p.extractAudioMessage(msgMap, &message); err != nil {
			return message, err
		}
	case entities.MessageTypeDocument:
		if err := p.extractDocumentMessage(msgMap, &message); err != nil {
			return message, err
		}
	case entities.MessageTypeLocation:
		if err := p.extractLocationMessage(msgMap, &message); err != nil {
			return message, err
		}
	case entities.MessageTypeContact:
		if err := p.extractContactMessage(msgMap, &message); err != nil {
			return message, err
		}
	case entities.MessageTypeSticker:
		if err := p.extractStickerMessage(msgMap, &message); err != nil {
			return message, err
		}
	default:
		return message, fmt.Errorf("unsupported message type: %d", message.MessageType)
	}

	return message, nil
}

// extractTextMessage extracts text message content
func (p *whatsAppProcessor) extractTextMessage(msgMap map[string]interface{}, message *entities.WebhookMessage) error {
	if text, exists := msgMap["text"]; exists {
		if textMap, ok := text.(map[string]interface{}); ok {
			if body, exists := textMap["body"]; exists {
				if bodyStr, ok := body.(string); ok {
					message.Content = bodyStr
				}
			}
		}
	}
	return nil
}

// extractImageMessage extracts image message content
func (p *whatsAppProcessor) extractImageMessage(msgMap map[string]interface{}, message *entities.WebhookMessage) error {
	if image, exists := msgMap["image"]; exists {
		if imageMap, ok := image.(map[string]interface{}); ok {
			// Extract media ID
			if id, exists := imageMap["id"]; exists {
				if idStr, ok := id.(string); ok {
					message.MediaID = idStr
				}
			}

			// Extract caption as content
			if caption, exists := imageMap["caption"]; exists {
				if captionStr, ok := caption.(string); ok {
					message.Content = captionStr
				}
			}

			// Extract metadata
			metadata := make(map[string]interface{})
			if mimeType, exists := imageMap["mime_type"]; exists {
				metadata["mime_type"] = mimeType
			}
			if sha256, exists := imageMap["sha256"]; exists {
				metadata["sha256"] = sha256
			}
			message.Metadata = metadata
		}
	}
	return nil
}

// extractVideoMessage extracts video message content
func (p *whatsAppProcessor) extractVideoMessage(msgMap map[string]interface{}, message *entities.WebhookMessage) error {
	if video, exists := msgMap["video"]; exists {
		if videoMap, ok := video.(map[string]interface{}); ok {
			// Extract media ID
			if id, exists := videoMap["id"]; exists {
				if idStr, ok := id.(string); ok {
					message.MediaID = idStr
				}
			}

			// Extract caption as content
			if caption, exists := videoMap["caption"]; exists {
				if captionStr, ok := caption.(string); ok {
					message.Content = captionStr
				}
			}

			// Extract metadata
			metadata := make(map[string]interface{})
			if mimeType, exists := videoMap["mime_type"]; exists {
				metadata["mime_type"] = mimeType
			}
			if sha256, exists := videoMap["sha256"]; exists {
				metadata["sha256"] = sha256
			}
			message.Metadata = metadata
		}
	}
	return nil
}

// extractAudioMessage extracts audio message content
func (p *whatsAppProcessor) extractAudioMessage(msgMap map[string]interface{}, message *entities.WebhookMessage) error {
	if audio, exists := msgMap["audio"]; exists {
		if audioMap, ok := audio.(map[string]interface{}); ok {
			// Extract media ID
			if id, exists := audioMap["id"]; exists {
				if idStr, ok := id.(string); ok {
					message.MediaID = idStr
				}
			}

			// For audio, content is typically empty or voice note description
			message.Content = "Voice message"

			// Extract metadata
			metadata := make(map[string]interface{})
			if mimeType, exists := audioMap["mime_type"]; exists {
				metadata["mime_type"] = mimeType
			}
			if sha256, exists := audioMap["sha256"]; exists {
				metadata["sha256"] = sha256
			}
			if voice, exists := audioMap["voice"]; exists {
				metadata["voice"] = voice
			}
			message.Metadata = metadata
		}
	}
	return nil
}

// extractDocumentMessage extracts document message content
func (p *whatsAppProcessor) extractDocumentMessage(msgMap map[string]interface{}, message *entities.WebhookMessage) error {
	if document, exists := msgMap["document"]; exists {
		if docMap, ok := document.(map[string]interface{}); ok {
			// Extract media ID
			if id, exists := docMap["id"]; exists {
				if idStr, ok := id.(string); ok {
					message.MediaID = idStr
				}
			}

			// Extract filename as content
			if filename, exists := docMap["filename"]; exists {
				if filenameStr, ok := filename.(string); ok {
					message.Content = fmt.Sprintf("Document: %s", filenameStr)
				}
			}

			// Extract caption if available
			if caption, exists := docMap["caption"]; exists {
				if captionStr, ok := caption.(string); ok {
					if message.Content != "" {
						message.Content += " - " + captionStr
					} else {
						message.Content = captionStr
					}
				}
			}

			// Extract metadata
			metadata := make(map[string]interface{})
			if mimeType, exists := docMap["mime_type"]; exists {
				metadata["mime_type"] = mimeType
			}
			if sha256, exists := docMap["sha256"]; exists {
				metadata["sha256"] = sha256
			}
			if filename, exists := docMap["filename"]; exists {
				metadata["filename"] = filename
			}
			message.Metadata = metadata
		}
	}
	return nil
}

// extractLocationMessage extracts location message content
func (p *whatsAppProcessor) extractLocationMessage(msgMap map[string]interface{}, message *entities.WebhookMessage) error {
	if location, exists := msgMap["location"]; exists {
		if locMap, ok := location.(map[string]interface{}); ok {
			var lat, lng float64
			var name, address string

			if latitude, exists := locMap["latitude"]; exists {
				if latFloat, ok := latitude.(float64); ok {
					lat = latFloat
				} else if latStr, ok := latitude.(string); ok {
					if parsedLat, err := strconv.ParseFloat(latStr, 64); err == nil {
						lat = parsedLat
					}
				}
			}

			if longitude, exists := locMap["longitude"]; exists {
				if lngFloat, ok := longitude.(float64); ok {
					lng = lngFloat
				} else if lngStr, ok := longitude.(string); ok {
					if parsedLng, err := strconv.ParseFloat(lngStr, 64); err == nil {
						lng = parsedLng
					}
				}
			}

			if locName, exists := locMap["name"]; exists {
				if nameStr, ok := locName.(string); ok {
					name = nameStr
				}
			}

			if locAddress, exists := locMap["address"]; exists {
				if addressStr, ok := locAddress.(string); ok {
					address = addressStr
				}
			}

			// Create content description
			content := fmt.Sprintf("Location: %f, %f", lat, lng)
			if name != "" {
				content = fmt.Sprintf("Location: %s (%f, %f)", name, lat, lng)
			}
			if address != "" {
				content += fmt.Sprintf(" - %s", address)
			}
			message.Content = content

			// Store location data as metadata
			metadata := map[string]interface{}{
				"latitude":  lat,
				"longitude": lng,
			}
			if name != "" {
				metadata["name"] = name
			}
			if address != "" {
				metadata["address"] = address
			}
			message.Metadata = metadata
		}
	}
	return nil
}

// extractContactMessage extracts contact message content
func (p *whatsAppProcessor) extractContactMessage(msgMap map[string]interface{}, message *entities.WebhookMessage) error {
	if contacts, exists := msgMap["contacts"]; exists {
		if contactsSlice, ok := contacts.([]interface{}); ok && len(contactsSlice) > 0 {
			if contact, ok := contactsSlice[0].(map[string]interface{}); ok {
				var contactName string
				var phones []string

				// Extract name
				if name, exists := contact["name"]; exists {
					if nameMap, ok := name.(map[string]interface{}); ok {
						if formattedName, exists := nameMap["formatted_name"]; exists {
							if nameStr, ok := formattedName.(string); ok {
								contactName = nameStr
							}
						}
					}
				}

				// Extract phone numbers
				if phonesList, exists := contact["phones"]; exists {
					if phonesSlice, ok := phonesList.([]interface{}); ok {
						for _, phone := range phonesSlice {
							if phoneMap, ok := phone.(map[string]interface{}); ok {
								if phoneNum, exists := phoneMap["phone"]; exists {
									if phoneStr, ok := phoneNum.(string); ok {
										phones = append(phones, phoneStr)
									}
								}
							}
						}
					}
				}

				// Create content description
				content := fmt.Sprintf("Contact: %s", contactName)
				if len(phones) > 0 {
					content += fmt.Sprintf(" (%s)", strings.Join(phones, ", "))
				}
				message.Content = content

				// Store contact data as metadata
				metadata := map[string]interface{}{
					"contact_name": contactName,
					"phones":       phones,
				}
				message.Metadata = metadata
			}
		}
	}
	return nil
}

// extractStickerMessage extracts sticker message content
func (p *whatsAppProcessor) extractStickerMessage(msgMap map[string]interface{}, message *entities.WebhookMessage) error {
	if sticker, exists := msgMap["sticker"]; exists {
		if stickerMap, ok := sticker.(map[string]interface{}); ok {
			// Extract media ID
			if id, exists := stickerMap["id"]; exists {
				if idStr, ok := id.(string); ok {
					message.MediaID = idStr
				}
			}

			message.Content = "Sticker"

			// Extract metadata
			metadata := make(map[string]interface{})
			if mimeType, exists := stickerMap["mime_type"]; exists {
				metadata["mime_type"] = mimeType
			}
			if sha256, exists := stickerMap["sha256"]; exists {
				metadata["sha256"] = sha256
			}
			if animated, exists := stickerMap["animated"]; exists {
				metadata["animated"] = animated
			}
			message.Metadata = metadata
		}
	}
	return nil
}

// getMessageTypeID converts WhatsApp message type string to internal message type ID
func (p *whatsAppProcessor) getMessageTypeID(msgType string) int {
	switch strings.ToLower(msgType) {
	case "text":
		return entities.MessageTypeText
	case "image":
		return entities.MessageTypeImage
	case "video":
		return entities.MessageTypeVideo
	case "audio":
		return entities.MessageTypeAudio
	case "document":
		return entities.MessageTypeDocument
	case "location":
		return entities.MessageTypeLocation
	case "contacts":
		return entities.MessageTypeContact
	case "sticker":
		return entities.MessageTypeSticker
	default:
		return entities.MessageTypeText // Default to text for unknown types
	}
}