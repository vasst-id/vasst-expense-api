# ExpenseTracker SaaS - Technical Documentation

## 1. System Architecture Overview

### 1.1 High-Level Architecture
```
┌─────────────────┐                            ┌─────────────────┐
│   WhatsApp      │                            │   Frontend      │
│   Business API  │                            │   (Next.js)     │
└─────────────────┘                            └─────────────────┘
        │                                                  │
        ▼                                                  │
┌────────────────────────────────────────────────────────┬─────────┐
│                 Golang Backend Services                          │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────────┐          │
│  │    Auth      │ │  Transaction │ │    Workspace     │          │
│  │   Service    │ │   Service    │ │    Service       │          │
│  └──────────────┘ └──────────────┘ └──────────────────┘          │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────────┐          │
│  │   Webhook    │ │      AI      │ │   Scheduler      │          │
│  │   Service    │ │   Service    │ │    Service       │          │
│  └──────────────┘ └──────────────┘ └──────────────────┘          │
└────────────────────────────────────────────────────────┬─────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────┬─────────┐
│                     Data Layer                                    │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────────┐          │
│  │  PostgreSQL  │ │    Redis     │ │      GCS.        │          │
│  │   Primary    │ │   Cache +    │ │   File Storage   │          │
│  │   Database   │ │    Queue     │ │                  │          │
│  └──────────────┘ └──────────────┘ └──────────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 Event-Driven Architecture
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Event Bus     │    │   Event Store    │    │   Dead Letter   │
(Google Pub/Sub)  │◄──►│   (PostgreSQL)   │    │     Queue       │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Event Processors                            │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────────┐        │
│  │   Budget     │ │  Settlement  │ │    Webhook       │        │
│  │   Monitor    │ │  Calculator  │ │   Dispatcher     │        │
│  └──────────────┘ └──────────────┘ └──────────────────┘        │
└─────────────────────────────────────────────────────────────────┘
```

## 2. Technology Stack Decisions

### 2.1 Backend Technology Stack
```yaml
Language: Go 1.21+
Framework: Gin Web Framework
Database: PostgreSQL 15+
Cache: Redis 7+
Queue: Google Pub Sub
File Storage: GCS
Search: PostgreSQL Full-Text Search
Monitoring: Prometheus + Grafana
Logging: Structured logging with logrus
Containerization: Docker + Docker Compose
```

### 2.2 Frontend Technology Stack
```yaml
Framework: Next.js 14 (App Router)
Language: TypeScript
Styling: Custom CSS
State Management: Zustand
Forms: React Hook Form + Zod
HTTP Client: Axios with React Query
UI Components: shadcn/ui + Headless UI
Charts: Recharts
Icons: Lucide React
Authentication: NextAuth.js
PWA: next-pwa
```


### 2.3 Third-Party Integrations
```yaml
AI/ML: Google Gemini API, openAI API, claude API
OCR: Google Gemini API
Payment: Midtrans
WhatsApp: Whatsapp Business API
Email: SendGrid
Push Notifications: Firebase Cloud Messaging
Analytics: Mixpanel + Google Analytics
Error Tracking: Sentry
```

## 3. Database Schema & Structure

### 3.1 Core Tables
```sql
-- Currency
CREATE TABLE "vasst_expense".currency (
    currency_id INT PRIMARY KEY,
    currency_code VARCHAR(3) NOT NULL
);

-- Subcription Plan
CREATE TABLE "vasst_expense".subscription_plan (
    subscription_plan_id INT PRIMARY KEY,
    subscription_plan_name VARCHAR(20) NOT NULL,
)

-- Users
CREATE TABLE "vasst_expense".users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    phone_number VARCHAR(20),
    timezone VARCHAR(50) DEFAULT 'Asia/Jakarta',
    currency_id INT DEFAULT 1,
    subscription_plan_id INT DEFAULT 1,
    email_verified_at TIMESTAMP,
    phone_verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    status INT DEFAULT 1
);

-- Workspaces
CREATE TABLE "vasst_expense".workspaces (
    workspace_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    workspace_type VARCHAR(50) NOT NULL, -- 'personal', 'business', 'event', 'travel', 'project', 'shared'
    icon VARCHAR(50) DEFAULT 'folder',
    color_code VARCHAR(7) DEFAULT '#3B82F6',
    currency_id INT DEFAULT 1,
    timezone VARCHAR(50) DEFAULT 'Asia/Jakarta',
    settings JSONB DEFAULT '{}', -- workspace-specific settings
    is_active BOOLEAN DEFAULT true,
    created_by UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Workspace Members (NEW - For collaborative workspaces)
-- CREATE TABLE "vasst_expense".workspace_members (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
--     user_id UUID REFERENCES users(id) ON DELETE CASCADE,
--     role VARCHAR(20) DEFAULT 'member', -- 'owner', 'admin', 'member', 'viewer'
--     permissions JSONB DEFAULT '{}', -- specific permissions
--     joined_at TIMESTAMP DEFAULT NOW(),
--     is_active BOOLEAN DEFAULT true,
--     UNIQUE(workspace_id, user_id)
-- );

-- Banks
CREATE TABLE "vasst_expense".banks (
    bank_id INT PRIMARY KEY NOT NULL,
    bank_name VARCHAR(50) NOT NULL
)

-- Accounts (Modified to support workspace-specific accounts)
CREATE TABLE "vasst_expense".accounts (
    account_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    account_name VARCHAR(100) NOT NULL,
    account_type VARCHAR(20) NOT NULL, -- 'debit', 'credit', 'savings', 'cash', 'shared'
    bank_id INT,
    account_number_masked VARCHAR(20),
    current_balance DECIMAL(15,2) DEFAULT 0,
    credit_limit DECIMAL(15,2), -- For credit cards
    due_date INTEGER, -- Day of month for credit card due date
    currency_id INT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Categories (Modified to support workspace-specific categories)
CREATE TABLE "vasst_expense".categories (
    category_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    color_code VARCHAR(7) DEFAULT '#3B82F6',
    icon VARCHAR(50) DEFAULT 'receipt',
    parent_category_id UUID REFERENCES categories(id),
    is_system_category BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(name) -- Unique category names per workspace
);

-- User Categories (Enhanced - User's custom categories)
-- Note: This replaces the previous categories table
CREATE TABLE "vasst_expense".user_categories (
    user_category_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(user_id),
    category_id UUID REFERENCES categories(user_id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    color_code VARCHAR(7) DEFAULT '#3B82F6',
    icon VARCHAR(50) DEFAULT 'receipt',
    is_custom BOOLEAN DEFAULT false, -- true if user created, false if from predefined
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
);

-- Budgets (Modified to support workspace-specific budgets)
CREATE TABLE "vasst_expense".budgets (
    budget_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    category_id UUID REFERENCES user_categories(id) ON DELETE CASCADE, -- Updated reference
    name VARCHAR(100) NOT NULL,
    budgeted_amount DECIMAL(15,2) NOT NULL,
    period_type VARCHAR(20) DEFAULT 'monthly', -- 'weekly', 'monthly', 'yearly', 'event'
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    spent_amount DECIMAL(15,2) DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Transactions (Enhanced for multi-workspace and bill splitting)
CREATE TABLE "vasst_expense".transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    account_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
    category_id UUID REFERENCES user_categories(id) ON DELETE SET NULL, -- Updated reference
    description TEXT NOT NULL,
    amount DECIMAL(15,2) NOT NULL, -- Total transaction amount
    transaction_type VARCHAR(20) NOT NULL, -- 'income', 'expense'
    payment_method VARCHAR(20) NOT NULL, -- 'debit', 'credit', 'cash', 'transfer'
    transaction_date DATE NOT NULL,
    merchant_name VARCHAR(255),
    location TEXT,
    notes TEXT,
    receipt_url TEXT,
    
    -- Bill splitting and group features
    is_split BOOLEAN DEFAULT false,
    split_type VARCHAR(20), -- 'equal', 'percentage', 'exact', 'by_share'
    split_count INTEGER DEFAULT 1,
    paid_by UUID REFERENCES users(id), -- Who paid the bill
    
    -- Recurring transactions
    is_recurring BOOLEAN DEFAULT false,
    recurring_pattern JSONB, -- recurring configuration
    parent_transaction_id UUID REFERENCES transactions(id),
    scheduler_task_id UUID REFERENCES scheduler_tasks(id), -- Link to recurring task
    
    -- AI and processing
    ai_confidence_score DECIMAL(3,2), -- 0.00 to 1.00
    ai_categorized BOOLEAN DEFAULT false,
    processing_status VARCHAR(20) DEFAULT 'completed', -- 'pending', 'completed', 'failed'
    
    -- Credit tracking
    credit_status VARCHAR(20), -- 'paid', 'unpaid' for credit transactions
    
    -- Metadata
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Transaction Splits (NEW - For bill splitting functionality)
CREATE TABLE "vasst_expense".transaction_splits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID REFERENCES transactions(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(15,2) NOT NULL, -- Amount this user owes/paid
    percentage DECIMAL(5,2), -- Percentage of total (for percentage splits)
    shares INTEGER, -- Number of shares (for share-based splits)
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'paid', 'settled'
    paid_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Settlements (NEW - For tracking who owes whom)
CREATE TABLE "vasst_expense".settlements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    from_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    to_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(15,2) NOT NULL,
    description TEXT,
    transaction_ids UUID[], -- Array of related transaction IDs
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'completed', 'cancelled'
    settled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Documents (Enhanced for workspace-specific documents)
CREATE TABLE "vasst_expense".documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    original_filename VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_size INTEGER NOT NULL,
    document_type VARCHAR(50) NOT NULL, -- 'receipt', 'statement', 'invoice'
    ai_analysis_result JSONB,
    processing_status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'processing', 'completed', 'failed'
    source_type VARCHAR(20) DEFAULT 'web', -- 'web', 'whatsapp', 'email'
    uploaded_at TIMESTAMP DEFAULT NOW(),
    processed_at TIMESTAMP
);

-- AI Analysis Logs (Enhanced)
CREATE TABLE "vasst_expense".ai_analysis_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    document_id UUID REFERENCES documents(id) ON DELETE CASCADE,
    analysis_type VARCHAR(50) NOT NULL,
    input_data JSONB,
    output_data JSONB,
    model_used VARCHAR(100),
    confidence_score DECIMAL(3,2),
    processing_time_ms INTEGER,
    status VARCHAR(20) NOT NULL,
    error_details JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- User Preferences (Enhanced with workspace preferences)
CREATE TABLE "vasst_expense".user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE, -- NULL for global preferences
    preference_type VARCHAR(50) NOT NULL, -- 'notification', 'display', 'ai', 'export'
    preferences JSONB NOT NULL DEFAULT '{}',
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, workspace_id, preference_type)
);

-- Audit Logs (Enhanced for workspace tracking)
CREATE TABLE "vasst_expense".audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    source_type VARCHAR(20) DEFAULT 'web', -- 'web', 'whatsapp', 'api'
    created_at TIMESTAMP DEFAULT NOW()
);

-- Workspace Invitations (NEW - For inviting users to workspaces)
CREATE TABLE "vasst_expense".workspace_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    invited_by UUID REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'member',
    invitation_token VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'accepted', 'declined', 'expired'
    expires_at TIMESTAMP NOT NULL,
    accepted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Webhook URLs (NEW - For VASST communication agent integration)
CREATE TABLE "vasst_expense".webhook_urls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE, -- NULL for global webhooks
    webhook_type VARCHAR(50) NOT NULL, -- 'vasst_callback', 'transaction_update', 'budget_alert', 'settlement_reminder'
    url TEXT NOT NULL,
    secret_token VARCHAR(255), -- For webhook verification
    is_active BOOLEAN DEFAULT true,
    retry_count INTEGER DEFAULT 3,
    timeout_seconds INTEGER DEFAULT 30,
    headers JSONB DEFAULT '{}', -- Custom headers
    events JSONB DEFAULT '[]', -- Array of events to trigger webhook
    last_triggered_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Request Callback Logs (NEW - For tracking VASST communication)
CREATE TABLE "vasst_expense".request_callback_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    webhook_url_id UUID REFERENCES webhook_urls(id) ON DELETE SET NULL,
    request_type VARCHAR(50) NOT NULL, -- 'transaction_created', 'budget_exceeded', 'settlement_reminder', etc.
    request_method VARCHAR(10) DEFAULT 'POST',
    request_url TEXT NOT NULL,
    request_headers JSONB,
    request_body JSONB,
    response_status INTEGER,
    response_headers JSONB,
    response_body JSONB,
    processing_time_ms INTEGER,
    retry_attempt INTEGER DEFAULT 0,
    status VARCHAR(20) NOT NULL, -- 'pending', 'success', 'failed', 'retrying'
    error_message TEXT,
    triggered_by VARCHAR(50), -- 'system', 'user_action', 'scheduler', 'ai_agent'
    event_id UUID, -- Reference to the event that triggered this callback
    created_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- Predefined Categories (NEW - System-wide category templates)
CREATE TABLE "vasst_expense".predefined_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    name_id VARCHAR(100) NOT NULL, -- Slug/identifier version
    description TEXT,
    workspace_type VARCHAR(50), -- NULL for all types, or specific to 'personal', 'business', etc.
    parent_category_id UUID REFERENCES predefined_categories(id),
    color_code VARCHAR(7) DEFAULT '#3B82F6',
    icon VARCHAR(50) DEFAULT 'receipt',
    keywords JSONB DEFAULT '[]', -- Keywords for AI categorization
    sort_order INTEGER DEFAULT 0,
    is_income_category BOOLEAN DEFAULT false,
    is_system_default BOOLEAN DEFAULT true,
    locale VARCHAR(10) DEFAULT 'id-ID',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(name_id, workspace_type, locale)
);

-- Predefined Tags (NEW - System-wide tag templates)
CREATE TABLE "vasst_expense".predefined_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    name_id VARCHAR(50) NOT NULL, -- Slug/identifier version
    description TEXT,
    tag_type VARCHAR(30) NOT NULL, -- 'expense_type', 'priority', 'payment_status', 'location', 'occasion'
    color_code VARCHAR(7) DEFAULT '#6B7280',
    icon VARCHAR(50),
    usage_context JSONB DEFAULT '[]', -- Array of contexts where this tag applies
    auto_apply_rules JSONB DEFAULT '{}', -- Rules for automatic tag application
    sort_order INTEGER DEFAULT 0,
    is_system_default BOOLEAN DEFAULT true,
    locale VARCHAR(10) DEFAULT 'id-ID',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(name_id, tag_type, locale)
);

-- User Tags (NEW - User's custom and applied tags)
CREATE TABLE "vasst_expense".user_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    predefined_tag_id UUID REFERENCES predefined_tags(id) ON DELETE SET NULL,
    name VARCHAR(50) NOT NULL,
    color_code VARCHAR(7) DEFAULT '#6B7280',
    icon VARCHAR(50),
    is_custom BOOLEAN DEFAULT false, -- true if user created, false if from predefined
    usage_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(workspace_id, name)
);

-- Transaction Tags (NEW - Many-to-many relationship for transaction tagging)
CREATE TABLE "vasst_expense".transaction_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID REFERENCES transactions(id) ON DELETE CASCADE,
    user_tag_id UUID REFERENCES user_tags(id) ON DELETE CASCADE,
    applied_by UUID REFERENCES users(id) ON DELETE SET NULL,
    applied_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(transaction_id, user_tag_id)
);

-- Scheduler Tasks (NEW - For recurring transactions and reminders)
CREATE TABLE "vasst_expense".scheduler_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    task_type VARCHAR(50) NOT NULL, -- 'recurring_transaction', 'budget_reminder', 'settlement_reminder', 'bill_due_reminder'
    task_name VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Scheduling configuration
    schedule_type VARCHAR(20) NOT NULL, -- 'once', 'daily', 'weekly', 'monthly', 'yearly', 'cron'
    schedule_config JSONB NOT NULL, -- Contains schedule details (cron expression, day of month, etc.)
    timezone VARCHAR(50) DEFAULT 'Asia/Jakarta',
    
    -- Execution details
    next_execution_at TIMESTAMP NOT NULL,
    last_execution_at TIMESTAMP,
    execution_count INTEGER DEFAULT 0,
    max_executions INTEGER, -- NULL for infinite
    
    -- Task payload
    task_payload JSONB NOT NULL, -- Contains the data needed to execute the task
    
    -- Status and error handling
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'paused', 'completed', 'failed', 'cancelled'
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    last_error TEXT,
    
    -- Metadata
    created_by UUID REFERENCES users(id),
    is_system_task BOOLEAN DEFAULT false, -- true for system-generated tasks
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Scheduler Execution Logs (NEW - For tracking task executions)
CREATE TABLE "vasst_expense".scheduler_execution_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID REFERENCES scheduler_tasks(id) ON DELETE CASCADE,
    execution_id UUID NOT NULL, -- Unique identifier for this execution attempt
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    status VARCHAR(20) NOT NULL, -- 'running', 'completed', 'failed', 'skipped'
    result JSONB, -- Execution result data
    error_message TEXT,
    execution_time_ms INTEGER,
    retry_attempt INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 3.2 Indexing Strategy
```sql
-- Performance-critical indexes
CREATE INDEX idx_transactions_workspace_date ON transactions(workspace_id, transaction_date DESC);
CREATE INDEX idx_transactions_user_created ON transactions(created_by, created_at DESC);
CREATE INDEX idx_workspace_members_user ON workspace_members(user_id, is_active);
CREATE INDEX idx_budgets_workspace_period ON budgets(workspace_id, period_start, period_end);
CREATE INDEX idx_settlements_workspace_status ON settlements(workspace_id, status);


## 4. Backend Architecture & Standards

### 4.1 Project Structure
```

- `cmd/api/` - API server entry point
- `cmd/worker/` - Worker service entry point  
- `internal/api/` - API application setup
- `internal/subscriber/` - Worker application setup
- `internal/controller/http/v1/` - Operational endpoints
- `internal/entities/` - Domain models and database entities
- `internal/events/` - Event publishing and handling
- `internal/workers/` - Background workers for AI processing
- `internal/services/` - Business logic services
- `internal/repositories/` - Database access layer
- `internal/middleware/` - HTTP middleware (auth, rate limiting, tenant isolation)
- `internal/pubsub/` - Google Cloud Pub/Sub integration
- `migrations/` - Database migrations
```

### 4.2 Clean Architecture Standards

#### 4.2.1 Entity Layer Example
```go
// internal/entities/transaction.go
package entities

import (
    "time"
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

type Transaction struct {
    ID                uuid.UUID       `json:"id"`
    WorkspaceID       uuid.UUID       `json:"workspace_id"`
    AccountID         *uuid.UUID      `json:"account_id,omitempty"`
    CategoryID        *uuid.UUID      `json:"category_id,omitempty"`
    Description       string          `json:"description"`
    Amount            decimal.Decimal `json:"amount"`
    TransactionType   TransactionType `json:"transaction_type"`
    PaymentMethod     PaymentMethod   `json:"payment_method"`
    TransactionDate   time.Time       `json:"transaction_date"`
    MerchantName      *string         `json:"merchant_name,omitempty"`
    Location          *string         `json:"location,omitempty"`
    Notes             *string         `json:"notes,omitempty"`
    ReceiptURL        *string         `json:"receipt_url,omitempty"`
    IsSplit           bool            `json:"is_split"`
    SplitType         *string         `json:"split_type,omitempty"`
    PaidBy            *uuid.UUID      `json:"paid_by,omitempty"`
    AIConfidenceScore *decimal.Decimal `json:"ai_confidence_score,omitempty"`
    AICategorized     bool            `json:"ai_categorized"`
    ProcessingStatus  string          `json:"processing_status"`
    CreditStatus      *string         `json:"credit_status,omitempty"`
    CreatedBy         uuid.UUID       `json:"created_by"`
    CreatedAt         time.Time       `json:"created_at"`
    UpdatedAt         time.Time       `json:"updated_at"`
}

type TransactionType string

const (
    TransactionTypeIncome  TransactionType = "income"
    TransactionTypeExpense TransactionType = "expense"
)

type PaymentMethod string

const (
    PaymentMethodCash    PaymentMethod = "cash"
    PaymentMethodDebit   PaymentMethod = "debit"
    PaymentMethodCredit  PaymentMethod = "credit"
    PaymentMethodDigital PaymentMethod = "digital"
)
```

#### 4.2.2 Repository Layer Example
```go
// internal/repositories/transaction_repository.go
package repositories

import (
    "context"
    "github.com/google/uuid"
    "github.com/expensetracker/internal/entities"
)

type TransactionRepository interface {
    Create(ctx context.Context, transaction *entities.Transaction) error
    GetByID(ctx context.Context, id uuid.UUID) (*entities.Transaction, error)
    GetByWorkspace(ctx context.Context, workspaceID uuid.UUID, filters TransactionFilters) ([]*entities.Transaction, error)
    Update(ctx context.Context, transaction *entities.Transaction) error
    Delete(ctx context.Context, id uuid.UUID) error
    GetByDateRange(ctx context.Context, workspaceID uuid.UUID, startDate, endDate time.Time) ([]*entities.Transaction, error)
    GetByCategoryAndPeriod(ctx context.Context, categoryID uuid.UUID, startDate, endDate time.Time) ([]*entities.Transaction, error)
}

type TransactionFilters struct {
    StartDate       *time.Time
    EndDate         *time.Time
    CategoryIDs     []uuid.UUID
    TransactionType *entities.TransactionType
    PaymentMethod   *entities.PaymentMethod
    MinAmount       *decimal.Decimal
    MaxAmount       *decimal.Decimal
    Limit           int
    Offset          int
    OrderBy         string
    OrderDirection  string
}
```

#### 4.2.3 Service Layer Example
```go
// services/transaction_service.go
package services

import (
    "context"
    "github.com/google/uuid"
    "github.com/expensetracker/internal/entities"
    "github.com/expensetracker/internal/repositories"
)

type TransactionService interface {
    CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*entities.Transaction, error)
    GetTransaction(ctx context.Context, id uuid.UUID) (*entities.Transaction, error)
    GetWorkspaceTransactions(ctx context.Context, workspaceID uuid.UUID, filters repositories.TransactionFilters) ([]*entities.Transaction, error)
    UpdateTransaction(ctx context.Context, id uuid.UUID, req UpdateTransactionRequest) (*entities.Transaction, error)
    DeleteTransaction(ctx context.Context, id uuid.UUID) error
    SplitTransaction(ctx context.Context, transactionID uuid.UUID, splits []SplitRequest) error
    ProcessReceiptImage(ctx context.Context, workspaceID uuid.UUID, imageURL string, userID uuid.UUID) (*entities.Transaction, error)
}

type transactionService struct {
    transactionRepo repositories.TransactionRepository
    workspaceRepo   repositories.WorkspaceRepository
    categoryRepo    repositories.CategoryRepository
    aiService       AIService
    eventBus        EventBus
}

func NewTransactionService(
    transactionRepo repositories.TransactionRepository,
    workspaceRepo repositories.WorkspaceRepository,
    categoryRepo repositories.CategoryRepository,
    aiService AIService,
    eventBus EventBus,
) TransactionService {
    return &transactionService{
        transactionRepo: transactionRepo,
        workspaceRepo:   workspaceRepo,
        categoryRepo:    categoryRepo,
        aiService:       aiService,
        eventBus:        eventBus,
    }
}

func (s *transactionService) CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*entities.Transaction, error) {
    // Validate request
    if err := s.validateCreateRequest(req); err != nil {
        return nil, err
    }

    // Check workspace access
    if err := s.checkWorkspaceAccess(ctx, req.WorkspaceID, req.UserID); err != nil {
        return nil, err
    }

    // Create transaction entity
    transaction := &entities.Transaction{
        ID:              uuid.New(),
        WorkspaceID:     req.WorkspaceID,
        Description:     req.Description,
        Amount:          req.Amount,
        TransactionType: req.TransactionType,
        PaymentMethod:   req.PaymentMethod,
        TransactionDate: req.TransactionDate,
        CreatedBy:       req.UserID,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }

    // Auto-categorize if not provided
    if req.CategoryID == nil && req.AutoCategorize {
        categoryID, confidence, err := s.aiService.CategorizeTransaction(ctx, req.Description, req.MerchantName)
        if err == nil && confidence > 0.7 {
            transaction.CategoryID = &categoryID
            transaction.AICategorized = true
            transaction.AIConfidenceScore = &confidence
        }
    } else if req.CategoryID != nil {
        transaction.CategoryID = req.CategoryID
    }

    // Save transaction
    if err := s.transactionRepo.Create(ctx, transaction); err != nil {
        return nil, err
    }

    // Publish event for budget tracking, webhooks, etc.
    event := &TransactionCreatedEvent{
        TransactionID: transaction.ID,
        WorkspaceID:   transaction.WorkspaceID,
        Amount:        transaction.Amount,
        CategoryID:    transaction.CategoryID,
        CreatedBy:     transaction.CreatedBy,
        CreatedAt:     transaction.CreatedAt,
    }
    s.eventBus.Publish(ctx, "transaction.created", event)

    return transaction, nil
}
```

#### 4.2.4 Controller/Handler Layer Example
```go
// api/handlers/transaction/handler.go
package transaction

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/expensetracker/internal/domain/services"
    "github.com/expensetracker/pkg/errors"
)

type Handler struct {
    transactionService services.TransactionService
    logger            logger.Logger
}

func NewHandler(transactionService services.TransactionService, logger logger.Logger) *Handler {
    return &Handler{
        transactionService: transactionService,
        logger:            logger,
    }
}

func (h *Handler) CreateTransaction(c *gin.Context) {
    var req CreateTransactionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Error("Invalid request body", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Get user ID from context (set by auth middleware)
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    req.UserID = userID.(uuid.UUID)

    transaction, err := h.transactionService.CreateTransaction(c.Request.Context(), req)
    if err != nil {
        h.logger.Error("Failed to create transaction", err)
        
        if appErr, ok := err.(*errors.AppError); ok {
            c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
            return
        }
        
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "data":    transaction,
    })
}

func (h *Handler) GetWorkspaceTransactions(c *gin.Context) {
    workspaceIDStr := c.Param("workspaceId")
    workspaceID, err := uuid.Parse(workspaceIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
        return
    }

    // Parse query parameters for filters
    filters := parseTransactionFilters(c)

    transactions, err := h.transactionService.GetWorkspaceTransactions(
        c.Request.Context(),
        workspaceID,
        filters,
    )
    if err != nil {
        h.logger.Error("Failed to get workspace transactions", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    transactions,
    })
}

type CreateTransactionRequest struct {
    WorkspaceID     uuid.UUID                    `json:"workspace_id" binding:"required"`
    AccountID       *uuid.UUID                   `json:"account_id,omitempty"`
    CategoryID      *uuid.UUID                   `json:"category_id,omitempty"`
    Description     string                       `json:"description" binding:"required,min=1,max=500"`
    Amount          decimal.Decimal              `json:"amount" binding:"required"`
    TransactionType entities.TransactionType     `json:"transaction_type" binding:"required,oneof=income expense"`
    PaymentMethod   entities.PaymentMethod       `json:"payment_method" binding:"required"`
    TransactionDate time.Time                    `json:"transaction_date" binding:"required"`
    MerchantName    *string                      `json:"merchant_name,omitempty"`
    Location        *string                      `json:"location,omitempty"`
    Notes           *string                      `json:"notes,omitempty"`
    AutoCategorize  bool                         `json:"auto_categorize"`
    UserID          uuid.UUID                    `json:"-"` // Set by middleware
}
```

### 4.3 Validation Standards
```go
// pkg/validation/validator.go
package validation

import (
    "github.com/go-playground/validator/v10"
    "github.com/shopspring/decimal"
)

type Validator struct {
    validator *validator.Validate
}

func NewValidator() *Validator {
    v := validator.New()
    
    // Register custom validators
    v.RegisterValidation("decimal_positive", validateDecimalPositive)
    v.RegisterValidation("decimal_range", validateDecimalRange)
    v.RegisterValidation("currency_code", validateCurrencyCode)
    v.RegisterValidation("workspace_type", validateWorkspaceType)
    
    return &Validator{validator: v}
}

func (v *Validator) Validate(i interface{}) error {
    return v.validator.Struct(i)
}

func validateDecimalPositive(fl validator.FieldLevel) bool {
    if dec, ok := fl.Field().Interface().(decimal.Decimal); ok {
        return dec.GreaterThan(decimal.Zero)
    }
    return false
}

func validateDecimalRange(fl validator.FieldLevel) bool {
    if dec, ok := fl.Field().Interface().(decimal.Decimal); ok {
        min := decimal.NewFromInt(100)        // Minimum Rp 100
        max := decimal.NewFromInt(1000000000) // Maximum Rp 1 billion
        return dec.GreaterThanOrEqual(min) && dec.LessThanOrEqual(max)
    }
    return false
}
```

### 4.4 Authorization Middleware
```go
// api/middleware/auth.go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/expensetracker/pkg/auth"
)

func AuthMiddleware(jwtService auth.JWTService) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        claims, err := jwtService.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        userID, err := uuid.Parse(claims.UserID)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
            c.Abort()
            return
        }

        c.Set("userID", userID)
        c.Set("userEmail", claims.Email)
        c.Next()
    }
}

func WorkspaceAuthMiddleware(workspaceService services.WorkspaceService) gin.HandlerFunc {
    return func(c *gin.Context) {
        workspaceIDStr := c.Param("workspaceId")
        if workspaceIDStr == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Workspace ID required"})
            c.Abort()
            return
        }

        workspaceID, err := uuid.Parse(workspaceIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
            c.Abort()
            return
        }

        userID, exists := c.Get("userID")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        member, err := workspaceService.GetWorkspaceMember(
            c.Request.Context(),
            workspaceID,
            userID.(uuid.UUID),
        )
        if err != nil {
            c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to workspace"})
            c.Abort()
            return
        }

        c.Set("workspaceID", workspaceID)
        c.Set("userRole", member.Role)
        c.Set("permissions", member.Permissions)
        c.Next()
    }
}
```

## 5. API Endpoints Specification

### 5.1 REST API Endpoints
```yaml
# Authentication
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
POST   /api/v1/auth/verify-phone-number
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password

# Workspaces
GET    /api/v1/workspaces
POST   /api/v1/workspaces
GET    /api/v1/workspaces/{id}
PUT    /api/v1/workspaces/{id}
DELETE /api/v1/workspaces/{id}
POST   /api/v1/workspaces/{id}/switch

# Workspace Members -- Skip for now
GET    /api/v1/workspaces/{id}/members
POST   /api/v1/workspaces/{id}/invite
PUT    /api/v1/workspaces/{id}/members/{userId}
DELETE /api/v1/workspaces/{id}/members/{userId}

# Transactions
GET    /api/v1/workspaces/{id}/transactions
POST   /api/v1/workspaces/{id}/transactions
GET    /api/v1/workspaces/{id}/transactions/{transactionId}
PUT    /api/v1/workspaces/{id}/transactions/{transactionId}
DELETE /api/v1/workspaces/{id}/transactions/{transactionId}
POST   /api/v1/workspaces/{id}/transactions/{transactionId}/split -- Skip for now

# Categories
GET    /api/v1/workspaces/{id}/categories
POST   /api/v1/workspaces/{id}/categories
PUT    /api/v1/workspaces/{id}/categories/{categoryId}
DELETE /api/v1/workspaces/{id}/categories/{categoryId}

# Budgets
GET    /api/v1/workspaces/{id}/budgets
POST   /api/v1/workspaces/{id}/budgets
PUT    /api/v1/workspaces/{id}/budgets/{budgetId}
DELETE /api/v1/workspaces/{id}/budgets/{budgetId}

# Settlements -- Skip for now
GET    /api/v1/workspaces/{id}/settlements
POST   /api/v1/workspaces/{id}/settlements/{id}/settle

# Webhooks
GET    /api/v1/webhooks/whatsapp
POST   /api/v1/webhooks/whatsapp

# AI & Processing
POST   /api/v1/ai/analyze-receipt
POST   /api/v1/ai/categorize-transaction
POST   /api/v1/ai/process-whatsapp-message
```

### 5.2 API Response Standards
```go
// pkg/response/response.go
package response

type SuccessResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
    Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorResponse struct {
    Success bool   `json:"success"`
    Error   string `json:"error"`
    Code    string `json:"code,omitempty"`
}

type Meta struct {
    Page       int   `json:"page,omitempty"`
    Limit      int   `json:"limit,omitempty"`
    Total      int64 `json:"total,omitempty"`
    TotalPages int   `json:"total_pages,omitempty"`
}

func Success(data interface{}) SuccessResponse {
    return SuccessResponse{
        Success: true,
        Data:    data,
    }
}

func SuccessWithMeta(data interface{}, meta *Meta) SuccessResponse {
    return SuccessResponse{
        Success: true,
        Data:    data,
        Meta:    meta,
    }
}

func Error(message string) ErrorResponse {
    return ErrorResponse{
        Success: false,
        Error:   message,
    }
}

func ErrorWithCode(message, code string) ErrorResponse {
    return ErrorResponse{
        Success: false,
        Error:   message,
        Code:    code,
    }
}
```

## How it Work in Backend
```
1. Whatsapp send a webhook message
2. VASST EXPENSE API will receive the webhook, publish an event webhook.whatsapp.received
3. Message Handler receive the topic and process it by storing the message from webhook (with various message type: text, image, audio, document)
4. After storing the message or even upload the files, publish new event message.webhook.created
5. AI Model Handler receive the topic and process it by sending the message to selected AI model service.
6. When receiving response from AI model, publish an event aimodel.response.received
7. Message Handler receive the topic and process it by storing the message from AI model
8. After storing the message or even do some actions such as create transactions, updte transactions, create workspace, create the budget, etc. Then publish new event message.aireply.created (This part need confirmation if we need to create MCP but it's still inside the same repository)
9. Whatsapp Handler receive the topic and process it by sending the message to whatsapp API
```

## 6. Frontend Architecture

### 6.1 Next.js Project Structure
```
src/
├── app/                        # App Router (Next.js 14)
│   ├── (auth)/                 # Route groups
│   │   ├── login/
│   │   │   └── page.tsx
│   │   └── register/
│   │       └── page.tsx
│   ├── (dashboard)/
│   │   ├── [workspaceId]/
│   │   │   ├── page.tsx        # Workspace Dashboard
│   │   │   ├── transactions/
│   │   │   │   ├── page.tsx
│   │   │   │   └── [id]/
│   │   │   │       └── page.tsx
│   │   │   ├── budgets/
│   │   │   │   └── page.tsx
│   │   │   ├── accounts/
│   │   │   │   └── page.tsx
│   │   │   ├── members/
│   │   │   │   └── page.tsx
│   │   │   └── settings/
│   │   │       └── page.tsx
│   │   └── workspaces/
│   │       ├── page.tsx        # Workspace Selection
│   │       └── create/
│   │           └── page.tsx
│   ├── api/                    # API Routes
│   │   ├── auth/
│   │   │   └── route.ts
│   │   └── upload/
│   │       └── route.ts
│   ├── globals.css
│   ├── layout.tsx              # Root Layout
│   ├── loading.tsx             # Global Loading UI
│   ├── error.tsx               # Global Error UI
│   └── not-found.tsx
│
├── components/                 # Reusable Components
│   ├── ui/                     # shadcn/ui components
│   │   ├── button.tsx
│   │   ├── input.tsx
│   │   ├── dialog.tsx
│   │   └── ...
│   ├── forms/                  # Form components
│   │   ├── transaction-form.tsx
│   │   ├── budget-form.tsx
│   │   └── workspace-form.tsx
│   ├── charts/                 # Chart components
│   │   ├── expense-chart.tsx
│   │   ├── budget-progress.tsx
│   │   └── spending-trends.tsx
│   ├── workspace/              # Workspace-specific
│   │   ├── workspace-switcher.tsx
│   │   ├── member-list.tsx
│   │   └── invite-modal.tsx
│   ├── transaction/            # Transaction-specific
│   │   ├── transaction-card.tsx
│   │   ├── split-modal.tsx
│   │   └── receipt-upload.tsx
│   └── layout/                 # Layout components
│       ├── sidebar.tsx
│       ├── header.tsx
│       └── mobile-nav.tsx
│
├── lib/                        # Utilities & Configuration
│   ├── api.ts                  # API client configuration
│   ├── auth.ts                 # Authentication utilities
│   ├── utils.ts                # General utilities
│   ├── validations.ts          # Zod schemas
│   ├── constants.ts            # App constants
│   └── hooks/                  # Custom React hooks
│       ├── use-workspace.ts
│       ├── use-transactions.ts
│       └── use-auth.ts
│
├── stores/                     # Zustand stores
│   ├── auth-store.ts
│   ├── workspace-store.ts
│   ├── transaction-store.ts
│   ├── budget-store.ts
│   └── ui-store.ts
│
├── types/                      # TypeScript type definitions
│   ├── api.ts                  # API response types
│   ├── workspace.ts
│   ├── transaction.ts
│   ├── budget.ts
│   └── index.ts
│
└── middleware.ts               # Next.js middleware
```

### 6.2 State Management with Zustand
```typescript
// stores/workspace-store.ts
import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';
import { Workspace, WorkspaceMember } from '@/types';

interface WorkspaceState {
  workspaces: Workspace[];
  activeWorkspace: Workspace | null;
  members: Record<string, WorkspaceMember[]>;
  loading: boolean;
  error: string | null;
  
  // Actions
  fetchWorkspaces: () => Promise<void>;
  setActiveWorkspace: (workspace: Workspace) => void;
  createWorkspace: (data: CreateWorkspaceInput) => Promise<Workspace>;
  updateWorkspace: (id: string, data: UpdateWorkspaceInput) => Promise<void>;
  deleteWorkspace: (id: string) => Promise<void>;
  inviteMember: (workspaceId: string, email: string, role: string) => Promise<void>;
  removeMember: (workspaceId: string, memberId: string) => Promise<void>;
}

export const useWorkspaceStore = create<WorkspaceState>()(
  devtools(
    persist(
      (set, get) => ({
        workspaces: [],
        activeWorkspace: null,
        members: {},
        loading: false,
        error: null,

        fetchWorkspaces: async () => {
          set({ loading: true, error: null });
          try {
            const response = await api.get('/workspaces');
            const workspaces = response.data.data;
            set({ 
              workspaces, 
              activeWorkspace: workspaces[0] || null,
              loading: false 
            });
          } catch (error) {
            set({ 
              error: error instanceof Error ? error.message : 'Failed to fetch workspaces',
              loading: false 
            });
          }
        },

        setActiveWorkspace: (workspace) => {
          set({ activeWorkspace: workspace });
        },

        createWorkspace: async (data) => {
          set({ loading: true, error: null });
          try {
            const response = await api.post('/workspaces', data);
            const newWorkspace = response.data.data;
            set((state) => ({
              workspaces: [...state.workspaces, newWorkspace],
              activeWorkspace: newWorkspace,
              loading: false
            }));
            return newWorkspace;
          } catch (error) {
            set({ 
              error: error instanceof Error ? error.message : 'Failed to create workspace',
              loading: false 
            });
            throw error;
          }
        },

        // ... other actions
      }),
      {
        name: 'workspace-store',
        partialize: (state) => ({ 
          activeWorkspace: state.activeWorkspace 
        }),
      }
    ),
    { name: 'workspace-store' }
  )
);
```

### 6.3 Form Validation with Zod
```typescript
// lib/validations.ts
import { z } from 'zod';

export const createTransactionSchema = z.object({
  workspaceId: z.string().uuid('Invalid workspace ID'),
  accountId: z.string().uuid().optional(),
  categoryId: z.string().uuid().optional(),
  description: z
    .string()
    .min(1, 'Description is required')
    .max(500, 'Description too long'),
  amount: z
    .number()
    .min(100, 'Minimum amount is Rp 100')
    .max(1000000000, 'Maximum amount is Rp 1 billion'),
  transactionType: z.enum(['income', 'expense']),
  paymentMethod: z.enum(['cash', 'debit', 'credit', 'digital']),
  transactionDate: z.date(),
  merchantName: z.string().max(255).optional(),
  location: z.string().max(500).optional(),
  notes: z.string().max(1000).optional(),
  autoCategorize: z.boolean().default(false),
});

export const createWorkspaceSchema = z.object({
  name: z
    .string()
    .min(1, 'Workspace name is required')
    .max(100, 'Name too long'),
  description: z.string().max(500).optional(),
  type: z.enum(['personal', 'business', 'event', 'travel', 'project', 'shared']),
  currency: z.string().length(3).default('IDR'),
  icon: z.string().default('folder'),
  colorCode: z.string().regex(/^#[0-9A-F]{6}$/i).default('#3B82F6'),
});

export const splitTransactionSchema = z.object({
  transactionId: z.string().uuid(),
  splitType: z.enum(['equal', 'percentage', 'exact']),
  participants: z.array(z.object({
    userId: z.string().uuid(),
    amount: z.number().optional(),
    percentage: z.number().min(0).max(100).optional(),
  })).min(2, 'At least 2 participants required'),
});

export type CreateTransactionInput = z.infer<typeof createTransactionSchema>;
export type CreateWorkspaceInput = z.infer<typeof createWorkspaceSchema>;
export type SplitTransactionInput = z.infer<typeof splitTransactionSchema>;
```

### 6.4 API Client Configuration
```typescript
// lib/api.ts
import axios, { AxiosResponse } from 'axios';
import { useAuthStore } from '@/stores/auth-store';

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
  timeout: 10000,
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response: AxiosResponse) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      try {
        const refreshToken = useAuthStore.getState().refreshToken;
        const response = await axios.post('/auth/refresh', {
          refresh_token: refreshToken,
        });
        
        const { token } = response.data.data;
        useAuthStore.getState().setToken(token);
        
        originalRequest.headers.Authorization = `Bearer ${token}`;
        return api(originalRequest);
      } catch (refreshError) {
        useAuthStore.getState().logout();
        window.location.href = '/login';
      }
    }

    return Promise.reject(error);
  }
);

export default api;

// API service functions
export const transactionApi = {
  getWorkspaceTransactions: (workspaceId: string, params?: any) =>
    api.get(`/workspaces/${workspaceId}/transactions`, { params }),
  
  createTransaction: (workspaceId: string, data: CreateTransactionInput) =>
    api.post(`/workspaces/${workspaceId}/transactions`, data),
  
  updateTransaction: (workspaceId: string, transactionId: string, data: any) =>
    api.put(`/workspaces/${workspaceId}/transactions/${transactionId}`, data),
  
  deleteTransaction: (workspaceId: string, transactionId: string) =>
    api.delete(`/workspaces/${workspaceId}/transactions/${transactionId}`),
  
  splitTransaction: (workspaceId: string, transactionId: string, data: SplitTransactionInput) =>
    api.post(`/workspaces/${workspaceId}/transactions/${transactionId}/split`, data),
};

export const workspaceApi = {
  getWorkspaces: () => api.get('/workspaces'),
  createWorkspace: (data: CreateWorkspaceInput) => api.post('/workspaces', data),
  updateWorkspace: (id: string, data: any) => api.put(`/workspaces/${id}`, data),
  deleteWorkspace: (id: string) => api.delete(`/workspaces/${id}`),
  getMembers: (id: string) => api.get(`/workspaces/${id}/members`),
  inviteMember: (id: string, data: { email: string; role: string }) =>
    api.post(`/workspaces/${id}/invite`, data),
};
```