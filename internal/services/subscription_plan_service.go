package services

import (
	"context"
	"errors"

	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=subscription_plan_service.go -package=mock -destination=mock/subscription_plan_service_mock.go
type (
	SubscriptionPlanService interface {
		CreateSubscriptionPlan(ctx context.Context, input *entities.CreateSubscriptionPlanInput) (*entities.SubscriptionPlan, error)
		UpdateSubscriptionPlan(ctx context.Context, subscriptionPlanID int, input *entities.UpdateSubscriptionPlanInput) (*entities.SubscriptionPlan, error)
		DeleteSubscriptionPlan(ctx context.Context, subscriptionPlanID int) error
		GetAllSubscriptionPlans(ctx context.Context) ([]*entities.SubscriptionPlanSimple, error)
		GetSubscriptionPlanByID(ctx context.Context, subscriptionPlanID int) (*entities.SubscriptionPlan, error)
	}

	subscriptionPlanService struct {
		subscriptionPlanRepo repositories.SubscriptionPlanRepository
	}
)

// NewSubscriptionPlanService creates a new subscription plan service
func NewSubscriptionPlanService(subscriptionPlanRepo repositories.SubscriptionPlanRepository) SubscriptionPlanService {
	return &subscriptionPlanService{
		subscriptionPlanRepo: subscriptionPlanRepo,
	}
}

// CreateSubscriptionPlan creates a new subscription plan
func (s *subscriptionPlanService) CreateSubscriptionPlan(ctx context.Context, input *entities.CreateSubscriptionPlanInput) (*entities.SubscriptionPlan, error) {
	// Validate required fields
	if input.SubscriptionPlanName == "" {
		return nil, errors.New("subscription plan name is required")
	}
	if input.SubscriptionPlanPrice == "" {
		return nil, errors.New("subscription plan price is required")
	}
	if input.SubscriptionPlanCurrencyID == 0 {
		return nil, errors.New("plan duration is required")
	}

	plan := &entities.SubscriptionPlan{
		SubscriptionPlanName:        input.SubscriptionPlanName,
		SubscriptionPlanDescription: input.SubscriptionPlanDescription,
		SubscriptionPlanPrice:       input.SubscriptionPlanPrice,
		SubscriptionPlanCurrencyID:  input.SubscriptionPlanCurrencyID,
		SubscriptionPlanFeatures:    input.SubscriptionPlanFeatures,
		SubscriptionPlanStatus:      input.SubscriptionPlanStatus,
	}

	// Create the plan - the repository will populate the struct with the actual data from DB
	createdPlan, err := s.subscriptionPlanRepo.Create(ctx, plan)
	if err != nil {
		return nil, err
	}

	// Return the plan with data populated from the database
	return &createdPlan, nil
}

// UpdateSubscriptionPlan updates an existing subscription plan
func (s *subscriptionPlanService) UpdateSubscriptionPlan(ctx context.Context, subscriptionPlanID int, input *entities.UpdateSubscriptionPlanInput) (*entities.SubscriptionPlan, error) {
	existingPlan, err := s.subscriptionPlanRepo.FindByID(ctx, subscriptionPlanID)
	if err != nil {
		return nil, err
	}
	if existingPlan == nil {
		return nil, errorsutil.New(404, "plan not found")
	}

	// Update fields
	if input.SubscriptionPlanName != "" {
		existingPlan.SubscriptionPlanName = input.SubscriptionPlanName
	}
	if input.SubscriptionPlanDescription != "" {
		existingPlan.SubscriptionPlanDescription = input.SubscriptionPlanDescription
	}
	if input.SubscriptionPlanPrice != "" {
		existingPlan.SubscriptionPlanPrice = input.SubscriptionPlanPrice
	}
	if input.SubscriptionPlanCurrencyID != 0 {
		existingPlan.SubscriptionPlanCurrencyID = input.SubscriptionPlanCurrencyID
	}
	if input.SubscriptionPlanFeatures != nil {
		existingPlan.SubscriptionPlanFeatures = input.SubscriptionPlanFeatures
	}
	if input.SubscriptionPlanStatus != false {
		existingPlan.SubscriptionPlanStatus = input.SubscriptionPlanStatus
	}

	// Update the plan - the repository will populate the struct with the actual data from DB
	updatedPlan, err := s.subscriptionPlanRepo.Update(ctx, existingPlan)
	if err != nil {
		return nil, err
	}

	// Return the plan with data populated from the database
	return &updatedPlan, nil
}

// DeletePlan deletes a plan
func (s *subscriptionPlanService) DeleteSubscriptionPlan(ctx context.Context, subscriptionPlanID int) error {
	existingPlan, err := s.subscriptionPlanRepo.FindByID(ctx, subscriptionPlanID)
	if err != nil {
		return err
	}
	if existingPlan == nil {
		return errorsutil.New(404, "plan not found")
	}
	return s.subscriptionPlanRepo.Delete(ctx, subscriptionPlanID)
}

// GetAllPlans returns all active plans in simple format
func (s *subscriptionPlanService) GetAllSubscriptionPlans(ctx context.Context) ([]*entities.SubscriptionPlanSimple, error) {
	return s.subscriptionPlanRepo.FindAll(ctx)
}

// GetSubscriptionPlanByID returns a subscription plan by ID
func (s *subscriptionPlanService) GetSubscriptionPlanByID(ctx context.Context, subscriptionPlanID int) (*entities.SubscriptionPlan, error) {
	subscriptionPlan, err := s.subscriptionPlanRepo.FindByID(ctx, subscriptionPlanID)
	if err != nil {
		return nil, err
	}
	if subscriptionPlan == nil {
		return nil, errorsutil.New(404, "subscription plan not found")
	}
	return subscriptionPlan, nil
}
