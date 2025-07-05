package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

//go:generate mockgen -source=google_storage_service.go -package=mock -destination=mock/google_storage_service_mock.go
type (
	GoogleStorageService interface {
		UploadFile(ctx context.Context, organizationCode string, conversationID uuid.UUID, file *multipart.FileHeader) (*FileUploadResult, error)
		UploadFileFromBytes(ctx context.Context, organizationCode string, conversationID uuid.UUID, fileName string, fileBytes []byte, contentType string) (*FileUploadResult, error)
		GetFileURL(ctx context.Context, organizationCode string, conversationID uuid.UUID, fileName string) (string, error)
		DeleteFile(ctx context.Context, organizationCode string, conversationID uuid.UUID, fileName string) error
		CreateOrganizationBucket(ctx context.Context, organizationCode string) error
		CreateConversationBucket(ctx context.Context, organizationCode string, conversationID uuid.UUID) error
	}

	googleStorageService struct {
		client       *storage.Client
		projectID    string
		bucketPrefix string
		region       string
	}

	FileUploadResult struct {
		FileName    string    `json:"file_name"`
		FileURL     string    `json:"file_url"`
		FileSize    int64     `json:"file_size"`
		ContentType string    `json:"content_type"`
		UploadedAt  time.Time `json:"uploaded_at"`
		BucketName  string    `json:"bucket_name"`
		ObjectName  string    `json:"object_name"`
	}
)

// NewGoogleStorageService creates a new Google Cloud Storage service
func NewGoogleStorageService(projectID, bucketPrefix string, credentialsFile string, region string) (GoogleStorageService, error) {
	ctx := context.Background()

	var client *storage.Client
	var err error

	if credentialsFile != "" {
		// Use service account credentials file
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	} else {
		// Use default credentials (Application Default Credentials)
		client, err = storage.NewClient(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	// Default to Singapore region if not specified
	if region == "" {
		region = "ASIA-SOUTHEAST1"
	}

	return &googleStorageService{
		client:       client,
		projectID:    projectID,
		bucketPrefix: bucketPrefix,
		region:       region,
	}, nil
}

// generateBucketName generates a bucket name for an organization
func (s *googleStorageService) generateOrganizationBucketName(organizationCode string) string {
	return fmt.Sprintf("%s-org-%s", s.bucketPrefix, strings.ToLower(organizationCode))
}

// generateConversationBucketName generates a bucket name for a conversation
func (s *googleStorageService) generateConversationBucketName(organizationCode string, conversationID uuid.UUID) string {
	return fmt.Sprintf("%s-conv-%s-%s", s.bucketPrefix, strings.ToLower(organizationCode), conversationID)
}

// generateObjectName generates an object name for a file
func (s *googleStorageService) generateObjectName(fileName string) string {
	timestamp := time.Now().Format("20060102-150405")
	ext := filepath.Ext(fileName)
	nameWithoutExt := strings.TrimSuffix(fileName, ext)
	uniqueID := uuid.New().String()[:8]

	return fmt.Sprintf("%s-%s-%s%s", nameWithoutExt, timestamp, uniqueID, ext)
}

// CreateOrganizationBucket creates a bucket for an organization
func (s *googleStorageService) CreateOrganizationBucket(ctx context.Context, organizationCode string) error {
	bucketName := s.generateOrganizationBucketName(organizationCode)

	// Check if bucket already exists
	bucket := s.client.Bucket(bucketName)
	_, err := bucket.Attrs(ctx)
	if err == nil {
		// Bucket already exists
		return nil
	}

	// Create the bucket with Singapore region and public access
	if err := bucket.Create(ctx, s.projectID, &storage.BucketAttrs{
		Location: s.region,
		Labels: map[string]string{
			"organization-code": strings.ToLower(organizationCode),
			"type":              "organization",
		},
		// Configure public access for files that need to be accessible
		// This allows files to be accessed via direct URLs
		UniformBucketLevelAccess: storage.UniformBucketLevelAccess{
			Enabled: false, // Disable uniform access to allow public access
		},
		// Set default object ACL to public read
		DefaultObjectACL: []storage.ACLRule{
			{Entity: storage.AllUsers, Role: storage.RoleReader},
		},
	}); err != nil {
		return fmt.Errorf("failed to create organization bucket %s: %w", bucketName, err)
	}

	return nil
}

// CreateConversationBucket creates a bucket for a conversation
func (s *googleStorageService) CreateConversationBucket(ctx context.Context, organizationCode string, conversationID uuid.UUID) error {
	bucketName := s.generateConversationBucketName(organizationCode, conversationID)

	// Check if bucket already exists
	bucket := s.client.Bucket(bucketName)
	_, err := bucket.Attrs(ctx)
	if err == nil {
		// Bucket already exists
		return nil
	}

	// Create the bucket with Singapore region and public access
	if err := bucket.Create(ctx, s.projectID, &storage.BucketAttrs{
		Location: s.region,
		Labels: map[string]string{
			"organization-code": strings.ToLower(organizationCode),
			"conversation-id":   conversationID.String(),
			"type":              "conversation",
		},
		// Configure public access for files that need to be accessible
		// This allows files to be accessed via direct URLs
		UniformBucketLevelAccess: storage.UniformBucketLevelAccess{
			Enabled: false, // Disable uniform access to allow public access
		},
		// Set default object ACL to public read
		DefaultObjectACL: []storage.ACLRule{
			{Entity: storage.AllUsers, Role: storage.RoleReader},
		},
	}); err != nil {
		return fmt.Errorf("failed to create conversation bucket %s: %w", bucketName, err)
	}

	return nil
}

// UploadFile uploads a file to Google Cloud Storage
func (s *googleStorageService) UploadFile(ctx context.Context, organizationCode string, conversationID uuid.UUID, file *multipart.FileHeader) (*FileUploadResult, error) {

	fmt.Println("organizationCode", organizationCode)
	fmt.Println("conversationID", conversationID)

	if conversationID == uuid.Nil {
		// Ensure organization bucket exists
		if err := s.CreateOrganizationBucket(ctx, organizationCode); err != nil {
			return nil, fmt.Errorf("failed to ensure organization bucket exists: %w", err)
		}
	} else {
		// Ensure conversation bucket exists
		if err := s.CreateConversationBucket(ctx, organizationCode, conversationID); err != nil {
			return nil, fmt.Errorf("failed to ensure conversation bucket exists: %w", err)
		}
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Generate bucket and object names
	var bucketName string
	if conversationID == uuid.Nil {
		bucketName = s.generateOrganizationBucketName(organizationCode)
	} else {
		bucketName = s.generateConversationBucketName(organizationCode, conversationID)
	}
	objectName := s.generateObjectName(file.Filename)

	// Get the bucket and create a new object
	bucket := s.client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	// Create a writer
	writer := obj.NewWriter(ctx)
	writer.ContentType = file.Header.Get("Content-Type")

	// Copy the file data to the writer
	if _, err := io.Copy(writer, src); err != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", err)
	}

	// Close the writer to finalize the upload
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize upload: %w", err)
	}

	// Set public read access for the uploaded file
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, fmt.Errorf("failed to set public access for file: %w", err)
	}

	// Get the uploaded object's attributes
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get object attributes: %w", err)
	}

	// Generate the public URL
	fileURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)

	return &FileUploadResult{
		FileName:    file.Filename,
		FileURL:     fileURL,
		FileSize:    attrs.Size,
		ContentType: attrs.ContentType,
		UploadedAt:  attrs.Created,
		BucketName:  bucketName,
		ObjectName:  objectName,
	}, nil
}

// UploadFileFromBytes uploads a file from bytes to Google Cloud Storage
func (s *googleStorageService) UploadFileFromBytes(ctx context.Context, organizationCode string, conversationID uuid.UUID, fileName string, fileBytes []byte, contentType string) (*FileUploadResult, error) {
	// Ensure conversation bucket exists
	if err := s.CreateConversationBucket(ctx, organizationCode, conversationID); err != nil {
		return nil, fmt.Errorf("failed to ensure conversation bucket exists: %w", err)
	}

	// Generate bucket and object names
	bucketName := s.generateConversationBucketName(organizationCode, conversationID)
	objectName := s.generateObjectName(fileName)

	// Get the bucket and create a new object
	bucket := s.client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	// Create a writer
	writer := obj.NewWriter(ctx)
	writer.ContentType = contentType

	// Write the bytes to the writer
	if _, err := writer.Write(fileBytes); err != nil {
		return nil, fmt.Errorf("failed to write file data: %w", err)
	}

	// Close the writer to finalize the upload
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize upload: %w", err)
	}

	// Set public read access for the uploaded file
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, fmt.Errorf("failed to set public access for file: %w", err)
	}

	// Get the uploaded object's attributes
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get object attributes: %w", err)
	}

	// Generate the public URL
	fileURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)

	return &FileUploadResult{
		FileName:    fileName,
		FileURL:     fileURL,
		FileSize:    attrs.Size,
		ContentType: attrs.ContentType,
		UploadedAt:  attrs.Created,
		BucketName:  bucketName,
		ObjectName:  objectName,
	}, nil
}

// GetFileURL generates a URL for a file
func (s *googleStorageService) GetFileURL(ctx context.Context, organizationCode string, conversationID uuid.UUID, fileName string) (string, error) {
	bucketName := s.generateConversationBucketName(organizationCode, conversationID)

	// Check if the object exists
	bucket := s.client.Bucket(bucketName)
	obj := bucket.Object(fileName)

	_, err := obj.Attrs(ctx)
	if err != nil {
		return "", fmt.Errorf("file not found: %w", err)
	}

	// Generate the public URL
	fileURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, fileName)
	return fileURL, nil
}

// DeleteFile deletes a file from Google Cloud Storage
func (s *googleStorageService) DeleteFile(ctx context.Context, organizationCode string, conversationID uuid.UUID, fileName string) error {
	bucketName := s.generateConversationBucketName(organizationCode, conversationID)

	bucket := s.client.Bucket(bucketName)
	obj := bucket.Object(fileName)

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete file %s from bucket %s: %w", fileName, bucketName, err)
	}

	return nil
}

// Close closes the storage client
func (s *googleStorageService) Close() error {
	return s.client.Close()
}
