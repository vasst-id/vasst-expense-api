package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// OrganizationKnowledge represents knowledge base entries for an organization
type OrganizationKnowledge struct {
	KnowledgeID    uuid.UUID       `json:"knowledge_id" db:"organization_knowledge_id"`
	OrganizationID uuid.UUID       `json:"organization_id" db:"organization_id"`
	KnowledgeType  int             `json:"knowledge_type" db:"knowledge_type"`
	Title          *string         `json:"title" db:"title"`
	Content        string          `json:"content" db:"knowledge_content"`
	Description    *string         `json:"description" db:"description"`
	SourceURL      *string         `json:"source_url" db:"source_url"`
	Metadata       json.RawMessage `json:"metadata" db:"metadata"`
	FileName       *string         `json:"file_name" db:"file_name"`
	FileSize       *int64          `json:"file_size" db:"file_size"`
	ContentType    *string         `json:"content_type" db:"content_type"`
	BucketName     *string         `json:"bucket_name" db:"bucket_name"`
	ObjectName     *string         `json:"object_name" db:"object_name"`
	IsActive       bool            `json:"is_active" db:"is_active"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}

// CreateOrganizationKnowledgeInput is used for creating new organization knowledge
type CreateOrganizationKnowledgeInput struct {
	OrganizationID   uuid.UUID          `json:"organization_id"`
	OrganizationCode string             `json:"organization_code"`
	KnowledgeType    int                `json:"knowledge_type" binding:"required"`
	Title            string             `json:"title" binding:"required"`
	Content          string             `json:"content" binding:"required"`
	Description      *string            `json:"description"`
	Metadata         *KnowledgeMetadata `json:"metadata"`
	IsActive         bool               `json:"is_active"`
}

// UpdateOrganizationKnowledgeInput is used for updating organization knowledge
type UpdateOrganizationKnowledgeInput struct {
	KnowledgeType int                `json:"knowledge_type"`
	Title         *string            `json:"title"`
	Content       string             `json:"content"`
	Description   *string            `json:"description"`
	Metadata      *KnowledgeMetadata `json:"metadata"`
	IsActive      *bool              `json:"is_active"`
}

// KnowledgeMetadata represents metadata for knowledge entries
type KnowledgeMetadata struct {
	Tags        []string `json:"tags,omitempty"`
	Category    string   `json:"category,omitempty"`
	Author      string   `json:"author,omitempty"`
	Version     string   `json:"version,omitempty"`
	Language    string   `json:"language,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
	LastUpdated string   `json:"last_updated,omitempty"`
	// Add more metadata fields as needed
}

// KnowledgeType constants
const (
	KnowledgeTypeProduct = 1
	KnowledgeTypeService = 2
	KnowledgeTypeFAQ     = 3
	KnowledgeTypeOther   = 4
)

// FileUploadResult represents the result of a file upload
type FileUploadResult struct {
	FileID      uuid.UUID `json:"file_id"`
	FileName    string    `json:"file_name"`
	FileSize    int64     `json:"file_size"`
	ContentType string    `json:"content_type"`
	FileURL     string    `json:"file_url"`
	BucketName  string    `json:"bucket_name"`
	ObjectName  string    `json:"object_name"`
	UploadedAt  time.Time `json:"uploaded_at"`
}
