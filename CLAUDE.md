# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

### Core Commands
- **Build API**: `make build` or `go build -o vasst-expense-api ./cmd/api/`
- **Run API**: `make run` or `./vasst-expense-api`
- **Build & Run**: `make all`
- **Run API directly**: `go run cmd/api/main.go`
- **Run Worker**: `go run cmd/worker/main.go`
- **Install dependencies**: `go mod download`
- **Run tests**: `go test ./...`

### Development Setup
1. Copy environment: `cp docs/sample.env .env`
2. Install dependencies: `go mod download`
3. Run migrations (manual - check migrations/ directory)
4. Start API: `go run cmd/api/main.go` (default port 8080)
5. Start worker: `go run cmd/worker/main.go`

## Architecture Overview

### Event-Driven Dual-Service Architecture
This is a **dual-service architecture** with event-driven communication:

- **API Service** (`cmd/api/`): HTTP REST API server handling client requests
- **Worker Service** (`cmd/worker/`): Background event processing and AI automation
- **Event Flow**: HTTP Request → Event Publisher → Google Cloud Pub/Sub → Worker → AI Processing

### Key Architectural Components

**Multi-Workspace Design**: Each user can create multiple workspaces (personal, business, shared) with isolated data and configurations. Always consider workspace isolation when modifying database queries or API endpoints.

**Event-Driven Pattern**: 
- Events published via `internal/events/publisher.go`
- Workers consume events in `internal/workers/`
- Event handlers in `internal/events/handlers/`

**AI Integration**: Dual provider support (OpenAI + Google Gemini) with user-level configuration. AI services are in `internal/services/openai_service.go` and `internal/services/gemini_service.go`.

### Directory Structure

- `cmd/api/` - API server entry point
- `cmd/worker/` - Worker service entry point  
- `internal/api/` - API application setup
- `internal/subscriber/` - Worker application setup
- `internal/controller/http/v0/` - Admin/setup endpoints
- `internal/controller/http/v1/` - Operational endpoints
- `internal/entities/` - Domain models and database entities
- `internal/events/` - Event publishing and handling
- `internal/workers/` - Background workers for AI processing
- `internal/services/` - Business logic services
- `internal/repositories/` - Database access layer
- `internal/middleware/` - HTTP middleware (auth, rate limiting, tenant isolation)
- `internal/pubsub/` - Google Cloud Pub/Sub integration
- `migrations/` - Database migrations

### Database Schema

Uses PostgreSQL with `vasst_expense` schema. Key tables:
- `users` - User accounts with authentication
- `workspaces` - Multi-workspace design for personal/business/shared expense tracking
- `workspace_members` - Collaborative workspace membership
- `accounts` - Financial accounts (bank, credit, cash) per workspace
- `transactions` - Expense/income transactions with AI categorization and bill splitting
- `user_categories` - Custom expense categories per workspace
- `budgets` - Budget management per workspace
- `transaction_splits` - Bill splitting functionality
- `settlements` - Who owes whom tracking
- `documents` - Receipt and document storage with AI analysis
- `scheduler_tasks` - Recurring transactions and reminders
- `webhook_urls` - VASST communication agent integration

### Authentication & Security

- **JWT Authentication**: Tokens contain user context for multi-workspace access
- **API Key Authentication**: User-level API keys for external integrations
- **Multi-workspace Isolation**: Middleware enforces workspace-level data separation
- **Rate Limiting**: Per-user rate limits via middleware

### Communication Channels

Each channel has dedicated processors in `internal/events/handlers/processors/`:
- **WhatsApp Business API**: Webhook handling + message sending
- **Email**: Email communication integration

### Testing

- Test files follow `*_test.go` pattern
- Found in: `internal/utils/xcontext/`, `internal/utils/redis/`, `internal/utils/httpclient/`, `internal/utils/healthcheck/`
- Run with: `go test ./...`

## API Documentation

- Swagger available at `/swagger/index.html` when running
- Health check at `/health-check`
- Base URLs:
  - Development: `http://localhost:8888/v1`
  - Production: `https://api.hi-emma.com/v1`

## Important Notes

- Always test both API and Worker services when making changes
- Consider user isolation for all database operations
- AI processing happens asynchronously via events
- Environment variables critical for Google Cloud services (Pub/Sub, Storage)
- Database uses UUID for primary keys in most tables
- Event-driven architecture means HTTP responses may be immediate while processing happens in background
- Auto create transaction by scanning documents, multi-workspace features, Bill splitting are core to the expense tracking functionality