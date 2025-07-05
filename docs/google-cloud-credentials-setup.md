# Google Cloud Credentials Setup Guide

This guide explains how to securely set up Google Cloud credentials for the Vasst Communication Agent.

## 1. Create a Google Cloud Service Account

### Step 1: Go to Google Cloud Console
1. Visit [Google Cloud Console](https://console.cloud.google.com/)
2. Select your project or create a new one
3. Navigate to "IAM & Admin" > "Service Accounts"

### Step 2: Create Service Account
1. Click "Create Service Account"
2. Fill in the details:
   - **Name**: `vasst-storage-service`
   - **Description**: `Service account for Vasst Communication Agent file storage`
3. Click "Create and Continue"

### Step 3: Assign Roles
Add the following roles:
- **Storage Object Admin** (`roles/storage.objectAdmin`)
- **Storage Object Creator** (`roles/storage.objectCreator`)
- **Storage Object Viewer** (`roles/storage.objectViewer`)

### Step 4: Create and Download Key
1. Click "Done"
2. Find your service account in the list
3. Click the three dots (⋮) > "Manage keys"
4. Click "Add Key" > "Create new key"
5. Choose "JSON" format
6. Click "Create"
7. **Important**: The JSON file will download automatically. Keep it secure!

## 2. Secure Storage Locations

### Option A: Outside Project Directory (Recommended)
```bash
# Create a secure credentials directory
mkdir -p ~/.credentials/vasst-expense-api

# Move the downloaded JSON file
mv ~/Downloads/your-project-123456-abc123.json ~/.credentials/vasst-expense-api/google-cloud-credentials.json

# Set proper permissions
chmod 600 ~/.credentials/vasst-expense-api/google-cloud-credentials.json
```

### Option B: System-wide Configuration
```bash
# macOS/Linux
sudo mkdir -p /etc/vasst-expense-api/credentials
sudo mv ~/Downloads/your-project-123456-abc123.json /etc/vasst-expense-api/credentials/google-cloud-credentials.json
sudo chown $USER:$USER /etc/vasst-expense-api/credentials/google-cloud-credentials.json
sudo chmod 600 /etc/vasst-expense-api/credentials/google-cloud-credentials.json
```

### Option C: User Config Directory
```bash
# Create user config directory
mkdir -p ~/.config/vasst-expense-api
mv ~/Downloads/your-project-123456-abc123.json ~/.config/vasst-expense-api/google-cloud-credentials.json
chmod 600 ~/.config/vasst-expense-api/google-cloud-credentials.json
```

## 3. Environment Configuration

### Create .env File
Create a `.env` file in your project root:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=vasst_ca
DB_SSL_MODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRATION=24h

# Google Cloud Storage Configuration
GOOGLE_CLOUD_PROJECT_ID=your-project-id
GOOGLE_CLOUD_BUCKET_PREFIX=vasst-comm
GOOGLE_CLOUD_CREDENTIALS_FILE=/Users/valent/.credentials/vasst-expense-api/google-cloud-credentials.json

# WhatsApp Configuration (if using)
WHATSAPP_PHONE_NUMBER_ID=your-phone-number-id
WHATSAPP_ACCESS_TOKEN=your-access-token
WHATSAPP_BASE_URL=https://graph.facebook.com/v18.0

# OpenAI Configuration (if using)
OPENAI_API_KEY=your-openai-api-key
OPENAI_MODEL=gpt-3.5-turbo

# Gemini Configuration (if using)
GEMINI_API_KEY=your-gemini-api-key
GEMINI_MODEL=gemini-pro
```

### Update .gitignore
Make sure your `.gitignore` includes:

```gitignore
# Environment files
.env
.env.local
.env.production
.env.staging

# Credentials
*.json
credentials/
.credentials/

# Logs
logs/
*.log

# Build artifacts
vasst-expense-api
dist/
build/

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db
```

## 4. Path Examples for Different Systems

### macOS
```bash
# Option 1: User credentials directory
GOOGLE_CLOUD_CREDENTIALS_FILE=/Users/valent/.credentials/vasst-expense-api/google-cloud-credentials.json

# Option 2: System-wide
GOOGLE_CLOUD_CREDENTIALS_FILE=/etc/vasst-expense-api/credentials/google-cloud-credentials.json

# Option 3: User config
GOOGLE_CLOUD_CREDENTIALS_FILE=/Users/valent/.config/vasst-expense-api/google-cloud-credentials.json
```

### Linux
```bash
# Option 1: User credentials directory
GOOGLE_CLOUD_CREDENTIALS_FILE=/home/username/.credentials/vasst-expense-api/google-cloud-credentials.json

# Option 2: System-wide
GOOGLE_CLOUD_CREDENTIALS_FILE=/etc/vasst-expense-api/credentials/google-cloud-credentials.json

# Option 3: User config
GOOGLE_CLOUD_CREDENTIALS_FILE=/home/username/.config/vasst-expense-api/google-cloud-credentials.json
```

### Windows
```bash
# Option 1: User credentials directory
GOOGLE_CLOUD_CREDENTIALS_FILE=C:\Users\username\.credentials\vasst-expense-api\google-cloud-credentials.json

# Option 2: System-wide
GOOGLE_CLOUD_CREDENTIALS_FILE=C:\ProgramData\vasst-expense-api\credentials\google-cloud-credentials.json

# Option 3: User config
GOOGLE_CLOUD_CREDENTIALS_FILE=C:\Users\username\AppData\Local\vasst-expense-api\google-cloud-credentials.json
```

## 5. Docker Configuration

If running in Docker, you can mount the credentials:

### Docker Compose Example
```yaml
version: '3.8'
services:
  vasst-expense-api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GOOGLE_CLOUD_PROJECT_ID=your-project-id
      - GOOGLE_CLOUD_BUCKET_PREFIX=vasst-comm
      - GOOGLE_CLOUD_CREDENTIALS_FILE=/app/credentials/google-cloud-credentials.json
    volumes:
      - ~/.credentials/vasst-expense-api:/app/credentials:ro
    env_file:
      - .env
```

### Docker Run Example
```bash
docker run -d \
  --name vasst-expense-api \
  -p 8080:8080 \
  -v ~/.credentials/vasst-expense-api:/app/credentials:ro \
  -e GOOGLE_CLOUD_CREDENTIALS_FILE=/app/credentials/google-cloud-credentials.json \
  vasst-expense-api:latest
```

## 6. Production Deployment

### Kubernetes Secret
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: google-cloud-credentials
type: Opaque
data:
  google-cloud-credentials.json: <base64-encoded-json-content>
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vasst-expense-api
spec:
  template:
    spec:
      containers:
      - name: vasst-expense-api
        image: vasst-expense-api:latest
        env:
        - name: GOOGLE_CLOUD_CREDENTIALS_FILE
          value: /app/credentials/google-cloud-credentials.json
        volumeMounts:
        - name: google-credentials
          mountPath: /app/credentials
          readOnly: true
      volumes:
      - name: google-credentials
        secret:
          secretName: google-cloud-credentials
```

### Environment Variables (Alternative)
Instead of using a file, you can set the credentials as an environment variable:

```bash
# Base64 encode the JSON content
cat ~/.credentials/vasst-expense-api/google-cloud-credentials.json | base64

# Set as environment variable
export GOOGLE_APPLICATION_CREDENTIALS_JSON="<base64-encoded-content>"
```

Then modify your application to decode and use the JSON content.

## 7. Security Best Practices

### File Permissions
```bash
# Set restrictive permissions
chmod 600 ~/.credentials/vasst-expense-api/google-cloud-credentials.json

# Set proper ownership
chown $USER:$USER ~/.credentials/vasst-expense-api/google-cloud-credentials.json
```

### Network Security
- Use HTTPS for all API communications
- Restrict service account permissions to minimum required
- Regularly rotate service account keys
- Monitor API usage and costs

### Access Control
- Limit who has access to the credentials file
- Use different service accounts for different environments
- Implement proper logging and monitoring

## 8. Testing the Configuration

### Test Connection
```bash
# Test with gcloud CLI
gcloud auth activate-service-account --key-file=~/.credentials/vasst-expense-api/google-cloud-credentials.json
gcloud config set project your-project-id
gsutil ls
```

### Test with Application
```bash
# Start your application
go run cmd/api/main.go

# Test file upload endpoint
curl -X POST "http://localhost:8080/v1/org/upload" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -F "file=@/path/to/test-file.pdf"
```

## 9. Troubleshooting

### Common Issues

1. **Permission Denied**
   ```bash
   # Check file permissions
   ls -la ~/.credentials/vasst-expense-api/google-cloud-credentials.json
   
   # Fix permissions
   chmod 600 ~/.credentials/vasst-expense-api/google-cloud-credentials.json
   ```

2. **File Not Found**
   ```bash
   # Verify file exists
   ls -la ~/.credentials/vasst-expense-api/
   
   # Check path in .env file
   cat .env | grep GOOGLE_CLOUD_CREDENTIALS_FILE
   ```

3. **Invalid Credentials**
   - Verify the JSON file is not corrupted
   - Check that the service account has the required permissions
   - Ensure the project ID is correct

4. **Bucket Creation Failed**
   - Verify billing is enabled for the project
   - Check that the service account has Storage Admin permissions
   - Ensure the bucket name follows Google Cloud naming conventions

### Debug Mode
Enable debug logging to get more information:

```bash
export LOG_LEVEL=debug
go run cmd/api/main.go
```

## 10. Cost Optimization

### Storage Classes
Consider using different storage classes for different file types:
- **Standard**: Frequently accessed files
- **Nearline**: Files accessed less than once per month
- **Coldline**: Files accessed less than once per quarter
- **Archive**: Long-term backup files

### Lifecycle Management
Set up automatic deletion of old files:
```bash
gsutil lifecycle set lifecycle.json gs://your-bucket-name
```

Example `lifecycle.json`:
```json
{
  "rule": [
    {
      "action": {"type": "Delete"},
      "condition": {
        "age": 365,
        "matchesStorageClass": ["STANDARD", "NEARLINE"]
      }
    }
  ]
}
```

## Support

If you encounter issues:
1. Check the Google Cloud Console for error messages
2. Review the application logs
3. Verify your service account permissions
4. Test with the gcloud CLI
5. Contact the development team 