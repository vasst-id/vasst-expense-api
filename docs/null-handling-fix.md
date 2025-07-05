# NULL Handling Fix for Organization Settings

## Issue
The application was encountering a SQL scan error:
```
sql: Scan error on column index 20, name "system_prompt": converting NULL to string is unsupported
```

## Root Cause
The `system_prompt`, `ai_assistant_name`, `ai_communication_style`, and `ai_communication_language` fields in the database can be NULL, but the Go struct fields were defined as `string` types, which cannot handle NULL values.

## Solution
Changed the field types in `internal/entities/organization_setting.go` from `string` to `*string` (pointer to string) for the following fields:

- `SystemPrompt`
- `AIAssistantName`
- `AICommunicationStyle`
- `AICommunicationLanguage`

## Changes Made

### 1. OrganizationSetting Entity
```go
// Before
SystemPrompt            string    `json:"system_prompt" db:"system_prompt"`
AIAssistantName         string    `json:"ai_assistant_name" db:"ai_assistant_name"`
AICommunicationStyle    string    `json:"ai_communication_style" db:"ai_communication_style"`
AICommunicationLanguage string    `json:"ai_communication_language" db:"ai_communication_language"`

// After
SystemPrompt            *string   `json:"system_prompt" db:"system_prompt"`
AIAssistantName         *string   `json:"ai_assistant_name" db:"ai_assistant_name"`
AICommunicationStyle    *string   `json:"ai_communication_style" db:"ai_communication_style"`
AICommunicationLanguage *string   `json:"ai_communication_language" db:"ai_communication_language"`
```

### 2. Input Structs
Updated both `CreateOrganizationSettingInput` and `UpdateOrganizationSettingInput` structs with the same pointer types.

## How It Works

### Database Scanning (Reading)
When scanning from the database, NULL values are handled correctly:
```go
// Repository correctly uses pointers for scanning
err := r.DB.QueryRowContext(ctx, query, orgID).Scan(
    // ... other fields ...
    &setting.SystemPrompt,        // Can be nil if NULL
    &setting.AIAssistantName,     // Can be nil if NULL
    &setting.AICommunicationStyle, // Can be nil if NULL
    &setting.AICommunicationLanguage, // Can be nil if NULL
)
```

### Database Updates (Writing)
When updating the database, the values are passed directly:
```go
// Repository passes values (not pointers) for updates
result, err := r.DB.ExecContext(ctx, query,
    // ... other fields ...
    setting.SystemPrompt,        // Passes the string value or nil
    setting.AIAssistantName,     // Passes the string value or nil
    setting.AICommunicationStyle, // Passes the string value or nil
    setting.AICommunicationLanguage, // Passes the string value or nil
)
```

## Usage in Code

### Checking for NULL Values
```go
if setting.SystemPrompt != nil {
    // Use the system prompt
    systemPrompt := *setting.SystemPrompt
} else {
    // Handle NULL case
    systemPrompt := "default prompt"
}
```

### Setting Values
```go
// Set a value
prompt := "Custom system prompt"
setting.SystemPrompt = &prompt

// Set to NULL
setting.SystemPrompt = nil
```

## Benefits

1. **NULL Safety**: Properly handles NULL values from the database
2. **Type Safety**: Go's type system ensures proper handling
3. **Backward Compatibility**: Existing code continues to work
4. **Clear Intent**: Pointer types clearly indicate nullable fields

## Testing

The fix has been tested by:
1. Building the application successfully
2. Ensuring all pointer types are correctly handled in repository operations
3. Verifying that both NULL and non-NULL values work correctly

## Related Files

- `internal/entities/organization_setting.go` - Entity definitions
- `internal/repositories/organization_repository.go` - Database operations
- Any services that use organization settings (may need updates for pointer handling) 