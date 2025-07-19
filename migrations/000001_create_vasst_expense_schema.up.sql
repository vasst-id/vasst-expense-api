-- Create vasst_expense database schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS "vasst_expense";

-- Create extension if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Currency
CREATE TABLE "vasst_expense".currency (
    currency_id SERIAL PRIMARY KEY,
    currency_code VARCHAR(3) NOT NULL,
    currency_name VARCHAR(100) NOT NULL,
    currency_symbol VARCHAR(10) NOT NULL,
    currency_decimal_places INT NOT NULL DEFAULT 2,
    currency_status INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Subcription Plan
CREATE TABLE "vasst_expense".subscription_plan (
    subscription_plan_id SERIAL PRIMARY KEY,
    subscription_plan_name VARCHAR(20) NOT NULL,
    subscription_plan_description TEXT,
    subscription_plan_features JSONB NOT NULL DEFAULT '{}',
    subscription_plan_price DECIMAL(10,2) NOT NULL,
    subscription_plan_currency_id INT NOT NULL REFERENCES "vasst_expense".currency(currency_id),
    subscription_plan_status INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Taxonomy
CREATE TABLE IF NOT EXISTS "vasst_expense".taxonomy (
    taxonomy_id SERIAL PRIMARY KEY,
    label VARCHAR(100) NOT NULL,
    value VARCHAR(100) NOT NULL,
    type VARCHAR(100) NOT NULL,
    type_label VARCHAR(100) NOT NULL,
    status INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Users
CREATE TABLE "vasst_expense".users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(100) UNIQUE,
    phone_number VARCHAR(20) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Jakarta',
    currency_id INT NOT NULL REFERENCES "vasst_expense".currency(currency_id),
    subscription_plan_id INT NOT NULL REFERENCES "vasst_expense".subscription_plan(subscription_plan_id),
    email_verified_at TIMESTAMPTZ,
    phone_verified_at TIMESTAMPTZ,
    status INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Workspaces
CREATE TABLE "vasst_expense".workspaces (
    workspace_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    workspace_type INT NOT NULL, -- '1 - personal', '2 - business', '3 - event', '4 - travel', '5 - project', '6 - shared'
    icon VARCHAR(50) DEFAULT 'folder',
    color_code VARCHAR(7) DEFAULT '#3B82F6',
    currency_id INT NOT NULL REFERENCES "vasst_expense".currency(currency_id),
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Jakarta',
    settings JSONB NOT NULL DEFAULT '{}', -- workspace-specific settings
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID REFERENCES "vasst_expense".users(user_id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
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
    bank_name VARCHAR(50) NOT NULL,
    bank_code VARCHAR(50) NOT NULL,
    bank_logo_url VARCHAR(255) NOT NULL,
    status INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Accounts (Modified to support workspace-specific accounts)
CREATE TABLE "vasst_expense".accounts (
    account_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES "vasst_expense".users(user_id),
    account_name VARCHAR(100) NOT NULL,
    account_type INT NOT NULL, -- '1 - debit', '2 - credit', '3 - savings', '4 - cash', '5 - digital wallet'
    bank_id INT REFERENCES "vasst_expense".banks(bank_id),
    account_number VARCHAR(20),
    current_balance DECIMAL(15,2) NOT NULL DEFAULT 0,
    credit_limit DECIMAL(15,2), -- For credit cards
    due_date INTEGER, -- Day of month for credit card due date
    currency_id INT NOT NULL REFERENCES "vasst_expense".currency(currency_id),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Categories (Modified to support workspace-specific categories)
CREATE TABLE "vasst_expense".categories (
    category_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon VARCHAR(50) DEFAULT 'receipt',
    parent_category_id UUID REFERENCES "vasst_expense".categories(category_id),
    is_system_category BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- User Categories (Enhanced - User's custom categories)
-- Note: This replaces the previous categories table
CREATE TABLE "vasst_expense".user_categories (
    user_category_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES "vasst_expense".users(user_id),
    category_id UUID REFERENCES "vasst_expense".categories(category_id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(50) DEFAULT 'receipt',
    is_custom BOOLEAN NOT NULL DEFAULT false, -- true if user created, false if from predefined
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Budgets (Modified to support workspace-specific budgets)
CREATE TABLE "vasst_expense".budgets (
    budget_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES "vasst_expense".workspaces(workspace_id) ON DELETE SET NULL,
    user_category_id UUID REFERENCES "vasst_expense".user_categories(user_category_id) ON DELETE SET NULL, -- Updated reference
    name VARCHAR(100) NOT NULL,
    budgeted_amount DECIMAL(15,2) NOT NULL,
    period_type INT NOT NULL, -- '1 - weekly', '2 - monthly', '3 - yearly', '4 - one time'
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    spent_amount DECIMAL(15,2) DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID REFERENCES "vasst_expense".users(user_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Transactions (Enhanced for multi-workspace and bill splitting)
CREATE TABLE "vasst_expense".transactions (
    transaction_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES "vasst_expense".workspaces(workspace_id) ON DELETE SET NULL,
    account_id UUID REFERENCES "vasst_expense".accounts(account_id) ON DELETE SET NULL,
    user_category_id UUID REFERENCES "vasst_expense".user_categories(user_category_id) ON DELETE SET NULL, -- Updated reference
    description TEXT NOT NULL,
    amount DECIMAL(15,2) NOT NULL, -- Total transaction amount
    transaction_type INT NOT NULL, -- '1 - income', '2 - expense'
    payment_method INT, -- '1 - debit/qris', '2 - credit', '3 - cash', '4 - transfer'
    transaction_date DATE NOT NULL,
    merchant_name VARCHAR(255),
    location VARCHAR(255),
    notes TEXT,
    receipt_url TEXT,
    
    -- Recurring transactions
    is_recurring BOOLEAN NOT NULL DEFAULT false,
    recurrence_interval INT, -- '1 - daily', '2 - weekly', '3 - monthly', '4 - yearly'
    recurrence_end_date DATE, -- When the recurring transaction ends
    parent_transaction_id UUID REFERENCES "vasst_expense".transactions(transaction_id),
    -- scheduler_task_id UUID REFERENCES scheduler_tasks(scheduler_task_id), -- Link to recurring task (table not defined)
    
    -- AI and processing
    ai_confidence_score DECIMAL(3,2), -- 0.00 to 1.00
    ai_categorized BOOLEAN DEFAULT false,
    
    -- Credit tracking
    credit_status INT, -- '1 - paid', '2 - unpaid' for credit transactions
    
    -- Metadata
    created_by UUID REFERENCES "vasst_expense".users(user_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Transaction Splits (NEW - For bill splitting functionality)
-- CREATE TABLE "vasst_expense".transaction_splits (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     transaction_id UUID REFERENCES transactions(id) ON DELETE CASCADE,
--     user_id UUID REFERENCES users(id) ON DELETE CASCADE,
--     amount DECIMAL(15,2) NOT NULL, -- Amount this user owes/paid
--     percentage DECIMAL(5,2), -- Percentage of total (for percentage splits)
--     shares INTEGER, -- Number of shares (for share-based splits)
--     status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'paid', 'settled'
--     paid_at TIMESTAMP,
--     notes TEXT,
--     created_at TIMESTAMP DEFAULT NOW(),
--     updated_at TIMESTAMP DEFAULT NOW()
-- );

-- Settlements (NEW - For tracking who owes whom)
-- CREATE TABLE "vasst_expense".settlements (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
--     from_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
--     to_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
--     amount DECIMAL(15,2) NOT NULL,
--     description TEXT,
--     transaction_ids UUID[], -- Array of related transaction IDs
--     status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'completed', 'cancelled'
--     settled_at TIMESTAMP,
--     created_at TIMESTAMP DEFAULT NOW(),
--     updated_at TIMESTAMP DEFAULT NOW()
-- );

-- Documents (Enhanced for workspace-specific documents)
CREATE TABLE "vasst_expense".documents (
    document_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES "vasst_expense".workspaces(workspace_id) ON DELETE SET NULL,
    user_id UUID REFERENCES "vasst_expense".users(user_id) ON DELETE SET NULL,
    original_filename VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_size INTEGER NOT NULL,
    document_type INT NOT NULL REFERENCES "vasst_expense".taxonomy(taxonomy_id), -- '1 - bank-statement', '2 - credit card bill', '3 - invoice'
    ai_analysis_result JSONB,
    processing_status INT NOT NULL REFERENCES "vasst_expense".taxonomy(taxonomy_id), -- '1 - pending', '2 - processing', '3 - completed', '4 - failed'
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- AI Analysis Logs (Enhanced)
CREATE TABLE "vasst_expense".ai_analysis_logs (
    ai_analysis_log_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES "vasst_expense".workspaces(workspace_id) ON DELETE SET NULL,
    user_id UUID REFERENCES "vasst_expense".users(user_id) ON DELETE SET NULL,
    document_id UUID REFERENCES "vasst_expense".documents(document_id) ON DELETE SET NULL,
    input_data JSONB,
    output_data JSONB,
    model_used VARCHAR(100),
    confidence_score DECIMAL(3,2),
    processing_time_ms INTEGER,
    error_details JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- User Preferences (Enhanced with workspace preferences)
-- CREATE TABLE "vasst_expense".user_preferences (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     user_id UUID REFERENCES users(id) ON DELETE CASCADE,
--     workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE, -- NULL for global preferences
--     preference_type VARCHAR(50) NOT NULL, -- 'notification', 'display', 'ai', 'export'
--     preferences JSONB NOT NULL DEFAULT '{}',
--     updated_at TIMESTAMP DEFAULT NOW(),
--     UNIQUE(user_id, workspace_id, preference_type)
-- );

-- Audit Logs (Enhanced for workspace tracking)
CREATE TABLE "vasst_expense".audit_logs (
    audit_log_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES "vasst_expense".workspaces(workspace_id) ON DELETE SET NULL,
    user_id UUID REFERENCES "vasst_expense".users(user_id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL, -- 'create', 'update', 'delete', 'view', etc.
    resource_type VARCHAR(50) NOT NULL, -- 'transaction', 'budget', 'document', 'user', 'workspace', etc.
    resource_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    source_type VARCHAR(20) DEFAULT 'web', -- 'web', 'whatsapp', 'api'
    created_at TIMESTAMP DEFAULT NOW()
);

-- Workspace Invitations (NEW - For inviting users to workspaces)
-- CREATE TABLE "vasst_expense".workspace_invitations (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
--     invited_by UUID REFERENCES users(id) ON DELETE CASCADE,
--     email VARCHAR(255) NOT NULL,
--     role VARCHAR(20) DEFAULT 'member',
--     invitation_token VARCHAR(255) UNIQUE NOT NULL,
--     status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'accepted', 'declined', 'expired'
--     expires_at TIMESTAMP NOT NULL,
--     accepted_at TIMESTAMP,
--     created_at TIMESTAMP DEFAULT NOW()
-- );


-- User Tags (NEW - User's custom and applied tags)
CREATE TABLE "vasst_expense".user_tags (
    user_tag_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES "vasst_expense".users(user_id) ON DELETE SET NULL,
    name VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Transaction Tags (NEW - Many-to-many relationship for transaction tagging)
CREATE TABLE "vasst_expense".transaction_tags (
    transaction_tag_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID REFERENCES "vasst_expense".transactions(transaction_id) ON DELETE SET NULL,
    user_tag_id UUID REFERENCES "vasst_expense".user_tags(user_tag_id) ON DELETE SET NULL,
    applied_by UUID REFERENCES "vasst_expense".users(user_id) ON DELETE SET NULL,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Conversations table: Represents a chat thread between the system and a user (1:1 for now)
CREATE TABLE "vasst_expense".conversations (
    conversation_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES "vasst_expense".users(user_id) ON DELETE CASCADE,
    channel VARCHAR(30) NOT NULL DEFAULT 'whatsapp', -- 'whatsapp', 'web', etc.
    is_active BOOLEAN NOT NULL DEFAULT true,
    context TEXT, -- Context of the conversation
    metadata JSONB, -- Metadata of the conversation
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Conversation Messages table: Stores all inbound/outbound messages in a conversation
CREATE TABLE "vasst_expense".messages (
    message_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID REFERENCES "vasst_expense".conversations(conversation_id) ON DELETE CASCADE,
    user_id UUID REFERENCES "vasst_expense".users(user_id), -- Sender (nullable for system/AI)
    sender_type INT NOT NULL, -- '1 - user', '2 - ai', '3 - system', '4 - scheduler'
    direction CHAR(1) NOT NULL, -- 'i - inbound', 'o - outbound'
    message_type INT NOT NULL REFERENCES "vasst_expense".taxonomy(taxonomy_id), -- 'text', 'image', 'document', 'audio', 'video', etc.
    content TEXT, -- Raw text content (if applicable)
    media_url TEXT, -- URL to image/document/audio if applicable
    attachments JSONB, -- Array of attachments (if applicable)
    media_mime_type VARCHAR(100), -- MIME type for media
    transcription TEXT, -- LLM/AI-generated transcription of media (if applicable)
    ai_processed BOOLEAN NOT NULL DEFAULT false, -- Whether LLM/AI has processed this message
    ai_model VARCHAR(50), -- e.g. 'gemini', 'gpt-4', etc.
    ai_confidence_score DECIMAL(3,2), -- 0.00 to 1.00
    related_transaction_id UUID REFERENCES "vasst_expense".transactions(transaction_id), -- If message is linked to a transaction
    scheduled_task_id UUID, -- If generated by a scheduler (future use)
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast lookup of messages by conversation
CREATE INDEX idx_messages_conversation_id ON "vasst_expense".messages(conversation_id);

-- Index for fast lookup of conversations by user
CREATE INDEX idx_conversations_user_id ON "vasst_expense".conversations(user_id);


CREATE TABLE "vasst_expense".verification_codes (
    verification_code_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(20) NOT NULL,
    code VARCHAR(6) NOT NULL,
    code_type VARCHAR(20) NOT NULL, -- 'phone_verification', 'password_reset', etc.
    expires_at TIMESTAMP NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    attempts_count INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for quick lookups
CREATE INDEX idx_verification_codes_phone_type ON "vasst_expense".verification_codes(phone_number, code_type);
CREATE INDEX idx_verification_codes_expires ON "vasst_expense".verification_codes(expires_at);