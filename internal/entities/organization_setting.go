package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// OrganizationSetting represents organization settings and limits
type OrganizationSetting struct {
	OrganizationID            uuid.UUID       `json:"organization_id" db:"organization_id"`
	PlanID                    int             `json:"plan_id" db:"plan_id"`
	PlanName                  string          `json:"plan_name" db:"plan_name"`
	PlanStartDate             time.Time       `json:"plan_start_date" db:"plan_start_date"`
	PlanEndDate               time.Time       `json:"plan_end_date" db:"plan_end_date"`
	PlanStatus                int             `json:"plan_status" db:"plan_status"`
	PlanAmount                float64         `json:"plan_amount" db:"plan_amount"`
	PlanCurrency              string          `json:"plan_currency" db:"plan_currency"`
	MaxContacts               int             `json:"max_contacts" db:"max_contacts"`
	MaxMessages               int             `json:"max_messages" db:"max_messages"`
	MaxBroadcasts             int             `json:"max_broadcasts" db:"max_broadcasts"`
	MaxUsers                  int             `json:"max_users" db:"max_users"`
	MaxTags                   int             `json:"max_tags" db:"max_tags"`
	MaxOrders                 int             `json:"max_orders" db:"max_orders"`
	CurrentContacts           int             `json:"current_contacts" db:"current_contacts"`
	CurrentMessages           int             `json:"current_messages" db:"current_messages"`
	CurrentBroadcasts         int             `json:"current_broadcasts" db:"current_broadcasts"`
	CurrentUsers              int             `json:"current_users" db:"current_users"`
	CurrentTags               int             `json:"current_tags" db:"current_tags"`
	CurrentOrders             int             `json:"current_orders" db:"current_orders"`
	SystemPrompt              *string         `json:"system_prompt" db:"system_prompt"`
	AIAssistantName           *string         `json:"ai_assistant_name" db:"ai_assistant_name"`
	AICommunicationStyle      *string         `json:"ai_communication_style" db:"ai_communication_style"`
	AICommunicationLanguage   *string         `json:"ai_communication_language" db:"ai_communication_language"`
	OrganizationInfoStructure json.RawMessage `json:"organization_info_structure" db:"organization_info_structure"`
	ContactInfoStructure      json.RawMessage `json:"contact_info_structure" db:"contact_info_structure"`
	CreatedAt                 time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time       `json:"updated_at" db:"updated_at"`
}

// CreateOrganizationSettingInput is used for creating new organization settings
type CreateOrganizationSettingInput struct {
	OrganizationID            uuid.UUID       `json:"organization_id" binding:"required"`
	PlanID                    int             `json:"plan_id" binding:"required"`
	PlanStartDate             time.Time       `json:"plan_start_date" binding:"required"`
	PlanEndDate               time.Time       `json:"plan_end_date" binding:"required"`
	PlanStatus                int             `json:"plan_status" binding:"required"`
	PlanAmount                float64         `json:"plan_amount" binding:"required"`
	PlanCurrency              string          `json:"plan_currency" binding:"required"`
	MaxContacts               int             `json:"max_contacts"`
	MaxMessages               int             `json:"max_messages"`
	MaxBroadcasts             int             `json:"max_broadcasts"`
	MaxUsers                  int             `json:"max_users"`
	MaxTags                   int             `json:"max_tags"`
	MaxOrders                 int             `json:"max_orders"`
	SystemPrompt              *string         `json:"system_prompt" db:"system_prompt"`
	AIAssistantName           *string         `json:"ai_assistant_name" db:"ai_assistant_name"`
	AICommunicationStyle      *string         `json:"ai_communication_style" db:"ai_communication_style"`
	AICommunicationLanguage   *string         `json:"ai_communication_language" db:"ai_communication_language"`
	OrganizationInfoStructure json.RawMessage `json:"organization_info_structure" db:"organization_info_structure"`
	ContactInfoStructure      json.RawMessage `json:"contact_info_structure" db:"contact_info_structure"`
}

// UpdateOrganizationSettingInput is used for updating organization settings
type UpdateOrganizationSettingInput struct {
	PlanID                    int             `json:"plan_id"`
	PlanStartDate             time.Time       `json:"plan_start_date"`
	PlanEndDate               time.Time       `json:"plan_end_date"`
	PlanStatus                int             `json:"plan_status"`
	PlanAmount                float64         `json:"plan_amount"`
	PlanCurrency              string          `json:"plan_currency"`
	MaxContacts               int             `json:"max_contacts"`
	MaxMessages               int             `json:"max_messages"`
	MaxBroadcasts             int             `json:"max_broadcasts"`
	MaxUsers                  int             `json:"max_users"`
	MaxTags                   int             `json:"max_tags"`
	MaxOrders                 int             `json:"max_orders"`
	CurrentContacts           int             `json:"current_contacts"`
	CurrentMessages           int             `json:"current_messages"`
	CurrentBroadcasts         int             `json:"current_broadcasts"`
	CurrentUsers              int             `json:"current_users"`
	CurrentTags               int             `json:"current_tags"`
	CurrentOrders             int             `json:"current_orders"`
	SystemPrompt              *string         `json:"system_prompt" db:"system_prompt"`
	AIAssistantName           *string         `json:"ai_assistant_name" db:"ai_assistant_name"`
	AICommunicationStyle      *string         `json:"ai_communication_style" db:"ai_communication_style"`
	AICommunicationLanguage   *string         `json:"ai_communication_language" db:"ai_communication_language"`
	OrganizationInfoStructure json.RawMessage `json:"organization_info_structure" db:"organization_info_structure"`
	ContactInfoStructure      json.RawMessage `json:"contact_info_structure" db:"contact_info_structure"`
}

// default organization info structure
var OrganizationInfoStructure = map[string]string{
	"name":     "string",
	"address":  "string",
	"phone":    "string",
	"email":    "string",
	"website":  "string",
	"industry": "string",
}
