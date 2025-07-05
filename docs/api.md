# API Documentation

This document provides comprehensive documentation for all API endpoints in the Vasst Communication Agent.

## Base URLs

- **v0 API**: `/v0`
- **v1 API**: `/v1`

## Authentication

Most endpoints require authentication using Bearer tokens. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## API Endpoints

### v0 API Endpoints

#### User Management (Superadmin Only)

##### Get All Users
```
GET /v0/admin/users
```

**Query Parameters:**
- `limit` (optional): Number of users to return (default: 10)
- `offset` (optional): Number of users to skip (default: 0)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "user_id": "uuid",
      "username": "string",
      "full_name": "string",
      "phone_number": "string",
      "organization_id": "uuid",
      "role_id": 1,
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  ]
}
```

##### Get User by ID
```
GET /v0/admin/users/{id}
```

**Path Parameters:**
- `id`: User UUID

**Response:**
```json
{
  "success": true,
  "data": {
    "user_id": "uuid",
    "username": "string",
    "full_name": "string",
    "phone_number": "string",
    "organization_id": "uuid",
    "role_id": 1,
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

##### Create User
```
POST /v0/admin/users
```

**Request Body:**
```json
{
  "phone_number": "string",
  "full_name": "string",
  "username": "string",
  "organization_id": "uuid",
  "role_id": 1
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user_id": "uuid",
    "username": "string",
    "full_name": "string",
    "phone_number": "string",
    "organization_id": "uuid",
    "role_id": 1,
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

##### Update User
```
PUT /v0/admin/users/{id}
```

**Path Parameters:**
- `id`: User UUID

**Request Body:**
```json
{
  "phone_number": "string",
  "full_name": "string",
  "username": "string",
  "role_id": 1
}
```

##### Delete User
```
DELETE /v0/admin/users/{id}
```

**Path Parameters:**
- `id`: User UUID

##### Get User by Phone Number
```
GET /v0/admin/users/phone/{phone}
```

**Path Parameters:**
- `phone`: Phone number

##### Get User by Username
```
GET /v0/admin/users/username/{username}
```

**Path Parameters:**
- `username`: Username

##### Reset User Password
```
POST /v0/admin/users/{id}/reset-password
```

**Path Parameters:**
- `id`: User UUID

**Request Body:**
```json
{
  "old_password": "string",
  "new_password": "string"
}
```

##### Generate User Password
```
POST /v0/admin/users/{id}/generate-password
```

**Path Parameters:**
- `id`: User UUID

**Request Body:**
```json
{
  "password": "string"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "password": "generated_password"
  }
}
```

##### Login
```
POST /v0/admin/users/login
```

**Request Body:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "jwt_token",
    "user": {
      "user_id": "uuid",
      "username": "string",
      "full_name": "string",
      "organization_id": "uuid",
      "role_id": 1
    }
  }
}
```

#### Organization Management

##### Get All Organizations
```
GET /v0/organizations
```

**Query Parameters:**
- `limit` (optional): Number of organizations to return (default: 10)
- `offset` (optional): Number of organizations to skip (default: 0)

##### Get Organization by ID
```
GET /v0/organizations/{id}
```

**Path Parameters:**
- `id`: Organization UUID

##### Get Organization by Code
```
GET /v0/organizations/code/{code}
```

**Path Parameters:**
- `code`: Organization code

##### Create Organization
```
POST /v0/organizations
```

**Request Body:**
```json
{
  "name": "string",
  "code": "string",
  "description": "string",
  "category_id": 1
}
```

##### Update Organization
```
PUT /v0/organizations/{id}
```

**Path Parameters:**
- `id`: Organization UUID

**Request Body:**
```json
{
  "name": "string",
  "code": "string",
  "description": "string",
  "category_id": 1
}
```

##### Delete Organization
```
DELETE /v0/organizations/{id}
```

**Path Parameters:**
- `id`: Organization UUID

#### Organization Categories

##### Get All Categories
```
GET /v0/organizations/categories
```

##### Get Category by ID
```
GET /v0/organizations/categories/{id}
```

**Path Parameters:**
- `id`: Category ID (integer)

##### Create Category
```
POST /v0/organizations/categories
```

**Request Body:**
```json
{
  "name": "string",
  "description": "string"
}
```

##### Update Category
```
PUT /v0/organizations/categories/{id}
```

**Path Parameters:**
- `id`: Category ID (integer)

**Request Body:**
```json
{
  "name": "string",
  "description": "string"
}
```

##### Delete Category
```
DELETE /v0/organizations/categories/{id}
```

**Path Parameters:**
- `id`: Category ID (integer)

#### Organization Settings

##### Get Organization Settings
```
GET /v0/organizations/{id}/settings
```

**Path Parameters:**
- `id`: Organization UUID

##### Update Organization Settings
```
PUT /v0/organizations/{id}/settings
```

**Path Parameters:**
- `id`: Organization UUID

**Request Body:**
```json
{
  "ai_enabled": true,
  "ai_model": "string",
  "ai_provider": "string",
  "ai_settings": {}
}
```

#### Organization Knowledge

##### Get Organization Knowledge
```
GET /v0/organizations/{id}/knowledge
```

**Path Parameters:**
- `id`: Organization UUID

##### Get Knowledge by ID
```
GET /v0/organizations/knowledge/{id}
```

**Path Parameters:**
- `id`: Knowledge UUID

##### Create Knowledge
```
POST /v0/organizations/{id}/knowledge
```

**Path Parameters:**
- `id`: Organization UUID

**Request Body:**
```json
{
  "title": "string",
  "content": "string",
  "category": "string",
  "tags": ["string"]
}
```

##### Update Knowledge
```
PUT /v0/organizations/knowledge/{id}
```

**Path Parameters:**
- `id`: Knowledge UUID

**Request Body:**
```json
{
  "title": "string",
  "content": "string",
  "category": "string",
  "tags": ["string"]
}
```

##### Delete Knowledge
```
DELETE /v0/organizations/knowledge/{id}
```

**Path Parameters:**
- `id`: Knowledge UUID

#### Organization Models

##### Get Organization Models
```
GET /v0/organizations/{id}/models
```

**Path Parameters:**
- `id`: Organization UUID

#### Organization Integrations

##### Get Organization Integrations
```
GET /v0/organizations/{id}/integrations
```

**Path Parameters:**
- `id`: Organization UUID

##### Get Integration by ID
```
GET /v0/organizations/integrations/{id}
```

**Path Parameters:**
- `id`: Integration UUID

##### Create Integration
```
POST /v0/organizations/{id}/integrations
```

**Path Parameters:**
- `id`: Organization UUID

**Request Body:**
```json
{
  "type": "string",
  "name": "string",
  "config": {},
  "is_enabled": true,
  "is_ai_enabled": true
}
```

##### Update Integration
```
PUT /v0/organizations/integrations/{id}
```

**Path Parameters:**
- `id`: Integration UUID

**Request Body:**
```json
{
  "type": "string",
  "name": "string",
  "config": {},
  "is_enabled": true,
  "is_ai_enabled": true
}
```

##### Delete Integration
```
DELETE /v0/organizations/integrations/{id}
```

**Path Parameters:**
- `id`: Integration UUID

#### WhatsApp Integration

##### Verify Webhook
```
GET /v0/whatsapp/webhook
```

**Query Parameters:**
- `hub.mode`: "subscribe"
- `hub.verify_token`: Verification token
- `hub.challenge`: Challenge string
- `key`: Organization key

##### Handle Webhook
```
POST /v0/whatsapp/webhook
```

**Query Parameters:**
- `key`: Organization key

**Request Body:**
```json
{
  "object": "whatsapp_business_account",
  "entry": [
    {
      "id": "string",
      "changes": [
        {
          "field": "messages",
          "value": {
            "messaging_product": "whatsapp",
            "metadata": {},
            "messages": [
              {
                "from": "string",
                "id": "string",
                "timestamp": "string",
                "text": {
                  "body": "string"
                }
              }
            ]
          }
        }
      ]
    }
  ]
}
```

##### Send Template Message
```
POST /v0/whatsapp/send/template
```

**Request Body:**
```json
{
  "to": "string",
  "template_name": "string",
  "language_code": "string"
}
```

##### Send Text Message
```
POST /v0/whatsapp/send/text
```

**Request Body (Simple Format):**
```json
{
  "to": "string",
  "message": "string"
}
```

**Request Body (WhatsApp API Format):**
```json
{
  "messaging_product": "whatsapp",
  "recipient_type": "individual",
  "to": "string",
  "type": "text",
  "text": {
    "preview_url": true,
    "body": "string"
  }
}
```


### v1 API Endpoints

#### Organization Management

##### Get Organization by ID
```
GET /v1/org/{id}
```

**Path Parameters:**
- `id`: Organization UUID

##### Update Organization
```
PUT /v1/org/{id}
```

**Path Parameters:**
- `id`: Organization UUID

**Request Body:**
```json
{
  "name": "string",
  "code": "string",
  "description": "string",
  "category_id": 1
}
```

#### Organization Settings

##### Get Organization Settings
```
GET /v1/org/{id}/settings
```

**Path Parameters:**
- `id`: Organization UUID

##### Update Organization Settings
```
PUT /v1/org/{id}/settings
```

**Path Parameters:**
- `id`: Organization UUID

**Request Body:**
```json
{
  "ai_enabled": true,
  "ai_model": "string",
  "ai_provider": "string",
  "ai_settings": {}
}
```

#### Organization Knowledge

##### Get Organization Knowledge
```
GET /v1/org/{id}/knowledge
```

**Path Parameters:**
- `id`: Organization UUID

##### Get Knowledge by ID
```
GET /v1/org/knowledge/{id}
```

**Path Parameters:**
- `id`: Knowledge UUID

##### Create Knowledge
```
POST /v1/org/{id}/knowledge
```

**Path Parameters:**
- `id`: Organization UUID

**Request Body:**
```json
{
  "title": "string",
  "content": "string",
  "category": "string",
  "tags": ["string"]
}
```

##### Update Knowledge
```
PUT /v1/org/knowledge/{id}
```

**Path Parameters:**
- `id`: Knowledge UUID

**Request Body:**
```json
{
  "title": "string",
  "content": "string",
  "category": "string",
  "tags": ["string"]
}
```

##### Delete Knowledge
```
DELETE /v1/org/knowledge/{id}
```

**Path Parameters:**
- `id`: Knowledge UUID

#### Organization Models

##### Get Organization Models
```
GET /v1/org/{id}/models
```

**Path Parameters:**
- `id`: Organization UUID

#### Organization Integrations

##### Get Organization Integrations
```
GET /v1/org/{id}/integrations
```

**Path Parameters:**
- `id`: Organization UUID

##### Get Integration by ID
```
GET /v1/org/integrations/{id}
```

**Path Parameters:**
- `id`: Integration UUID

##### Create Integration
```
POST /v1/org/{id}/integrations
```

**Path Parameters:**
- `id`: Organization UUID

**Request Body:**
```json
{
  "type": "string",
  "name": "string",
  "config": {},
  "is_enabled": true,
  "is_ai_enabled": true
}
```

##### Update Integration
```
PUT /v1/org/integrations/{id}
```

**Path Parameters:**
- `id`: Integration UUID

**Request Body:**
```json
{
  "type": "string",
  "name": "string",
  "config": {},
  "is_enabled": true,
  "is_ai_enabled": true
}
```

##### Delete Integration
```
DELETE /v1/org/integrations/{id}
```

**Path Parameters:**
- `id`: Integration UUID

#### Contact Management

##### Get All Contacts
```
GET /v1/contacts
```

**Query Parameters:**
- `limit` (optional): Number of contacts to return (default: 10)
- `offset` (optional): Number of contacts to skip (default: 0)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "contact_id": "uuid",
      "phone_number": "string",
      "full_name": "string",
      "email": "string",
      "organization_id": "uuid",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  ]
}
```

##### Get Contact by ID
```
GET /v1/contacts/{id}
```

**Path Parameters:**
- `id`: Contact UUID

**Response:**
```json
{
  "success": true,
  "data": {
    "contact_id": "uuid",
    "phone_number": "string",
    "full_name": "string",
    "email": "string",
    "organization_id": "uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

##### Get Contact by Phone Number
```
GET /v1/contacts/phone/{phone}
```

**Path Parameters:**
- `phone`: Phone number

##### Get Contacts by Organization
```
GET /v1/contacts/organization/{id}
```

**Path Parameters:**
- `id`: Organization UUID

**Query Parameters:**
- `limit` (optional): Number of contacts to return (default: 10)
- `offset` (optional): Number of contacts to skip (default: 0)

##### Create Contact
```
POST /v1/contacts
```

**Request Body:**
```json
{
  "phone_number": "string",
  "full_name": "string",
  "email": "string"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "contact_id": "uuid",
    "phone_number": "string",
    "full_name": "string",
    "email": "string",
    "organization_id": "uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

##### Update Contact
```
PUT /v1/contacts/{id}
```

**Path Parameters:**
- `id`: Contact UUID

**Request Body:**
```json
{
  "phone_number": "string",
  "full_name": "string",
  "email": "string"
}
```

##### Delete Contact
```
DELETE /v1/contacts/{id}
```

**Path Parameters:**
- `id`: Contact UUID

#### Conversation Management

##### Get Conversations by Organization
```
GET /v1/conversations
```

**Query Parameters:**
- `limit` (optional): Number of conversations to return (default: 10)
- `offset` (optional): Number of conversations to skip (default: 0)
- `status` (optional): Filter by status (0: Open, 1: Closed, 2: Pending, 3: Resolved)
- `priority` (optional): Filter by priority (0: Low, 1: Medium, 2: High, 3: Urgent)
- `is_active` (optional): Filter by active status (true/false)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "conversation_id": "uuid",
      "user_id": "uuid",
      "contact_id": "uuid",
      "medium_id": 1,
      "status": 0,
      "priority": 1,
      "is_active": true,
      "organization_id": "uuid",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  ]
}
```

##### Get Conversation Count by Organization
```
GET /v1/conversations/count
```

**Response:**
```json
{
  "success": true,
  "data": {
    "count": 42
  }
}
```

##### Get Conversation by ID
```
GET /v1/conversations/{id}
```

**Path Parameters:**
- `id`: Conversation UUID

##### Get Conversation Detail with Messages
```
GET /v1/conversations/{id}/detail
```

**Path Parameters:**
- `id`: Conversation UUID

**Query Parameters:**
- `message_limit` (optional): Message limit (default: 50)
- `message_offset` (optional): Message offset (default: 0)

**Response:**
```json
{
  "success": true,
  "data": {
    "conversation": {
      "conversation_id": "uuid",
      "user_id": "uuid",
      "contact_id": "uuid",
      "medium_id": 1,
      "status": 0,
      "priority": 1,
      "is_active": true,
      "organization_id": "uuid",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    },
    "messages": [
      {
        "message_id": "uuid",
        "conversation_id": "uuid",
        "sender_type_id": 1,
        "sender_id": "uuid",
        "direction": "i",
        "message_type_id": 1,
        "content": "string",
        "status": 1,
        "created_at": "timestamp"
      }
    ]
  }
}
```

##### Create Conversation
```
POST /v1/conversations
```

**Request Body:**
```json
{
  "user_id": "uuid",
  "contact_id": "uuid",
  "medium_id": 1,
  "status": 0,
  "priority": 1,
  "is_active": true
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "conversation_id": "uuid",
    "user_id": "uuid",
    "contact_id": "uuid",
    "medium_id": 1,
    "status": 0,
    "priority": 1,
    "is_active": true,
    "organization_id": "uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

##### Update Conversation
```
PUT /v1/conversations/{id}
```

**Path Parameters:**
- `id`: Conversation UUID

**Request Body:**
```json
{
  "status": 1,
  "priority": 2,
  "is_active": false
}
```

##### Delete Conversation
```
DELETE /v1/conversations/{id}
```

**Path Parameters:**
- `id`: Conversation UUID

##### Get Conversations by User ID
```
GET /v1/conversations/user/{user_id}
```

**Path Parameters:**
- `user_id`: User UUID

**Query Parameters:**
- `limit` (optional): Number of conversations to return (default: 10)
- `offset` (optional): Number of conversations to skip (default: 0)

##### Get Conversations by Contact ID
```
GET /v1/conversations/contact/{contact_id}
```

**Path Parameters:**
- `contact_id`: Contact UUID

**Query Parameters:**
- `limit` (optional): Number of conversations to return (default: 10)
- `offset` (optional): Number of conversations to skip (default: 0)

##### Get Conversations by Status
```
GET /v1/conversations/status/{status}
```

**Path Parameters:**
- `status`: Status value (0: Open, 1: Closed, 2: Pending, 3: Resolved)

**Query Parameters:**
- `limit` (optional): Number of conversations to return (default: 10)
- `offset` (optional): Number of conversations to skip (default: 0)

##### Get Conversations by Priority
```
GET /v1/conversations/priority/{priority}
```

**Path Parameters:**
- `priority`: Priority value (0: Low, 1: Medium, 2: High, 3: Urgent)

**Query Parameters:**
- `limit` (optional): Number of conversations to return (default: 10)
- `offset` (optional): Number of conversations to skip (default: 0)

##### Get Active Conversation
```
GET /v1/conversations/active/{user_id}/{contact_id}/{medium_id}
```

**Path Parameters:**
- `user_id`: User UUID
- `contact_id`: Contact UUID
- `medium_id`: Medium ID (integer)

#### Message Management

##### Get Messages by Organization
```
GET /v1/messages
```

**Query Parameters:**
- `limit` (optional): Number of messages to return (default: 50)
- `offset` (optional): Number of messages to skip (default: 0)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "message_id": "uuid",
      "conversation_id": "uuid",
      "organization_id": "uuid",
      "sender_type_id": 1,
      "sender_id": "uuid",
      "direction": "i",
      "message_type_id": 1,
      "content": "string",
      "status": 1,
      "ai_generated": false,
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  ]
}
```

##### Get Message by ID
```
GET /v1/messages/{id}
```

**Path Parameters:**
- `id`: Message UUID

##### Create Message
```
POST /v1/messages
```

**Request Body:**
```json
{
  "conversation_id": "uuid",
  "sender_type_id": 1,
  "sender_id": "uuid",
  "direction": "o",
  "message_type_id": 1,
  "content": "string",
  "is_broadcast": false,
  "is_order_message": false,
  "ai_generated": false
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message_id": "uuid",
    "conversation_id": "uuid",
    "organization_id": "uuid",
    "sender_type_id": 1,
    "sender_id": "uuid",
    "direction": "o",
    "message_type_id": 1,
    "content": "string",
    "status": 0,
    "ai_generated": false,
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

##### Create Message with Media
```
POST /v1/messages/with-media
```

**Content-Type:** `multipart/form-data`

**Form Data:**
- `file`: Media file (required)
- `conversation_id`: Conversation UUID (required)
- `organization_id`: Organization UUID (required)
- `direction`: Message direction "i" or "o" (required)
- `message_type_id`: Message type ID (required)
- `content`: Message content (optional)
- `is_broadcast`: Is broadcast message (optional, default: false)
- `is_order_message`: Is order message (optional, default: false)

##### Update Message
```
PUT /v1/messages/{id}
```

**Path Parameters:**
- `id`: Message UUID

**Request Body:**
```json
{
  "content": "string",
  "status": 1,
  "ai_generated": true
}
```

##### Update Message Status
```
PUT /v1/messages/{id}/status
```

**Path Parameters:**
- `id`: Message UUID

**Request Body:**
```json
{
  "status": 1
}
```

##### Delete Message
```
DELETE /v1/messages/{id}
```

**Path Parameters:**
- `id`: Message UUID

##### Get Messages by Conversation ID
```
GET /v1/messages/conversation/{conversation_id}
```

**Path Parameters:**
- `conversation_id`: Conversation UUID

**Query Parameters:**
- `limit` (optional): Number of messages to return (default: 50)
- `offset` (optional): Number of messages to skip (default: 0)

##### Get Pending Messages
```
GET /v1/messages/pending
```

**Query Parameters:**
- `limit` (optional): Number of messages to return (default: 50)
- `offset` (optional): Number of messages to skip (default: 0)

##### Get Messages by Status
```
GET /v1/messages/status/{status}
```

**Path Parameters:**
- `status`: Message status (0: Pending, 1: Sent, 2: Delivered, 3: Read, 4: Failed)

**Query Parameters:**
- `limit` (optional): Number of messages to return (default: 50)
- `offset` (optional): Number of messages to skip (default: 0)

#### User Management (Organization-scoped)

##### Get Users by Organization
```
GET /v1/users
```

**Query Parameters:**
- `limit` (optional): Number of users to return (default: 10)
- `offset` (optional): Number of users to skip (default: 0)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "user_id": "uuid",
      "username": "string",
      "full_name": "string",
      "phone_number": "string",
      "organization_id": "uuid",
      "role_id": 1,
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  ]
}
```

##### Get User by ID and Organization
```
GET /v1/users/{id}
```

**Path Parameters:**
- `id`: User UUID

##### Get User by Phone Number and Organization
```
GET /v1/users/phone/{phone}
```

**Path Parameters:**
- `phone`: Phone number

##### Get User by Username and Organization
```
GET /v1/users/username/{username}
```

**Path Parameters:**
- `username`: Username

##### Login
```
POST /v1/users/login
```

**Request Body:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "jwt_token",
    "user": {
      "user_id": "uuid",
      "username": "string",
      "full_name": "string",
      "organization_id": "uuid",
      "role_id": 1
    }
  }
}
```

## Error Responses

All endpoints return errors in the following format:

```json
{
  "success": false,
  "error": "Error message description"
}
```

## Common HTTP Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict (e.g., duplicate entry)
- `500 Internal Server Error`: Server error

## Rate Limiting

API requests are subject to rate limiting. When rate limited, the API will return:

```json
{
  "success": false,
  "error": "Rate limit exceeded"
}
```

## Pagination

For endpoints that return lists, pagination is supported using `limit` and `offset` query parameters. The response will include pagination metadata when applicable.

## Webhooks

Webhook endpoints (like WhatsApp) expect specific formats and require proper authentication. Always include the organization key as a query parameter for webhook endpoints. 