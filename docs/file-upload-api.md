# File Upload API Documentation

This document describes the file upload functionality available in the Vasst Communication Agent API.

## Overview

The file upload API allows you to upload files to Google Cloud Storage and associate them with organizations, conversations, and knowledge base entries. Files are automatically organized into buckets based on organization and conversation IDs.

## Authentication

All file upload endpoints require authentication. Include your API key in the Authorization header:

```
Authorization: Bearer YOUR_API_KEY
```

## File Upload Endpoints

### 1. General File Upload

#### Single File Upload

**Endpoint:** `POST /v1/org/upload`  
**Content-Type:** `multipart/form-data`

Upload a single file to your organization's storage.

**Form Data:**
- `file` (required): The file to upload

**Response:**
```json
{
  "success": true,
  "data": {
    "file_id": "123e4567-e89b-12d3-a456-426614174000",
    "file_name": "document.pdf",
    "file_size": 1024000,
    "content_type": "application/pdf",
    "file_url": "https://storage.googleapis.com/vasst-comm-org-123e4567-e89b-12d3-a456-426614174000/document-20231201-143022-abc123.pdf",
    "bucket_name": "vasst-comm-org-123e4567-e89b-12d3-a456-426614174000",
    "object_name": "document-20231201-143022-abc123.pdf",
    "uploaded_at": "2023-12-01T14:30:22Z"
  }
}
```

**Example (cURL):**
```bash
curl -X POST "http://localhost:8080/v1/org/upload" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -F "file=@/path/to/document.pdf"
```

#### Multiple File Upload

**Endpoint:** `POST /v1/org/upload/multiple`  
**Content-Type:** `multipart/form-data`

Upload multiple files at once (up to 10 files).

**Form Data:**
- `files` (required): Array of files to upload

**Response:**
```json
{
  "success": true,
  "data": {
    "uploaded_files": [
      {
        "file_id": "123e4567-e89b-12d3-a456-426614174000",
        "file_name": "document1.pdf",
        "file_size": 1024000,
        "content_type": "application/pdf",
        "file_url": "https://storage.googleapis.com/...",
        "bucket_name": "vasst-comm-org-123e4567-e89b-12d3-a456-426614174000",
        "object_name": "document1-20231201-143022-abc123.pdf",
        "uploaded_at": "2023-12-01T14:30:22Z"
      }
    ],
    "total_files": 1,
    "success_count": 1,
    "error_count": 0
  }
}
```

**Example (cURL):**
```bash
curl -X POST "http://localhost:8080/v1/org/upload/multiple" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -F "files=@/path/to/document1.pdf" \
  -F "files=@/path/to/document2.pdf"
```

### 2. Knowledge Base with File Upload

#### Create Knowledge with File

**Endpoint:** `POST /v1/org/knowledge/with-file`  
**Content-Type:** `multipart/form-data`

Create a knowledge base entry with an attached file.

**Form Data:**
- `file` (required): The file to upload
- `knowledge_type` (required): Type of knowledge (1-4)
- `title` (required): Knowledge title
- `content` (required): Knowledge content
- `description` (optional): Knowledge description
- `is_active` (optional): Whether the knowledge is active (true/false)

**Knowledge Types:**
- `1`: Product
- `2`: Service
- `3`: FAQ
- `4`: Other

**Response:**
```json
{
  "success": true,
  "data": {
    "knowledge_id": "123e4567-e89b-12d3-a456-426614174000",
    "organization_id": "123e4567-e89b-12d3-a456-426614174000",
    "knowledge_type": 1,
    "title": "Product Manual",
    "content": "This is the product manual content",
    "description": "Product manual for customers",
    "source_url": "https://storage.googleapis.com/...",
    "file_name": "manual.pdf",
    "file_size": 1024000,
    "content_type": "application/pdf",
    "bucket_name": "vasst-comm-org-123e4567-e89b-12d3-a456-426614174000",
    "object_name": "manual-20231201-143022-abc123.pdf",
    "is_active": true,
    "created_at": "2023-12-01T14:30:22Z",
    "updated_at": "2023-12-01T14:30:22Z"
  }
}
```

**Example (cURL):**
```bash
curl -X POST "http://localhost:8080/v1/org/knowledge/with-file" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -F "file=@/path/to/manual.pdf" \
  -F "knowledge_type=1" \
  -F "title=Product Manual" \
  -F "content=This is the product manual content" \
  -F "description=Product manual for customers" \
  -F "is_active=true"
```

#### Update Knowledge with File

**Endpoint:** `PUT /v1/org/knowledge/:id/with-file`  
**Content-Type:** `multipart/form-data`

Update a knowledge base entry with a new file.

**Form Data:**
- `file` (required): The new file to upload
- `knowledge_type` (optional): Type of knowledge (1-4)
- `title` (optional): Knowledge title
- `content` (optional): Knowledge content
- `description` (optional): Knowledge description
- `is_active` (optional): Whether the knowledge is active (true/false)

**Example (cURL):**
```bash
curl -X PUT "http://localhost:8080/v1/org/knowledge/123e4567-e89b-12d3-a456-426614174000/with-file" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -F "file=@/path/to/updated-manual.pdf" \
  -F "title=Updated Product Manual" \
  -F "content=Updated content"
```

### 3. V0 API Endpoints

The V0 API provides similar functionality with different URL patterns:

- **Single File Upload:** `POST /v0/organizations/:id/upload`
- **Multiple File Upload:** `POST /v0/organizations/:id/upload/multiple`
- **Create Knowledge with File:** `POST /v0/organizations/:id/knowledge/with-file`
- **Update Knowledge with File:** `PUT /v0/organizations/knowledge/:id/with-file`

## File Validation

### Supported File Types

The following file types are supported:

- **Documents:** PDF, DOC, DOCX, TXT
- **Images:** JPEG, PNG, GIF
- **Videos:** MP4, AVI, MOV (for media messages)
- **Audio:** MP3, WAV, M4A (for media messages)

### File Size Limits

- **Single File:** Maximum 10MB
- **Multiple Files:** Maximum 10MB per file
- **Total Upload:** Maximum 32MB for multiple file uploads

### File Naming

Files are automatically renamed to prevent conflicts:
- Format: `{original-name}-{timestamp}-{unique-id}.{extension}`
- Example: `document-20231201-143022-abc123.pdf`

## Storage Organization

### Bucket Structure

Files are organized into Google Cloud Storage buckets:

- **Organization Files:** `{prefix}-org-{organization-id-short}`
- **Conversation Files:** `{prefix}-conv-{organization-id-short}-{conversation-id-short}`
- **Knowledge Files:** Stored in organization buckets

**Note**: UUIDs are truncated to stay within Google Cloud Storage's 63-character bucket name limit.

### File URLs

Files are accessible via public URLs:
```
https://storage.googleapis.com/{bucket-name}/{object-name}
```

## Error Handling

### Common Error Responses

**File Too Large:**
```json
{
  "success": false,
  "error": "file size too large, maximum 10MB allowed"
}
```

**Invalid File Type:**
```json
{
  "success": false,
  "error": "file type not allowed. Allowed types: PDF, JPEG, PNG, GIF, TXT, DOC, DOCX"
}
```

**Missing File:**
```json
{
  "success": false,
  "error": "file is required"
}
```

**Authentication Error:**
```json
{
  "success": false,
  "error": "organization ID not found in context"
}
```

## Best Practices

### 1. File Preparation

- Compress large files before upload
- Use appropriate file formats for your use case
- Ensure files are not corrupted

### 2. Error Handling

- Always check the response status
- Handle partial upload failures for multiple files
- Implement retry logic for network issues

### 3. Security

- Validate file types on the client side
- Don't trust file extensions alone
- Implement proper access controls

### 4. Performance

- Use multiple file upload for batch operations
- Consider file compression for large uploads
- Monitor upload progress for large files

## Integration Examples

### JavaScript/Node.js

```javascript
const FormData = require('form-data');
const fs = require('fs');

async function uploadFile(apiKey, filePath) {
  const form = new FormData();
  form.append('file', fs.createReadStream(filePath));

  const response = await fetch('http://localhost:8080/v1/org/upload', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${apiKey}`,
      ...form.getHeaders()
    },
    body: form
  });

  return response.json();
}
```

### Python

```python
import requests

def upload_file(api_key, file_path):
    url = 'http://localhost:8080/v1/org/upload'
    headers = {'Authorization': f'Bearer {api_key}'}
    
    with open(file_path, 'rb') as f:
        files = {'file': f}
        response = requests.post(url, headers=headers, files=files)
    
    return response.json()
```

### PHP

```php
function uploadFile($apiKey, $filePath) {
    $url = 'http://localhost:8080/v1/org/upload';
    $headers = ['Authorization: Bearer ' . $apiKey];
    
    $postData = ['file' => new CURLFile($filePath)];
    
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_POST, true);
    curl_setopt($ch, CURLOPT_POSTFIELDS, $postData);
    curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    
    $response = curl_exec($ch);
    curl_close($ch);
    
    return json_decode($response, true);
}
```

## Troubleshooting

### Common Issues

1. **Upload Fails with 401 Error**
   - Check your API key is valid
   - Ensure the organization exists and is active

2. **File Type Not Allowed**
   - Verify the file's MIME type
   - Check the allowed file types list

3. **File Size Too Large**
   - Compress the file before upload
   - Split large files if possible

4. **Network Timeout**
   - Check your internet connection
   - Implement retry logic
   - Consider using smaller files

### Debug Information

Enable debug logging to get more information about upload failures:

```bash
# Set debug level in your environment
export LOG_LEVEL=debug
```

## Rate Limits

- **Single File Upload:** 100 requests per minute per organization
- **Multiple File Upload:** 50 requests per minute per organization
- **Knowledge Upload:** 200 requests per minute per organization

## Support

For additional support or questions about the file upload API:

1. Check the API documentation
2. Review the error messages
3. Contact the development team
4. Check the system logs for detailed error information 