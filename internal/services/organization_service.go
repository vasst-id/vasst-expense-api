package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/vasst-id/vasst-expense-api/internal/entities"
	"github.com/vasst-id/vasst-expense-api/internal/repositories"
	"github.com/vasst-id/vasst-expense-api/internal/utils"
	errorsutil "github.com/vasst-id/vasst-expense-api/internal/utils/errors"
)

//go:generate mockgen -source=organization_service.go -package=mock -destination=mock/organization_service_mock.go
type (
	OrganizationService interface {
		// Organization methods
		CreateOrganization(ctx context.Context, input *entities.CreateOrganizationInput) (*entities.Organization, error)
		UpdateOrganization(ctx context.Context, orgID uuid.UUID, input *entities.UpdateOrganizationInput) (*entities.Organization, error)
		DeleteOrganization(ctx context.Context, orgID uuid.UUID) error
		ListOrganizations(ctx context.Context, limit, offset int) ([]*entities.Organization, error)
		GetOrganizationByID(ctx context.Context, orgID uuid.UUID) (*entities.Organization, error)
		GetOrganizationByCode(ctx context.Context, code string) (*entities.Organization, error)
		GetOrganizationByKey(ctx context.Context, key string) (*entities.Organization, error)
		GetDefaultUserByOrgID(ctx context.Context, orgID uuid.UUID) (*entities.User, error)

		// OrganizationCategory methods
		CreateCategory(ctx context.Context, input *entities.CreateOrganizationCategoryInput) (*entities.OrganizationCategory, error)
		UpdateCategory(ctx context.Context, categoryID int, input *entities.UpdateOrganizationCategoryInput) (*entities.OrganizationCategory, error)
		DeleteCategory(ctx context.Context, categoryID int) error
		ListCategories(ctx context.Context) ([]*entities.OrganizationCategory, error)
		GetCategoryByID(ctx context.Context, categoryID int) (*entities.OrganizationCategory, error)

		// OrganizationSetting methods
		UpdateSetting(ctx context.Context, orgID uuid.UUID, input *entities.UpdateOrganizationSettingInput) (*entities.OrganizationSetting, error)
		GetSettingByOrgID(ctx context.Context, orgID uuid.UUID) (*entities.OrganizationSetting, error)

		// OrganizationKnowledge methods
		CreateKnowledge(ctx context.Context, input *entities.CreateOrganizationKnowledgeInput) (*entities.OrganizationKnowledge, error)
		CreateKnowledgeWithFile(ctx context.Context, input *entities.CreateOrganizationKnowledgeInput, file *multipart.FileHeader) (*entities.OrganizationKnowledge, error)
		UpdateKnowledge(ctx context.Context, knowledgeID uuid.UUID, input *entities.UpdateOrganizationKnowledgeInput) (*entities.OrganizationKnowledge, error)
		UpdateKnowledgeWithFile(ctx context.Context, knowledgeID uuid.UUID, input *entities.UpdateOrganizationKnowledgeInput, file *multipart.FileHeader) (*entities.OrganizationKnowledge, error)
		DeleteKnowledge(ctx context.Context, knowledgeID uuid.UUID) error
		ListKnowledgeByOrgID(ctx context.Context, orgID uuid.UUID) ([]*entities.OrganizationKnowledge, error)
		GetKnowledgeByID(ctx context.Context, knowledgeID uuid.UUID) (*entities.OrganizationKnowledge, error)

		// OrganizationModel methods
		ListModelsByOrgID(ctx context.Context, orgID uuid.UUID) ([]*entities.OrganizationModel, error)

		// OrganizationIntegration methods
		CreateIntegration(ctx context.Context, input *entities.CreateOrganizationIntegrationInput) (*entities.OrganizationIntegration, error)
		UpdateIntegration(ctx context.Context, integrationID uuid.UUID, input *entities.UpdateOrganizationIntegrationInput) (*entities.OrganizationIntegration, error)
		DeleteIntegration(ctx context.Context, integrationID uuid.UUID) error
		ListIntegrationsByOrgID(ctx context.Context, orgID uuid.UUID) ([]*entities.OrganizationIntegration, error)
		GetIntegrationByID(ctx context.Context, integrationID uuid.UUID) (*entities.OrganizationIntegration, error)
		GetIntegrationTokenByOrgIDAndType(ctx context.Context, orgID uuid.UUID, integrationType string) (string, error)
		GetIntegrationByOrgIDAndType(ctx context.Context, orgID uuid.UUID, integrationType string) (*entities.OrganizationIntegration, error)

		// File upload methods
		UploadFile(ctx context.Context, orgID uuid.UUID, fileID uuid.UUID, file *multipart.FileHeader) (*entities.FileUploadResult, error)
	}

	organizationService struct {
		orgRepo    repositories.OrganizationRepository
		storageSvc GoogleStorageService
	}
)

// NewOrganizationService creates a new organization service
func NewOrganizationService(orgRepo repositories.OrganizationRepository, storageSvc GoogleStorageService) OrganizationService {
	return &organizationService{
		orgRepo:    orgRepo,
		storageSvc: storageSvc,
	}
}

// Organization methods
func (s *organizationService) CreateOrganization(ctx context.Context, input *entities.CreateOrganizationInput) (*entities.Organization, error) {
	// Validate input
	if input.OrganizationCode == "" {
		return nil, errors.New("organization code is required")
	}

	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	if input.ContactName == "" {
		return nil, errors.New("contact name is required")
	}

	if input.PhoneNumber == "" {
		return nil, errors.New("phone number is required")
	}

	// Check if organization already exists
	existingOrg, err := s.orgRepo.FindOrganizationByCode(ctx, input.OrganizationCode)
	if err != nil {
		return nil, err
	}

	if existingOrg != nil {
		return nil, errorsutil.New(409, "organization with this code already exists")
	}

	// Generate API key
	apiKey := utils.GenerateAPIKey()

	// Create new organization
	org := &entities.Organization{
		OrganizationID:   uuid.New(),
		OrganizationCode: input.OrganizationCode,
		Name:             input.Name,
		ContactName:      input.ContactName,
		PhoneNumber:      input.PhoneNumber,
		Email:            input.Email,
		WhatsappNumber:   input.WhatsappNumber,
		OrganizationType: input.OrganizationType,
		CategoryID:       input.CategoryID,
		Address:          input.Address,
		City:             input.City,
		Province:         input.Province,
		PostalCode:       input.PostalCode,
		Country:          input.Country,
		Status:           input.Status,
		APIKey:           apiKey,
	}

	if err := s.orgRepo.CreateOrganization(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

func (s *organizationService) UpdateOrganization(ctx context.Context, orgID uuid.UUID, input *entities.UpdateOrganizationInput) (*entities.Organization, error) {
	// Check if organization exists
	existingOrg, err := s.orgRepo.FindOrganizationByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingOrg == nil {
		return nil, errorsutil.New(404, "organization not found")
	}

	// Update organization fields
	if input.Name != "" {
		existingOrg.Name = input.Name
	}
	if input.ContactName != "" {
		existingOrg.ContactName = input.ContactName
	}
	if input.PhoneNumber != "" {
		existingOrg.PhoneNumber = input.PhoneNumber
	}
	if input.Email != "" {
		existingOrg.Email = input.Email
	}
	if input.WhatsappNumber != "" {
		existingOrg.WhatsappNumber = input.WhatsappNumber
	}
	if input.OrganizationType != 0 {
		existingOrg.OrganizationType = input.OrganizationType
	}
	if input.CategoryID != 0 {
		existingOrg.CategoryID = input.CategoryID
	}
	if input.Address != "" {
		existingOrg.Address = input.Address
	}
	if input.City != "" {
		existingOrg.City = input.City
	}
	if input.Province != "" {
		existingOrg.Province = input.Province
	}
	if input.PostalCode != "" {
		existingOrg.PostalCode = input.PostalCode
	}
	if input.Country != "" {
		existingOrg.Country = input.Country
	}
	if input.Status != nil {
		existingOrg.Status = *input.Status
	}

	if err := s.orgRepo.UpdateOrganization(ctx, existingOrg); err != nil {
		return nil, err
	}

	return existingOrg, nil
}

func (s *organizationService) DeleteOrganization(ctx context.Context, orgID uuid.UUID) error {
	// Check if organization exists
	existingOrg, err := s.orgRepo.FindOrganizationByID(ctx, orgID)
	if err != nil {
		return err
	}

	if existingOrg == nil {
		return errorsutil.New(404, "organization not found")
	}

	return s.orgRepo.DeleteOrganization(ctx, orgID)
}

func (s *organizationService) ListOrganizations(ctx context.Context, limit, offset int) ([]*entities.Organization, error) {
	return s.orgRepo.ListOrganizations(ctx, limit, offset)
}

func (s *organizationService) GetOrganizationByID(ctx context.Context, orgID uuid.UUID) (*entities.Organization, error) {
	org, err := s.orgRepo.FindOrganizationByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if org == nil {
		return nil, errorsutil.New(404, "organization not found")
	}

	return org, nil
}

func (s *organizationService) GetOrganizationByCode(ctx context.Context, code string) (*entities.Organization, error) {
	org, err := s.orgRepo.FindOrganizationByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if org == nil {
		return nil, errorsutil.New(404, "organization not found")
	}

	return org, nil
}

// OrganizationCategory methods
func (s *organizationService) CreateCategory(ctx context.Context, input *entities.CreateOrganizationCategoryInput) (*entities.OrganizationCategory, error) {
	// Validate input
	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	// Create new category
	category := &entities.OrganizationCategory{
		CategoryID:  0, // Will be set by the database
		Name:        input.Name,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		IsActive:    input.IsActive,
	}

	if err := s.orgRepo.CreateCategory(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *organizationService) UpdateCategory(ctx context.Context, categoryID int, input *entities.UpdateOrganizationCategoryInput) (*entities.OrganizationCategory, error) {
	// Check if category exists
	existingCategory, err := s.orgRepo.FindCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	if existingCategory == nil {
		return nil, errorsutil.New(404, "category not found")
	}

	// Update category fields
	if input.Name != "" {
		existingCategory.Name = input.Name
	}
	if input.Description != "" {
		existingCategory.Description = input.Description
	}
	if input.ImageURL != "" {
		existingCategory.ImageURL = input.ImageURL
	}
	if input.IsActive != nil {
		existingCategory.IsActive = *input.IsActive
	}

	if err := s.orgRepo.UpdateCategory(ctx, existingCategory); err != nil {
		return nil, err
	}

	return existingCategory, nil
}

func (s *organizationService) DeleteCategory(ctx context.Context, categoryID int) error {
	// Check if category exists
	existingCategory, err := s.orgRepo.FindCategoryByID(ctx, categoryID)
	if err != nil {
		return err
	}

	if existingCategory == nil {
		return errorsutil.New(404, "category not found")
	}

	return s.orgRepo.DeleteCategory(ctx, categoryID)
}

func (s *organizationService) ListCategories(ctx context.Context) ([]*entities.OrganizationCategory, error) {
	return s.orgRepo.ListCategories(ctx)
}

func (s *organizationService) GetCategoryByID(ctx context.Context, categoryID int) (*entities.OrganizationCategory, error) {
	category, err := s.orgRepo.FindCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errorsutil.New(404, "category not found")
	}

	return category, nil
}

// OrganizationSetting methods
func (s *organizationService) UpdateSetting(ctx context.Context, orgID uuid.UUID, input *entities.UpdateOrganizationSettingInput) (*entities.OrganizationSetting, error) {
	// Check if setting exists
	existingSetting, err := s.orgRepo.FindSettingByOrgID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingSetting == nil {
		return nil, errorsutil.New(404, "setting not found")
	}

	// Update setting fields
	if input.PlanID != 0 {
		existingSetting.PlanID = input.PlanID
	}
	if !input.PlanStartDate.IsZero() {
		existingSetting.PlanStartDate = input.PlanStartDate
	}
	if !input.PlanEndDate.IsZero() {
		existingSetting.PlanEndDate = input.PlanEndDate
	}
	if input.PlanStatus != 0 {
		existingSetting.PlanStatus = input.PlanStatus
	}
	if input.PlanAmount != 0 {
		existingSetting.PlanAmount = input.PlanAmount
	}
	if input.PlanCurrency != "" {
		existingSetting.PlanCurrency = input.PlanCurrency
	}
	if input.MaxContacts != 0 {
		existingSetting.MaxContacts = input.MaxContacts
	}
	if input.MaxMessages != 0 {
		existingSetting.MaxMessages = input.MaxMessages
	}
	if input.MaxBroadcasts != 0 {
		existingSetting.MaxBroadcasts = input.MaxBroadcasts
	}
	if input.MaxUsers != 0 {
		existingSetting.MaxUsers = input.MaxUsers
	}
	if input.MaxTags != 0 {
		existingSetting.MaxTags = input.MaxTags
	}
	if input.MaxOrders != 0 {
		existingSetting.MaxOrders = input.MaxOrders
	}
	if input.CurrentContacts != 0 {
		existingSetting.CurrentContacts = input.CurrentContacts
	}
	if input.CurrentMessages != 0 {
		existingSetting.CurrentMessages = input.CurrentMessages
	}
	if input.CurrentBroadcasts != 0 {
		existingSetting.CurrentBroadcasts = input.CurrentBroadcasts
	}
	if input.CurrentUsers != 0 {
		existingSetting.CurrentUsers = input.CurrentUsers
	}
	if input.CurrentTags != 0 {
		existingSetting.CurrentTags = input.CurrentTags
	}
	if input.CurrentOrders != 0 {
		existingSetting.CurrentOrders = input.CurrentOrders
	}
	if input.SystemPrompt != nil {
		existingSetting.SystemPrompt = input.SystemPrompt
	}
	if input.AIAssistantName != nil {
		existingSetting.AIAssistantName = input.AIAssistantName
	}
	if input.AICommunicationStyle != nil {
		existingSetting.AICommunicationStyle = input.AICommunicationStyle
	}
	if input.AICommunicationLanguage != nil {
		existingSetting.AICommunicationLanguage = input.AICommunicationLanguage
	}

	if err := s.orgRepo.UpdateSetting(ctx, existingSetting); err != nil {
		return nil, err
	}

	return existingSetting, nil
}

func (s *organizationService) GetSettingByOrgID(ctx context.Context, orgID uuid.UUID) (*entities.OrganizationSetting, error) {
	setting, err := s.orgRepo.FindSettingByOrgID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if setting == nil {
		return nil, errorsutil.New(404, "setting not found")
	}

	return setting, nil
}

// OrganizationKnowledge methods
func (s *organizationService) CreateKnowledge(ctx context.Context, input *entities.CreateOrganizationKnowledgeInput) (*entities.OrganizationKnowledge, error) {
	// Validate input
	if input.OrganizationID == uuid.Nil {
		return nil, errors.New("organization ID is required")
	}

	if input.KnowledgeType == 0 {
		return nil, errors.New("knowledge type is required")
	}

	if input.Title == "" {
		return nil, errors.New("title is required")
	}

	if input.Content == "" {
		return nil, errors.New("content is required")
	}

	// Prepare metadata
	var metadataJSON json.RawMessage
	if input.Metadata != nil {
		metadataBytes, err := json.Marshal(input.Metadata)
		if err != nil {
			return nil, err
		}
		metadataJSON = metadataBytes
	} else {
		metadataJSON = json.RawMessage("{}")
	}

	// Create new knowledge
	knowledge := &entities.OrganizationKnowledge{
		KnowledgeID:    uuid.New(),
		OrganizationID: input.OrganizationID,
		KnowledgeType:  input.KnowledgeType,
		Title:          &input.Title,
		Content:        input.Content,
		Description:    input.Description,
		Metadata:       metadataJSON,
		IsActive:       input.IsActive,
	}

	if err := s.orgRepo.CreateKnowledge(ctx, knowledge); err != nil {
		return nil, err
	}

	return knowledge, nil
}

func (s *organizationService) CreateKnowledgeWithFile(ctx context.Context, input *entities.CreateOrganizationKnowledgeInput, file *multipart.FileHeader) (*entities.OrganizationKnowledge, error) {
	// First create the knowledge entry
	knowledge, err := s.CreateKnowledge(ctx, input)
	if err != nil {
		return nil, err
	}

	// Upload file to Google Cloud Storage
	uploadResult, err := s.storageSvc.UploadFile(ctx, input.OrganizationCode, uuid.Nil, file)
	if err != nil {
		return nil, err
	}

	fmt.Println("uploadResult", uploadResult)

	if uploadResult == nil {
		return nil, errors.New("upload result is empty")
	}

	// Update knowledge with file information
	knowledge.SourceURL = &uploadResult.FileURL
	knowledge.FileName = &uploadResult.FileName
	knowledge.FileSize = &uploadResult.FileSize
	knowledge.ContentType = &uploadResult.ContentType
	knowledge.BucketName = &uploadResult.BucketName
	knowledge.ObjectName = &uploadResult.ObjectName

	// Update the knowledge entry in the database
	if err := s.orgRepo.UpdateKnowledge(ctx, knowledge); err != nil {
		return nil, err
	}

	return knowledge, nil
}

func (s *organizationService) UpdateKnowledge(ctx context.Context, knowledgeID uuid.UUID, input *entities.UpdateOrganizationKnowledgeInput) (*entities.OrganizationKnowledge, error) {
	// Check if knowledge exists
	existingKnowledge, err := s.orgRepo.FindKnowledgeByID(ctx, knowledgeID)
	if err != nil {
		return nil, err
	}

	if existingKnowledge == nil {
		return nil, errorsutil.New(404, "knowledge not found")
	}

	// Update knowledge fields
	if input.KnowledgeType != 0 {
		existingKnowledge.KnowledgeType = input.KnowledgeType
	}
	if input.Title != nil && *input.Title != "" {
		existingKnowledge.Title = input.Title
	}
	if input.Content != "" {
		existingKnowledge.Content = input.Content
	}
	if input.Description != nil && *input.Description != "" {
		existingKnowledge.Description = input.Description
	}
	if input.Metadata != nil {
		metadataBytes, err := json.Marshal(input.Metadata)
		if err != nil {
			return nil, err
		}
		existingKnowledge.Metadata = metadataBytes
	}
	if input.IsActive != nil {
		existingKnowledge.IsActive = *input.IsActive
	}

	if err := s.orgRepo.UpdateKnowledge(ctx, existingKnowledge); err != nil {
		return nil, err
	}

	return existingKnowledge, nil
}

func (s *organizationService) UpdateKnowledgeWithFile(ctx context.Context, knowledgeID uuid.UUID, input *entities.UpdateOrganizationKnowledgeInput, file *multipart.FileHeader) (*entities.OrganizationKnowledge, error) {
	// First update the knowledge entry
	knowledge, err := s.UpdateKnowledge(ctx, knowledgeID, input)
	if err != nil {
		return nil, err
	}

	// Get organization code
	organization, err := s.orgRepo.FindOrganizationByID(ctx, knowledge.OrganizationID)
	if err != nil {
		return nil, err
	}

	if organization == nil {
		return nil, errorsutil.New(404, "organization not found")
	}

	// Upload new file to Google Cloud Storage
	uploadResult, err := s.storageSvc.UploadFile(ctx, organization.OrganizationCode, uuid.Nil, file)
	if err != nil {
		return nil, err
	}

	// Update knowledge with new file information
	knowledge.SourceURL = &uploadResult.FileURL
	knowledge.FileName = &uploadResult.FileName
	knowledge.FileSize = &uploadResult.FileSize
	knowledge.ContentType = &uploadResult.ContentType
	knowledge.BucketName = &uploadResult.BucketName
	knowledge.ObjectName = &uploadResult.ObjectName

	// Update the knowledge entry in the database
	if err := s.orgRepo.UpdateKnowledge(ctx, knowledge); err != nil {
		return nil, err
	}

	return knowledge, nil
}

func (s *organizationService) DeleteKnowledge(ctx context.Context, knowledgeID uuid.UUID) error {
	// Check if knowledge exists
	existingKnowledge, err := s.orgRepo.FindKnowledgeByID(ctx, knowledgeID)
	if err != nil {
		return err
	}

	if existingKnowledge == nil {
		return errorsutil.New(404, "knowledge not found")
	}

	return s.orgRepo.DeleteKnowledge(ctx, knowledgeID)
}

func (s *organizationService) ListKnowledgeByOrgID(ctx context.Context, orgID uuid.UUID) ([]*entities.OrganizationKnowledge, error) {
	return s.orgRepo.ListKnowledgeByOrgID(ctx, orgID)
}

func (s *organizationService) GetKnowledgeByID(ctx context.Context, knowledgeID uuid.UUID) (*entities.OrganizationKnowledge, error) {
	knowledge, err := s.orgRepo.FindKnowledgeByID(ctx, knowledgeID)
	if err != nil {
		return nil, err
	}

	if knowledge == nil {
		return nil, errorsutil.New(404, "knowledge not found")
	}

	return knowledge, nil
}

// OrganizationModel methods
func (s *organizationService) CreateModel(ctx context.Context, input *entities.CreateOrganizationModelInput) (*entities.OrganizationModel, error) {
	// Validate input
	if input.OrganizationID == uuid.Nil {
		return nil, errors.New("organization ID is required")
	}

	if input.ModelID == uuid.Nil {
		return nil, errors.New("model ID is required")
	}

	// Create new model
	model := &entities.OrganizationModel{
		OrganizationModelID: uuid.New(),
		OrganizationID:      input.OrganizationID,
		ModelID:             input.ModelID,
		IsActive:            input.IsActive,
	}

	if err := s.orgRepo.CreateModel(ctx, model); err != nil {
		return nil, err
	}

	return model, nil
}

func (s *organizationService) DeleteModel(ctx context.Context, modelID uuid.UUID) error {
	return s.orgRepo.DeleteModel(ctx, modelID)
}

func (s *organizationService) ListModelsByOrgID(ctx context.Context, orgID uuid.UUID) ([]*entities.OrganizationModel, error) {
	return s.orgRepo.ListModelsByOrgID(ctx, orgID)
}

// OrganizationIntegration methods
func (s *organizationService) CreateIntegration(ctx context.Context, input *entities.CreateOrganizationIntegrationInput) (*entities.OrganizationIntegration, error) {
	// Validate input
	if input.OrganizationID == uuid.Nil {
		return nil, errors.New("organization ID is required")
	}

	if input.IntegrationID == uuid.Nil {
		return nil, errors.New("integration ID is required")
	}

	if input.Token == "" {
		return nil, errors.New("token is required")
	}

	// Create new integration
	integration := &entities.OrganizationIntegration{
		OrganizationIntegrationID: uuid.New(),
		OrganizationID:            input.OrganizationID,
		IntegrationID:             input.IntegrationID,
		Token:                     input.Token,
		LastUsedAt:                nil,
		IsActive:                  input.IsActive,
	}

	if err := s.orgRepo.CreateIntegration(ctx, integration); err != nil {
		return nil, err
	}

	return integration, nil
}

func (s *organizationService) UpdateIntegration(ctx context.Context, integrationID uuid.UUID, input *entities.UpdateOrganizationIntegrationInput) (*entities.OrganizationIntegration, error) {
	// Check if integration exists
	existingIntegration, err := s.orgRepo.FindIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	if existingIntegration == nil {
		return nil, errorsutil.New(404, "integration not found")
	}

	// Update integration fields
	if input.Token != "" {
		existingIntegration.Token = input.Token
	}
	if input.LastUsedAt != nil {
		existingIntegration.LastUsedAt = input.LastUsedAt
	}
	if input.IsActive != nil {
		existingIntegration.IsActive = *input.IsActive
	}

	if err := s.orgRepo.UpdateIntegration(ctx, existingIntegration); err != nil {
		return nil, err
	}

	return existingIntegration, nil
}

func (s *organizationService) DeleteIntegration(ctx context.Context, integrationID uuid.UUID) error {
	// Check if integration exists
	existingIntegration, err := s.orgRepo.FindIntegrationByID(ctx, integrationID)
	if err != nil {
		return err
	}

	if existingIntegration == nil {
		return errorsutil.New(404, "integration not found")
	}

	return s.orgRepo.DeleteIntegration(ctx, integrationID)
}

func (s *organizationService) ListIntegrationsByOrgID(ctx context.Context, orgID uuid.UUID) ([]*entities.OrganizationIntegration, error) {
	return s.orgRepo.ListIntegrationsByOrgID(ctx, orgID)
}

func (s *organizationService) GetIntegrationByID(ctx context.Context, integrationID uuid.UUID) (*entities.OrganizationIntegration, error) {
	integration, err := s.orgRepo.FindIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	if integration == nil {
		return nil, errorsutil.New(404, "integration not found")
	}

	return integration, nil
}

func (s *organizationService) GetOrganizationByKey(ctx context.Context, key string) (*entities.Organization, error) {
	if key == "" {
		return nil, errors.New("key is required")
	}

	org, err := s.orgRepo.FindOrganizationByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	if org == nil {
		return nil, errorsutil.New(404, "org not found")
	}

	return org, nil
}

func (s *organizationService) GetIntegrationTokenByOrgIDAndType(ctx context.Context, orgID uuid.UUID, integrationType string) (string, error) {
	token, err := s.orgRepo.FindIntegrationTokenByOrgIDAndType(ctx, orgID, integrationType)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *organizationService) GetIntegrationByOrgIDAndType(ctx context.Context, orgID uuid.UUID, integrationType string) (*entities.OrganizationIntegration, error) {
	integration, err := s.orgRepo.FindIntegrationByOrgIDAndType(ctx, orgID, integrationType)
	if err != nil {
		return nil, err
	}

	return integration, nil
}

func (s *organizationService) GetDefaultUserByOrgID(ctx context.Context, orgID uuid.UUID) (*entities.User, error) {
	user, err := s.orgRepo.FindDefaultUserByOrgID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// File upload methods
func (s *organizationService) UploadFile(ctx context.Context, orgID uuid.UUID, fileID uuid.UUID, file *multipart.FileHeader) (*entities.FileUploadResult, error) {
	// Get organization code
	organization, err := s.orgRepo.FindOrganizationByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, errorsutil.New(404, "organization not found")
	}

	// Upload file to Google Cloud Storage using organization bucket
	uploadResult, err := s.storageSvc.UploadFile(ctx, organization.OrganizationCode, fileID, file)
	if err != nil {
		return nil, err
	}

	// Map the storage service result to entities FileUploadResult
	return &entities.FileUploadResult{
		FileID:      fileID,
		FileName:    uploadResult.FileName,
		FileSize:    uploadResult.FileSize,
		ContentType: uploadResult.ContentType,
		FileURL:     uploadResult.FileURL,
		BucketName:  uploadResult.BucketName,
		ObjectName:  uploadResult.ObjectName,
		UploadedAt:  uploadResult.UploadedAt,
	}, nil
}
