# Google Cloud Storage Setup

This document provides instructions for setting up Google Cloud Storage integration for the Vasst Communication Agent.

## Prerequisites

1. Google Cloud Platform account
2. A Google Cloud project
3. Google Cloud Storage API enabled
4. Service account with appropriate permissions

## Setup Steps

### 1. Create a Google Cloud Project

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Note down your Project ID

### 2. Enable Google Cloud Storage API

1. In the Google Cloud Console, go to "APIs & Services" > "Library"
2. Search for "Cloud Storage"
3. Click on "Cloud Storage" and enable it

### 3. Create a Service Account

1. Go to "IAM & Admin" > "Service Accounts"
2. Click "Create Service Account"
3. Give it a name (e.g., "vasst-storage-service")
4. Add the following roles:
   - Storage Object Admin
   - Storage Object Creator
   - Storage Object Viewer
5. Create and download the JSON key file

### 4. Environment Configuration

Add the following environment variables to your `.env` file:

```bash
# Google Cloud Storage
GOOGLE_CLOUD_PROJECT_ID=your-project-id
GOOGLE_CLOUD_BUCKET_PREFIX=vasst-comm
GOOGLE_CLOUD_CREDENTIALS_FILE=path/to/service-account-key.json
```

### 5. Bucket Structure

The system automatically creates buckets with the following naming convention:

- **Organization Buckets**: `{prefix}-org-{organization-id-short}`
- **Conversation Buckets**: `{prefix}-conv-{organization-id-short}-{conversation-id-short}`

**Note**: UUIDs are truncated to stay within Google Cloud Storage's 63-character bucket name limit.

Example:
- Organization bucket: `vasst-comm-org-123e4567e89b12d3`
- Conversation bucket: `vasst-comm-conv-123e4567e89b-987fcdeb51a2`

### 6. File Upload Process

When a file is uploaded:

1. The system checks if the conversation bucket exists
2. If not, it creates the bucket automatically
3. Files are uploaded with unique names: `{original-name}-{timestamp}-{unique-id}.{extension}`
4. File metadata is stored in the message attachments
5. Public URLs are generated for file access

### 7. Supported File Types

The system supports the following file types:
- Images (image/*)
- Videos (video/*)
- Audio (audio/*)
- Documents (application/*)
- PDFs (application/pdf)

### 8. Security Considerations

1. **Service Account Permissions**: Use the principle of least privilege
2. **Bucket Access**: Buckets are created with default settings
3. **File Access**: Files are publicly accessible via generated URLs
4. **Credentials**: Store service account keys securely

### 9. Troubleshooting

#### Common Issues

1. **Authentication Error**: Check service account credentials and permissions
2. **Bucket Creation Failed**: Verify project ID and billing is enabled
3. **File Upload Failed**: Check file size limits and network connectivity

#### Debug Commands

```bash
# Test Google Cloud Storage connection
gcloud auth application-default login
gcloud config set project YOUR_PROJECT_ID
gsutil ls
```

### 10. Cost Optimization

1. **Storage Class**: Consider using different storage classes for different file types
2. **Lifecycle Management**: Set up automatic deletion of old files
3. **Monitoring**: Use Google Cloud Monitoring to track usage

## Integration with WhatsApp

When WhatsApp media messages are received:

1. Media is downloaded from WhatsApp API
2. Files are uploaded to the conversation bucket
3. Message is updated with file metadata
4. File URLs are stored for future access

## API Usage

The Google Cloud Storage service provides the following methods:

- `UploadFile()`: Upload multipart files
- `UploadFileFromBytes()`: Upload file bytes
- `GetFileURL()`: Get file access URL
- `DeleteFile()`: Delete files
- `CreateOrganizationBucket()`: Create organization bucket
- `CreateConversationBucket()`: Create conversation bucket 