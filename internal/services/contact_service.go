package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=contact_service.go -package=mock -destination=mock/contact_service_mock.go
type (
	ContactService interface {
		CreateContact(ctx context.Context, input *entities.CreateContactInput) (*entities.Contact, error)
		UpdateContact(ctx context.Context, contactID uuid.UUID, input *entities.UpdateContactInput) (*entities.Contact, error)
		DeleteContact(ctx context.Context, contactID uuid.UUID) error
		ListContacts(ctx context.Context, limit, offset int) ([]*entities.Contact, error)
		GetContactByID(ctx context.Context, contactID uuid.UUID) (*entities.Contact, error)
		GetContactByIDAndOrganization(ctx context.Context, contactID, organizationID uuid.UUID) (*entities.Contact, error)
		GetContactByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.Contact, error)
		ListContactsByOrganizationID(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.Contact, error)
	}

	contactService struct {
		contactRepo repositories.ContactRepository
	}
)

// NewContactService creates a new contact service
func NewContactService(contactRepo repositories.ContactRepository) ContactService {
	return &contactService{
		contactRepo: contactRepo,
	}
}

// CreateContact creates a new contact
func (s *contactService) CreateContact(ctx context.Context, input *entities.CreateContactInput) (*entities.Contact, error) {
	// Validate input
	if input.OrganizationID == uuid.Nil {
		return nil, errors.New("organization ID is required")
	}

	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	if input.PhoneNumber == "" {
		return nil, errors.New("phone number is required")
	}

	// Check if contact already exists
	existingContact, err := s.contactRepo.FindByPhoneNumber(ctx, input.PhoneNumber)
	if err != nil {
		return nil, err
	}

	if existingContact != nil {
		return nil, errorsutil.New(409, "contact with this phone number already exists")
	}

	// Create new contact
	contact := &entities.Contact{
		ContactID:      uuid.New(),
		OrganizationID: input.OrganizationID,
		Name:           input.Name,
		PhoneNumber:    input.PhoneNumber,
		Email:          input.Email,
	}

	if err := s.contactRepo.Create(ctx, contact); err != nil {
		return nil, err
	}

	return contact, nil
}

// UpdateContact updates an existing contact
func (s *contactService) UpdateContact(ctx context.Context, contactID uuid.UUID, input *entities.UpdateContactInput) (*entities.Contact, error) {
	// Check if contact exists
	existingContact, err := s.contactRepo.FindByID(ctx, contactID)
	if err != nil {
		return nil, err
	}

	if existingContact == nil {
		return nil, errorsutil.New(404, "contact not found")
	}

	// Check if phone number is already taken by another contact
	if input.PhoneNumber != "" && input.PhoneNumber != existingContact.PhoneNumber {
		contactWithPhone, err := s.contactRepo.FindByPhoneNumber(ctx, input.PhoneNumber)
		if err != nil {
			return nil, err
		}
		if contactWithPhone != nil && contactWithPhone.ContactID != contactID {
			return nil, errorsutil.New(409, "phone number already in use")
		}
	}

	// Update contact fields
	if input.Name != "" {
		existingContact.Name = input.Name
	}
	if input.PhoneNumber != "" {
		existingContact.PhoneNumber = input.PhoneNumber
	}
	if input.Email != "" {
		existingContact.Email = input.Email
	}
	if input.Salutation != "" {
		existingContact.Salutation = input.Salutation
	}
	if input.Notes != "" {
		existingContact.Notes = input.Notes
	}
	if input.CustomSystemPrompt != "" {
		existingContact.CustomSystemPrompt = input.CustomSystemPrompt
	}
	if len(input.Context) > 0 {
		existingContact.Context = input.Context
	}

	if err := s.contactRepo.Update(ctx, existingContact); err != nil {
		return nil, err
	}

	return existingContact, nil
}

// DeleteContact deletes a contact by their ID
func (s *contactService) DeleteContact(ctx context.Context, contactID uuid.UUID) error {
	// Check if contact exists
	existingContact, err := s.contactRepo.FindByID(ctx, contactID)
	if err != nil {
		return err
	}

	if existingContact == nil {
		return errorsutil.New(404, "contact not found")
	}

	return s.contactRepo.Delete(ctx, contactID)
}

// ListContacts returns all contacts with optional filtering
func (s *contactService) ListContacts(ctx context.Context, limit, offset int) ([]*entities.Contact, error) {
	return s.contactRepo.List(ctx, limit, offset)
}

// GetContactByID returns a contact by their ID
func (s *contactService) GetContactByID(ctx context.Context, contactID uuid.UUID) (*entities.Contact, error) {
	contact, err := s.contactRepo.FindByID(ctx, contactID)
	if err != nil {
		return nil, err
	}

	if contact == nil {
		return nil, errorsutil.New(404, "contact not found")
	}

	return contact, nil
}

// GetContactByIDAndOrganization returns a contact by their ID with organization validation for tenant isolation
func (s *contactService) GetContactByIDAndOrganization(ctx context.Context, contactID, organizationID uuid.UUID) (*entities.Contact, error) {
	contact, err := s.contactRepo.FindByID(ctx, contactID)
	if err != nil {
		return nil, err
	}

	if contact == nil {
		return nil, errorsutil.New(404, "contact not found")
	}

	// Ensure tenant isolation - contact must belong to the organization
	if contact.OrganizationID != organizationID {
		return nil, errorsutil.New(404, "contact not found")
	}

	return contact, nil
}

// GetContactByPhoneNumber returns a contact by their phone number
func (s *contactService) GetContactByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.Contact, error) {
	contact, err := s.contactRepo.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, err
	}

	if contact == nil {
		return nil, errorsutil.New(404, "contact not found")
	}

	return contact, nil
}

// ListContactsByOrganizationID returns all contacts for an organization
func (s *contactService) ListContactsByOrganizationID(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.Contact, error) {
	return s.contactRepo.ListByOrganizationID(ctx, organizationID, limit, offset)
}
