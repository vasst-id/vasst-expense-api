-- Create vasst_expense database schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS "vasst_expense";

-- Create extension if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

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