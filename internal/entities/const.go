package entities

import "errors"

const (
	StatusActive   = 1
	StatusInactive = 0
	StatusDeleted  = 2
)

// ErrInvalidInput returns an error indicating invalid input
func ErrInvalidInput(msg string) error {
	return errors.New(msg)
}

var ErrNotFound = errors.New("record not found")
var SqlNoRows = errors.New("no rows in result set")

const (
	RoleAdmin      = 1
	RoleUser       = 2
	RoleCustomer   = 3
	RoleSuperAdmin = 99
)

const (
	PlanStarter    = 1
	PlanBusiness   = 2
	PlanPro        = 3
	PlanEnterprise = 4
)

const (
	MediumWhatsApp   = 1
	MediumEmail      = 2
	MediumSMS        = 3
	MediumCustomChat = 4
	MediumInstagram  = 5
	MediumFacebook   = 6
)

const (
	PlatformWhatsApp      = "whatsapp"
	PlatformInstagram     = "instagram"
	PlatformFacebook      = "facebook"
	PlatformEmail         = "email"
	PlatformVasstOrder    = "vasst-order"
	PlatformVasstSchedule = "vasst-schedule"
)

const (
	WebhookStatusPending   = 0
	WebhookStatusProcessed = 1
	WebhookStatusFailed    = 2
)

const (
	MessageTypeText     int = 1
	MessageTypeImage    int = 2
	MessageTypeVideo    int = 3
	MessageTypeAudio    int = 4
	MessageTypeDocument int = 5
	MessageTypeLocation int = 6
	MessageTypeContact  int = 7
	MessageTypeSticker  int = 8
)

// Rate limiting constants
const (
	RateLimitRequestsPerMinute = 100
	RateLimitBurstSize         = 10
)

// Google Pub/Sub Topics
const (
	TopicWebhookReceived    = "webhook-received"
	TopicMessageCreated     = "message-created"
	TopicAIResponseReceived = "ai-response-received"
	TopicMessageDelivery    = "message-delivery"
)

// Google Pub/Sub Subscriptions
const (
	SubscriptionWebhookProcessor    = "webhook-processor-sub"
	SubscriptionMessageProcessor    = "message-processor-sub"
	SubscriptionAIProcessor         = "ai-processor-sub"
	SubscriptionAIResponseProcessor = "ai-response-processor-sub"
	SubscriptionMessageDelivery     = "message-delivery-sub"
)
