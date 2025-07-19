package entities

import (
	"time"

	"github.com/google/uuid"
)

// Budget represents a spending budget
type Budget struct {
	BudgetID       uuid.UUID `json:"budget_id" db:"budget_id"`
	WorkspaceID    uuid.UUID `json:"workspace_id" db:"workspace_id"`
	UserCategoryID uuid.UUID `json:"user_category_id" db:"user_category_id"`
	Name           string    `json:"name" db:"name"`
	BudgetedAmount float64   `json:"budgeted_amount" db:"budgeted_amount"`
	PeriodType     int       `json:"period_type" db:"period_type"`
	PeriodStart    time.Time `json:"period_start" db:"period_start"`
	PeriodEnd      time.Time `json:"period_end" db:"period_end"`
	SpentAmount    float64   `json:"spent_amount" db:"spent_amount"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedBy      uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type BudgetSimple struct {
	BudgetID         uuid.UUID `json:"budget_id" db:"budget_id"`
	UserCategoryName string    `json:"user_category_name" db:"user_category_name"`
	Name             string    `json:"name" db:"name"`
	BudgetedAmount   float64   `json:"budgeted_amount" db:"budgeted_amount"`
	PeriodTypeLabel  string    `json:"period_type_label" db:"period_type_label"`
	PeriodStart      time.Time `json:"period_start" db:"period_start"`
	PeriodEnd        time.Time `json:"period_end" db:"period_end"`
	SpentAmount      float64   `json:"spent_amount" db:"spent_amount"`
	RemainingAmount  float64   `json:"remaining_amount" db:"remaining_amount"`
	PercentageUsed   float64   `json:"percentage_used" db:"percentage_used"`
	DaysRemaining    int       `json:"days_remaining" db:"days_remaining"`
	IsOverspent      bool      `json:"is_overspent" db:"is_overspent"`
}

// CreateBudgetRequest represents the create budget request
type CreateBudgetRequest struct {
	WorkspaceID    uuid.UUID `json:"workspace_id"`
	UserCategoryID uuid.UUID `json:"user_category_id"`
	Name           string    `json:"name" binding:"required"`
	BudgetedAmount float64   `json:"budgeted_amount" binding:"required"`
	PeriodType     int       `json:"period_type" binding:"required"`
	PeriodStart    time.Time `json:"period_start" binding:"required" time:"2006-01-02"`
	PeriodEnd      time.Time `json:"period_end" binding:"required" time:"2006-01-02"`
	CreatedBy      uuid.UUID `json:"created_by" db:"created_by"`
}

type UpdateBudgetRequest struct {
	UserCategoryID uuid.UUID `json:"user_category_id"`
	Name           string    `json:"name" binding:"required"`
	BudgetedAmount float64   `json:"budgeted_amount" binding:"required"`
	PeriodType     int       `json:"period_type" binding:"required"`
	PeriodStart    time.Time `json:"period_start" binding:"required" time:"2006-01-02"`
	PeriodEnd      time.Time `json:"period_end" binding:"required" time:"2006-01-02"`
	SpentAmount    float64   `json:"spent_amount" binding:"required"`
	IsActive       bool      `json:"is_active" binding:"required"`
}

// Constants for period types
const (
	PeriodTypeWeekly  = 1
	PeriodTypeMonthly = 2
	PeriodTypeYearly  = 3
	PeriodTypeOneTime = 4
)
