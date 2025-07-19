# VASST Expense API v1 Documentation

## Base URL
```
https://api.vasst.id/v1
```

## Authentication
Most endpoints require authentication using Bearer token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## Response Format
All API responses follow this standard format:
```json
{
  "success": true|false,
  "data": <response_data>,
  "message": "optional message",
  "error": "error message if success is false"
}
```

---

## Table of Contents
1. [Authentication](#authentication-endpoints)
2. [Users](#user-endpoints)
3. [Workspaces](#workspace-endpoints)
4. [Accounts](#account-endpoints)
5. [Categories](#category-endpoints)
6. [Banks](#bank-endpoints)
7. [Currencies](#currency-endpoints)
8. [Plans](#plan-endpoints)
9. [Budgets](#budget-endpoints)
10. [Transactions](#transaction-endpoints)
11. [Conversations](#conversation-endpoints)
12. [Messages](#message-endpoints)
13. [Taxonomies](#taxonomy-endpoints)
14. [User Tags](#user-tags-endpoints)
15. [Transaction Tags](#transaction-tags-endpoints)
16. [Verification Codes](#verification-code-endpoints)

---

## Authentication Endpoints

### Register User
**POST** `/auth/register`

Create a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "phone": "+1234567890",
  "first_name": "John",
  "last_name": "Doe",
  "password": "123456",
  "pin": "123456"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "phone": "+1234567890",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### Login
**POST** `/auth/login`

Authenticate user and get access token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "123456"
}
```
OR
```json
{
  "phone": "+1234567890",
  "password": "123456"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "phone": "+1234567890",
      "first_name": "John",
      "last_name": "Doe"
    },
    "token": "jwt_access_token"
  }
}
```

### Forgot Password
**POST** `/auth/forgot-password`

Initiate password reset process.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

### Reset Password
**POST** `/auth/reset-password`

Reset password using token.

**Request Body:**
```json
{
  "token": "reset_token",
  "new_password": "123456"
}
```

### Change Password
**POST** `/auth/change-password`

Change user's password (requires authentication).

**Request Body:**
```json
{
  "current_password": "123456",
  "new_password": "123456"
}
```

### Verify Phone
**POST** `/auth/verify-phone`

Verify phone number with code.

**Request Body:**
```json
{
  "phone": "+1234567890",
  "code": "123456"
}
```

### Resend Verification Code
**POST** `/auth/resend-verification-code`

Resend SMS verification code.

**Request Body:**
```json
{
  "phone": "+1234567890"
}
```

### Verify Email
**POST** `/auth/verify-email`

Verify email address with token.

**Request Body:**
```json
{
  "token": "email_verification_token"
}
```

### Resend Verification Email
**POST** `/auth/resend-verification-email`

Resend email verification link.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

---

## User Endpoints

### Get User Profile
**GET** `/users/`

Get authenticated user's profile information.

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "phone": "+1234567890",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### Update User Profile
**PUT** `/users/`

Update authenticated user's profile information.

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "newemail@example.com",
  "phone": "+1234567890"
}
```

---

## Workspace Endpoints

### List All Workspaces
**GET** `/workspaces`

Get all workspaces for the authenticated user.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `limit` (optional): Number of items per page (default: 10)
- `offset` (optional): Number of items to skip (default: 0)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "Personal Finance",
      "type": "personal",
      "currency_id": "uuid",
      "created_by": "uuid",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### Get Workspace by ID
**GET** `/workspaces/{id}`

Get a specific workspace by ID.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Workspace UUID

### Create Workspace
**POST** `/workspaces`

Create a new workspace.

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "Personal Finance",
  "type": "personal",
  "currency_id": "uuid",
  "description": "My personal finance workspace"
}
```

### Update Workspace
**PUT** `/workspaces/{id}`

Update an existing workspace.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Workspace UUID

**Request Body:**
```json
{
  "name": "Updated Workspace Name",
  "type": "business",
  "currency_id": "uuid",
  "description": "Updated description"
}
```

### Delete Workspace
**DELETE** `/workspaces/{id}`

Delete a workspace.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Workspace UUID

---

## Account Endpoints

### Get Accounts by User
**GET** `/accounts`

Get all accounts for the authenticated user.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `limit` (optional): Number of items per page (default: 10)
- `offset` (optional): Number of items to skip (default: 0)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "Main Bank Account",
      "type": "checking",
      "currency_id": "uuid",
      "balance": 1000.00,
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### Create Account
**POST** `/accounts`

Create a new account.

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "Savings Account",
  "type": "savings",
  "currency_id": "uuid",
  "initial_balance": 500.00
}
```

### Get Account by ID
**GET** `/accounts/{id}`

Get a specific account by ID.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Account UUID

### Update Account
**PUT** `/accounts/{id}`

Update an existing account.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Account UUID

**Request Body:**
```json
{
  "name": "Updated Account Name",
  "type": "checking",
  "currency_id": "uuid",
  "balance": 1500.00
}
```

### Delete Account
**DELETE** `/accounts/{id}`

Delete an account (soft delete).

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Account UUID

---

## Category Endpoints

### Get System Categories
**GET** `/system-categories`

Get all system categories.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 20)

### Create System Category
**POST** `/system-categories`

Create a new system category.

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "Food & Dining",
  "description": "Expenses related to food and dining",
  "icon": "🍽️",
  "color": "#FF6B6B"
}
```

### Get System Category by ID
**GET** `/system-categories/{id}`

Get a specific system category.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Category UUID

### Add System Category to User
**POST** `/system-categories/{id}/add-to-user`

Add a system category to the authenticated user's categories.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Category UUID

### Get User Categories
**GET** `/user-categories`

Get all user categories for the authenticated user.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 20)
- `search` (optional): Search term
- `sort_by` (optional): Sort field
- `sort_order` (optional): Sort order (asc/desc)

### Create User Category
**POST** `/user-categories`

Create a new user category.

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "Custom Category",
  "description": "My custom category",
  "icon": "🎯",
  "color": "#4ECDC4"
}
```

### Get Active User Categories
**GET** `/user-categories/active`

Get all active user categories.

**Headers:**
```
Authorization: Bearer <token>
```

### Get Categories with Transaction Count
**GET** `/user-categories/with-transaction-count`

Get user categories with their transaction counts.

**Headers:**
```
Authorization: Bearer <token>
```

### Get User Category by ID
**GET** `/user-categories/{id}`

Get a specific user category.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Category UUID

### Update User Category
**PUT** `/user-categories/{id}`

Update an existing user category.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Category UUID

**Request Body:**
```json
{
  "name": "Updated Category Name",
  "description": "Updated description",
  "icon": "🎯",
  "color": "#4ECDC4"
}
```

### Delete User Category
**DELETE** `/user-categories/{id}`

Delete a user category (soft delete).

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Category UUID

---

## Bank Endpoints

### Get All Banks
**GET** `/banks`

Get all active banks (public endpoint, no authentication required).

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "Bank of America",
      "code": "BOA",
      "country": "US",
      "is_active": true
    }
  ]
}
```

---

## Currency Endpoints

### Get All Currencies
**GET** `/currencies`

Get all active currencies (public endpoint, no authentication required).

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "code": "USD",
      "name": "US Dollar",
      "symbol": "$",
      "is_active": true
    }
  ]
}
```

---

## Plan Endpoints

### Get All Plans
**GET** `/plans`

Get all active plans (public endpoint, no authentication required).

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "Basic Plan",
      "description": "Basic features",
      "price": 9.99,
      "currency": "USD",
      "features": ["feature1", "feature2"]
    }
  ]
}
```

---

## Budget Endpoints

### Get All Budgets
**GET** `/budgets`

Get all budgets for a workspace.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `workspace_id` (required): Workspace UUID
- `limit` (optional): Number of items per page (default: 10)
- `offset` (optional): Number of items to skip (default: 0)

### Create Budget
**POST** `/budgets`

Create a new budget.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `workspace_id` (required): Workspace UUID

**Request Body:**
```json
{
  "name": "Monthly Groceries",
  "budgeted_amount": 500.00,
  "period_type": "monthly",
  "period_start": "2024-01-01",
  "period_end": "2024-01-31",
  "user_category_id": "uuid"
}
```

### Get Budget by ID
**GET** `/budgets/{id}`

Get a specific budget.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Budget UUID

**Query Parameters:**
- `workspace_id` (required): Workspace UUID

### Update Budget
**PUT** `/budgets/{id}`

Update an existing budget.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Budget UUID

**Query Parameters:**
- `workspace_id` (required): Workspace UUID

**Request Body:**
```json
{
  "name": "Updated Budget Name",
  "budgeted_amount": 600.00,
  "period_type": "monthly",
  "period_start": "2024-01-01",
  "period_end": "2024-01-31",
  "user_category_id": "uuid"
}
```

### Delete Budget
**DELETE** `/budgets/{id}`

Delete a budget.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Budget UUID

**Query Parameters:**
- `workspace_id` (required): Workspace UUID

---

## Transaction Endpoints

### Get Transactions by Workspace
**GET** `/transactions`

Get all transactions for a workspace with filtering and pagination.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `workspace_id` (required): Workspace UUID
- `account_id` (optional): Filter by account ID
- `category_id` (optional): Filter by category ID
- `start_date` (optional): Start date filter (YYYY-MM-DD)
- `end_date` (optional): End date filter (YYYY-MM-DD)
- `payment_method` (optional): Filter by payment method
- `description` (optional): Filter by description (partial match)
- `merchant_name` (optional): Filter by merchant name (partial match)
- `amount` (optional): Filter by exact amount
- `is_recurring` (optional): Filter by recurring status
- `credit_status` (optional): Filter by credit status
- `limit` (optional): Limit for pagination (default: 10)
- `offset` (optional): Offset for pagination (default: 0)

**Response:**
```json
{
  "success": true,
  "data": {
    "transactions": [
      {
        "id": "uuid",
        "description": "Grocery shopping",
        "amount": 50.00,
        "transaction_type": "expense",
        "payment_method": "card",
        "transaction_date": "2024-01-15",
        "merchant_name": "Walmart",
        "account_id": "uuid",
        "category_id": "uuid",
        "is_recurring": false,
        "created_at": "2024-01-15T10:30:00Z"
      }
    ],
    "total": 100,
    "limit": 10,
    "offset": 0
  }
}
```

### Create Transaction
**POST** `/transactions`

Create a new transaction.

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "description": "Grocery shopping",
  "amount": 50.00,
  "transaction_type": "expense",
  "payment_method": "card",
  "transaction_date": "2024-01-15",
  "merchant_name": "Walmart",
  "workspace_id": "uuid",
  "account_id": "uuid",
  "category_id": "uuid",
  "is_recurring": false,
  "credit_status": 0
}
```

### Get Transaction by ID
**GET** `/transactions/{id}`

Get a specific transaction.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Transaction UUID

### Update Transaction
**PUT** `/transactions/{id}`

Update an existing transaction.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Transaction UUID

**Request Body:**
```json
{
  "description": "Updated description",
  "amount": 55.00,
  "transaction_type": "expense",
  "payment_method": "card",
  "transaction_date": "2024-01-15",
  "merchant_name": "Walmart",
  "account_id": "uuid",
  "category_id": "uuid"
}
```

### Delete Transaction
**DELETE** `/transactions/{id}`

Delete a transaction.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Transaction UUID

---

## Conversation Endpoints

### Get Active Conversations
**GET** `/conversations/active`

Get all active conversations for the authenticated user.

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "channel": "whatsapp",
      "is_active": true,
      "context": "expense tracking",
      "metadata": {},
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

## Message Endpoints

### Get Messages by Conversation
**GET** `/messages/conversation/{conversation_id}`

Get all messages for a specific conversation.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `conversation_id`: Conversation UUID

**Query Parameters:**
- `limit` (optional): Limit for pagination (default: 10)
- `offset` (optional): Offset for pagination (default: 0)

**Response:**
```json
{
  "success": true,
  "data": {
    "messages": [
      {
        "id": "uuid",
        "conversation_id": "uuid",
        "user_id": "uuid",
        "sender_type": "user",
        "direction": "inbound",
        "message_type": "text",
        "content": "Hello",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "limit": 10,
    "offset": 0
  }
}
```

---

## Taxonomy Endpoints

### Get Taxonomies by Type
**GET** `/taxonomies/type/{type}`

Get taxonomies by type with pagination.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `type`: Taxonomy type

**Query Parameters:**
- `limit` (optional): Limit for pagination (default: 10)
- `offset` (optional): Offset for pagination (default: 0)

---

## User Tags Endpoints

### Get User Tags
**GET** `/user-tags`

Get all user tags for the authenticated user.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `limit` (optional): Limit for pagination (default: 10)
- `offset` (optional): Offset for pagination (default: 0)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "name": "Important",
      "color": "#FF6B6B",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### Create User Tag
**POST** `/user-tags`

Create a new user tag.

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "Important",
  "color": "#FF6B6B",
  "description": "Important transactions"
}
```

### Get Active User Tags
**GET** `/user-tags/active`

Get all active user tags.

**Headers:**
```
Authorization: Bearer <token>
```

### Get User Tag by ID
**GET** `/user-tags/{id}`

Get a specific user tag.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: User Tag UUID

### Update User Tag
**PUT** `/user-tags/{id}`

Update an existing user tag.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: User Tag UUID

**Request Body:**
```json
{
  "name": "Updated Tag Name",
  "color": "#4ECDC4",
  "description": "Updated description"
}
```

### Delete User Tag
**DELETE** `/user-tags/{id}`

Delete a user tag (soft delete).

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: User Tag UUID

---

## Transaction Tags Endpoints

### Create Transaction Tag
**POST** `/transaction-tags`

Create a new transaction tag.

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "transaction_id": "uuid",
  "user_tag_id": "uuid"
}
```

### Delete Transaction Tag
**DELETE** `/transaction-tags/{id}`

Delete a transaction tag.

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Transaction Tag UUID

---

## Verification Code Endpoints

### Create Verification Code
**POST** `/verification-codes/create`

Create a new verification code for phone number verification.

**Request Body:**
```json
{
  "phone_number": "+1234567890",
  "code_type": "phone_verification"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "verification_code_id": "uuid",
    "phone_number": "+1234567890",
    "code": "123456",
    "code_type": "phone_verification",
    "expires_at": "2024-01-01T00:10:00Z",
    "is_used": false,
    "attempts_count": 0,
    "max_attempts": 3,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "message": "Verification code sent successfully"
}
```

### Verify Code
**POST** `/verification-codes/verify`

Verify a verification code for phone number.

**Request Body:**
```json
{
  "phone_number": "+1234567890",
  "code": "123456"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Phone number verified successfully"
}
```

### Resend Verification Code
**POST** `/verification-codes/resend`

Resend a verification code for phone number.

**Query Parameters:**
- `phone_number` (required): Phone number to send code to
- `code_type` (required): Type of verification code (e.g., "phone_verification")

**Example:**
```
POST /verification-codes/resend?phone_number=+1234567890&code_type=phone_verification
```

**Response:**
```json
{
  "success": true,
  "message": "Verification code resent successfully"
}
```

---

## Error Responses

### Common Error Codes

**400 Bad Request**
```json
{
  "success": false,
  "error": "Invalid request parameters"
}
```

**401 Unauthorized**
```json
{
  "success": false,
  "error": "Authentication required"
}
```

**403 Forbidden**
```json
{
  "success": false,
  "error": "Access denied"
}
```

**404 Not Found**
```json
{
  "success": false,
  "error": "Resource not found"
}
```

**409 Conflict**
```json
{
  "success": false,
  "error": "Resource already exists"
}
```

**500 Internal Server Error**
```json
{
  "success": false,
  "error": "Internal server error"
}
```

---

## Rate Limiting

The API implements rate limiting to ensure fair usage. Limits are applied per IP address and per authenticated user.

- **Public endpoints**: 100 requests per minute
- **Authenticated endpoints**: 1000 requests per minute

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640995200
```

---

## Pagination

For endpoints that return lists, pagination is supported using `limit` and `offset` query parameters:

- `limit`: Number of items to return (default: 10, max: 100)
- `offset`: Number of items to skip (default: 0)

Example:
```
GET /transactions?limit=20&offset=40
```

---

## Data Types

### UUID
All IDs are UUID v4 strings.

### Date/Time
All date and time fields use ISO 8601 format: `YYYY-MM-DDTHH:MM:SSZ`

### Currency
All monetary amounts are decimal numbers with up to 2 decimal places.

### Boolean
Boolean values are represented as `true` or `false`.

---

## TypeScript Data Structures

### Base Types

```typescript
// UUID type
type UUID = string;

// API Response wrapper
interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
}

// Pagination types
interface PaginationParams {
  limit?: number;
  offset?: number;
}

interface PaginatedResponse<T> {
  data: T[];
  total: number;
  limit: number;
  offset: number;
}
```

### Authentication Types

```typescript
// User registration
interface CreateUserRequest {
  email: string;
  phone: string;
  first_name: string;
  last_name: string;
  password: string;
  pin: string;
}

interface User {
  id: UUID;
  email: string;
  phone: string;
  first_name: string;
  last_name: string;
  created_at: string;
}

// Login
interface LoginRequest {
  email?: string;
  phone?: string;
  password: string;
}

interface LoginResponse {
  user: User;
  token: string;
}

// Password management
interface ForgotPasswordRequest {
  email: string;
}

interface ResetPasswordRequest {
  token: string;
  new_password: string;
}

interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

// Verification
interface VerifyPhoneRequest {
  phone: string;
  code: string;
}

interface ResendVerificationCodeRequest {
  phone: string;
}

interface VerifyEmailRequest {
  token: string;
}

interface ResendVerificationEmailRequest {
  email: string;
}
```

### Workspace Types

```typescript
interface Workspace {
  id: UUID;
  name: string;
  type: number;
  currency_id: UUID;
  created_by: UUID;
  created_at: string;
  updated_at: string;
}

interface CreateWorkspaceRequest {
  name: string;
  type: number;
  currency_id: UUID;
  description?: string;
}

interface UpdateWorkspaceRequest {
  name?: string;
  type?: number;
  currency_id?: UUID;
  description?: string;
}

// Workspace type constants
enum WorkspaceType {
  Personal = 1,
  Business = 2,
  Event = 3,
  Travel = 4,
  Project = 5,
  Shared = 6
}
```

### Account Types

```typescript
interface Account {
  id: UUID;
  name: string;
  type: number;
  currency_id: UUID;
  balance: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

interface CreateAccountRequest {
  name: string;
  type: number;
  currency_id: UUID;
  initial_balance?: number;
}

interface UpdateAccountRequest {
  name: string;
  type: number;
  currency_id: UUID;
  balance: number;
}

// Account type constants
enum AccountType {
  Debit = 1,
  Credit = 2,
  Savings = 3,
  Cash = 4,
  Shared = 5
}
```

### Category Types

```typescript
interface Category {
  id: UUID;
  name: string;
  description?: string;
  icon?: string;
  is_system_category: boolean;
  created_at: string;
  updated_at: string;
}

interface CreateCategoryRequest {
  name: string;
  description?: string;
  icon?: string;
  color?: string;
}

interface UpdateCategoryRequest {
  name: string;
  description?: string;
  icon?: string;
  color?: string;
}

interface UserCategory extends Category {
  user_id: UUID;
  is_custom: boolean;
  is_active: boolean;
}
```

### Bank Types

```typescript
interface Bank {
  id: number;
  name: string;
  code: string;
  logo_url?: string;
  is_active: boolean;
}
```

### Currency Types

```typescript
interface Currency {
  id: number;
  code: string;
  name: string;
  symbol: string;
  is_active: boolean;
}
```

### Plan Types

```typescript
interface Plan {
  id: number;
  name: string;
  description: string;
  price: string;
  currency: string;
  features: string[];
  is_active: boolean;
}
```

### Budget Types

```typescript
interface Budget {
  id: UUID;
  workspace_id: UUID;
  user_category_id: UUID;
  name: string;
  budgeted_amount: number;
  period_type: number;
  period_start: string;
  period_end: string;
  spent_amount: number;
  is_active: boolean;
  created_by: UUID;
  created_at: string;
  updated_at: string;
}

interface CreateBudgetRequest {
  workspace_id: UUID;
  user_category_id: UUID;
  name: string;
  budgeted_amount: number;
  period_type: number;
  period_start: string;
  period_end: string;
}

interface UpdateBudgetRequest {
  user_category_id: UUID;
  name: string;
  budgeted_amount: number;
  period_type: number;
  period_start: string;
  period_end: string;
  spent_amount: number;
  is_active: boolean;
}

// Period type constants
enum PeriodType {
  Weekly = 1,
  Monthly = 2,
  Yearly = 3,
  Event = 4
}
```

### Transaction Types

```typescript
interface Transaction {
  id: UUID;
  description: string;
  amount: number;
  transaction_type: number;
  payment_method: number;
  transaction_date: string;
  merchant_name?: string;
  account_id: UUID;
  category_id?: UUID;
  is_recurring: boolean;
  created_at: string;
}

interface CreateTransactionRequest {
  description: string;
  amount: number;
  transaction_type: number;
  payment_method: number;
  transaction_date: string;
  merchant_name?: string;
  workspace_id: UUID;
  account_id: UUID;
  category_id?: UUID;
  is_recurring?: boolean;
  credit_status?: number;
}

interface UpdateTransactionRequest {
  description: string;
  amount: number;
  transaction_type: number;
  payment_method: number;
  transaction_date: string;
  merchant_name?: string;
  account_id: UUID;
  category_id?: UUID;
}

interface TransactionListParams {
  workspace_id: UUID;
  account_id?: UUID;
  category_id?: UUID;
  start_date?: string;
  end_date?: string;
  payment_method?: number;
  description?: string;
  merchant_name?: string;
  amount?: number;
  is_recurring?: boolean;
  credit_status?: number;
  limit?: number;
  offset?: number;
}

// Transaction type constants
enum TransactionType {
  Income = 1,
  Expense = 2
}

enum PaymentMethod {
  DebitQRIS = 1,
  Credit = 2,
  Cash = 3,
  Transfer = 4
}

enum CreditStatus {
  Paid = 1,
  Unpaid = 2
}
```

### Conversation Types

```typescript
interface Conversation {
  id: UUID;
  user_id: UUID;
  channel: string;
  is_active: boolean;
  context?: string;
  metadata?: Record<string, any>;
  created_at: string;
  updated_at: string;
}
```

### Message Types

```typescript
interface Message {
  id: UUID;
  conversation_id: UUID;
  user_id: UUID;
  sender_type: string;
  direction: string;
  message_type: string;
  content: string;
  created_at: string;
}

interface MessageListResponse {
  messages: Message[];
  total: number;
  limit: number;
  offset: number;
}
```

### Taxonomy Types

```typescript
interface Taxonomy {
  id: number;
  label: string;
  value: string;
  type: string;
  type_label: string;
  status: number;
  created_at: string;
  updated_at: string;
}
```

### User Tag Types

```typescript
interface UserTag {
  id: UUID;
  user_id: UUID;
  name: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

interface CreateUserTagRequest {
  name: string;
}

interface UpdateUserTagRequest {
  name: string;
  is_active: boolean;
}
```

### Transaction Tag Types

```typescript
interface TransactionTag {
  id: UUID;
  transaction_id: UUID;
  user_tag_id: UUID;
  applied_by: UUID;
  applied_at: string;
}

interface CreateTransactionTagRequest {
  transaction_id: UUID;
  user_tag_id: UUID;
}

interface TransactionTagSimple {
  id: UUID;
  transaction_id: UUID;
  user_tag_id: UUID;
  tag_name: string;
  applied_by: UUID;
  applied_at: string;
}
```

### Verification Code Types

```typescript
interface VerificationCode {
  verification_code_id: UUID;
  phone_number: string;
  code: string;
  code_type: string;
  expires_at: string;
  is_used: boolean;
  attempts_count: number;
  max_attempts: number;
  created_at: string;
  updated_at: string;
}

interface CreateVerificationCodeRequest {
  phone_number: string;
  code_type: string;
}

interface VerifyVerificationCodeRequest {
  phone_number: string;
  code: string;
}

// Code type constants
enum CodeType {
  PhoneVerification = "phone_verification",
  PasswordReset = "password_reset",
  EmailVerification = "email_verification"
}
```

### API Client Example

```typescript
class VasstExpenseAPI {
  private baseURL: string;
  private token?: string;

  constructor(baseURL: string = 'https://api.vasst.id/v1') {
    this.baseURL = baseURL;
  }

  setToken(token: string) {
    this.token = token;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const url = `${this.baseURL}${endpoint}`;
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (this.token) {
      headers.Authorization = `Bearer ${this.token}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    return response.json();
  }

  // Authentication
  async login(request: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    return this.request<LoginResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  async register(request: CreateUserRequest): Promise<ApiResponse<User>> {
    return this.request<User>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  // Workspaces
  async getWorkspaces(params?: PaginationParams): Promise<ApiResponse<Workspace[]>> {
    const query = new URLSearchParams();
    if (params?.limit) query.append('limit', params.limit.toString());
    if (params?.offset) query.append('offset', params.offset.toString());
    
    return this.request<Workspace[]>(`/workspaces?${query}`);
  }

  async createWorkspace(request: CreateWorkspaceRequest): Promise<ApiResponse<Workspace>> {
    return this.request<Workspace>('/workspaces', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  // Transactions
  async getTransactions(params: TransactionListParams): Promise<ApiResponse<PaginatedResponse<Transaction>>> {
    const query = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined) {
        query.append(key, value.toString());
      }
    });
    
    return this.request<PaginatedResponse<Transaction>>(`/transactions?${query}`);
  }

  async createTransaction(request: CreateTransactionRequest): Promise<ApiResponse<Transaction>> {
    return this.request<Transaction>('/transactions', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  // Categories
  async getUserCategories(params?: PaginationParams): Promise<ApiResponse<UserCategory[]>> {
    const query = new URLSearchParams();
    if (params?.limit) query.append('limit', params.limit.toString());
    if (params?.offset) query.append('offset', params.offset.toString());
    
    return this.request<UserCategory[]>(`/user-categories?${query}`);
  }

  // User Tags
  async getUserTags(params?: PaginationParams): Promise<ApiResponse<UserTag[]>> {
    const query = new URLSearchParams();
    if (params?.limit) query.append('limit', params.limit.toString());
    if (params?.offset) query.append('offset', params.offset.toString());
    
    return this.request<UserTag[]>(`/user-tags?${query}`);
  }

  async createUserTag(request: CreateUserTagRequest): Promise<ApiResponse<UserTag>> {
    return this.request<UserTag>('/user-tags', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }
}

// Usage example
const api = new VasstExpenseAPI();

// Login
const loginResponse = await api.login({
  email: 'user@example.com',
  password: 'password123'
});

if (loginResponse.success && loginResponse.data) {
  api.setToken(loginResponse.data.token);
  
  // Get workspaces
  const workspaces = await api.getWorkspaces();
  
  // Create transaction
  const transaction = await api.createTransaction({
    description: 'Grocery shopping',
    amount: 50.00,
    transaction_type: TransactionType.Expense,
    payment_method: PaymentMethod.Cash,
    transaction_date: '2024-01-15',
    workspace_id: 'workspace-uuid',
    account_id: 'account-uuid'
  });
}
```

---

## SDKs and Libraries

Official SDKs are available for:
- JavaScript/TypeScript
- Python
- Go
- Java

For more information, visit our [Developer Portal](https://developers.vasst.id). 