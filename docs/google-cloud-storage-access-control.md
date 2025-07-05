# Google Cloud Storage Access Control Strategy

## Overview

This document explains the access control strategy for Google Cloud Storage (GCS) files in the communication agent system, specifically designed for files that need to be accessible in system prompts and sent via WhatsApp.

## Access Control Strategy

### Public Read Access

Files uploaded to GCS are configured with **public read access** for the following reasons:

1. **System Prompt Accessibility**: Files need to be accessible via direct URLs when included in AI system prompts
2. **WhatsApp Integration**: WhatsApp requires publicly accessible URLs for media files
3. **Simplicity**: Eliminates the need for signed URLs or complex authentication for file access

### Security Considerations

While files are publicly readable, the following security measures are in place:

1. **Unique File Names**: Each uploaded file gets a unique name with timestamp and UUID
2. **Organization Isolation**: Files are stored in organization-specific buckets
3. **Conversation Isolation**: Media files are stored in conversation-specific buckets
4. **No Public Write Access**: Only the application can upload/delete files

## Configuration

### Environment Variables

Add the following environment variables to your `.env` file:

```bash
# Google Cloud Storage Configuration
GOOGLE_CLOUD_PROJECT_ID=your-project-id
GOOGLE_CLOUD_BUCKET_PREFIX=vasst
GOOGLE_CLOUD_CREDENTIALS_FILE=path/to/service-account-key.json
GOOGLE_CLOUD_REGION=ASIA-SOUTHEAST1
```

### Region Configuration

- **Default Region**: `ASIA-SOUTHEAST1` (Singapore)
- **Purpose**: Optimize latency for Southeast Asian users
- **Configurable**: Can be changed via `GOOGLE_CLOUD_REGION` environment variable

## Bucket Structure

### Organization Buckets
- **Pattern**: `{prefix}-org-{organization-code}`
- **Purpose**: Store knowledge content files
- **Access**: Public read for system prompts

### Conversation Buckets
- **Pattern**: `{prefix}-conv-{organization-code}-{conversation-id}`
- **Purpose**: Store media files (images, documents, etc.)
- **Access**: Public read for WhatsApp integration

## File Access Patterns

### For System Prompts
```go
// Files are accessible via direct URLs
fileURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)
```

### For WhatsApp
```go
// WhatsApp can directly access the public URLs
// No additional authentication required
```

## Alternative Access Control Options

If you need more restrictive access control in the future, consider these alternatives:

### 1. Signed URLs (Recommended for Production)
```go
// Generate signed URLs with expiration
url, err := obj.SignedURL(opts)
```

### 2. IAM-based Access
```go
// Use IAM roles and service accounts
// More complex but more secure
```

### 3. VPC Service Controls
```go
// Restrict access to specific VPC networks
// Requires additional GCP setup
```

## Best Practices

1. **Monitor Usage**: Set up Cloud Monitoring for bucket access
2. **Lifecycle Policies**: Configure automatic deletion of old files
3. **Cost Optimization**: Use appropriate storage classes
4. **Backup Strategy**: Consider cross-region replication for critical files

## Troubleshooting

### Common Issues

1. **Bucket Name Too Long**
   - Solution: Bucket names are automatically truncated to meet GCS limits

2. **Access Denied**
   - Check service account permissions
   - Verify bucket ACL settings

3. **Region Issues**
   - Ensure region is supported in your GCP project
   - Check for any regional restrictions

### Debugging

Enable debug logging in the Google Storage service:

```go
// Add debug prints in upload methods
fmt.Printf("Uploading to bucket: %s, object: %s\n", bucketName, objectName)
```

## Migration from Private to Public Access

If you need to change existing buckets from private to public access:

```bash
# Update bucket ACL
gsutil iam ch allUsers:objectViewer gs://your-bucket-name

# Update default object ACL
gsutil defacl ch -u AllUsers:R gs://your-bucket-name
```

## Security Recommendations

1. **Regular Audits**: Monitor bucket access logs
2. **File Validation**: Validate file types and sizes before upload
3. **Rate Limiting**: Implement upload rate limits
4. **Monitoring**: Set up alerts for unusual access patterns

## Cost Considerations

- **Storage**: Standard storage class for frequently accessed files
- **Network**: Egress costs for file downloads
- **Operations**: PUT/GET operations are charged per request
- **Lifecycle**: Consider moving old files to cheaper storage classes 