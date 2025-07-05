### 9.1 Enhanced MCP Server Configuration

```json
{
  "name": "expense-tracker-mcp-v2",
  "version": "2.0.0",
  "description": "Enhanced MCP server for Multi-Workspace ExpenseTracker with VASST integration",
  "tools": [
    {
      "name": "register_vasst_webhook",
      "description": "Register webhook with VASST communication platform",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"},
          "callback_url": {"type": "string"},
          "events": {"type": "array", "items": {"type": "string"}}
        },
        "required": ["user_id", "callback_url"]
      }
    },
    {
      "name": "send_vasst_callback",
      "description": "Send callback to VASST platform",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "event_type": {"type": "string"},
          "payload": {"type": "object"},
          "priority": {"type": "string", "enum": ["low", "normal", "high"]}
        },
        "required": ["user_id", "event_type", "payload"]
      }
    },
    {
      "name": "create_recurring_transaction",
      "description": "Create a recurring transaction with scheduler",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"},
          "transaction_template": {"type": "object"},
          "schedule_config": {
            "type": "object",
            "properties": {
              "schedule_type": {"type": "string", "enum": ["daily", "weekly", "monthly", "yearly"]},
              "interval": {"type": "integer"},
              "day_of_month": {"type": "integer"},
              "day_of_week": {"type": "integer"},
              "end_date": {"type": "string"}
            }
          }
        },
        "required": ["user_id", "workspace_id", "transaction_template", "schedule_config"]
      }
    },
    {
      "name": "create_budget_reminder",
      "description": "Create automated budget reminder",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"},
          "budget_id": {"type": "string"},
          "reminder_thresholds": {"type": "array", "items": {"type": "number"}},
          "notification_channels": {"type": "array", "items": {"type": "string"}}
        },
        "required": ["user_id", "workspace_id", "budget_id"]
      }
    },
    {
      "name": "apply_predefined_categories",
      "description": "Apply predefined categories to workspace",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"},
          "workspace_type": {"type": "string"},
          "locale": {"type": "string", "default": "id-ID"}
        },
        "required": ["user_id", "workspace_id", "workspace_type"]
      }
    },
    {
      "name": "auto_tag_transaction",
      "description": "Automatically tag transaction based on predefined rules",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "transaction_id": {"type": "string"},
          "description": {"type": "string"},
          "amount": {"type": "number"},
          "merchant": {"type": "string"},
          "category": {"type": "string"}
        },
        "required": ["user_id", "transaction_id"]
      }
    },
    {
      "name": "get_user_workspaces",
      "description": "Get all workspaces for a user",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"}
        },
        "required": ["user_id"]
      }
    },
    {
      "name": "switch_workspace",
      "description": "Change user's active workspace",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"}
        },
        "required": ["user_id", "workspace_id"]
      }
    },
    {
      "name": "create_workspace",
      "description": "Create a new workspace",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "name": {"type": "string"},
          "type": {"type": "string", "enum": ["personal", "business", "event", "travel", "project", "shared"]},
          "description": {"type": "string"}
        },
        "required": ["user_id", "name", "type"]
      }
    },
    {
      "name": "create_workspace_transaction",
      "description": "Create transaction in specific workspace",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"},
          "amount": {"type": "number"},
          "description": {"type": "string"},
          "category": {"type": "string"},
          "type": {"type": "string", "enum": ["income", "expense"]},
          "payment_method": {"type": "string"},
          "date": {"type": "string"},
          "tags": {"type": "array", "items": {"type": "string"}},
          "is_recurring": {"type": "boolean", "default": false},
          "recurring_config": {"type": "object"}
        },
        "required": ["user_id", "workspace_id", "amount", "description", "type"]
      }
    },
    {
      "name": "split_transaction",
      "description": "Split a transaction among multiple people",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"},
          "transaction_id": {"type": "string"},
          "split_type": {"type": "string", "enum": ["equal", "percentage", "exact"]},
          "participants": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "user_id": {"type": "string"},
                "amount": {"type": "number"},
                "percentage": {"type": "number"}
              }
            }
          }
        },
        "required": ["user_id", "workspace_id", "transaction_id", "split_type", "participants"]
      }
    },
    {
      "name": "get_workspace_summary",
      "description": "Get financial summary for a workspace",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"},
          "month": {"type": "string"},
          "include_tags": {"type": "boolean", "default": false},
          "include_settlements": {"type": "boolean", "default": false}
        },
        "required": ["user_id", "workspace_id"]
      }
    },
    {
      "name": "get_settlement_status",
      "description": "Get who owes whom in a shared workspace",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"}
        },
        "required": ["user_id", "workspace_id"]
      }
    },
    {
      "name": "mark_settlement_paid",
      "description": "Mark a settlement as paid",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "settlement_id": {"type": "string"}
        },
        "required": ["user_id", "settlement_id"]
      }
    },
    {
      "name": "invite_to_workspace",
      "description": "Invite someone to a workspace",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"},
          "email": {"type": "string"},
          "role": {"type": "string", "enum": ["admin", "member", "viewer"]}
        },
        "required": ["user_id", "workspace_id", "email"]
      }
    },
    {
      "name": "process_receipt_in_workspace",
      "description": "Process receipt image for specific workspace",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"},
          "image_url": {"type": "string"},
          "auto_apply_tags": {"type": "boolean", "default": true},
          "auto_categorize": {"type": "boolean", "default": true}
        },
        "required": ["user_id", "workspace_id", "image_url"]
      }
    }
  ]
}
```

### 9.2 Enhanced WhatsApp Webhook Handler

```go
type EnhancedWhatsAppWebhookHandler struct {
    aiService         AIService
    userService       UserService
    vasst             VASSTClient
    webhookService    WebhookService
    schedulerService  SchedulerService
    tagService        TagService
    categoryService   CategoryService
}

func (h *EnhancedWhatsAppWebhookHandler) HandleMessage(ctx context.Context, payload WhatsAppPayload) error {
    // 1. Validate webhook signature
    if !h.validateWebhookSignature(payload) {
        return errors.New("invalid webhook signature")
    }
    
    // 2. Extract user info from phone number
    user, err := h.userService.GetByPhoneNumber(ctx, payload.From)
    if err != nil {
        return h.sendWelcomeMessage(ctx, payload.From)
    }
    
    // 3. Log the callback request
    logID, err := h.logCallbackRequest(ctx, user.ID, payload)
    if err != nil {
        log.Error("Failed to log callback request", err)
    }
    
    // 4. Process message through enhanced AI agent
    response, err := h.aiService.ProcessWhatsAppMessage(ctx, WhatsAppMessage{
        UserID:      user.ID,
        PhoneNumber: payload.From,
        Message:     payload.Body,
        MessageType: payload.Type,
        MediaURL:    payload.MediaURL,
        Context:     h.getUserContext(ctx, user.ID),
    })
    
    if err != nil {
        h.updateCallbackLog(ctx, logID, "failed", err.Error())
        return h.sendErrorMessage(ctx, payload.From)
    }
    
    // 5. Execute any scheduled tasks triggered by the message
    if response.TriggerScheduledTasks {
        for _, taskID := range response.ScheduledTaskIDs {
            go h.schedulerService.ExecuteTaskNow(ctx, taskID)
        }
    }
    
    // 6. Send response back to WhatsApp
    err = h.vasst.SendMessage(ctx, payload.From, response.Message)
    if err != nil {
        h.updateCallbackLog(ctx, logID, "failed", err.Error())
        return err
    }
    
    h.updateCallbackLog(ctx, logID, "success", "")
    return nil
}

func (h *EnhancedWhatsAppWebhookHandler) logCallbackRequest(ctx context.Context, userID string, payload WhatsAppPayload) (string, error) {
    log := &RequestCallbackLog{
        UserID:       userID,
        RequestType:  "whatsapp_message",
        RequestURL:   "/api/v1/whatsapp/webhook",
        RequestBody:  payload,
        Status:       "pending",
        TriggeredBy:  "whatsapp",
    }
    
    return h.webhookService.CreateCallbackLog(ctx, log)
}
```

### 9.3 VASST Callback System

```go
type VASSTCallbackManager struct {
    webhookService WebhookService
    httpClient     *http.Client
}

func (v *VASSTCallbackManager) SendCallback(ctx context.Context, userID string, event string, payload interface{}) error {
    // Get user's registered webhooks for this event
    webhooks, err := v.webhookService.GetActiveWebhooksForEvent(ctx, userID, event)
    if err != nil {
        return err
    }
    
    for _, webhook := range webhooks {
        go v.executeWebhook(ctx, webhook, event, payload)
    }
    
    return nil
}

func (v *VASSTCallbackManager) executeWebhook(ctx context.Context, webhook *WebhookURL, event string, payload interface{}) {
    logID := v.createExecutionLog(ctx, webhook, event, payload)
    
    // Prepare request
    requestBody, _ := json.Marshal(map[string]interface{}{
        "event":     event,
        "payload":   payload,
        "timestamp": time.Now().Unix(),
        "user_id":   webhook.UserID,
    })
    
    // Create HTTP request
    req, err := http.NewRequestWithContext(ctx, "POST", webhook.URL, bytes.NewBuffer(requestBody))
    if err != nil {
        v.updateExecutionLog(ctx, logID, "failed", err.Error(), 0, nil)
        return
    }
    
    // Add headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Webhook-Secret", webhook.SecretToken)
    req.Header.Set("User-Agent", "ExpenseTracker-Webhook/1.0")
    
    // Add custom headers
    for key, value := range webhook.Headers {
        req.Header.Set(key, value.(string))
    }
    
    // Execute request with timeout
    client := &http.Client{Timeout: time.Duration(webhook.TimeoutSeconds) * time.Second}
    start := time.Now()
    resp, err := client.Do(req)
    duration := time.Since(start)
    
    if err != nil {
        v.handleWebhookRetry(ctx, webhook, event, payload, logID, err.Error())
        return
    }
    defer resp.Body.Close()
    
    // Read response
    responseBody, _ := ioutil.ReadAll(resp.Body)
    
    // Update log
    status := "success"
    if resp.StatusCode >= 400 {
        status = "failed"
    }
    
    v.updateExecutionLog(ctx, logID, status, "", resp.StatusCode, responseBody)
    
    // Update webhook last triggered time
    v.webhookService.UpdateLastTriggered(ctx, webhook.ID)
}
```

## 10. Data Seeding Specifications

### 10.1 Predefined Categories (Indonesian Market)

```sql
-- Personal/Home Categories
INSERT INTO predefined_categories (name, name_id, workspace_type, color_code, icon, keywords, locale) VALUES
('Makanan & Minuman', 'food-drinks', 'personal', '#10B981', 'utensils', '["makan", "minum", "restoran", "kafe", "makanan"]', 'id-ID'),
('Transportasi', 'transportation', 'personal', '#3B82F6', 'car', '["bensin", "ojek", "bus", "kereta", "parkir", "toll"]', 'id-ID'),
('Belanja & Lifestyle', 'shopping-lifestyle', 'personal', '#8B5CF6', 'shopping-bag', '["baju", "sepatu", "kosmetik", "fashion"]', 'id-ID'),
('Hiburan', 'entertainment', 'personal', '#F59E0B', 'film', '["bioskop", "konser", "game", "spotify", "netflix"]', 'id-ID'),
('Kesehatan', 'health', 'personal', '#EF4444', 'heart', '["dokter", "obat", "rumah sakit", "vitamin"]', 'id-ID'),
('Pendidikan', 'education', 'personal', '#06B6D4', 'book', '["kursus", "buku", "sekolah", "training"]', 'id-ID'),
('Utilitas & Tagihan', 'utilities-bills', 'personal', '#84CC16', 'zap', '["listrik", "air", "internet", "pulsa", "gas"]', 'id-ID'),
('Investasi & Tabungan', 'investment-savings', 'personal', '#10B981', 'piggy-bank', '["saham", "reksadana", "deposito", "tabungan"]', 'id-ID'),
('Hadiah & Donasi', 'gifts-donations', 'personal', '#F97316', 'gift', '["angpao", "hadiah", "donasi", "sedekah"]', 'id-ID'),

-- Business Categories  
('Gaji & Upah', 'salary-wages', 'business', '#10B981', 'users', '["gaji", "upah", "bonus", "tunjangan"]', 'id-ID'),
('Sewa & Utilitas Kantor', 'office-rent-utilities', 'business', '#6B7280', 'building', '["sewa kantor", "listrik kantor", "internet kantor"]', 'id-ID'),
('Marketing & Promosi', 'marketing-promotion', 'business', '#F59E0B', 'megaphone', '["iklan", "promosi", "sosial media", "google ads"]', 'id-ID'),
('Operasional', 'operational', 'business', '#3B82F6', 'settings', '["alat tulis", "maintenance", "software", "lisensi"]', 'id-ID'),
('Perjalanan Bisnis', 'business-travel', 'business', '#8B5CF6', 'plane', '["tiket", "hotel bisnis", "transport bisnis"]', 'id-ID'),

-- Event Categories
('Venue & Sewa Tempat', 'venue-rental', 'event', '#EF4444', 'map-pin', '["venue", "gedung", "sewa tempat", "dekorasi tempat"]', 'id-ID'),
('Katering & Konsumsi', 'catering-food', 'event', '#F59E0B', 'utensils', '["katering", "snack", "minuman", "kue"]', 'id-ID'),
('Dokumentasi', 'documentation', 'event', '#8B5CF6', 'camera', '["foto", "video", "fotografer", "videografer"]', 'id-ID'),
('Dekorasi & Hiasan', 'decoration', 'event', '#EC4899', 'sparkles', '["bunga", "balon", "dekorasi", "hiasan"]', 'id-ID'),
('Entertainment', 'event-entertainment', 'event', '#10B981', 'music', '["band", "dj", "penyanyi", "entertainment"]', 'id-ID'),

-- Travel Categories
('Transportasi Perjalanan', 'travel-transport', 'travel', '#3B82F6', 'plane', '["tiket pesawat", "tiket kereta", "rental mobil"]', 'id-ID'),
('Akomodasi', 'accommodation', 'travel', '#10B981', 'bed', '["hotel", "penginapan", "airbnb", "villa"]', 'id-ID'),
('Wisata & Aktivitas', 'tourism-activities', 'travel', '#F59E0B', 'camera', '["tiket wisata", "tour", "aktivitas", "museum"]', 'id-ID'),
('Makanan Perjalanan', 'travel-food', 'travel', '#EF4444', 'utensils', '["makan di perjalanan", "street food", "restoran lokal"]', 'id-ID'),
('Oleh-oleh & Souvenir', 'souvenirs', 'travel', '#8B5CF6', 'gift', '["oleh-oleh", "souvenir", "kerajinan", "merchandise"]', 'id-ID');
```

### 10.2 Predefined Tags (Indonesian Context)

```sql
-- Expense Type Tags
INSERT INTO predefined_tags (name, name_id, tag_type, color_code, icon, usage_context, locale) VALUES
('Mendesak', 'urgent', 'priority', '#EF4444', 'alert-triangle', '["all"]', 'id-ID'),
('Bisa Ditunda', 'can-postpone', 'priority', '#84CC16', 'clock', '["all"]', 'id-ID'),
('Investasi Jangka Panjang', 'long-term-investment', 'expense_type', '#10B981', 'trending-up', '["personal", "business"]', 'id-ID'),
('Kebutuhan Pokok', 'basic-needs', 'expense_type', '#F59E0B', 'home', '["personal"]', 'id-ID'),

-- Payment Status Tags
('Lunas', 'paid', 'payment_status', '#10B981', 'check-circle', '["all"]', 'id-ID'),
('Belum Lunas', 'unpaid', 'payment_status', '#EF4444', 'x-circle', '["all"]', 'id-ID'),
('Cicilan', 'installment', 'payment_status', '#F59E0B', 'calendar', '["all"]', 'id-ID'),

-- Location Tags
('Jakarta', 'jakarta', 'location', '#3B82F6', 'map-pin', '["travel", "business"]', 'id-ID'),
('Bandung', 'bandung', 'location', '#8B5CF6', 'map-pin', '["travel", "business"]', 'id-ID'),
('Bali', 'bali', 'location', '#10B981', 'map-pin', '["travel", "event"]', 'id-ID'),
('Luar Negeri', 'international', 'location', '#F59E0B', 'globe', '["travel", "business"]', 'id-ID'),

-- Occasion Tags
('Lebaran', 'lebaran', 'occasion', '#10B981', 'moon', '["personal", "event"]', 'id-ID'),
('Natal', 'christmas', 'occasion', '#EF4444', 'gift', '["personal", "event"]', 'id-ID'),
('Ulang Tahun', 'birthday', 'occasion', '#8B5CF6', 'cake', '["personal", "event"]', 'id-ID'),
('Pernikahan', 'wedding', 'occasion', '#EC4899', 'heart', '["event"]', 'id-ID'),
('Wisuda', 'graduation', 'occasion', '#F59E0B', 'graduation-cap', '["education", "event"]', 'id-ID');
```

### 10.3 Default Scheduler Tasks Templates

```sql
-- Common recurring transaction templates
INSERT INTO scheduler_tasks (task_type, task_name, schedule_type, schedule_config, task_payload, is_system_task) VALUES
('recurring_transaction', 'Monthly Salary Template', 'monthly', '{"day_of_month": 25}', '{"type": "income", "category": "salary", "description": "Monthly Salary"}', true),
('recurring_transaction', 'Monthly Rent Template', 'monthly', '{"day_of_month": 1}', '{"type": "expense", "category": "housing", "description": "Monthly Rent"}', true),
('recurring_transaction', 'Weekly Groceries Template', 'weekly', '{"day_of_week": 6}', '{"type": "expense", "category": "food", "description": "Weekly Groceries"}', true),

-- Budget reminder templates
('budget_reminder', 'Monthly Budget Check', 'monthly', '{"day_of_month": 15}', '{"reminder_type": "budget_review", "thresholds": [0.8, 0.9, 1.0]}', true),
('budget_reminder', 'Weekly Budget Alert', 'weekly', '{"day_of_week": 0}', '{"reminder_type": "weekly_summary", "include_projections": true}', true),

-- Settlement reminder templates
('settlement_reminder', 'Weekly Settlement Check', 'weekly', '{"day_of_week": 0}', '{"reminder_type": "outstanding_settlements", "min_amount": 50000}', true),
('settlement_reminder', 'Monthly Settlement Summary', 'monthly', '{"day_of_month": 1}', '{"reminder_type": "monthly_settlement_summary"}', true);
```

## 11. Implementation Timeline (Updated)

### Phase 1 (Weeks 1-6): Enhanced Foundation
- Database setup with all enhanced tables
- Webhook management system
- Predefined data seeding (categories, tags)
- Basic scheduler infrastructure
- VASST integration setup

### Phase 2 (Weeks 7-12): Core Features + Enhanced Collaboration
- Multi-workspace transaction management
- Enhanced tagging and categorization system
- Bill splitting with automatic settlement calculation
- Scheduler task management
- Webhook callback logging

### Phase 3 (Weeks 13-18): AI Integration + Advanced WhatsApp
- Enhanced AI agent with tag and category suggestions
- Receipt processing with auto-tagging
- VASST platform integration with full callback support
- WhatsApp bot with scheduler task creation
- Automated budget and settlement reminders

### Phase 4 (Weeks 19-24): Advanced Features + Production Ready
- Advanced analytics with tag-based filtering
- Optimal settlement calculations
- Performance optimization and caching
- Security hardening and audit logging
- Advanced notification system with multiple channels

## 12. Success Metrics (Enhanced)

### Technical Metrics
- API response time < 200ms (95th percentile)
- Webhook delivery success rate > 95%
- Scheduler task execution accuracy > 99%
- AI auto-tagging accuracy > 80%
- System uptime > 99.9%

### Business Metrics
- User activation rate > 70%
- Multi-workspace adoption > 40%
- WhatsApp engagement rate > 60%
- Recurring transaction setup rate > 30%
- Tag usage adoption > 50%

### Integration Metrics
- VASST callback success rate > 95%
- Webhook response time < 5 seconds
- Scheduler task completion rate > 98%
- Auto-categorization acceptance rate > 75%

This enhanced PRD now includes comprehensive webhook management for VASST integration, detailed logging capabilities, predefined data structures for the Indonesian market, and a robust scheduler system for automation - making it a complete, production-ready specification for building the expense tracker SaaS platform.
      "# Expense Tracker SaaS - Product Requirements Document (Revised)

## 1. Project Overview

### 1.1 Product Vision
A multi-workspace expense tracking SaaS platform with AI-powered transaction processing and WhatsApp integration for seamless expense management across different contexts (home, business, events, travel, etc.).

### 1.2 Key Objectives
- Build a mobile-optimized expense tracking platform with workspace management
- Support multiple expense contexts (personal, business, events, travel)
- Enable collaborative expense tracking with bill splitting and group management
- Integrate AI-powered receipt processing and transaction categorization
- Provide WhatsApp-based expense management through VASST communication platform
- Support real-time budget tracking and financial insights per workspace

### 1.3 Core Innovation: Multi-Workspace Architecture
Users can create separate "workspaces" for different expense tracking needs:
- **Personal/Home**: Daily household expenses and personal finance
- **Business**: Company expenses, invoices, and business transactions
- **Events/Travel**: Holiday expenses, wedding planning, group trips
- **Projects**: Specific project budgets and expense tracking
- **Shared Groups**: Collaborative spaces for roommates, family, or teams

### 1.4 Architecture Overview
```
WhatsApp Business API ↔ VASST Communication Platform ↔ AI Agent (Gemini) ↔ Golang API Server ↔ Next.js Frontend
                                                                              ↕
                                                                         PostgreSQL Database
```

## 2. Technical Stack

### 2.1 Frontend (Next.js)
- **Framework**: Next.js 14 with App Router
- **Styling**: Tailwind CSS
- **State Management**: Zustand
- **Forms**: React Hook Form + Zod validation
- **Charts**: Chart.js or Recharts
- **Icons**: Lucide React
- **Authentication**: NextAuth.js
- **HTTP Client**: Axios
- **PWA**: next-pwa

### 2.2 Backend (Golang)
- **Framework**: Gin or Fiber
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT tokens
- **File Storage**: AWS S3 or MinIO
- **Queue**: Redis for background jobs
- **Caching**: Redis
- **Validation**: go-playground/validator
- **PDF Processing**: unidoc/unioffice
- **OCR**: Tesseract or cloud OCR APIs

### 2.3 AI Integration
- **Primary AI**: Google Gemini via VASST platform
- **Document Processing**: PDF parsing + OCR
- **NLP**: Intent recognition and entity extraction
- **Image Processing**: Receipt analysis and data extraction

## 3. Database Schema (Revised)

### 3.1 Core Tables

```sql
-- Users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    whatsapp_number VARCHAR(20),
    timezone VARCHAR(50) DEFAULT 'Asia/Jakarta',
    currency VARCHAR(3) DEFAULT 'IDR',
    subscription_tier VARCHAR(20) DEFAULT 'free',
    email_verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    is_active BOOLEAN DEFAULT true
);

-- Workspaces (NEW - Core feature for multi-context tracking)
CREATE TABLE workspaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    workspace_type VARCHAR(50) NOT NULL, -- 'personal', 'business', 'event', 'travel', 'project', 'shared'
    icon VARCHAR(50) DEFAULT 'folder',
    color_code VARCHAR(7) DEFAULT '#3B82F6',
    currency VARCHAR(3) DEFAULT 'IDR',
    timezone VARCHAR(50) DEFAULT 'Asia/Jakarta',
    settings JSONB DEFAULT '{}', -- workspace-specific settings
    is_active BOOLEAN DEFAULT true,
    created_by UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Workspace Members (NEW - For collaborative workspaces)
CREATE TABLE workspace_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'member', -- 'owner', 'admin', 'member', 'viewer'
    permissions JSONB DEFAULT '{}', -- specific permissions
    joined_at TIMESTAMP DEFAULT NOW(),
    is_active BOOLEAN DEFAULT true,
    UNIQUE(workspace_id, user_id)
);

-- Accounts (Modified to support workspace-specific accounts)
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    account_name VARCHAR(100) NOT NULL,
    account_type VARCHAR(20) NOT NULL, -- 'debit', 'credit', 'savings', 'cash', 'shared'
    bank_name VARCHAR(100),
    account_number_masked VARCHAR(20),
    current_balance DECIMAL(15,2) DEFAULT 0,
    credit_limit DECIMAL(15,2), -- For credit cards
    due_date INTEGER, -- Day of month for credit card due date
    currency VARCHAR(3) DEFAULT 'IDR',
    owner_id UUID REFERENCES users(id), -- Account owner (for shared workspaces)
    is_shared BOOLEAN DEFAULT false, -- For group accounts
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Categories (Modified to support workspace-specific categories)
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    color_code VARCHAR(7) DEFAULT '#3B82F6',
    icon VARCHAR(50) DEFAULT 'receipt',
    parent_category_id UUID REFERENCES categories(id),
    is_system_category BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(workspace_id, name) -- Unique category names per workspace
);

-- Budgets (Modified to support workspace-specific budgets)
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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
CREATE TABLE transactions (
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
CREATE TABLE transaction_splits (
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
CREATE TABLE settlements (
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
CREATE TABLE documents (
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
CREATE TABLE ai_analysis_logs (
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
CREATE TABLE user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE, -- NULL for global preferences
    preference_type VARCHAR(50) NOT NULL, -- 'notification', 'display', 'ai', 'export'
    preferences JSONB NOT NULL DEFAULT '{}',
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, workspace_id, preference_type)
);

-- Audit Logs (Enhanced for workspace tracking)
CREATE TABLE audit_logs (
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
CREATE TABLE workspace_invitations (
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
CREATE TABLE webhook_urls (
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
CREATE TABLE request_callback_logs (
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
CREATE TABLE predefined_categories (
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
CREATE TABLE predefined_tags (
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

-- User Categories (Enhanced - User's custom categories)
-- Note: This replaces the previous categories table
CREATE TABLE user_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    predefined_category_id UUID REFERENCES predefined_categories(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    color_code VARCHAR(7) DEFAULT '#3B82F6',
    icon VARCHAR(50) DEFAULT 'receipt',
    parent_category_id UUID REFERENCES user_categories(id),
    is_custom BOOLEAN DEFAULT false, -- true if user created, false if from predefined
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(workspace_id, name)
);

-- User Tags (NEW - User's custom and applied tags)
CREATE TABLE user_tags (
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
CREATE TABLE transaction_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID REFERENCES transactions(id) ON DELETE CASCADE,
    user_tag_id UUID REFERENCES user_tags(id) ON DELETE CASCADE,
    applied_by UUID REFERENCES users(id) ON DELETE SET NULL,
    applied_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(transaction_id, user_tag_id)
);

-- Scheduler Tasks (NEW - For recurring transactions and reminders)
CREATE TABLE scheduler_tasks (
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
CREATE TABLE scheduler_execution_logs (
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

## 4. API Specifications (Enhanced)

### 4.1 Authentication Endpoints

```yaml
# Authentication
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
POST /api/v1/auth/verify-email
POST /api/v1/auth/forgot-password
POST /api/v1/auth/reset-password
```

### 4.2 Workspace Management (NEW)

```yaml
# Workspaces
GET /api/v1/workspaces                          # Get user's workspaces
POST /api/v1/workspaces                         # Create new workspace
GET /api/v1/workspaces/{id}                     # Get workspace details
PUT /api/v1/workspaces/{id}                     # Update workspace
DELETE /api/v1/workspaces/{id}                  # Delete workspace
POST /api/v1/workspaces/{id}/switch             # Switch active workspace

# Workspace Members
GET /api/v1/workspaces/{id}/members             # Get workspace members
POST /api/v1/workspaces/{id}/invite             # Invite user to workspace
PUT /api/v1/workspaces/{id}/members/{user_id}   # Update member role
DELETE /api/v1/workspaces/{id}/members/{user_id} # Remove member

# Workspace Invitations
GET /api/v1/invitations                         # Get user invitations
POST /api/v1/invitations/{token}/accept         # Accept invitation
POST /api/v1/invitations/{token}/decline        # Decline invitation
```

### 4.3 Core API Endpoints (Enhanced)

```yaml
# Dashboard (workspace-specific)
GET /api/v1/workspaces/{id}/dashboard/summary?month=2025-01
GET /api/v1/workspaces/{id}/dashboard/spending-by-category?month=2025-01
GET /api/v1/workspaces/{id}/dashboard/budget-overview?month=2025-01
GET /api/v1/workspaces/{id}/dashboard/group-summary     # For shared workspaces

# Transactions (workspace-specific)
GET /api/v1/workspaces/{id}/transactions?page=1&limit=20&from=2025-01-01&to=2025-01-31
POST /api/v1/workspaces/{id}/transactions
PUT /api/v1/workspaces/{id}/transactions/{transaction_id}
DELETE /api/v1/workspaces/{id}/transactions/{transaction_id}
POST /api/v1/workspaces/{id}/transactions/{transaction_id}/split
GET /api/v1/workspaces/{id}/transactions/{transaction_id}/splits

# Bill Splitting (NEW)
POST /api/v1/workspaces/{id}/transactions/{transaction_id}/split-equal
POST /api/v1/workspaces/{id}/transactions/{transaction_id}/split-custom
PUT /api/v1/workspaces/{id}/transaction-splits/{split_id}/mark-paid
GET /api/v1/workspaces/{id}/settlements                 # Who owes whom
POST /api/v1/workspaces/{id}/settlements/{id}/settle    # Mark settlement as paid

# Categories (workspace-specific)
GET /api/v1/workspaces/{id}/categories
POST /api/v1/workspaces/{id}/categories
PUT /api/v1/workspaces/{id}/categories/{category_id}
DELETE /api/v1/workspaces/{id}/categories/{category_id}

# Budgets (workspace-specific)
GET /api/v1/workspaces/{id}/budgets?month=2025-01
POST /api/v1/workspaces/{id}/budgets
PUT /api/v1/workspaces/{id}/budgets/{budget_id}
DELETE /api/v1/workspaces/{id}/budgets/{budget_id}
GET /api/v1/workspaces/{id}/budgets/{budget_id}/progress

# Accounts (workspace-specific)
GET /api/v1/workspaces/{id}/accounts
POST /api/v1/workspaces/{id}/accounts
PUT /api/v1/workspaces/{id}/accounts/{account_id}
DELETE /api/v1/workspaces/{id}/accounts/{account_id}
GET /api/v1/workspaces/{id}/accounts/{account_id}/transactions

# Documents (workspace-specific)
POST /api/v1/workspaces/{id}/documents/upload
GET /api/v1/workspaces/{id}/documents
GET /api/v1/workspaces/{id}/documents/{document_id}
DELETE /api/v1/workspaces/{id}/documents/{document_id}
POST /api/v1/workspaces/{id}/documents/{document_id}/process

# AI Integration
POST /api/v1/ai/analyze-receipt
POST /api/v1/ai/categorize-transaction
POST /api/v1/ai/process-whatsapp-message

# WhatsApp Integration
POST /api/v1/whatsapp/webhook
POST /api/v1/whatsapp/send-message
```

## 5. AI Agent System Prompt (Enhanced)

### 5.1 Core System Prompt

```
You are ExpenseBot, an intelligent financial assistant integrated with the ExpenseTracker SaaS platform. You help users manage their expenses, budgets, and financial tracking across multiple workspaces through WhatsApp conversations.

## Your Enhanced Capabilities:
1. Manage multiple expense tracking workspaces (personal, business, events, travel)
2. Record transactions in specific workspaces
3. Process receipt images and extract transaction data
4. Handle bill splitting and group expense management
5. Create and manage budgets per workspace
6. Provide workspace-specific financial summaries and insights
7. Track settlements between group members
8. Send budget alerts and notifications per workspace
9. Switch between workspaces during conversations

## Multi-Workspace Context:
Users can have multiple workspaces for different purposes:
- Personal/Home expenses
- Business/Company expenses  
- Event planning (weddings, parties)
- Travel and holidays
- Project-specific budgets
- Shared group expenses

## Workspace Management Commands:
- /workspaces - List all user workspaces
- /switch [workspace] - Change active workspace
- /create-workspace [name] [type] - Create new workspace
- /invite [email] [workspace] - Invite someone to workspace
- /group-summary - Show group expense overview

## Bill Splitting Features:
- Process group receipts and split automatically
- Track who paid what and who owes whom
- Handle different split types (equal, percentage, custom amounts)
- Manage settlements between group members

## API Integration:
You communicate with the ExpenseTracker API at {API_BASE_URL} using workspace-specific endpoints:
- POST /api/v1/workspaces/{id}/transactions - Create workspace transaction
- GET /api/v1/workspaces/{id}/dashboard/summary - Get workspace summary
- POST /api/v1/workspaces/{id}/budgets - Create workspace budget
- GET /api/v1/workspaces - Get user workspaces
- POST /api/v1/workspaces/{id}/transactions/{id}/split - Split transaction

## Context Management:
- Always confirm which workspace user wants to work with
- Remember the active workspace during conversation
- Ask for clarification when workspace context is unclear
- Provide workspace-specific insights and recommendations

## Response Format:
- Always respond in Indonesian (Bahasa Indonesia)
- Use clear, friendly, and professional tone
- Include relevant emojis for better engagement
- Format amounts with proper Indonesian currency formatting (Rp X.XXX.XXX)
- Clearly indicate which workspace you're working with
- Provide actionable insights and suggestions

## User Context:
- User ID: {user_id}
- Phone Number: {phone_number}
- Active Workspace: {active_workspace}
- Available Workspaces: {user_workspaces}
- Timezone: Asia/Jakarta
- Currency: IDR (Indonesian Rupiah)
- Current Month: {current_month}

## Enhanced Commands You Handle:
- /summary [workspace] - Show workspace financial dashboard
- /new-transaction [workspace] - Create new expense/income in workspace
- /new-budget [workspace] - Create new budget in workspace
- /split-bill [amount] [people] - Split a bill among group members
- /who-owes - Show settlement status in shared workspaces
- /switch [workspace] - Change active workspace
- /workspaces - List all workspaces
- /accounts [workspace] - Show workspace accounts
- /budgets [workspace] - Show workspace budgets
- /help - Show available commands
- Receipt images - Process and create transactions in active workspace
- Natural language expense entries with workspace context

## Example Multi-Workspace Interactions:

User: "Bayar makan 150rb"
Bot: "🤔 Di workspace mana transaksi ini?\n\n📁 Workspace tersedia:\n1️⃣ Personal\n2️⃣ Liburan Bali\n3️⃣ Kantor\n\nKetik nomor atau nama workspace, atau gunakan /switch [workspace]"

User: "2"
Bot: "✅ Transaksi berhasil dicatat di workspace Liburan Bali!\n🏖️ Pengeluaran: Rp 150.000\n🍽️ Kategori: Makanan\n📅 Tanggal: {today}\n\nBudget Makanan Liburan: Rp 850.000 tersisa dari Rp 2.000.000"

User: "Split bill restoran 800rb untuk 4 orang"
Bot: "🧾 Bill Splitting - Workspace: Liburan Bali\n\n💰 Total: Rp 800.000\n👥 Dibagi: 4 orang\n💵 Per orang: Rp 200.000\n\n👤 Siapa yang bayar?\n1️⃣ Saya\n2️⃣ Orang lain\n\nSetelah itu saya akan catat siapa saja yang terlibat."

Always verify transaction details and workspace context with users before finalizing, and provide helpful financial insights specific to the workspace being used.
```

### 5.2 Workspace-Specific Intent Handlers

```
## Intent: WORKSPACE_MANAGEMENT
When user wants to manage workspaces:
1. List available workspaces
2. Create new workspace with proper type
3. Switch active workspace context
4. Invite members to shared workspaces

## Intent: SPLIT_TRANSACTION
When user wants to split a bill:
1. Extract total amount and participant count
2. Determine split type (equal, custom, percentage)
3. Identify who paid the bill
4. Create split transaction records
5. Update settlement balances

## Intent: GROUP_SUMMARY
When user requests group financial overview:
1. Get workspace group summary
2. Show total expenses by member
3. Display settlement status (who owes whom)
4. Provide group spending insights

## Intent: SETTLEMENT_TRACKING
When user asks about settlements:
1. Query settlement status for workspace
2. Show outstanding balances between members
3. Provide payment instructions
4. Track settlement completion
```

## 6. WhatsApp Command Templates (Enhanced)

### 6.1 Multi-Workspace Commands

```
/workspaces
- Shows all user workspaces
- Indicates active workspace
- Provides workspace switching options

Example Response:
📁 *Workspace Anda*

✅ *Personal* (Aktif)
💰 Saldo: Rp 5.200.000
📊 Budget bulan ini: 85% terpakai

🏢 *Kantor*  
💰 Saldo: Rp 12.500.000
📊 Budget bulan ini: 60% terpakai

🏖️ *Liburan Bali*
👥 4 anggota
💰 Total pengeluaran: Rp 8.500.000
🧾 Hutang tertunggak: Rp 650.000

🎉 *Wedding Planning*
👥 6 anggota  
💰 Budget: Rp 150.000.000
📊 Terpakai: 35%

Ketik `/switch [nama workspace]` untuk pindah workspace
```

```
/switch [workspace]
- Changes active workspace context
- Shows workspace overview after switching

Example Flow:
User: "/switch liburan bali"

Bot: "🏖️ *Workspace berhasil diubah ke 'Liburan Bali'*\n\n📊 *Overview Januari 2025*:\n💰 Total pengeluaran: Rp 8.500.000\n👥 Anggota: Anda, Sarah, Budi, Dini\n\n🧾 *Status Hutang*:\n• Sarah → Anda: Rp 250.000\n• Budi → Anda: Rp 150.000\n• Anda → Dini: Rp 180.000\n\n💡 Semua transaksi selanjutnya akan dicatat di workspace ini."
```

```
/create-workspace [name] [type]
- Creates new workspace
- Sets up basic categories and structure
- Explains workspace features

Example Flow:
User: "/create-workspace Weekend Bandung travel"

Bot: "🆕 *Workspace Baru Berhasil Dibuat!*\n\n🏔️ Nama: Weekend Bandung\n🧳 Jenis: Travel\n📅 Periode: Januari 2025\n\n📋 *Kategori Default*:\n• Transportasi\n• Akomodasi  \n• Makanan & Minuman\n• Wisata & Hiburan\n• Belanja Oleh-oleh\n\n👥 *Fitur Available*:\n• Bill splitting otomatis\n• Tracking hutang antar member\n• Budget per kategori\n• Laporan pengeluaran grup\n\nKetik `/invite [email]` untuk mengundang teman!"
```

### 6.2 Bill Splitting Commands

```
/split-bill [amount] [people]
- Splits bill equally among specified people
- Tracks who paid and who owes
- Updates settlement balances

Example Flow:
User: "/split-bill 800000 4"

Bot: "🧾 *Split Bill - Workspace: Liburan Bali*\n\n💰 Total: Rp 800.000\n👥 Dibagi: 4 orang\n💵 Per orang: Rp 200.000\n\n👤 *Siapa yang bayar?*\n1️⃣ Saya\n2️⃣ Sarah  \n3️⃣ Budi\n4️⃣ Dini\n\nKetik nomor atau nama yang bayar"

User: "1"

Bot: "👥 *Pilih anggota yang terlibat*:\n\n✅ Anda (yang bayar)\n☐ Sarah\n☐ Budi  \n☐ Dini\n\nKetik nama untuk toggle, atau 'semua' untuk pilih semua"

User: "semua"

Bot: "✅ *Bill berhasil di-split!*\n\n🧾 Restoran - Rp 800.000\n👤 Dibayar: Anda\n\n💰 *Yang harus bayar Anda*:\n• Sarah: Rp 200.000\n• Budi: Rp 200.000  \n• Dini: Rp 200.000\n\n📊 *Update Hutang*:\n• Sarah total hutang: Rp 450.000\n• Budi total hutang: Rp 350.000"
```

```
/who-owes
- Shows settlement status between all members
- Displays net balances
- Provides payment suggestions

Example Response:
🧾 *Status Hutang - Liburan Bali*

💰 *Net Balances*:
• Anda: +Rp 850.000 (harus terima)
• Sarah: -Rp 450.000 (harus bayar)
• Budi: -Rp 350.000 (harus bayar)  
• Dini: -Rp 50.000 (harus bayar)

📋 *Rekomendasi Settlement*:
1. Sarah transfer Rp 450.000 → Anda
2. Budi transfer Rp 350.000 → Anda
3. Dini transfer Rp 50.000 → Anda

💡 *Tip*: Gunakan '/settle [nama] [jumlah]' setelah menerima pembayaran
```

### 6.3 Workspace-Specific Commands

```
/summary [workspace]
- Shows financial overview for specific workspace
- Adapts content based on workspace type
- Includes collaboration features for shared workspaces

Business Workspace Example:
📊 *Ringkasan Bisnis - Januari 2025*

💰 *Keuangan*:
• Pendapatan: Rp 85.000.000
• Pengeluaran: Rp 62.000.000
• Profit: Rp 23.000.000 (37%)

📈 *Pengeluaran Terbesar*:
• Gaji Karyawan: Rp 35.000.000 (56%)
• Sewa Kantor: Rp 8.000.000 (13%)
• Marketing: Rp 6.500.000 (10%)

🎯 *Budget Status*:
• Operasional: 78% (Rp 48M/62M)
• Marketing: 65% (Rp 6.5M/10M)

Event Workspace Example:
🎉 *Ringkasan Wedding Planning - Januari*

💰 *Budget Overview*:
• Total Budget: Rp 150.000.000
• Terpakai: Rp 52.500.000 (35%)
• Tersisa: Rp 97.500.000

👥 *Kontributor*:
• Anda: Rp 30.000.000
• Keluarga A: Rp 15.000.000  
• Keluarga B: Rp 7.500.000

📋 *Progress Kategori*:
• Venue: 80% (Rp 40M/50M) ✅
• Catering: 25% (Rp 7.5M/30M)
• Dekorasi: 0% (Rp 0/20M)
• Photography: 50% (Rp 5M/10M)

⏰ *Timeline*: 3 bulan tersisa
```

### 6.4 Enhanced Natural Language Processing

```
## Multi-Workspace Expense Recognition:
- "bayar makan 150rb di workspace liburan"
- "catat pengeluaran 500rb untuk kantor"
- "split bill restoran 800rb untuk 4 orang"
- "masukkan ke budget wedding planning"

## Workspace Context Patterns:
- "pindah ke workspace [name]"
- "buat workspace baru untuk [purpose]"
- "undang [email] ke grup [workspace]"
- "siapa yang masih hutang di [workspace]?"

## Settlement Patterns:
- "[name] sudah bayar 200rb"
- "settle hutang dengan [name]"
- "mark [name] sudah lunas"
- "siapa yang harus bayar berapa?"

## Group Management Patterns:
- "tambah anggota [name] ke grup"
- "remove [name] dari workspace"
- "transfer ownership ke [name]"
- "tutup workspace [name]"
```

## 7. Frontend Specifications (Enhanced for Multi-Workspace)

### 7.1 Enhanced Project Structure

```
src/
├── app/
│   ├── (auth)/
│   │   ├── login/
│   │   └── register/
│   ├── (workspace)/
│   │   ├── [workspaceId]/
│   │   │   ├── page.tsx (Workspace Dashboard)
│   │   │   ├── transactions/
│   │   │   ├── budgets/
│   │   │   ├── accounts/
│   │   │   ├── members/
│   │   │   ├── settlements/
│   │   │   └── settings/
│   │   ├── select/page.tsx (Workspace Selection)
│   │   └── create/page.tsx (Create Workspace)
│   ├── invitations/
│   │   └── [token]/page.tsx (Accept Invitation)
│   ├── globals.css
│   ├── layout.tsx
│   └── loading.tsx
├── components/
│   ├── ui/ (shadcn components)
│   ├── workspace/
│   │   ├── WorkspaceSelector.tsx
│   │   ├── WorkspaceSwitcher.tsx
│   │   ├── MembersList.tsx
│   │   └── InviteModal.tsx
│   ├── splits/
│   │   ├── SplitBillModal.tsx
│   │   ├── SettlementCard.tsx
│   │   └── SplitHistory.tsx
│   ├── forms/
│   ├── charts/
│   └── layouts/
├── lib/
│   ├── api.ts
│   ├── auth.ts
│   ├── workspace.ts
│   ├── utils.ts
│   └── validations.ts
├── stores/
│   ├── auth-store.ts
│   ├── workspace-store.ts
│   ├── transaction-store.ts
│   ├── budget-store.ts
│   └── settlement-store.ts
└── types/
    ├── workspace.ts
    ├── transaction.ts
    └── index.ts
```

### 7.2 Key Enhanced Components

```typescript
// Workspace Selector Component
interface WorkspaceSelectorProps {
  workspaces: Workspace[];
  activeWorkspace: Workspace;
  onWorkspaceChange: (workspace: Workspace) => void;
  onCreateNew: () => void;
}

// Split Bill Component
interface SplitBillProps {
  transaction: Transaction;
  members: WorkspaceMember[];
  onSplitComplete: (splits: TransactionSplit[]) => void;
}

// Settlement Dashboard Component
interface SettlementDashboardProps {
  workspaceId: string;
  settlements: Settlement[];
  members: WorkspaceMember[];
  onSettleComplete: (settlementId: string) => void;
}

// Workspace Dashboard Component
interface WorkspaceDashboardProps {
  workspace: Workspace;
  isShared: boolean;
  members?: WorkspaceMember[];
  settlements?: Settlement[];
}

// Member Management Component
interface MemberManagementProps {
  workspace: Workspace;
  members: WorkspaceMember[];
  canInvite: boolean;
  canManageRoles: boolean;
  onInvite: (email: string, role: string) => void;
  onRoleChange: (memberId: string, role: string) => void;
}
```

### 7.3 Enhanced State Management (Zustand)

```typescript
// Workspace Store
interface WorkspaceStore {
  workspaces: Workspace[];
  activeWorkspace: Workspace | null;
  members: Record<string, WorkspaceMember[]>;
  loading: boolean;
  
  fetchWorkspaces: () => Promise<void>;
  createWorkspace: (data: CreateWorkspaceInput) => Promise<Workspace>;
  updateWorkspace: (id: string, data: UpdateWorkspaceInput) => Promise<void>;
  deleteWorkspace: (id: string) => Promise<void>;
  switchWorkspace: (workspace: Workspace) => void;
  
  // Member management
  inviteMember: (workspaceId: string, email: string, role: string) => Promise<void>;
  removeMember: (workspaceId: string, memberId: string) => Promise<void>;
  updateMemberRole: (workspaceId: string, memberId: string, role: string) => Promise<void>;
}

// Settlement Store
interface SettlementStore {
  settlements: Record<string, Settlement[]>; // by workspace ID
  loading: boolean;
  
  fetchSettlements: (workspaceId: string) => Promise<void>;
  createSettlement: (data: CreateSettlementInput) => Promise<void>;
  markSettled: (settlementId: string) => Promise<void>;
  getBalanceSummary: (workspaceId: string) => BalanceSummary;
}

// Enhanced Transaction Store
interface TransactionStore {
  transactions: Record<string, Transaction[]>; // by workspace ID
  splits: Record<string, TransactionSplit[]>; // by transaction ID
  loading: boolean;
  
  fetchTransactions: (workspaceId: string) => Promise<void>;
  createTransaction: (workspaceId: string, data: TransactionInput) => Promise<void>;
  splitTransaction: (transactionId: string, splits: SplitInput[]) => Promise<void>;
  markSplitPaid: (splitId: string) => Promise<void>;
  
  // Filters
  filters: Record<string, TransactionFilters>; // by workspace ID
  setFilters: (workspaceId: string, filters: TransactionFilters) => void;
}
```

## 8. Backend Specifications (Enhanced for Multi-Workspace)

### 8.1 Enhanced Project Structure

```
cmd/
├── api/
│   └── main.go
internal/
├── api/
│   ├── handlers/
│   │   ├── auth/
│   │   ├── workspace/
│   │   ├── transaction/
│   │   ├── budget/
│   │   ├── settlement/
│   │   └── whatsapp/
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── workspace.go
│   │   └── ratelimit.go
│   └── routes/
├── domain/
│   ├── entities/
│   │   ├── workspace.go
│   │   ├── transaction.go
│   │   ├── settlement.go
│   │   └── split.go
│   └── repositories/
│       ├── workspace_repo.go
│       ├── transaction_repo.go
│       └── settlement_repo.go
├── infrastructure/
│   ├── database/
│   ├── cache/
│   └── storage/
├── services/
│   ├── workspace_service.go
│   ├── transaction_service.go
│   ├── settlement_service.go
│   ├── ai_service.go
│   └── whatsapp_service.go
├── utils/
└── config/
```

### 8.2 Enhanced Key Services

```go
// Workspace Service
type WorkspaceService interface {
    CreateWorkspace(ctx context.Context, req CreateWorkspaceRequest) (*Workspace, error)
    GetUserWorkspaces(ctx context.Context, userID string) ([]*Workspace, error)
    GetWorkspaceDetails(ctx context.Context, workspaceID string, userID string) (*WorkspaceDetails, error)
    UpdateWorkspace(ctx context.Context, workspaceID string, req UpdateWorkspaceRequest) (*Workspace, error)
    DeleteWorkspace(ctx context.Context, workspaceID string, userID string) error
    
    // Member management
    InviteMember(ctx context.Context, workspaceID string, email string, role string, invitedBy string) error
    AcceptInvitation(ctx context.Context, token string, userID string) error
    RemoveMember(ctx context.Context, workspaceID string, memberID string, removedBy string) error
    UpdateMemberRole(ctx context.Context, workspaceID string, memberID string, role string) error
    GetWorkspaceMembers(ctx context.Context, workspaceID string) ([]*WorkspaceMember, error)
}

```go
// Webhook Service (NEW)
type WebhookService interface {
    CreateWebhook(ctx context.Context, req CreateWebhookRequest) (*WebhookURL, error)
    GetUserWebhooks(ctx context.Context, userID string) ([]*WebhookURL, error)
    UpdateWebhook(ctx context.Context, webhookID string, req UpdateWebhookRequest) (*WebhookURL, error)
    DeleteWebhook(ctx context.Context, webhookID string) error
    TestWebhook(ctx context.Context, webhookID string) error
    TriggerWebhook(ctx context.Context, webhookID string, payload interface{}) error
    GetWebhookLogs(ctx context.Context, webhookID string, filters LogFilters) ([]*RequestCallbackLog, error)
}

// VASST Communication Service (NEW)
type VASSTCommService interface {
    RegisterWebhook(ctx context.Context, userID string, workspaceID string, callbackURL string) error
    SendCallback(ctx context.Context, userID string, event string, payload interface{}) error
    ProcessVASSTCallback(ctx context.Context, payload VASSTCallbackPayload) error
    UpdateWebhookConfig(ctx context.Context, userID string, config VASSTWebhookConfig) error
}

// Predefined Data Service (NEW)
type PredefinedDataService interface {
    GetPredefinedCategories(ctx context.Context, workspaceType string, locale string) ([]*PredefinedCategory, error)
    GetPredefinedTags(ctx context.Context, tagType string, locale string) ([]*PredefinedTag, error)
    CreateUserCategoriesFromPredefined(ctx context.Context, workspaceID string, categoryIDs []string) ([]*UserCategory, error)
    CreateUserTagsFromPredefined(ctx context.Context, workspaceID string, tagIDs []string) ([]*UserTag, error)
    SuggestCategoriesForWorkspace(ctx context.Context, workspaceType string) ([]*PredefinedCategory, error)
    SuggestTagsForTransaction(ctx context.Context, description string, amount float64) ([]*PredefinedTag, error)
}

// Enhanced Category Service
type CategoryService interface {
    GetWorkspaceCategories(ctx context.Context, workspaceID string) ([]*UserCategory, error)
    CreateCategory(ctx context.Context, req CreateCategoryRequest) (*UserCategory, error)
    UpdateCategory(ctx context.Context, categoryID string, req UpdateCategoryRequest) (*UserCategory, error)
    DeleteCategory(ctx context.Context, categoryID string) error
    GetCategoryUsageStats(ctx context.Context, workspaceID string) ([]*CategoryUsage, error)
}

// Tag Service (NEW)
type TagService interface {
    GetWorkspaceTags(ctx context.Context, workspaceID string) ([]*UserTag, error)
    CreateTag(ctx context.Context, req CreateTagRequest) (*UserTag, error)
    UpdateTag(ctx context.Context, tagID string, req UpdateTagRequest) (*UserTag, error)
    DeleteTag(ctx context.Context, tagID string) error
    ApplyTagsToTransaction(ctx context.Context, transactionID string, tagIDs []string) error
    RemoveTagFromTransaction(ctx context.Context, transactionID string, tagID string) error
    GetTagUsageStats(ctx context.Context, workspaceID string) ([]*TagUsage, error)
}

// Scheduler Service (NEW)
type SchedulerService interface {
    CreateSchedulerTask(ctx context.Context, req CreateSchedulerTaskRequest) (*SchedulerTask, error)
    GetWorkspaceSchedulerTasks(ctx context.Context, workspaceID string) ([]*SchedulerTask, error)
    UpdateSchedulerTask(ctx context.Context, taskID string, req UpdateSchedulerTaskRequest) (*SchedulerTask, error)
    DeleteSchedulerTask(ctx context.Context, taskID string) error
    PauseTask(ctx context.Context, taskID string) error
    ResumeTask(ctx context.Context, taskID string) error
    ExecuteTaskNow(ctx context.Context, taskID string) (*SchedulerExecutionLog, error)
    GetTaskExecutions(ctx context.Context, taskID string, filters ExecutionFilters) ([]*SchedulerExecutionLog, error)
    
    // System scheduler methods
    ProcessPendingTasks(ctx context.Context) error
    ExecuteTask(ctx context.Context, task *SchedulerTask) (*SchedulerExecutionLog, error)
}

// Enhanced Transaction Service
type TransactionService interface {
    CreateTransaction(ctx context.Context, workspaceID string, req CreateTransactionRequest) (*Transaction, error)
    GetWorkspaceTransactions(ctx context.Context, workspaceID string, filters TransactionFilters) ([]*Transaction, error)
    UpdateTransaction(ctx context.Context, transactionID string, req UpdateTransactionRequest) (*Transaction, error)
    DeleteTransaction(ctx context.Context, transactionID string) error
    
    // Enhanced tagging
    ApplyTagsToTransaction(ctx context.Context, transactionID string, tagIDs []string) error
    GetTransactionTags(ctx context.Context, transactionID string) ([]*UserTag, error)
    
    // Bill splitting
    SplitTransaction(ctx context.Context, transactionID string, splits []SplitInput) ([]*TransactionSplit, error)
    MarkSplitPaid(ctx context.Context, splitID string) error
    GetTransactionSplits(ctx context.Context, transactionID string) ([]*TransactionSplit, error)
    
    // Recurring transactions
    CreateRecurringTransaction(ctx context.Context, req CreateRecurringTransactionRequest) (*SchedulerTask, error)
    
    // Receipt processing
    ProcessReceiptImage(ctx context.Context, workspaceID string, imageURL string, userID string) (*Transaction, error)
}

// Settlement Service (NEW)
type SettlementService interface {
    GetWorkspaceSettlements(ctx context.Context, workspaceID string) ([]*Settlement, error)
    CreateSettlement(ctx context.Context, req CreateSettlementRequest) (*Settlement, error)
    MarkSettled(ctx context.Context, settlementID string) error
    GetBalanceSummary(ctx context.Context, workspaceID string) (*BalanceSummary, error)
    CalculateOptimalSettlements(ctx context.Context, workspaceID string) ([]*OptimalSettlement, error)
}

// Enhanced AI Service
type AIService interface {
    AnalyzeReceipt(ctx context.Context, imageURL string) (*ReceiptAnalysis, error)
    CategorizeTransaction(ctx context.Context, description string, workspaceType string) (string, float64, error)
    ProcessWhatsAppMessage(ctx context.Context, message WhatsAppMessage) (*AIResponse, error)
    
    // Workspace-specific AI features
    SuggestBudgetCategories(ctx context.Context, workspaceType string) ([]string, error)
    AnalyzeSpendingPattern(ctx context.Context, workspaceID string) (*SpendingAnalysis, error)
    GenerateGroupInsights(ctx context.Context, workspaceID string) (*GroupInsights, error)
}
```

### 8.3 Workspace Authorization Middleware

```go
func WorkspaceAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        workspaceID := c.Param("workspaceId")
        userID := c.GetString("userID")
        
        // Check if user has access to workspace
        member, err := workspaceService.GetWorkspaceMember(c, workspaceID, userID)
        if err != nil {
            c.JSON(403, gin.H{"error": "Access denied to workspace"})
            c.Abort()
            return
        }
        
        // Set workspace context
        c.Set("workspaceID", workspaceID)
        c.Set("userRole", member.Role)
        c.Set("permissions", member.Permissions)
        
        c.Next()
    }
}
```

## 9. VASST Integration Specifications (Enhanced)

### 9.1 Enhanced MCP Server Configuration

```json
{
  "name": "expense-tracker-mcp-v2",
  "version": "2.0.0",
  "description": "Enhanced MCP server for Multi-Workspace ExpenseTracker",
  "tools": [
    {
      "name": "get_user_workspaces",
      "description": "Get all workspaces for a user",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"}
        },
        "required": ["user_id"]
      }
    },
    {
      "name": "switch_workspace",
      "description": "Change user's active workspace",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"}
        },
        "required": ["user_id", "workspace_id"]
      }
    },
    {
      "name": "create_workspace",
      "description": "Create a new workspace",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "name": {"type": "string"},
          "type": {"type": "string", "enum": ["personal", "business", "event", "travel", "project", "shared"]},
          "description": {"type": "string"}
        },
        "required": ["user_id", "name", "type"]
      }
    },
    {
      "name": "create_workspace_transaction",
      "description": "Create transaction in specific workspace",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"},
          "amount": {"type": "number"},
          "description": {"type": "string"},
          "category": {"type": "string"},
          "type": {"type": "string", "enum": ["income", "expense"]},
          "payment_method": {"type": "string"},
          "date": {"type": "string"}
        },
        "required": ["user_id", "workspace_id", "amount", "description", "type"]
      }
    },
    {
      "name": "split_transaction",
      "description": "Split a transaction among multiple people",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"},
          "transaction_id": {"type": "string"},
          "split_type": {"type": "string", "enum": ["equal", "percentage", "exact"]},
          "participants": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "user_id": {"type": "string"},
                "amount": {"type": "number"},
                "percentage": {"type": "number"}
              }
            }
          }
        },
        "required": ["user_id", "workspace_id", "transaction_id", "split_type", "participants"]
      }
    },
    {
      "name": "get_workspace_summary",
      "description": "Get financial summary for a workspace",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"},
          "month": {"type": "string"}
        },
        "required": ["user_id", "workspace_id"]
      }
    },
    {
      "name": "get_settlement_status",
      "description": "Get who owes whom in a shared workspace",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"}
        },
        "required": ["user_id", "workspace_id"]
      }
    },
    {
      "name": "mark_settlement_paid",
      "description": "Mark a settlement as paid",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "settlement_id": {"type": "string"}
        },
        "required": ["user_id", "settlement_id"]
      }
    },
    {
      "name": "invite_to_workspace",
      "description": "Invite someone to a workspace",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"},
          "email": {"type": "string"},
          "role": {"type": "string", "enum": ["admin", "member", "viewer"]}
        },
        "required": ["user_id", "workspace_id", "email"]
      }
    },
    {
      "name": "process_receipt_in_workspace",
      "description": "Process receipt image for specific workspace",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {"type": "string"},
          "workspace_id": {"type": "string"},
          "image_url": {"type": "string"}
        },
        "required": ["user_id", "workspace_id", "image_url"]
      }
    }
  ]
}
```

## 10. Implementation Timeline (Revised)

### Phase 1 (Weeks 1-5): Multi-Workspace Foundation
- Database setup with enhanced schema
- Workspace management API endpoints
- User authentication and workspace authorization
- Basic workspace CRUD operations
- Member invitation system

### Phase 2 (Weeks 6-10): Core Features + Collaboration
- Transaction management per workspace
- Category and budget management per workspace
- Bill splitting functionality
- Settlement tracking system
- Enhanced mobile UI for workspaces

### Phase 3 (Weeks 11-15): AI Integration + WhatsApp
- Enhanced AI agent with workspace context
- Receipt processing per workspace
- VASST platform integration with workspace support
- WhatsApp bot with multi-workspace commands
- Group expense management features

### Phase 4 (Weeks 16-20): Advanced Features + Optimization
- Advanced analytics per workspace type
- Optimal settlement calculations
- Performance optimization
- Security hardening
- Advanced notification system

## 11. Success Metrics (Enhanced)

### Technical Metrics
- API response time < 200ms (95th percentile)
- Mobile app performance score > 90
- AI categorization accuracy > 85% per workspace type
- System uptime > 99.9%
- Workspace switching time < 100ms

### Business Metrics
- User activation rate > 70%
- Multi-workspace adoption > 40%
- WhatsApp engagement rate > 60%
- Bill splitting feature usage > 30%
- Average workspaces per user > 2.5

### Collaboration Metrics
- Workspace invitation acceptance rate > 65%
- Settlement completion rate > 80%
- Group expense tracking engagement > 50%
- Average members per shared workspace > 3

## 12. Competitive Advantages (Enhanced)

### Multi-Context Expense Tracking
- **First platform** to offer workspace-based expense management
- **Seamless switching** between personal, business, and event expenses
- **Context-aware AI** that understands different spending patterns

### Advanced Collaboration Features
- **Real-time bill splitting** with automatic settlement tracking
- **Group budget management** for events and projects
- **Multi-role workspace** management with proper permissions

### WhatsApp-Native Group Management
- **Group expense coordination** through WhatsApp
- **Automatic settlement reminders** via chat
- **Collaborative budget tracking** in group chats

### Indonesian Market Focus
- **Local banking integration** and currency handling
- **Cultural spending patterns** (angpao, group dining, etc.)
- **Indonesian language** natural language processing

This enhanced PRD provides a comprehensive foundation for building a multi-workspace expense tracker that revolutionizes how people manage finances across different contexts, with seamless collaboration and AI-powered automation through WhatsApp integration.