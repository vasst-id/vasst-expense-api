package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=leads_service.go -package=mock -destination=mock/leads_service_mock.go
type (
	LeadsService interface {
		CreateLead(ctx context.Context, input *entities.CreateLeadInput) (*entities.Lead, error)
		UpdateLead(ctx context.Context, leadID uuid.UUID, input *entities.UpdateLeadInput) (*entities.Lead, error)
		GetLeadByID(ctx context.Context, leadID uuid.UUID) (*entities.Lead, error)
		ListAllLeads(ctx context.Context, limit, offset int) ([]*entities.Lead, error)
		GetLeadByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.Lead, error)
		GetLeadByEmail(ctx context.Context, email string) (*entities.Lead, error)
	}

	leadsService struct {
		leadsRepo repositories.LeadsRepository
	}
)

// NewLeadsService creates a new leads service
func NewLeadsService(leadsRepo repositories.LeadsRepository) LeadsService {
	return &leadsService{
		leadsRepo: leadsRepo,
	}
}

func (s *leadsService) CreateLead(ctx context.Context, input *entities.CreateLeadInput) (*entities.Lead, error) {
	if input.Name == "" {
		return nil, errors.New("name is required")
	}
	if input.PhoneNumber == "" {
		return nil, errors.New("phone number is required")
	}

	// Check if lead with this phone number already exists
	existingLead, err := s.leadsRepo.FindByPhoneNumber(ctx, input.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if existingLead != nil {
		return nil, errorsutil.New(409, "lead with this phone number already exists")
	}

	// Check if lead with this email already exists (if email is provided)
	if input.Email != "" {
		existingEmailLead, err := s.leadsRepo.FindByEmail(ctx, input.Email)
		if err != nil {
			return nil, err
		}
		if existingEmailLead != nil {
			return nil, errorsutil.New(409, "lead with this email already exists")
		}
	}

	lead := &entities.Lead{
		LeadID:              uuid.New(),
		Name:                input.Name,
		PhoneNumber:         input.PhoneNumber,
		Email:               input.Email,
		BusinessName:        input.BusinessName,
		BusinessAddress:     input.BusinessAddress,
		BusinessPhoneNumber: input.BusinessPhoneNumber,
		BusinessEmail:       input.BusinessEmail,
		BusinessWebsite:     input.BusinessWebsite,
		BusinessIndustry:    input.BusinessIndustry,
		BusinessSize:        input.BusinessSize,
	}

	if err := s.leadsRepo.Create(ctx, lead); err != nil {
		return nil, err
	}

	return lead, nil
}

func (s *leadsService) UpdateLead(ctx context.Context, leadID uuid.UUID, input *entities.UpdateLeadInput) (*entities.Lead, error) {
	existingLead, err := s.leadsRepo.FindByID(ctx, leadID)
	if err != nil {
		return nil, err
	}
	if existingLead == nil {
		return nil, errorsutil.New(404, "lead not found")
	}

	// Check if phone number is being changed and if it's already in use
	if input.PhoneNumber != "" && input.PhoneNumber != existingLead.PhoneNumber {
		leadWithPhone, err := s.leadsRepo.FindByPhoneNumber(ctx, input.PhoneNumber)
		if err != nil {
			return nil, err
		}
		if leadWithPhone != nil && leadWithPhone.LeadID != leadID {
			return nil, errorsutil.New(409, "phone number already in use")
		}
	}

	// Check if email is being changed and if it's already in use
	if input.Email != "" && input.Email != existingLead.Email {
		leadWithEmail, err := s.leadsRepo.FindByEmail(ctx, input.Email)
		if err != nil {
			return nil, err
		}
		if leadWithEmail != nil && leadWithEmail.LeadID != leadID {
			return nil, errorsutil.New(409, "email already in use")
		}
	}

	// Update fields if provided
	if input.Name != "" {
		existingLead.Name = input.Name
	}
	if input.PhoneNumber != "" {
		existingLead.PhoneNumber = input.PhoneNumber
	}
	if input.Email != "" {
		existingLead.Email = input.Email
	}
	if input.BusinessName != "" {
		existingLead.BusinessName = input.BusinessName
	}
	if input.BusinessAddress != "" {
		existingLead.BusinessAddress = input.BusinessAddress
	}
	if input.BusinessPhoneNumber != "" {
		existingLead.BusinessPhoneNumber = input.BusinessPhoneNumber
	}
	if input.BusinessEmail != "" {
		existingLead.BusinessEmail = input.BusinessEmail
	}
	if input.BusinessWebsite != "" {
		existingLead.BusinessWebsite = input.BusinessWebsite
	}
	if input.BusinessIndustry != "" {
		existingLead.BusinessIndustry = input.BusinessIndustry
	}
	if input.BusinessSize != "" {
		existingLead.BusinessSize = input.BusinessSize
	}

	if err := s.leadsRepo.Update(ctx, existingLead); err != nil {
		return nil, err
	}

	return existingLead, nil
}

func (s *leadsService) GetLeadByID(ctx context.Context, leadID uuid.UUID) (*entities.Lead, error) {
	lead, err := s.leadsRepo.FindByID(ctx, leadID)
	if err != nil {
		return nil, err
	}
	if lead == nil {
		return nil, errorsutil.New(404, "lead not found")
	}
	return lead, nil
}

func (s *leadsService) ListAllLeads(ctx context.Context, limit, offset int) ([]*entities.Lead, error) {
	return s.leadsRepo.ListAll(ctx, limit, offset)
}

func (s *leadsService) GetLeadByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.Lead, error) {
	lead, err := s.leadsRepo.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, err
	}
	if lead == nil {
		return nil, errorsutil.New(404, "lead not found")
	}
	return lead, nil
}

func (s *leadsService) GetLeadByEmail(ctx context.Context, email string) (*entities.Lead, error) {
	lead, err := s.leadsRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if lead == nil {
		return nil, errorsutil.New(404, "lead not found")
	}
	return lead, nil
}
