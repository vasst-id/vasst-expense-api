package entities

import (
	"encoding/json"
	"time"
)

// Plan represents a subscription plan in the system
type SubscriptionPlan struct {
	SubscriptionPlanID          int             `json:"subscription_plan_id" db:"subscription_plan_id"`
	SubscriptionPlanName        string          `json:"subscription_plan_name" db:"subscription_plan_name"`
	SubscriptionPlanDescription string          `json:"subscription_plan_description" db:"subscription_plan_description"`
	SubscriptionPlanPrice       string          `json:"subscription_plan_price" db:"subscription_plan_price"`
	SubscriptionPlanCurrencyID  int             `json:"subscription_plan_currency_id" db:"subscription_plan_currency_id"`
	SubscriptionPlanFeatures    json.RawMessage `json:"subscription_plan_features" db:"subscription_plan_features"`
	SubscriptionPlanStatus      bool            `json:"subscription_plan_status" db:"subscription_plan_status"`
	CreatedAt                   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt                   time.Time       `json:"updated_at" db:"updated_at"`
}

// CreatePlanInput is used for creating a new plan
type CreateSubscriptionPlanInput struct {
	SubscriptionPlanName        string          `json:"subscription_plan_name" binding:"required"`
	SubscriptionPlanDescription string          `json:"subscription_plan_description"`
	SubscriptionPlanPrice       string          `json:"subscription_plan_price" binding:"required"`
	SubscriptionPlanCurrencyID  int             `json:"subscription_plan_currency_id"`
	SubscriptionPlanFeatures    json.RawMessage `json:"subscription_plan_features" binding:"required"`
	SubscriptionPlanStatus      bool            `json:"subscription_plan_status"`
}

// UpdatePlanInput is used for updating an existing plan
type UpdateSubscriptionPlanInput struct {
	SubscriptionPlanName        string          `json:"subscription_plan_name"`
	SubscriptionPlanDescription string          `json:"subscription_plan_description"`
	SubscriptionPlanPrice       string          `json:"subscription_plan_price"`
	SubscriptionPlanCurrencyID  int             `json:"subscription_plan_currency_id"`
	SubscriptionPlanFeatures    json.RawMessage `json:"subscription_plan_features"`
	SubscriptionPlanStatus      bool            `json:"subscription_plan_status"`
}

type SubscriptionPlanSimple struct {
	SubscriptionPlanID          int             `json:"subscription_plan_id" db:"subscription_plan_id"`
	SubscriptionPlanName        string          `json:"subscription_plan_name" db:"subscription_plan_name"`
	SubscriptionPlanDescription string          `json:"subscription_plan_description" db:"subscription_plan_description"`
	SubscriptionPlanPrice       string          `json:"subscription_plan_price" db:"subscription_plan_price"`
	SubscriptionPlanCurrencyID  int             `json:"subscription_plan_currency_id" db:"subscription_plan_currency_id"`
	SubscriptionPlanFeatures    json.RawMessage `json:"subscription_plan_features" db:"subscription_plan_features"`
	SubscriptionPlanStatus      bool            `json:"subscription_plan_status" db:"subscription_plan_status"`
}
