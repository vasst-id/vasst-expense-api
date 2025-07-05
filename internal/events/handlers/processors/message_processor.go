package processors

// Package processors provides platform-specific message processing

// MessageInfo represents extracted message information from any platform
type MessageInfo struct {
	PhoneNumber     string
	Content         string
	MediaURL        string
	MessageType     int
	Metadata        map[string]interface{}
	MessageOriginId string
}

// MessageProcessor defines the interface for platform-specific message processing
type MessageProcessor interface {
	ExtractMessages(payload map[string]interface{}) ([]*MessageInfo, error)
	ValidatePayload(payload map[string]interface{}) error
	GetMediumID() int
}

// GetMessageProcessor returns the appropriate processor for the platform
func GetMessageProcessor(platform string) MessageProcessor {
	switch platform {
	case "whatsapp":
		return &WhatsAppProcessor{}
	case "instagram":
		return &InstagramProcessor{}
	case "facebook":
		return &FacebookProcessor{}
	case "email":
		return &EmailProcessor{}
	default:
		return nil
	}
}
